// Copyright 2021 - 2022 Matrix Origin
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package plan

import (
	"context"
	"github.com/matrixorigin/matrixone/pkg/catalog"
	"github.com/matrixorigin/matrixone/pkg/common/moerr"
	"github.com/matrixorigin/matrixone/pkg/container/types"
	"github.com/matrixorigin/matrixone/pkg/pb/plan"
	"github.com/matrixorigin/matrixone/pkg/sql/parsers/dialect"
	"github.com/matrixorigin/matrixone/pkg/sql/parsers/tree"
)

func buildInsert(stmt *tree.Insert, ctx CompilerContext) (p *Plan, err error) {
	if stmt.OnDuplicateUpdate != nil {
		return nil, moerr.NewNotSupported(ctx.GetContext(), "INSERT ... ON DUPLICATE KEY UPDATE ...")
	}
	rows := stmt.Rows
	switch rows.Select.(type) {
	case *tree.ValuesClause:
		return buildInsertValues(stmt, ctx)
	case *tree.SelectClause, *tree.ParenSelect:
		return buildInsertSelect(stmt, ctx)
	default:
		return nil, moerr.NewInvalidInput(ctx.GetContext(), "insert has unknown select statement")
	}
}

func buildInsertValues(stmt *tree.Insert, ctx CompilerContext) (p *Plan, err error) {
	// get table source
	tbl, ok := stmt.Table.(*tree.TableName)
	if !ok {
		return nil, moerr.NewInvalidInput(ctx.GetContext(), "insert table is invalid '%s'", tree.String(stmt.Table, dialect.MYSQL))
	}
	tblName := string(tbl.ObjectName)
	dbName := string(tbl.SchemaName)
	if dbName == "" {
		dbName = ctx.DefaultDatabase()
	}
	_, tblRef := ctx.Resolve(dbName, tblName)
	if tblRef == nil {
		return nil, moerr.NewInvalidInput(ctx.GetContext(), "insert table is invalid '%s'", tree.String(stmt.Table, dialect.MYSQL))
	}
	if tblRef.TableType == catalog.SystemExternalRel {
		return nil, moerr.NewInvalidInput(ctx.GetContext(), "cannot insert into external table '%s'", tblName)
	} else if tblRef.TableType == catalog.SystemViewRel {
		return nil, moerr.NewInvalidInput(ctx.GetContext(), "cannot insert into view '%s'", tblName)
	}

	// build columns
	colCount := len(tblRef.Cols)

	hasExplicitCols := false
	if stmt.Columns != nil {
		hasExplicitCols = true
	}

	var explicitCols []*ColDef
	if stmt.Columns == nil {
		explicitCols = append(explicitCols, tblRef.Cols...)
	} else {
		for _, attr := range stmt.Columns {
			hasAttr := false
			for _, col := range tblRef.Cols {
				if string(attr) == col.Name {
					explicitCols = append(explicitCols, col)
					hasAttr = true
					break
				}
			}
			if !hasAttr {
				return nil, moerr.NewInvalidInput(ctx.GetContext(), "insert value into unknown column '%s'", string(attr))
			}
		}
	}
	explicitCount := len(explicitCols)

	orderAttrs := make([]string, 0, colCount)
	for _, col := range tblRef.Cols {
		orderAttrs = append(orderAttrs, col.Name)
	}

	var otherCols []*ColDef
	if len(explicitCols) < colCount {
		for _, c1 := range tblRef.Cols {
			hasCol := false
			for _, c2 := range explicitCols {
				if c1.Name == c2.Name {
					hasCol = true
					break
				}
			}
			if !hasCol {
				otherCols = append(otherCols, c1)
			}
		}
	}

	rows := stmt.Rows.Select.(*tree.ValuesClause).Rows
	isAllDefault := false
	if rows[0] == nil {
		isAllDefault = true
	}

	if isAllDefault && hasExplicitCols {
		return nil, moerr.NewInvalidInput(ctx.GetContext(), "insert values does not match number of columns")
	}

	rowCount := len(rows)
	columns := make([]*plan.Column, colCount)
	for i := range columns {
		columns[i] = &plan.Column{
			Column: make([]*plan.Expr, 0, rowCount),
		}
	}

	if isAllDefault {
		// hasExplicitCols must be false
		for _, row := range rows {
			if row != nil {
				return nil, moerr.NewInvalidInput(ctx.GetContext(), "insert values does not match number of columns")
			}
			// build column
			for j, col := range explicitCols {
				expr, err := getDefaultExpr(ctx.GetContext(), col)
				if err != nil {
					return nil, err
				}
				columns[j].Column = append(columns[j].Column, expr)
			}
		}
	} else {
		// hasExplicitCols maybe true or false
		binders := make([]*DefaultBinder, 0, len(explicitCols))
		for _, col := range explicitCols {
			binders = append(binders, NewDefaultBinder(ctx.GetContext(), nil, nil, col.Typ, nil))
		}
		for i, row := range rows {
			if row == nil || explicitCount != len(row) {
				return nil, moerr.NewInvalidInput(ctx.GetContext(), "insert values does not match the number of columns")
			}

			idx := 0
			for j, col := range explicitCols {
				if _, ok := row[idx].(*tree.DefaultVal); ok {
					expr, err := getDefaultExpr(ctx.GetContext(), col)
					if err != nil {
						return nil, err
					}
					columns[idx].Column = append(columns[idx].Column, expr)
				} else {
					planExpr, err := binders[j].BindExpr(row[idx], 0, false)
					if err != nil {
						err = MakeInsertError(ctx.GetContext(), types.T(col.Typ.Id), col, rows, j, i, err)
						return nil, err
					}
					resExpr, err := makePlan2CastExpr(ctx.GetContext(), planExpr, col.Typ)
					if err != nil {
						err = MakeInsertError(ctx.GetContext(), types.T(col.Typ.Id), col, rows, j, i, err)
						return nil, err
					}
					columns[idx].Column = append(columns[idx].Column, resExpr)
				}
				idx++
			}

			for _, col := range otherCols {
				expr, err := getDefaultExpr(ctx.GetContext(), col)
				if err != nil {
					return nil, err
				}
				columns[idx].Column = append(columns[idx].Column, expr)
				idx++
			}
		}
	}

	uDef, sDef := buildIndexDefs(tblRef.Defs)

	return &Plan{
		Plan: &plan.Plan_Ins{
			Ins: &plan.InsertValues{
				DbName:            dbName,
				TblName:           tblName,
				ExplicitCols:      explicitCols,
				OtherCols:         otherCols,
				OrderAttrs:        orderAttrs,
				Columns:           columns,
				CompositePkey:     tblRef.CompositePkey,
				UniqueIndexDef:    uDef,
				SecondaryIndexDef: sDef,
			},
		},
	}, nil
}

func MakeInsertError(ctx context.Context, id types.T, col *ColDef, rows []tree.Exprs, colIdx, rowIdx int, err error) error {
	if moerr.IsMoErrCode(err, moerr.ErrFileNotFound) {
		return err
	}
	var str string
	if rows[rowIdx] == nil || len(rows[rowIdx]) < colIdx {
		str = col.Default.OriginString
	} else if _, ok := rows[rowIdx][colIdx].(*tree.DefaultVal); ok {
		str = col.Default.OriginString
	} else {
		str = tree.String(rows[rowIdx][colIdx], dialect.MYSQL)
	}
	if id == types.T_json {
		return moerr.NewInvalidInput(ctx, "Invalid %s text: '%s' for column '%s' at row '%d'", id.String(), str, col.Name, rowIdx+1)
	}
	return moerr.NewTruncatedValueForField(ctx, id.String(), str, col.Name, rowIdx+1)
}

func SetPlanLoadTag(pn *Plan) {
	pn2, ok := pn.Plan.(*plan.Plan_Query)
	if !ok {
		return
	}
	nodes := pn2.Query.Nodes
	for i := 0; i < len(nodes); i++ {
		if nodes[i].NodeType == plan.Node_EXTERNAL_SCAN {
			pn2.Query.LoadTag = true
			return
		}
	}
}

func buildInsertSelect(stmt *tree.Insert, ctx CompilerContext) (p *Plan, err error) {
	pn, err := runBuildSelectByBinder(plan.Query_SELECT, ctx, stmt.Rows)
	if err != nil {
		return nil, err
	}
	SetPlanLoadTag(pn)
	cols := GetResultColumnsFromPlan(pn)
	pn.Plan.(*plan.Plan_Query).Query.StmtType = plan.Query_INSERT
	if len(stmt.Columns) != 0 && len(stmt.Columns) != len(cols) {
		return nil, moerr.NewInvalidInput(ctx.GetContext(), "insert statement column count does not match")
	}

	objRef, tableDef, err := getInsertTable(stmt.Table, ctx)
	if err != nil {
		return nil, err
	}
	if tableDef.TableType == catalog.SystemExternalRel {
		return nil, moerr.NewInvalidInput(ctx.GetContext(), "cannot insert into external table")
	} else if tableDef.TableType == catalog.SystemViewRel {
		return nil, moerr.NewInvalidInput(ctx.GetContext(), "cannot insert into view")
	}

	valueCount := len(stmt.Columns)
	if len(stmt.Columns) == 0 {
		valueCount = len(tableDef.Cols)
	}
	if valueCount != len(cols) {
		return nil, moerr.NewInvalidInput(ctx.GetContext(), "insert statement column count does not match value count")
	}

	// generate values expr
	exprs, err := getInsertExprs(ctx, stmt, cols, tableDef)
	if err != nil {
		return nil, err
	}

	// do type cast if needed
	for i := range tableDef.Cols {
		exprs[i], err = makePlan2CastExpr(ctx.GetContext(), exprs[i], tableDef.Cols[i].Typ)
		if err != nil {
			return nil, err
		}
	}
	qry := pn.Plan.(*plan.Plan_Query).Query
	n := &Node{
		ObjRef:      objRef,
		TableDef:    tableDef,
		NodeType:    plan.Node_INSERT,
		NodeId:      int32(len(qry.Nodes)),
		Children:    []int32{qry.Steps[len(qry.Steps)-1]},
		ProjectList: exprs,
	}
	appendQueryNode(qry, n)
	qry.Steps[len(qry.Steps)-1] = n.NodeId
	return pn, nil
}

func getInsertExprs(ctx CompilerContext, stmt *tree.Insert, cols []*ColDef, tableDef *TableDef) ([]*Expr, error) {
	var exprs []*Expr

	if len(stmt.Columns) == 0 {
		exprs = make([]*Expr, len(cols))
		for i := range exprs {
			exprs[i] = &plan.Expr{
				Typ: cols[i].Typ,
				Expr: &plan.Expr_Col{
					Col: &plan.ColRef{
						ColPos: int32(i),
					},
				},
			}
		}
	} else {
		exprs = make([]*Expr, len(tableDef.Cols))
		tableColMap := make(map[string]int)
		targetMap := make(map[string]int)
		for i, col := range stmt.Columns {
			targetMap[string(col)] = i
		}
		for i, col := range tableDef.Cols {
			tableColMap[col.GetName()] = i
		}
		// check if the column name is legal
		for k := range targetMap {
			if _, ok := tableColMap[k]; !ok {
				return nil, moerr.NewInvalidInput(ctx.GetContext(), "insert column '%s' does not exist", k)
			}
		}
		for i := range exprs {
			if ref, ok := targetMap[tableDef.Cols[i].GetName()]; ok {
				exprs[i] = &plan.Expr{
					Typ: cols[ref].Typ,
					Expr: &plan.Expr_Col{
						Col: &plan.ColRef{
							ColPos: int32(ref),
						},
					},
				}
			} else {
				var err error
				exprs[i], err = getDefaultExpr(ctx.GetContext(), tableDef.Cols[i])
				if err != nil {
					return nil, err
				}
			}
		}
	}
	return exprs, nil
}

func getInsertTable(stmt tree.TableExpr, ctx CompilerContext) (*ObjectRef, *TableDef, error) {
	switch tbl := stmt.(type) {
	case *tree.TableName:
		tblName := string(tbl.ObjectName)
		dbName := string(tbl.SchemaName)
		objRef, tableDef := ctx.Resolve(dbName, tblName)
		if tableDef == nil {
			return nil, nil, moerr.NewInvalidInput(ctx.GetContext(), "insert target table '%s' does not exist", tblName)
		}
		return objRef, tableDef, nil
	case *tree.ParenTableExpr:
		return getInsertTable(tbl.Expr, ctx)
	case *tree.AliasedTableExpr:
		return getInsertTable(tbl.Expr, ctx)
	case *tree.Select:
		return nil, nil, moerr.NewNotSupported(ctx.GetContext(), "insert table expr %v", stmt)
	case *tree.StatementSource:
		return nil, nil, moerr.NewNotSupported(ctx.GetContext(), "insert table expr %v", stmt)
	default:
		return nil, nil, moerr.NewNotSupported(ctx.GetContext(), "insert table expr %v", stmt)
	}
}

func buildIndexDefs(defs []*plan.TableDef_DefType) (*UniqueIndexDef, *SecondaryIndexDef) {
	var uIdxDef *UniqueIndexDef = nil
	var sIdxDef *SecondaryIndexDef = nil
	for _, def := range defs {
		if idxDef, ok := def.Def.(*plan.TableDef_DefType_UIdx); ok {
			uIdxDef = idxDef.UIdx
		}
		if idxDef, ok := def.Def.(*plan.TableDef_DefType_SIdx); ok {
			sIdxDef = idxDef.SIdx
		}
	}
	return uIdxDef, sIdxDef
}

// Copyright 2021 Matrix Origin
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

package compile

import (
	"context"
	"github.com/matrixorigin/matrixone/pkg/sql/util"

	"github.com/matrixorigin/matrixone/pkg/common/moerr"
	"github.com/matrixorigin/matrixone/pkg/compress"
	"github.com/matrixorigin/matrixone/pkg/container/types"
	"github.com/matrixorigin/matrixone/pkg/defines"
	"github.com/matrixorigin/matrixone/pkg/pb/plan"
	"github.com/matrixorigin/matrixone/pkg/sql/colexec"
	"github.com/matrixorigin/matrixone/pkg/util/trace"
	"github.com/matrixorigin/matrixone/pkg/vm/engine"
)

func (s *Scope) CreateDatabase(c *Compile) error {
	var span trace.Span
	c.ctx, span = trace.Start(c.ctx, "CreateDatabase")
	defer span.End()
	dbName := s.Plan.GetDdl().GetCreateDatabase().GetDatabase()
	if _, err := c.e.Database(c.ctx, dbName, c.proc.TxnOperator); err == nil {
		if s.Plan.GetDdl().GetCreateDatabase().GetIfNotExists() {
			return nil
		}
		return moerr.NewDBAlreadyExists(c.ctx, dbName)
	}
	err := c.e.Create(context.WithValue(c.ctx, defines.SqlKey{}, c.sql),
		dbName, c.proc.TxnOperator)
	if err != nil {
		return err
	}
	return colexec.CreateAutoIncrTable(c.e, c.ctx, c.proc, dbName)
}

func (s *Scope) DropDatabase(c *Compile) error {
	dbName := s.Plan.GetDdl().GetDropDatabase().GetDatabase()
	if _, err := c.e.Database(c.ctx, dbName, c.proc.TxnOperator); err != nil {
		if s.Plan.GetDdl().GetDropDatabase().GetIfExists() {
			return nil
		}
		return moerr.NewErrDropNonExistsDB(c.ctx, dbName)
	}
	return c.e.Delete(c.ctx, dbName, c.proc.TxnOperator)
}

func (s *Scope) CreateTable(c *Compile) error {
	qry := s.Plan.GetDdl().GetCreateTable()
	// convert the plan's cols to the execution's cols
	planCols := qry.GetTableDef().GetCols()
	tableCols := planCols
	exeCols := planColsToExeCols(planCols)

	// convert the plan's defs to the execution's defs
	exeDefs, err := planDefsToExeDefs(qry.GetTableDef())
	if err != nil {
		return err
	}

	dbName := c.db
	if qry.GetDatabase() != "" {
		dbName = qry.GetDatabase()
	}
	dbSource, err := c.e.Database(c.ctx, dbName, c.proc.TxnOperator)
	if err != nil {
		if dbName == "" {
			return moerr.NewNoDB(c.ctx)
		}
		return err
	}
	tblName := qry.GetTableDef().GetName()
	if _, err := dbSource.Relation(c.ctx, tblName); err == nil {
		if qry.GetIfNotExists() {
			return nil
		}
		return moerr.NewTableAlreadyExists(c.ctx, tblName)
	}
	if err := dbSource.Create(context.WithValue(c.ctx, defines.SqlKey{}, c.sql), tblName, append(exeCols, exeDefs...)); err != nil {
		return err
	}
	// build index table
	for _, def := range qry.IndexTables {
		planCols = def.GetCols()
		exeCols = planColsToExeCols(planCols)
		exeDefs, err = planDefsToExeDefs(def)
		if err != nil {
			return err
		}
		if _, err := dbSource.Relation(c.ctx, def.Name); err == nil {
			return moerr.NewTableAlreadyExists(c.ctx, def.Name)
		}
		if err := dbSource.Create(c.ctx, def.Name, append(exeCols, exeDefs...)); err != nil {
			return err
		}
	}
	return colexec.CreateAutoIncrCol(c.e, c.ctx, dbSource, c.proc, tableCols, dbName, tblName)
}

func (s *Scope) CreateIndex(c *Compile) error {
	qry := s.Plan.GetDdl().GetCreateIndex()
	d, err := c.e.Database(c.ctx, qry.Database, c.proc.TxnOperator)
	if err != nil {
		return err
	}
	r, err := d.Relation(c.ctx, qry.Table)
	if err != nil {
		return err
	}

	// build and create index table
	def := qry.GetIndex().GetIndexTables()[0]
	planCols := def.GetCols()
	exeCols := planColsToExeCols(planCols)
	exeDefs, err := planDefsToExeDefs(def)
	if err != nil {
		return err
	}
	if _, err := d.Relation(c.ctx, def.Name); err == nil {
		return moerr.NewTableAlreadyExists(c.ctx, def.Name)
	}
	if err := d.Create(c.ctx, def.Name, append(exeCols, exeDefs...)); err != nil {
		return err
	}

	// build and update constraint def
	defs, err := planDefsToExeDefs(qry.GetIndex().GetTableDef())
	if err != nil {
		return err
	}
	ct := defs[0].(*engine.ConstraintDef)

	tblDefs, err := r.TableDefs(c.ctx)
	if err != nil {
		return err
	}
	var oldCt *engine.ConstraintDef
	for _, def := range tblDefs {
		if ct, ok := def.(*engine.ConstraintDef); ok {
			oldCt = ct
			break
		}
	}
	newCt, err := makeNewCreateConstraint(oldCt, ct.Cts[0])
	if err != nil {
		return err
	}
	err = r.UpdateConstraint(c.ctx, newCt)
	if err != nil {
		return err
	}

	// TODO: implement by insert ... select ...
	// insert data into index table
	switch t := qry.GetIndex().GetTableDef().Defs[0].Def.(type) {
	case *plan.TableDef_DefType_UIdx:
		ret, err := r.Ranges(c.ctx, nil)
		if err != nil {
			return err
		}
		rds, err := r.NewReader(c.ctx, 1, nil, ret)
		if err != nil {
			return err
		}
		bat, err := rds[0].Read(c.ctx, t.UIdx.Fields[0].Parts, nil, c.proc.Mp())
		if err != nil {
			return err
		}
		err = rds[0].Close()
		if err != nil {
			return err
		}
		indexBat, cnt := util.BuildUniqueKeyBatch(bat.Vecs, t.UIdx.Fields[0].Parts, t.UIdx.Fields[0].Cols, c.proc)
		indexR, err := d.Relation(c.ctx, t.UIdx.TableNames[0])
		if err != nil {
			return err
		}
		if cnt != 0 {
			if err := indexR.Write(c.ctx, indexBat); err != nil {
				return err
			}
		}
		indexBat.Clean(c.proc.Mp())
		// other situation is not supported now and check in plan
	}

	return nil
}

func (s *Scope) DropIndex(c *Compile) error {
	qry := s.Plan.GetDdl().GetDropIndex()
	d, err := c.e.Database(c.ctx, qry.Database, c.proc.TxnOperator)
	if err != nil {
		return err
	}
	r, err := d.Relation(c.ctx, qry.Table)
	if err != nil {
		return err
	}

	// build and update constraint def
	tblDefs, err := r.TableDefs(c.ctx)
	if err != nil {
		return err
	}
	var oldCt *engine.ConstraintDef
	for _, def := range tblDefs {
		if ct, ok := def.(*engine.ConstraintDef); ok {
			oldCt = ct
			break
		}
	}
	newCt, err := makeNewDropConstraint(oldCt, qry.GetIndexName())
	if err != nil {
		return err
	}
	err = r.UpdateConstraint(c.ctx, newCt)
	if err != nil {
		return err
	}

	// drop index table
	if _, err = d.Relation(c.ctx, qry.IndexTableName); err != nil {
		return err
	}
	if err = d.Delete(c.ctx, qry.IndexTableName); err != nil {
		return err
	}
	return nil
}

// TODO:
// func makeNewUpdateConstraint()
func makeNewDropConstraint(oldCt *engine.ConstraintDef, dropName string) (*engine.ConstraintDef, error) {
	// must fount dropName because of being checked in plan
	for j, ct := range oldCt.Cts {
		switch c := ct.(type) {
		case *engine.UniqueIndexDef:
			u := &plan.UniqueIndexDef{}
			err := u.UnMarshalUniqueIndexDef(([]byte)(c.UniqueIndex))
			if err != nil {
				return nil, err
			}
			for i, name := range u.IndexNames {
				if dropName == name {
					// If all indexes of a table are not defined in plan.UniqueIndexDef, the code will be much simpler
					u.IndexNames = append(u.IndexNames[:i], u.IndexNames[i+1:]...)
					u.TableNames = append(u.TableNames[:i], u.TableNames[i+1:]...)
					u.Fields = append(u.Fields[:i], u.Fields[i+1:]...)
					u.TableExists = append(u.TableExists[:i], u.TableExists[i+1:]...)

					oldCt.Cts = append(oldCt.Cts[:j], oldCt.Cts[j+1:]...)
					b, err := u.MarshalUniqueIndexDef()
					if err != nil {
						return nil, err
					}
					oldCt.Cts = append(oldCt.Cts, &engine.SecondaryIndexDef{
						SecondaryIndex: string(b),
					})
					break
				}
			}

		}
	}
	return oldCt, nil
}

func makeNewCreateConstraint(oldCt *engine.ConstraintDef, c engine.Constraint) (*engine.ConstraintDef, error) {
	// duplication has checked in plan
	if oldCt == nil {
		return &engine.ConstraintDef{
			Cts: []engine.Constraint{c},
		}, nil
	}
	switch t := c.(type) {
	case *engine.UniqueIndexDef:
		d := &plan.UniqueIndexDef{}
		err := d.UnMarshalUniqueIndexDef([]byte(t.UniqueIndex))
		if err != nil {
			return nil, err
		}

		ok := false
		var idx *engine.UniqueIndexDef
		for i, ct := range oldCt.Cts {
			if idx, ok = ct.(*engine.UniqueIndexDef); ok {
				u := &plan.UniqueIndexDef{}
				err := u.UnMarshalUniqueIndexDef([]byte(idx.UniqueIndex))
				if err != nil {
					return nil, err
				}
				u.IndexNames = append(u.IndexNames, d.IndexNames[0])
				u.TableNames = append(u.TableNames, d.TableNames[0])
				u.TableExists = append(u.TableExists, d.TableExists[0])
				u.Fields = append(u.Fields, d.Fields[0])

				oldCt.Cts = append(oldCt.Cts[:i], oldCt.Cts[i+1:]...)

				bytes, err := u.MarshalUniqueIndexDef()
				if err != nil {
					return nil, err
				}
				oldCt.Cts = append(oldCt.Cts, &engine.UniqueIndexDef{
					UniqueIndex: string(bytes),
				})
				break
			}
		}
		if !ok {
			oldCt.Cts = append(oldCt.Cts, c)
		}
	}
	return oldCt, nil
}

// Truncation operations cannot be performed if the session holds an active table lock.
func (s *Scope) TruncateTable(c *Compile) error {
	tqry := s.Plan.GetDdl().GetTruncateTable()
	dbName := tqry.GetDatabase()
	dbSource, err := c.e.Database(c.ctx, dbName, c.proc.TxnOperator)
	if err != nil {
		return err
	}
	tblName := tqry.GetTable()
	var rel engine.Relation
	if rel, err = dbSource.Relation(c.ctx, tblName); err != nil {
		return err
	}
	id := rel.GetTableID(c.ctx)
	err = dbSource.Truncate(c.ctx, tblName)
	if err != nil {
		return err
	}

	// Truncate Index Tables if needed
	for _, name := range tqry.IndexTableNames {
		err := dbSource.Truncate(c.ctx, name)
		if err != nil {
			return err
		}
	}

	err = colexec.ResetAutoInsrCol(c.e, c.ctx, tblName, dbSource, c.proc, id, dbName)
	if err != nil {
		return err
	}
	return nil
}

func (s *Scope) DropTable(c *Compile) error {
	qry := s.Plan.GetDdl().GetDropTable()

	dbName := qry.GetDatabase()
	dbSource, err := c.e.Database(c.ctx, dbName, c.proc.TxnOperator)
	if err != nil {
		if qry.GetIfExists() {
			return nil
		}
		return err
	}
	tblName := qry.GetTable()
	var rel engine.Relation
	if rel, err = dbSource.Relation(c.ctx, tblName); err != nil {
		if qry.GetIfExists() {
			return nil
		}
		return err
	}
	if err := dbSource.Delete(c.ctx, tblName); err != nil {
		return err
	}
	for _, name := range qry.IndexTableNames {
		if err := dbSource.Delete(c.ctx, name); err != nil {
			return err
		}
	}
	return colexec.DeleteAutoIncrCol(c.e, c.ctx, rel, c.proc, dbName, rel.GetTableID(c.ctx))
}

func planDefsToExeDefs(tableDef *plan.TableDef) ([]engine.TableDef, error) {
	planDefs := tableDef.GetDefs()
	var exeDefs []engine.TableDef
	c := new(engine.ConstraintDef)
	for _, def := range planDefs {
		switch defVal := def.GetDef().(type) {
		case *plan.TableDef_DefType_Cb:
			exeDefs = append(exeDefs, &engine.ClusterByDef{
				Name: defVal.Cb.Name,
			})
		case *plan.TableDef_DefType_Pk:
			exeDefs = append(exeDefs, &engine.PrimaryIndexDef{
				Names: defVal.Pk.GetNames(),
			})
		case *plan.TableDef_DefType_Properties:
			properties := make([]engine.Property, len(defVal.Properties.GetProperties()))
			for i, p := range defVal.Properties.GetProperties() {
				properties[i] = engine.Property{
					Key:   p.GetKey(),
					Value: p.GetValue(),
				}
			}
			exeDefs = append(exeDefs, &engine.PropertiesDef{
				Properties: properties,
			})
		case *plan.TableDef_DefType_Partition:
			bytes, err := defVal.Partition.MarshalPartitionInfo()
			if err != nil {
				return nil, err
			}
			exeDefs = append(exeDefs, &engine.PartitionDef{
				Partition: string(bytes),
			})
		case *plan.TableDef_DefType_UIdx:
			bytes, err := defVal.UIdx.MarshalUniqueIndexDef()
			if err != nil {
				return nil, err
			}
			c.Cts = append(c.Cts, &engine.UniqueIndexDef{
				UniqueIndex: string(bytes),
			})
		case *plan.TableDef_DefType_SIdx:
			bytes, err := defVal.SIdx.MarshalSecondaryIndexDef()
			if err != nil {
				return nil, err
			}
			c.Cts = append(c.Cts, &engine.SecondaryIndexDef{
				SecondaryIndex: string(bytes),
			})
		}
	}

	if tableDef.ViewSql != nil {
		exeDefs = append(exeDefs, &engine.ViewDef{
			View: tableDef.ViewSql.View,
		})
	}

	if len(c.Cts) > 0 {
		exeDefs = append(exeDefs, c)
	}
	return exeDefs, nil
}

func planColsToExeCols(planCols []*plan.ColDef) []engine.TableDef {
	exeCols := make([]engine.TableDef, len(planCols))
	for i, col := range planCols {
		var alg compress.T
		switch col.Alg {
		case plan.CompressType_None:
			alg = compress.None
		case plan.CompressType_Lz4:
			alg = compress.Lz4
		}
		colTyp := col.GetTyp()
		exeCols[i] = &engine.AttributeDef{
			Attr: engine.Attribute{
				Name: col.Name,
				Alg:  alg,
				Type: types.Type{
					Oid:       types.T(colTyp.GetId()),
					Width:     colTyp.GetWidth(),
					Precision: colTyp.GetPrecision(),
					Scale:     colTyp.GetScale(),
					Size:      colTyp.GetSize(),
				},
				Default:       planCols[i].GetDefault(),
				OnUpdate:      planCols[i].GetOnUpdate(),
				Primary:       col.GetPrimary(),
				Comment:       col.GetComment(),
				ClusterBy:     col.ClusterBy,
				AutoIncrement: col.Typ.GetAutoIncr(),
			},
		}
	}
	return exeCols
}

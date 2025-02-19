// Copyright 2022 Matrix Origin
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

package disttae

import (
	"context"
	"math/rand"
	"strconv"

	"github.com/google/uuid"
	"github.com/matrixorigin/matrixone/pkg/catalog"
	"github.com/matrixorigin/matrixone/pkg/container/batch"
	"github.com/matrixorigin/matrixone/pkg/container/types"
	"github.com/matrixorigin/matrixone/pkg/container/vector"
	"github.com/matrixorigin/matrixone/pkg/pb/plan"
	"github.com/matrixorigin/matrixone/pkg/txn/storage/memorystorage/memtable"
	"github.com/matrixorigin/matrixone/pkg/vm/engine"
)

var _ engine.Relation = new(table)

func (tbl *table) FilteredStats(ctx context.Context, expr *plan.Expr) (int32, int64, error) {
	switch tbl.tableId {
	case catalog.MO_DATABASE_ID:
		return 1, 1000, nil
	case catalog.MO_TABLES_ID:
		return 10, 10000, nil
	case catalog.MO_COLUMNS_ID:
		return 10, 10000, nil
	}

	if expr == nil {
		return tbl.Stats(ctx)
	}
	var blockNum, totalBlockCnt int
	var outcnt int64

	exprMono := checkExprIsMonotonic(ctx, expr)
	columnMap, columns, maxCol := getColumnsByExpr(expr, tbl.getTableDef())

	for _, blockmetas := range tbl.meta.blocks {
		totalBlockCnt += len(blockmetas)
		for _, blk := range blockmetas {
			if !exprMono || needRead(ctx, expr, blk, tbl.getTableDef(), columnMap, columns, maxCol, tbl.db.txn.proc) {
				outcnt += blockRows(blk)
				blockNum++
			}
		}
	}
	// before first execution, no metadata.
	if totalBlockCnt == 0 {
		return 100, 1000000, nil
	}
	return int32(blockNum), outcnt, nil
}

func (tbl *table) Stats(ctx context.Context) (int32, int64, error) {
	var rows int64
	var totalBlockCnt int
	for _, blks := range tbl.meta.blocks {
		totalBlockCnt += len(blks)
		for _, blk := range blks {
			rows += blockRows(blk)
		}
	}
	// before first execution, no metadata.
	if totalBlockCnt == 0 {
		return 100, 1000000, nil
	}
	return int32(totalBlockCnt), rows, nil
}

func (tbl *table) Rows(ctx context.Context) (int64, error) {
	var rows int64
	writes := make([]Entry, 0, len(tbl.db.txn.writes))
	tbl.db.txn.Lock()
	for i := range tbl.db.txn.writes {
		for _, entry := range tbl.db.txn.writes[i] {
			if entry.databaseId == tbl.db.databaseId &&
				entry.tableId == tbl.tableId {
				writes = append(writes, entry)
			}
		}
	}
	tbl.db.txn.Unlock()

	deletes := make(map[types.Rowid]uint8)
	for _, entry := range writes {
		if entry.typ == INSERT {
			rows = rows + int64(entry.bat.Length())
		} else {
			if entry.bat.GetVector(0).GetType().Oid == types.T_Rowid {
				vs := vector.MustTCols[types.Rowid](entry.bat.GetVector(0))
				for _, v := range vs {
					deletes[v] = 0
				}
			}
		}
	}

	t := memtable.Time{
		Timestamp: tbl.db.txn.meta.SnapshotTS,
	}
	tx := memtable.NewTransaction(
		uuid.NewString(),
		t,
		memtable.SnapshotIsolation,
	)
	for _, partition := range tbl.parts {
		pRows, err := partition.Rows(tx, deletes, tbl.skipBlocks)
		if err != nil {
			return 0, err
		}
		rows = rows + pRows
	}

	if tbl.meta == nil {
		return rows, nil
	}
	for _, blks := range tbl.meta.blocks {
		for _, blk := range blks {
			rows += blockRows(blk)
		}
	}
	return rows, nil
}

func (tbl *table) Size(ctx context.Context, name string) (int64, error) {
	// TODO
	return 0, nil
}

// return all unmodified blocks
func (tbl *table) Ranges(ctx context.Context, expr *plan.Expr) ([][]byte, error) {
	/*
		// consider halloween problem
		if int64(tbl.db.txn.statementId)-1 > 0 {
			writes = tbl.db.txn.writes[:tbl.db.txn.statementId-1]
		}
	*/
	writes := make([]Entry, 0, len(tbl.db.txn.writes))
	for i := range tbl.db.txn.writes {
		for _, entry := range tbl.db.txn.writes[i] {
			if entry.databaseId == tbl.db.databaseId &&
				entry.tableId == tbl.tableId {
				writes = append(writes, entry)
			}
		}
	}
	priKeys := make([]*engine.Attribute, 0, 1)
	if tbl.primaryIdx >= 0 {
		for _, def := range tbl.defs {
			if attr, ok := def.(*engine.AttributeDef); ok {
				if attr.Attr.Primary {
					priKeys = append(priKeys, &attr.Attr)
				}
			}
		}
	}
	dnList := needSyncDnStores(ctx, expr, tbl.tableDef, priKeys, tbl.db.txn.dnStores, tbl.db.txn.proc)
	switch {
	case tbl.tableId == catalog.MO_DATABASE_ID:
		tbl.dnList = []int{0}
	case tbl.tableId == catalog.MO_TABLES_ID:
		tbl.dnList = []int{0}
	case tbl.tableId == catalog.MO_COLUMNS_ID:
		tbl.dnList = []int{0}
	default:
		tbl.dnList = dnList
	}
	dnStores := make([]DNStore, 0, len(dnList))
	for _, i := range dnList {
		dnStores = append(dnStores, tbl.db.txn.dnStores[i])
	}
	_, ok := tbl.db.txn.createTableMap[tbl.tableId]
	if !ok && !tbl.updated {
		if err := tbl.db.txn.db.Update(ctx, dnStores, tbl, tbl.db.txn.op, tbl.primaryIdx,
			tbl.db.databaseId, tbl.tableId, tbl.db.txn.meta.SnapshotTS); err != nil {
			return nil, err
		}
		// update meta
		columnLength := len(tbl.getTableDef().Cols) - 1 //we use this data to fetch zonemap, but row_id has no zonemap
		meta, err := tbl.db.txn.getTableMeta(ctx, tbl.db.databaseId,
			genMetaTableName(tbl.tableId), !ok, columnLength)
		if err != nil {
			return nil, err
		}
		tbl.meta = meta
		tbl.updated = true
	}
	ranges := make([][]byte, 0, 1)
	ranges = append(ranges, []byte{})
	tbl.skipBlocks = make(map[uint64]uint8)
	if tbl.meta == nil {
		return ranges, nil
	}
	tbl.meta.modifedBlocks = make([][]ModifyBlockMeta, len(tbl.meta.blocks))

	exprMono := checkExprIsMonotonic(tbl.db.txn.proc.Ctx, expr)
	columnMap, columns, maxCol := getColumnsByExpr(expr, tbl.getTableDef())
	for _, i := range dnList {
		blks, deletes := tbl.parts[i].BlockList(ctx, tbl.db.txn.meta.SnapshotTS,
			tbl.meta.blocks[i], writes)
		for _, blk := range blks {
			tbl.skipBlocks[blk.Info.BlockID] = 0
			if !exprMono || needRead(ctx, expr, blk, tbl.getTableDef(), columnMap, columns, maxCol, tbl.db.txn.proc) {
				ranges = append(ranges, blockMarshal(blk))
			}
		}
		tbl.meta.modifedBlocks[i] = genModifedBlocks(ctx, deletes,
			tbl.meta.blocks[i], blks, expr, tbl.getTableDef(), tbl.db.txn.proc)
	}
	return ranges, nil
}

// getTableDef only return all cols and their index.
func (tbl *table) getTableDef() *plan.TableDef {
	if tbl.tableDef == nil {
		var cols []*plan.ColDef
		i := int32(0)
		name2index := make(map[string]int32)
		for _, def := range tbl.defs {
			if attr, ok := def.(*engine.AttributeDef); ok {
				name2index[attr.Attr.Name] = i
				cols = append(cols, &plan.ColDef{
					Name: attr.Attr.Name,
					Typ: &plan.Type{
						Id:        int32(attr.Attr.Type.Oid),
						Width:     attr.Attr.Type.Width,
						Size:      attr.Attr.Type.Size,
						Precision: attr.Attr.Type.Precision,
						Scale:     attr.Attr.Type.Scale,
						AutoIncr:  attr.Attr.AutoIncrement,
					},
					Primary:   attr.Attr.Primary,
					Default:   attr.Attr.Default,
					OnUpdate:  attr.Attr.OnUpdate,
					Comment:   attr.Attr.Comment,
					ClusterBy: attr.Attr.ClusterBy,
				})
				i++
			}
		}
		tbl.tableDef = &plan.TableDef{
			Name:          tbl.tableName,
			Cols:          cols,
			Name2ColIndex: name2index,
		}
	}
	return tbl.tableDef
}

func (tbl *table) TableDefs(ctx context.Context) ([]engine.TableDef, error) {
	//return tbl.defs, nil
	// I don't understand why the logic now is not to get all the tableDef. Don't understand.
	// copy from tae's logic
	defs := make([]engine.TableDef, 0, len(tbl.defs))
	if tbl.comment != "" {
		commentDef := new(engine.CommentDef)
		commentDef.Comment = tbl.comment
		defs = append(defs, commentDef)
	}
	if tbl.partition != "" {
		partitionDef := new(engine.PartitionDef)
		partitionDef.Partition = tbl.partition
		defs = append(defs, partitionDef)
	}

	if tbl.viewdef != "" {
		viewDef := new(engine.ViewDef)
		viewDef.View = tbl.viewdef
		defs = append(defs, viewDef)
	}
	if len(tbl.constraint) > 0 {
		c := &engine.ConstraintDef{}
		err := c.UnmarshalBinary(tbl.constraint)
		if err != nil {
			return nil, err
		}
		defs = append(defs, c)
	}
	for i, def := range tbl.defs {
		if attr, ok := def.(*engine.AttributeDef); ok {
			if attr.Attr.Name != catalog.Row_ID {
				defs = append(defs, tbl.defs[i])
			}

		}
	}
	pro := new(engine.PropertiesDef)
	pro.Properties = append(pro.Properties, engine.Property{
		Key:   catalog.SystemRelAttr_Kind,
		Value: string(tbl.relKind),
	})
	if tbl.createSql != "" {
		pro.Properties = append(pro.Properties, engine.Property{
			Key:   catalog.SystemRelAttr_CreateSQL,
			Value: tbl.createSql,
		})
	}
	defs = append(defs, pro)
	return defs, nil

}

func (tbl *table) UpdateConstraint(ctx context.Context, c *engine.ConstraintDef) error {
	ct, err := c.MarshalBinary()
	if err != nil {
		return err
	}
	bat, err := genTableConstraintTuple(tbl.tableId, tbl.db.databaseId, tbl.tableName, tbl.db.databaseName, ct, tbl.db.txn.proc.Mp())
	if err != nil {
		return err
	}
	if err = tbl.db.txn.WriteBatch(UPDATE, catalog.MO_CATALOG_ID, catalog.MO_TABLES_ID,
		catalog.MO_CATALOG, catalog.MO_TABLES, bat, tbl.db.txn.dnStores[0], -1); err != nil {
		return err
	}
	tbl.constraint = ct
	return nil
}

func (tbl *table) TableColumns(ctx context.Context) ([]*engine.Attribute, error) {
	var attrs []*engine.Attribute
	for _, def := range tbl.defs {
		if attr, ok := def.(*engine.AttributeDef); ok {
			attrs = append(attrs, &attr.Attr)
		}
	}
	return attrs, nil

}

func (tbl *table) GetPrimaryKeys(ctx context.Context) ([]*engine.Attribute, error) {
	attrs := make([]*engine.Attribute, 0, 1)
	for _, def := range tbl.defs {
		if attr, ok := def.(*engine.AttributeDef); ok {
			if attr.Attr.Primary {
				attrs = append(attrs, &attr.Attr)
			}
		}
	}
	return attrs, nil
}

func (tbl *table) GetHideKeys(ctx context.Context) ([]*engine.Attribute, error) {
	attrs := make([]*engine.Attribute, 0, 1)
	attrs = append(attrs, &engine.Attribute{
		IsHidden: true,
		IsRowId:  true,
		Name:     catalog.Row_ID,
		Type:     types.New(types.T_Rowid, 0, 0, 0),
		Primary:  true,
	})
	return attrs, nil
}

func (tbl *table) Write(ctx context.Context, bat *batch.Batch) error {
	if tbl.insertExpr == nil {
		ibat := batch.New(true, bat.Attrs)
		for j := range bat.Vecs {
			ibat.SetVector(int32(j), vector.New(bat.GetVector(int32(j)).GetType()))
		}
		if _, err := ibat.Append(ctx, tbl.db.txn.proc.Mp(), bat); err != nil {
			return err
		}
		i := rand.Int() % len(tbl.db.txn.dnStores)
		return tbl.db.txn.WriteBatch(INSERT, tbl.db.databaseId, tbl.tableId,
			tbl.db.databaseName, tbl.tableName, ibat, tbl.db.txn.dnStores[i], tbl.primaryIdx)
	}
	bats, err := partitionBatch(bat, tbl.insertExpr, tbl.db.txn.proc, len(tbl.parts))
	if err != nil {
		return err
	}
	for i := range bats {
		if bats[i].Length() == 0 {
			continue
		}
		if err := tbl.db.txn.WriteBatch(INSERT, tbl.db.databaseId, tbl.tableId,
			tbl.db.databaseName, tbl.tableName, bats[i], tbl.db.txn.dnStores[i], tbl.primaryIdx); err != nil {
			return err
		}
	}
	return nil
}

func (tbl *table) Update(ctx context.Context, bat *batch.Batch) error {
	return nil
}

func (tbl *table) Delete(ctx context.Context, bat *batch.Batch, name string) error {

	bat.SetAttributes([]string{catalog.Row_ID})
	bat = tbl.db.txn.deleteBatch(bat, tbl.db.databaseId, tbl.tableId)
	if bat.Length() == 0 {
		return nil
	}
	bats, err := partitionDeleteBatch(tbl, bat)
	if err != nil {
		return err
	}
	for i := range bats {
		if bats[i].Length() == 0 {
			continue
		}
		if err := tbl.db.txn.WriteBatch(DELETE, tbl.db.databaseId, tbl.tableId,
			tbl.db.databaseName, tbl.tableName, bats[i], tbl.db.txn.dnStores[i], tbl.primaryIdx); err != nil {
			return err
		}
	}
	return nil
}

func (tbl *table) AddTableDef(ctx context.Context, def engine.TableDef) error {
	return nil
}

func (tbl *table) DelTableDef(ctx context.Context, def engine.TableDef) error {
	return nil
}

func (tbl *table) GetTableID(ctx context.Context) string {
	return strconv.FormatUint(tbl.tableId, 10)
}

func (tbl *table) NewReader(ctx context.Context, num int, expr *plan.Expr, ranges [][]byte) ([]engine.Reader, error) {
	if len(ranges) == 0 {
		return tbl.newMergeReader(ctx, num, expr)
	}
	if len(ranges) == 1 && len(ranges[0]) == 0 {
		return tbl.newMergeReader(ctx, num, expr)
	}
	rds := make([]engine.Reader, num)
	blks := make([]BlockMeta, len(ranges))
	for i := range ranges {
		blks[i] = blockUnmarshal(ranges[i])
	}
	ts := tbl.db.txn.meta.SnapshotTS
	tableDef := tbl.getTableDef()

	if len(ranges) < num {
		for i := range ranges {
			rds[i] = &blockReader{
				fs:         tbl.db.fs,
				tableDef:   tableDef,
				primaryIdx: tbl.primaryIdx,
				expr:       expr,
				ts:         ts,
				ctx:        ctx,
				blks:       []BlockMeta{blks[i]},
			}
		}
		for j := len(ranges); j < num; j++ {
			rds[j] = &emptyReader{}
		}
		return rds, nil
	}
	step := (len(ranges)) / num
	if step < 1 {
		step = 1
	}
	for i := 0; i < num; i++ {
		if i == num-1 {
			rds[i] = &blockReader{
				fs:         tbl.db.fs,
				tableDef:   tableDef,
				primaryIdx: tbl.primaryIdx,
				expr:       expr,
				ts:         ts,
				ctx:        ctx,
				blks:       blks[i*step:],
			}
		} else {
			rds[i] = &blockReader{
				fs:         tbl.db.fs,
				tableDef:   tableDef,
				primaryIdx: tbl.primaryIdx,
				expr:       expr,
				ts:         ts,
				ctx:        ctx,
				blks:       blks[i*step : (i+1)*step],
			}
		}
	}
	return rds, nil
}

func (tbl *table) newMergeReader(ctx context.Context, num int,
	expr *plan.Expr) ([]engine.Reader, error) {
	var index memtable.Tuple
	/*
		// consider halloween problem
		if int64(tbl.db.txn.statementId)-1 > 0 {
			writes = tbl.db.txn.writes[:tbl.db.txn.statementId-1]
		}
	*/
	if tbl.primaryIdx >= 0 && expr != nil {
		pkColumn := tbl.tableDef.Cols[tbl.primaryIdx]
		ok, v := getPkValueByExpr(expr, pkColumn.Name, types.T(pkColumn.Typ.Id))
		if ok {
			index = memtable.Tuple{
				index_PrimaryKey,
				memtable.ToOrdered(v),
			}
		}
	}
	writes := make([]Entry, 0, len(tbl.db.txn.writes))
	tbl.db.txn.Lock()
	for i := range tbl.db.txn.writes {
		for _, entry := range tbl.db.txn.writes[i] {
			if entry.databaseId == tbl.db.databaseId &&
				entry.tableId == tbl.tableId {
				writes = append(writes, entry)
			}
		}
	}
	tbl.db.txn.Unlock()
	rds := make([]engine.Reader, num)
	mrds := make([]mergeReader, num)
	for _, i := range tbl.dnList {
		var blks []ModifyBlockMeta

		if tbl.meta != nil {
			blks = tbl.meta.modifedBlocks[i]
		}
		rds0, err := tbl.parts[i].NewReader(ctx, num, index, tbl.defs, tbl.tableDef,
			tbl.skipBlocks, blks, tbl.db.txn.meta.SnapshotTS, tbl.db.fs, writes)
		if err != nil {
			return nil, err
		}
		for i := range rds0 {
			mrds[i].rds = append(mrds[i].rds, rds0[i])
		}
	}
	for i := range rds {
		rds[i] = &mrds[i]
	}
	return rds, nil
}

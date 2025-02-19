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

import "github.com/matrixorigin/matrixone/pkg/txn/storage/memorystorage/memtable"

const (
	index_Time               = memtable.Text("time")
	index_Table              = memtable.Text("table")
	index_Column             = memtable.Text("column")
	index_PrimaryKey         = memtable.Text("primary key")
	index_BlockID_Time_OP    = memtable.Text("block id, time, op")
	index_TableID_PrimaryKey = memtable.Text("table id, primary key")
	index_Database           = memtable.Text("database")
)

type ColumnsIndexDef struct {
	Name    memtable.Text
	Columns []int
}

func NewColumnsIndexDef(name memtable.Text, cols ...int) ColumnsIndexDef {
	return ColumnsIndexDef{
		Name:    name,
		Columns: cols,
	}
}

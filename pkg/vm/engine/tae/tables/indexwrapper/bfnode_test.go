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

package indexwrapper

import (
	"context"
	"fmt"
	"path"
	"testing"

	"github.com/matrixorigin/matrixone/pkg/container/types"
	"github.com/matrixorigin/matrixone/pkg/defines"
	"github.com/matrixorigin/matrixone/pkg/fileservice"
	"github.com/matrixorigin/matrixone/pkg/objectio"
	"github.com/matrixorigin/matrixone/pkg/vm/engine/tae/common"
	"github.com/matrixorigin/matrixone/pkg/vm/engine/tae/containers"
	"github.com/matrixorigin/matrixone/pkg/vm/engine/tae/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStaticFilterIndex(t *testing.T) {
	defer testutils.AfterTest(t)()
	//bufManager := buffer.NewNodeManager(1024*1024, nil)
	var err error
	//var res bool
	//var exist bool
	//var ans *roaring.Bitmap

	// var res bool
	dir := testutils.InitTestEnv(ModuleName, t)
	dir = path.Join(dir, "/local")
	id := 1
	name := fmt.Sprintf("%d.blk", id)
	bat := newBatch()
	c := fileservice.Config{
		Name:    defines.LocalFileServiceName,
		Backend: "DISK",
		DataDir: dir,
	}
	service, err := fileservice.NewFileService(c)
	assert.Nil(t, err)

	objectWriter, err := objectio.NewObjectWriter(name, service)
	assert.Nil(t, err)
	/*fd*/ block, err := objectWriter.Write(bat)
	assert.Nil(t, err)

	cType := common.Plain
	typ := types.Type{Oid: types.T_int32}
	colIdx := uint16(0)
	interIdx := uint16(0)

	writer := NewBFWriter()
	err = writer.Init(objectWriter, block, cType, colIdx, interIdx)
	require.NoError(t, err)

	keys := containers.MockVector2(typ, 1000, 0)
	err = writer.AddValues(keys)
	require.NoError(t, err)

	_, err = writer.Finalize()
	require.NoError(t, err)
	blocks, err := objectWriter.WriteEnd(context.Background())
	assert.Nil(t, err)
	assert.Equal(t, 1, len(blocks))

	/*reader := NewBFReader(bufManager, file, new(common.ID))

	res, err = reader.MayContainsKey(int32(500))
	require.NoError(t, err)
	require.True(t, res)

	res, err = reader.MayContainsKey(int32(2000))
	require.NoError(t, err)
	require.False(t, res)

	query := containers.MockVector2(typ, 1000, 1500)
	exist, ans, err = reader.MayContainsAnyKeys(query, nil)
	require.NoError(t, err)
	require.True(t, ans.GetCardinality() < uint64(10))
	require.True(t, exist)*/

	//t.Log(bufManager.String())
}

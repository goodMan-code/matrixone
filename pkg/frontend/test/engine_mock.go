// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/vm/engine/types.go

// Package mock_engine is a generated GoMock package.
package mock_frontend

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	mpool "github.com/matrixorigin/matrixone/pkg/common/mpool"
	batch "github.com/matrixorigin/matrixone/pkg/container/batch"
	plan "github.com/matrixorigin/matrixone/pkg/pb/plan"
	timestamp "github.com/matrixorigin/matrixone/pkg/pb/timestamp"
	client "github.com/matrixorigin/matrixone/pkg/txn/client"
	engine "github.com/matrixorigin/matrixone/pkg/vm/engine"
)

// MockStatistics is a mock of Statistics interface.
type MockStatistics struct {
	ctrl     *gomock.Controller
	recorder *MockStatisticsMockRecorder
}

// MockStatisticsMockRecorder is the mock recorder for MockStatistics.
type MockStatisticsMockRecorder struct {
	mock *MockStatistics
}

// NewMockStatistics creates a new mock instance.
func NewMockStatistics(ctrl *gomock.Controller) *MockStatistics {
	mock := &MockStatistics{ctrl: ctrl}
	mock.recorder = &MockStatisticsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStatistics) EXPECT() *MockStatisticsMockRecorder {
	return m.recorder
}

// FilteredStats mocks base method.
func (m *MockStatistics) FilteredStats(ctx context.Context, expr *plan.Expr) (int32, int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FilteredStats", ctx, expr)
	ret0, _ := ret[0].(int32)
	ret1, _ := ret[1].(int64)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// FilteredStats indicates an expected call of FilteredStats.
func (mr *MockStatisticsMockRecorder) FilteredStats(ctx, expr interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FilteredStats", reflect.TypeOf((*MockStatistics)(nil).FilteredStats), ctx, expr)
}

// Size mocks base method.
func (m *MockStatistics) Size(ctx context.Context, columnName string) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Size", ctx, columnName)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Size indicates an expected call of Size.
func (mr *MockStatisticsMockRecorder) Size(ctx, columnName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Size", reflect.TypeOf((*MockStatistics)(nil).Size), ctx, columnName)
}

// Stats mocks base method.
func (m *MockStatistics) Stats(ctx context.Context) (int32, int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stats", ctx)
	ret0, _ := ret[0].(int32)
	ret1, _ := ret[1].(int64)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Stats indicates an expected call of Stats.
func (mr *MockStatisticsMockRecorder) Stats(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stats", reflect.TypeOf((*MockStatistics)(nil).Stats), ctx)
}

// MockTableDef is a mock of TableDef interface.
type MockTableDef struct {
	ctrl     *gomock.Controller
	recorder *MockTableDefMockRecorder
}

// MockTableDefMockRecorder is the mock recorder for MockTableDef.
type MockTableDefMockRecorder struct {
	mock *MockTableDef
}

// NewMockTableDef creates a new mock instance.
func NewMockTableDef(ctrl *gomock.Controller) *MockTableDef {
	mock := &MockTableDef{ctrl: ctrl}
	mock.recorder = &MockTableDefMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTableDef) EXPECT() *MockTableDefMockRecorder {
	return m.recorder
}

// tableDef mocks base method.
func (m *MockTableDef) tableDef() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "tableDef")
}

// tableDef indicates an expected call of tableDef.
func (mr *MockTableDefMockRecorder) tableDef() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "tableDef", reflect.TypeOf((*MockTableDef)(nil).tableDef))
}

// MockConstraint is a mock of Constraint interface.
type MockConstraint struct {
	ctrl     *gomock.Controller
	recorder *MockConstraintMockRecorder
}

// MockConstraintMockRecorder is the mock recorder for MockConstraint.
type MockConstraintMockRecorder struct {
	mock *MockConstraint
}

// NewMockConstraint creates a new mock instance.
func NewMockConstraint(ctrl *gomock.Controller) *MockConstraint {
	mock := &MockConstraint{ctrl: ctrl}
	mock.recorder = &MockConstraintMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockConstraint) EXPECT() *MockConstraintMockRecorder {
	return m.recorder
}

// constraint mocks base method.
func (m *MockConstraint) constraint() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "constraint")
}

// constraint indicates an expected call of constraint.
func (mr *MockConstraintMockRecorder) constraint() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "constraint", reflect.TypeOf((*MockConstraint)(nil).constraint))
}

// MockRelation is a mock of Relation interface.
type MockRelation struct {
	ctrl     *gomock.Controller
	recorder *MockRelationMockRecorder
}

// MockRelationMockRecorder is the mock recorder for MockRelation.
type MockRelationMockRecorder struct {
	mock *MockRelation
}

// NewMockRelation creates a new mock instance.
func NewMockRelation(ctrl *gomock.Controller) *MockRelation {
	mock := &MockRelation{ctrl: ctrl}
	mock.recorder = &MockRelationMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRelation) EXPECT() *MockRelationMockRecorder {
	return m.recorder
}

// AddTableDef mocks base method.
func (m *MockRelation) AddTableDef(arg0 context.Context, arg1 engine.TableDef) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddTableDef", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddTableDef indicates an expected call of AddTableDef.
func (mr *MockRelationMockRecorder) AddTableDef(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddTableDef", reflect.TypeOf((*MockRelation)(nil).AddTableDef), arg0, arg1)
}

// DelTableDef mocks base method.
func (m *MockRelation) DelTableDef(arg0 context.Context, arg1 engine.TableDef) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DelTableDef", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DelTableDef indicates an expected call of DelTableDef.
func (mr *MockRelationMockRecorder) DelTableDef(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DelTableDef", reflect.TypeOf((*MockRelation)(nil).DelTableDef), arg0, arg1)
}

// Delete mocks base method.
func (m *MockRelation) Delete(arg0 context.Context, arg1 *batch.Batch, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockRelationMockRecorder) Delete(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockRelation)(nil).Delete), arg0, arg1, arg2)
}

// FilteredStats mocks base method.
func (m *MockRelation) FilteredStats(ctx context.Context, expr *plan.Expr) (int32, int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FilteredStats", ctx, expr)
	ret0, _ := ret[0].(int32)
	ret1, _ := ret[1].(int64)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// FilteredStats indicates an expected call of FilteredStats.
func (mr *MockRelationMockRecorder) FilteredStats(ctx, expr interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FilteredStats", reflect.TypeOf((*MockRelation)(nil).FilteredStats), ctx, expr)
}

// GetHideKeys mocks base method.
func (m *MockRelation) GetHideKeys(arg0 context.Context) ([]*engine.Attribute, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetHideKeys", arg0)
	ret0, _ := ret[0].([]*engine.Attribute)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetHideKeys indicates an expected call of GetHideKeys.
func (mr *MockRelationMockRecorder) GetHideKeys(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetHideKeys", reflect.TypeOf((*MockRelation)(nil).GetHideKeys), arg0)
}

// GetPrimaryKeys mocks base method.
func (m *MockRelation) GetPrimaryKeys(arg0 context.Context) ([]*engine.Attribute, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPrimaryKeys", arg0)
	ret0, _ := ret[0].([]*engine.Attribute)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPrimaryKeys indicates an expected call of GetPrimaryKeys.
func (mr *MockRelationMockRecorder) GetPrimaryKeys(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPrimaryKeys", reflect.TypeOf((*MockRelation)(nil).GetPrimaryKeys), arg0)
}

// GetTableID mocks base method.
func (m *MockRelation) GetTableID(arg0 context.Context) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTableID", arg0)
	ret0, _ := ret[0].(string)
	return ret0
}

// GetTableID indicates an expected call of GetTableID.
func (mr *MockRelationMockRecorder) GetTableID(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTableID", reflect.TypeOf((*MockRelation)(nil).GetTableID), arg0)
}

// NewReader mocks base method.
func (m *MockRelation) NewReader(arg0 context.Context, arg1 int, arg2 *plan.Expr, arg3 [][]byte) ([]engine.Reader, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewReader", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].([]engine.Reader)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewReader indicates an expected call of NewReader.
func (mr *MockRelationMockRecorder) NewReader(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewReader", reflect.TypeOf((*MockRelation)(nil).NewReader), arg0, arg1, arg2, arg3)
}

// Ranges mocks base method.
func (m *MockRelation) Ranges(arg0 context.Context, arg1 *plan.Expr) ([][]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ranges", arg0, arg1)
	ret0, _ := ret[0].([][]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Ranges indicates an expected call of Ranges.
func (mr *MockRelationMockRecorder) Ranges(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ranges", reflect.TypeOf((*MockRelation)(nil).Ranges), arg0, arg1)
}

// Size mocks base method.
func (m *MockRelation) Size(ctx context.Context, columnName string) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Size", ctx, columnName)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Size indicates an expected call of Size.
func (mr *MockRelationMockRecorder) Size(ctx, columnName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Size", reflect.TypeOf((*MockRelation)(nil).Size), ctx, columnName)
}

// Stats mocks base method.
func (m *MockRelation) Stats(ctx context.Context) (int32, int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stats", ctx)
	ret0, _ := ret[0].(int32)
	ret1, _ := ret[1].(int64)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

func (m *MockRelation) Rows(ctx context.Context) (int64, error) {
	_, rows, err := m.Stats(ctx)
	return rows, err
}

// Stats indicates an expected call of Stats.
func (mr *MockRelationMockRecorder) Stats(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stats", reflect.TypeOf((*MockRelation)(nil).Stats), ctx)
}

// TableColumns mocks base method.
func (m *MockRelation) TableColumns(ctx context.Context) ([]*engine.Attribute, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TableColumns", ctx)
	ret0, _ := ret[0].([]*engine.Attribute)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// TableColumns indicates an expected call of TableColumns.
func (mr *MockRelationMockRecorder) TableColumns(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TableColumns", reflect.TypeOf((*MockRelation)(nil).TableColumns), ctx)
}

// TableDefs mocks base method.
func (m *MockRelation) TableDefs(arg0 context.Context) ([]engine.TableDef, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TableDefs", arg0)
	ret0, _ := ret[0].([]engine.TableDef)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// TableDefs indicates an expected call of TableDefs.
func (mr *MockRelationMockRecorder) TableDefs(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TableDefs", reflect.TypeOf((*MockRelation)(nil).TableDefs), arg0)
}

// Update mocks base method.
func (m *MockRelation) Update(arg0 context.Context, arg1 *batch.Batch) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockRelationMockRecorder) Update(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockRelation)(nil).Update), arg0, arg1)
}

// UpdateConstraint mocks base method.
func (m *MockRelation) UpdateConstraint(arg0 context.Context, arg1 *engine.ConstraintDef) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateConstraint", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateConstraint indicates an expected call of UpdateConstraint.
func (mr *MockRelationMockRecorder) UpdateConstraint(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateConstraint", reflect.TypeOf((*MockRelation)(nil).UpdateConstraint), arg0, arg1)
}

// Write mocks base method.
func (m *MockRelation) Write(arg0 context.Context, arg1 *batch.Batch) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Write", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Write indicates an expected call of Write.
func (mr *MockRelationMockRecorder) Write(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Write", reflect.TypeOf((*MockRelation)(nil).Write), arg0, arg1)
}

// MockReader is a mock of Reader interface.
type MockReader struct {
	ctrl     *gomock.Controller
	recorder *MockReaderMockRecorder
}

// MockReaderMockRecorder is the mock recorder for MockReader.
type MockReaderMockRecorder struct {
	mock *MockReader
}

// NewMockReader creates a new mock instance.
func NewMockReader(ctrl *gomock.Controller) *MockReader {
	mock := &MockReader{ctrl: ctrl}
	mock.recorder = &MockReaderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockReader) EXPECT() *MockReaderMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockReader) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockReaderMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockReader)(nil).Close))
}

// Read mocks base method.
func (m *MockReader) Read(arg0 context.Context, arg1 []string, arg2 *plan.Expr, arg3 *mpool.MPool) (*batch.Batch, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Read", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(*batch.Batch)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Read indicates an expected call of Read.
func (mr *MockReaderMockRecorder) Read(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Read", reflect.TypeOf((*MockReader)(nil).Read), arg0, arg1, arg2, arg3)
}

// MockDatabase is a mock of Database interface.
type MockDatabase struct {
	ctrl     *gomock.Controller
	recorder *MockDatabaseMockRecorder
}

// MockDatabaseMockRecorder is the mock recorder for MockDatabase.
type MockDatabaseMockRecorder struct {
	mock *MockDatabase
}

// NewMockDatabase creates a new mock instance.
func NewMockDatabase(ctrl *gomock.Controller) *MockDatabase {
	mock := &MockDatabase{ctrl: ctrl}
	mock.recorder = &MockDatabaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDatabase) EXPECT() *MockDatabaseMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockDatabase) Create(arg0 context.Context, arg1 string, arg2 []engine.TableDef) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockDatabaseMockRecorder) Create(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockDatabase)(nil).Create), arg0, arg1, arg2)
}

// Delete mocks base method.
func (m *MockDatabase) Delete(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockDatabaseMockRecorder) Delete(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockDatabase)(nil).Delete), arg0, arg1)
}

// GetDatabaseId mocks base method.
func (m *MockDatabase) GetDatabaseId(arg0 context.Context) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDatabaseId", arg0)
	ret0, _ := ret[0].(string)
	return ret0
}

// GetDatabaseId indicates an expected call of GetDatabaseId.
func (mr *MockDatabaseMockRecorder) GetDatabaseId(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDatabaseId", reflect.TypeOf((*MockDatabase)(nil).GetDatabaseId), arg0)
}

// Relation mocks base method.
func (m *MockDatabase) Relation(arg0 context.Context, arg1 string) (engine.Relation, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Relation", arg0, arg1)
	ret0, _ := ret[0].(engine.Relation)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Relation indicates an expected call of Relation.
func (mr *MockDatabaseMockRecorder) Relation(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Relation", reflect.TypeOf((*MockDatabase)(nil).Relation), arg0, arg1)
}

// Relations mocks base method.
func (m *MockDatabase) Relations(arg0 context.Context) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Relations", arg0)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Relations indicates an expected call of Relations.
func (mr *MockDatabaseMockRecorder) Relations(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Relations", reflect.TypeOf((*MockDatabase)(nil).Relations), arg0)
}

// Truncate mocks base method.
func (m *MockDatabase) Truncate(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Truncate", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Truncate indicates an expected call of Truncate.
func (mr *MockDatabaseMockRecorder) Truncate(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Truncate", reflect.TypeOf((*MockDatabase)(nil).Truncate), arg0, arg1)
}

// MockEngine is a mock of Engine interface.
type MockEngine struct {
	ctrl     *gomock.Controller
	recorder *MockEngineMockRecorder
}

// MockEngineMockRecorder is the mock recorder for MockEngine.
type MockEngineMockRecorder struct {
	mock *MockEngine
}

// NewMockEngine creates a new mock instance.
func NewMockEngine(ctrl *gomock.Controller) *MockEngine {
	mock := &MockEngine{ctrl: ctrl}
	mock.recorder = &MockEngineMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEngine) EXPECT() *MockEngineMockRecorder {
	return m.recorder
}

// Commit mocks base method.
func (m *MockEngine) Commit(ctx context.Context, op client.TxnOperator) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Commit", ctx, op)
	ret0, _ := ret[0].(error)
	return ret0
}

// Commit indicates an expected call of Commit.
func (mr *MockEngineMockRecorder) Commit(ctx, op interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Commit", reflect.TypeOf((*MockEngine)(nil).Commit), ctx, op)
}

// Create mocks base method.
func (m *MockEngine) Create(ctx context.Context, databaseName string, op client.TxnOperator) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, databaseName, op)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockEngineMockRecorder) Create(ctx, databaseName, op interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockEngine)(nil).Create), ctx, databaseName, op)
}

// Database mocks base method.
func (m *MockEngine) Database(ctx context.Context, databaseName string, op client.TxnOperator) (engine.Database, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Database", ctx, databaseName, op)
	ret0, _ := ret[0].(engine.Database)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Database indicates an expected call of Database.
func (mr *MockEngineMockRecorder) Database(ctx, databaseName, op interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Database", reflect.TypeOf((*MockEngine)(nil).Database), ctx, databaseName, op)
}

// Databases mocks base method.
func (m *MockEngine) Databases(ctx context.Context, op client.TxnOperator) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Databases", ctx, op)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Databases indicates an expected call of Databases.
func (mr *MockEngineMockRecorder) Databases(ctx, op interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Databases", reflect.TypeOf((*MockEngine)(nil).Databases), ctx, op)
}

// Delete mocks base method.
func (m *MockEngine) Delete(ctx context.Context, databaseName string, op client.TxnOperator) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, databaseName, op)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockEngineMockRecorder) Delete(ctx, databaseName, op interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockEngine)(nil).Delete), ctx, databaseName, op)
}

// Hints mocks base method.
func (m *MockEngine) Hints() engine.Hints {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Hints")
	ret0, _ := ret[0].(engine.Hints)
	return ret0
}

// Hints indicates an expected call of Hints.
func (mr *MockEngineMockRecorder) Hints() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Hints", reflect.TypeOf((*MockEngine)(nil).Hints))
}

// New mocks base method.
func (m *MockEngine) New(ctx context.Context, op client.TxnOperator) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "New", ctx, op)
	ret0, _ := ret[0].(error)
	return ret0
}

// New indicates an expected call of New.
func (mr *MockEngineMockRecorder) New(ctx, op interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "New", reflect.TypeOf((*MockEngine)(nil).New), ctx, op)
}

// NewBlockReader mocks base method.
func (m *MockEngine) NewBlockReader(ctx context.Context, num int, ts timestamp.Timestamp, expr *plan.Expr, ranges [][]byte, tblDef *plan.TableDef) ([]engine.Reader, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewBlockReader", ctx, num, ts, expr, ranges, tblDef)
	ret0, _ := ret[0].([]engine.Reader)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewBlockReader indicates an expected call of NewBlockReader.
func (mr *MockEngineMockRecorder) NewBlockReader(ctx, num, ts, expr, ranges, tblDef interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewBlockReader", reflect.TypeOf((*MockEngine)(nil).NewBlockReader), ctx, num, ts, expr, ranges, tblDef)
}

// Nodes mocks base method.
func (m *MockEngine) Nodes() (engine.Nodes, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Nodes")
	ret0, _ := ret[0].(engine.Nodes)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Nodes indicates an expected call of Nodes.
func (mr *MockEngineMockRecorder) Nodes() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Nodes", reflect.TypeOf((*MockEngine)(nil).Nodes))
}

// Rollback mocks base method.
func (m *MockEngine) Rollback(ctx context.Context, op client.TxnOperator) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Rollback", ctx, op)
	ret0, _ := ret[0].(error)
	return ret0
}

// Rollback indicates an expected call of Rollback.
func (mr *MockEngineMockRecorder) Rollback(ctx, op interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Rollback", reflect.TypeOf((*MockEngine)(nil).Rollback), ctx, op)
}

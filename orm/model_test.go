package orm

import (
	"github.com/stretchr/testify/assert"
	"leason-toy-orm/orm/internal/errs"
	"reflect"
	"testing"
)

func Test_ParseModel(t *testing.T) {
	testCases := []struct {
		name      string
		entity    any
		wantModel *Model
		wantErr   error
	}{
		{
			name:      "test Model",
			entity:    TestModel{},
			wantModel: nil,
			wantErr:   errs.ErrPointerOnly,
		},
		{
			name:   "pointer",
			entity: &TestModel{},
			wantModel: &Model{
				tableName: "test_model",
				fields: map[string]*Field{
					"Id": {
						colName: "id",
					},
					"FirstName": {
						colName: "first_name",
					},
					"LastName": {
						colName: "last_name",
					},
					"Age": {
						colName: "age",
					},
				},
			},
			wantErr: nil,
		},
		{
			name:      "map",
			entity:    map[string]string{},
			wantModel: nil,
			wantErr:   errs.ErrPointerOnly,
		},
		{
			name:      "slice",
			entity:    []string{},
			wantModel: nil,
			wantErr:   errs.ErrPointerOnly,
		},
	}

	r := newRegistry()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m, err := r.Get(tc.entity)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantModel, m)
		})
	}
}

func TestRegistr_get(t *testing.T) {
	testCases := []struct {
		name string

		entity    any
		wantModel *Model
		wantErr   error
		cacheSize int
	}{
		{
			name:   "pointer",
			entity: &TestModel{},
			wantModel: &Model{
				tableName: "test_model",
				fields: map[string]*Field{
					"Id": {
						colName: "id",
					},
					"FirstName": {
						colName: "first_name",
					},
					"LastName": {
						colName: "last_name",
					},
					"Age": {
						colName: "age",
					},
				},
			},
			wantErr:   nil,
			cacheSize: 1,
		},
		{
			name: "tag",
			entity: func() any {
				type TagTable struct {
					FirstName string `orm:"column=first_name_t"`
				}
				return &TagTable{}
			}(),
			wantModel: &Model{
				tableName: "tag_table",
				fields: map[string]*Field{
					"FirstName": {
						colName: "first_name_t",
					},
				},
			},
			wantErr:   nil,
			cacheSize: 1,
		},
		{
			name: "empty column",
			entity: func() any {
				type TagTable struct {
					FirstName string `orm:"column="`
				}
				return &TagTable{}
			}(),
			wantModel: &Model{
				tableName: "tag_table",
				fields: map[string]*Field{
					"FirstName": {
						colName: "first_name",
					},
				},
			},
			wantErr:   nil,
			cacheSize: 1,
		},
		{
			name: "invalid column",
			entity: func() any {
				type TagTable struct {
					FirstName string `orm:"column"`
				}
				return &TagTable{}
			}(),
			wantModel: nil,
			wantErr:   errs.NewErrInvalidTagContent("column"),
		},
		{
			name: "ignore tag",
			entity: func() any {
				type TagTable struct {
					FirstName string `orm:"abc=abc"`
				}
				return &TagTable{}
			}(),
			wantModel: &Model{
				tableName: "tag_table",
				fields: map[string]*Field{
					"FirstName": {
						colName: "first_name",
					},
				},
			},
			wantErr:   nil,
			cacheSize: 1,
		},
		{
			name:   "table name",
			entity: &CustomTableName{},
			wantModel: &Model{
				tableName: "custom_table_name_t",
				fields: map[string]*Field{
					"FirstName": {
						colName: "first_name",
					},
				},
			},
			wantErr:   nil,
			cacheSize: 1,
		},
		{
			name:   "table name ptr",
			entity: &CustomTableNamePtr{},
			wantModel: &Model{
				tableName: "custom_table_name_ptr_t",
				fields: map[string]*Field{
					"FirstName": {
						colName: "first_name",
					},
				},
			},
			wantErr:   nil,
			cacheSize: 1,
		},
		{
			name:   "empty table name",
			entity: &EmptyTableName{},
			wantModel: &Model{
				tableName: "empty_table_name",
				fields: map[string]*Field{
					"FirstName": {
						colName: "first_name",
					},
				},
			},
			wantErr:   nil,
			cacheSize: 1,
		},
	}

	r := newRegistry()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m, err := r.Get(tc.entity)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}

			assert.Equal(t, tc.wantModel, m)

			typ := reflect.TypeOf(tc.entity)
			cache, ok := r.models.Load(typ)
			assert.True(t, ok)
			assert.Equal(t, tc.wantModel, cache)
		})
	}
}

type CustomTableName struct {
	FirstName string
}

func (c CustomTableName) TableName() string {
	return "custom_table_name_t"
}

type CustomTableNamePtr struct {
	FirstName string
}

func (c *CustomTableNamePtr) TableName() string {
	return "custom_table_name_ptr_t"
}

type EmptyTableName struct {
	FirstName string
}

func (c EmptyTableName) TableName() string {
	return ""
}

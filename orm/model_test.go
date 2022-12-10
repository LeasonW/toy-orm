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
		wantModel *model
		wantErr   error
	}{
		{
			name:      "test model",
			entity:    TestModel{},
			wantModel: nil,
			wantErr:   errs.ErrPointerOnly,
		},
		{
			name:   "pointer",
			entity: &TestModel{},
			wantModel: &model{
				tableName: "test_model",
				fields: map[string]*field{
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
			m, err := r.get(tc.entity)
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
		wantModel *model
		wantErr   error
		cacheSize int
	}{
		{
			name:   "pointer",
			entity: &TestModel{},
			wantModel: &model{
				tableName: "test_model",
				fields: map[string]*field{
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
	}

	r := newRegistry()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m, err := r.get(tc.entity)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}

			assert.Equal(t, tc.wantModel, m)
			assert.Equal(t, tc.cacheSize, len(r.models))

			typ := reflect.TypeOf(tc.entity)
			m, ok := r.models[typ]
			assert.True(t, ok)
			assert.Equal(t, tc.wantModel, m)
		})
	}
}

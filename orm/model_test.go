package orm

import (
	"github.com/stretchr/testify/assert"
	"leason-toy-orm/orm/internal/errs"
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

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m, err := parseModel(tc.entity)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantModel, m)
		})
	}
}

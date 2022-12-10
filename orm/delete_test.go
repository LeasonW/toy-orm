package orm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeleter_Build(t *testing.T) {
	testCases := []struct {
		name      string
		builder   QueryBuilder
		wantErr   error
		wantQuery *Query
	}{
		{
			name:    "no where",
			builder: (&Deleter[TestModel]{}).From("`test_model`"),
			wantQuery: &Query{
				SQL: "DELETE FROM `test_model`;",
			},
		},
		{
			name:    "where",
			builder: (&Deleter[TestModel]{}).Where(C("Id").EQ(16)),
			wantQuery: &Query{
				SQL:  "DELETE FROM `test_model` WHERE `id` = ?;",
				Args: []any{16},
			},
		},
		{
			name:    "from",
			builder: (&Deleter[TestModel]{}).From("`test_model`").Where(C("Id").EQ(16)),
			wantQuery: &Query{
				SQL:  "DELETE FROM `test_model` WHERE `id` = ?;",
				Args: []any{16},
			},
		},
		{
			name:    "not",
			builder: (&Deleter[TestModel]{}).From("`test_model`").Where(Not(C("Id").EQ(16))),
			wantQuery: &Query{
				SQL:  "DELETE FROM `test_model` WHERE NOT (`id` = ?);",
				Args: []any{16},
			},
		},
		{
			name:    "and",
			builder: (&Deleter[TestModel]{}).From("`test_model`").Where(C("Id").GT(16).And(C("Id").LT(35))),
			wantQuery: &Query{
				SQL:  "DELETE FROM `test_model` WHERE (`id` > ?) AND (`id` < ?);",
				Args: []any{16, 35},
			},
		},
	}

	for _, tc := range testCases {
		c := tc
		t.Run(c.name, func(t *testing.T) {
			query, err := c.builder.Build()
			assert.Equal(t, c.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantQuery, query)
		})
	}
}

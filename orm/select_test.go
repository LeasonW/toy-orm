package orm

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSelector_Build(t *testing.T) {
	testCases := []struct {
		name      string
		builder   QueryBuilder
		wantQuery *Query
		wantErr   error
	}{
		{
			name:    "no from",
			builder: &Selector[TestModel]{},
			wantQuery: &Query{
				SQL:  "SELECT * FROM `TestModel`;",
				Args: nil,
			},
		},
		{
			name:    "from",
			builder: (&Selector[TestModel]{}).From("`test_model`"),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model`;",
				Args: nil,
			},
		},
		{
			name:    "empty from",
			builder: (&Selector[TestModel]{}).From(""),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `TestModel`;",
				Args: nil,
			},
		},
		{
			name:    "db table",
			builder: (&Selector[TestModel]{}).From("`test_db`.`test_model`"),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_db`.`test_model`;",
				Args: nil,
			},
		},
		{
			name:    "where",
			builder: (&Selector[TestModel]{}).Where(C("Age").Eq(18)),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `TestModel` WHERE `Age` = ?;",
				Args: []any{18},
			},
		},
		{
			// 使用 OR
			name:    "or",
			builder: (&Selector[TestModel]{}).Where(C("Age").GT(18).Or(C("Age").LT(35))),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `TestModel` WHERE (`Age` > ?) OR (`Age` < ?);",
				Args: []any{18, 35},
			},
		},
		{
			// 使用 NOT
			name:    "not",
			builder: (&Selector[TestModel]{}).Where(Not(C("Age").GT(18))),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `TestModel` WHERE NOT (`Age` > ?);",
				Args: []any{18},
			},
		},
		{
			// 使用 AND
			name:    "and",
			builder: (&Selector[TestModel]{}).Where(C("Age").GT(18).And(C("Age").LT(35))),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `TestModel` WHERE (`Age` > ?) AND (`Age` < ?);",
				Args: []any{18, 35},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			q, err := tc.builder.Build()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantQuery, q)
		})
	}
}

type TestModel struct {
	Id        int64
	FirstName string
	Age       int8
	LastName  *sql.NullString
}
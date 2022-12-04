package reflect

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIterateArrayOrSlice(t *testing.T) {
	testCases := []struct {
		name   string
		entity any

		wantVals []any
		wantErr  error
	}{
		{
			name:     "array",
			entity:   [3]int{1, 2, 3},
			wantErr:  nil,
			wantVals: []any{1, 2, 3},
		},
		{
			name:     "slice",
			entity:   []int{1, 2, 3},
			wantErr:  nil,
			wantVals: []any{1, 2, 3},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			vals, err := IterateArrayOrSlice(tt.entity)
			assert.Equal(t, tt.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tt.wantVals, vals)
		})
	}
}

func TestIterateMap(t *testing.T) {
	testCases := []struct {
		name   string
		entity any

		wantKeys   []any
		wantValues []any
		wantErr    error
	}{
		{
			name: "map",
			entity: map[string]string{
				"A": "a",
				"B": "b",
			},
			wantKeys:   []any{"A", "B"},
			wantValues: []any{"a", "b"},
			wantErr:    nil,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			keys, values, err := IterateMap(tt.entity)
			assert.Equal(t, tt.wantErr, err)
			if err != nil {
				return
			}
			assert.EqualValues(t, tt.wantKeys, keys)
			assert.EqualValues(t, tt.wantValues, values)
		})
	}
}

package reflect

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

type User struct {
	Name string
	age  int
}

func Test_IterateFields(t *testing.T) {
	testCases := []struct {
		name    string
		entity  any
		wantErr error
		wantRes map[string]any
	}{
		{
			name: "structure",
			entity: User{
				Name: "Tom",
				age:  18,
			},
			wantErr: nil,
			wantRes: map[string]any{
				"Name": "Tom",
				"age":  0,
			},
		},
		{
			name: "pointer",
			entity: &User{
				Name: "Tom",
				age:  18,
			},
			wantErr: nil,
			wantRes: map[string]any{
				"Name": "Tom",
				"age":  0,
			},
		},
		{
			name:    "basic type",
			entity:  18,
			wantErr: errors.New("unsupported type"),
		},
		{
			name: "multiple pointer",
			entity: func() **User {
				res := &User{
					Name: "Tom",
					age:  18,
				}
				return &res
			}(),
			wantErr: nil,
			wantRes: map[string]any{
				"Name": "Tom",
				"age":  0,
			},
		},
		{
			name:    "nil",
			entity:  nil,
			wantErr: errors.New("unsupported type"),
		},
		{
			name:    "type nil",
			entity:  (*User)(nil),
			wantErr: errors.New("unsupported zero value"),
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			fields, err := IterateFields(tt.entity)
			assert.Equal(t, tt.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tt.wantRes, fields)
		})
	}
}

func TestSetField(t *testing.T) {
	testCases := []struct {
		name string

		entity any
		field  string
		newVal any

		wantErr    error
		wantEntity any
	}{
		{
			name: "structure",
			entity: User{
				Name: "Tom",
			},
			field:  "Name",
			newVal: "Jerry",

			wantErr: errors.New("cannot set field"),
			wantEntity: User{
				Name: "Jerry",
			},
		},
		{
			name: "pointer",
			entity: &User{
				Name: "Tom",
			},
			field:  "Name",
			newVal: "Jerry",

			wantErr: nil,
			wantEntity: &User{
				Name: "Jerry",
			},
		},
		{
			name: "pointer unexported",
			entity: &User{
				Name: "Tom",
				age:  18,
			},
			field:  "age",
			newVal: 19,

			wantErr: errors.New("cannot set field"),
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			err := SetField(tt.entity, tt.field, tt.newVal)
			assert.Equal(t, tt.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tt.wantEntity, tt.entity)
		})
	}
}

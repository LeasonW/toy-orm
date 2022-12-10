package orm

import (
	"leason-toy-orm/orm/internal/errs"
	"reflect"
	"sync"
	"unicode"
)

type model struct {
	tableName string
	fields    map[string]*field
}

type field struct {
	colName string
}

type registry struct {
	models sync.Map
}

func newRegistry() *registry {
	return &registry{}
}

func (r *registry) get(val any) (*model, error) {
	typ := reflect.TypeOf(val)

	m, ok := r.models.Load(typ)
	if ok {
		return m.(*model), nil
	}
	m, err := r.parseModel(val)
	if err != nil {
		return nil, err
	}
	r.models.Store(typ, m)
	return m.(*model), nil
}

// parseModel 限制只能使用一级指针
func (r *registry) parseModel(entity any) (*model, error) {
	typ := reflect.TypeOf(entity)
	if typ.Kind() != reflect.Pointer || typ.Elem().Kind() != reflect.Struct {
		return nil, errs.ErrPointerOnly
	}

	typ = typ.Elem()
	numField := typ.NumField()
	m := &model{
		tableName: formattedName(typ.Name()),
		fields:    make(map[string]*field, numField),
	}
	for i := 0; i < numField; i++ {
		m.fields[typ.Field(i).Name] = &field{
			colName: formattedName(typ.Field(i).Name),
		}
	}
	return m, nil
}

func formattedName(str string) string {
	var buf []byte
	for i, v := range str {
		if unicode.IsUpper(v) && i > 0 {
			buf = append(buf, '_')
		}
		buf = append(buf, byte(unicode.ToLower(v)))
	}
	return string(buf)
}

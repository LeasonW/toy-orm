package orm

import (
	"leason-toy-orm/orm/internal/errs"
	"reflect"
	"strings"
	"sync"
	"unicode"
)

const (
	tagColumn = "column"
)

type Registry interface {
	Get(val any) (*Model, error)
	Registry(val any) (*Model, error)
}

type Model struct {
	tableName string
	fields    map[string]*Field
}

type Field struct {
	colName string
}

type registry struct {
	models sync.Map
}

func newRegistry() *registry {
	return &registry{}
}

func (r *registry) Get(val any) (*Model, error) {
	typ := reflect.TypeOf(val)

	m, ok := r.models.Load(typ)
	if ok {
		return m.(*Model), nil
	}
	m, err := r.Register(val)
	if err != nil {
		return nil, err
	}
	r.models.Store(typ, m)
	return m.(*Model), nil
}

// Register 限制只能使用一级指针
func (r *registry) Register(entity any) (*Model, error) {
	typ := reflect.TypeOf(entity)
	if typ.Kind() != reflect.Pointer || typ.Elem().Kind() != reflect.Struct {
		return nil, errs.ErrPointerOnly
	}

	typ = typ.Elem()
	numField := typ.NumField()

	var tableName string
	if val, ok := entity.(TableName); ok {
		tableName = val.TableName()
	}
	if tableName == "" {
		tableName = formattedName(typ.Name())
	}

	m := &Model{
		tableName: tableName,
		fields:    make(map[string]*Field, numField),
	}
	for i := 0; i < numField; i++ {
		f := typ.Field(i)
		pairs, err := r.parseTag(f.Tag)
		if err != nil {
			return nil, err
		}
		columnName := pairs[tagColumn]
		if columnName == "" {
			columnName = formattedName(f.Name)
		}
		m.fields[f.Name] = &Field{
			colName: columnName,
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

func (r *registry) parseTag(tag reflect.StructTag) (map[string]string, error) {
	ormTag, ok := tag.Lookup("orm")
	if !ok {
		return map[string]string{}, nil
	}

	pairs := strings.Split(ormTag, ",")
	res := make(map[string]string, len(pairs))
	for _, pair := range pairs {
		segs := strings.Split(pair, "=")
		if len(segs) != 2 {
			return nil, errs.NewErrInvalidTagContent(pair)
		}
		key := segs[0]
		val := segs[1]
		res[key] = val
	}
	return res, nil
}

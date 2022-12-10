package orm

import (
	"leason-toy-orm/orm/internal/errs"
	"strings"
)

type Deleter[T any] struct {
	tableName string
	sb        *strings.Builder
	model     *model
	args      []any
	wheres    []Predicate

	db *DB
}

func NewDeleter[T any](db *DB) *Deleter[T] {
	return &Deleter[T]{
		sb: &strings.Builder{},
		db: db,
	}
}

func (d *Deleter[T]) Build() (*Query, error) {
	var err error
	d.model, err = d.db.r.get(new(T))
	if err != nil {
		return nil, err
	}

	d.sb.WriteString("DELETE FROM ")

	if d.tableName != "" {
		d.sb.WriteString(d.tableName)
	} else {
		d.sb.WriteByte('`')
		d.sb.WriteString(d.model.tableName)
		d.sb.WriteByte('`')
	}

	if len(d.wheres) > 0 {
		d.sb.WriteString(" WHERE ")
		p := d.wheres[0]
		for i := 1; i < len(d.wheres); i++ {
			p = p.And(d.wheres[i])
		}
		if err := d.buildExpression(p); err != nil {
			return nil, err
		}
	}

	d.sb.WriteByte(';')
	return &Query{
		SQL:  strings.ReplaceAll(d.sb.String(), "  ", " "),
		Args: d.args,
	}, nil
}

// From accepts model definition
func (d *Deleter[T]) From(table string) *Deleter[T] {
	d.tableName = table
	return d
}

// Where accepts predicates
func (d *Deleter[T]) Where(predicates ...Predicate) *Deleter[T] {
	d.wheres = predicates
	return d
}

func (d *Deleter[T]) buildExpression(e Expression) error {
	switch exp := e.(type) {
	case nil:
		return nil
	case Predicate:
		_, ok := exp.left.(Predicate)
		if ok {
			d.sb.WriteByte('(')
		}
		if err := d.buildExpression(exp.left); err != nil {
			return err
		}
		if ok {
			d.sb.WriteByte(')')
		}

		d.sb.WriteByte(' ')
		d.sb.WriteString(exp.op.String())
		d.sb.WriteByte(' ')

		_, ok = exp.right.(Predicate)
		if ok {
			d.sb.WriteByte('(')
		}

		if err := d.buildExpression(exp.right); err != nil {
			return err
		}

		if ok {
			d.sb.WriteByte(')')
		}

	case Column:
		fd, ok := d.model.fields[exp.name]
		if !ok {
			return errs.NewErrUnknownField(exp.name)
		}
		d.sb.WriteByte('`')
		d.sb.WriteString(fd.colName)
		d.sb.WriteByte('`')

	case Value:
		d.sb.WriteByte('?')
		d.addArgs(exp.val)

	default:
		return errs.NewErrUnsupportedExpression(exp)
	}
	return nil
}

func (d *Deleter[T]) addArgs(arg any) {
	if d.args == nil {
		d.args = make([]any, 0, 4)
	}
	d.args = append(d.args, arg)
}

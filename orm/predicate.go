package orm

type op string

const (
	opEq  op = "="
	opGT  op = ">"
	opLT  op = "<"
	opNot op = "NOT"
	opAnd op = "AND"
	opOr  op = "OR"
)

type Expression interface {
	expr()
}

func (o op) String() string {
	return string(o)
}

type Predicate struct {
	left  Expression
	op    op
	right Expression
}

func (p Predicate) expr() {}

type Column struct {
	name string
}

func C(name string) Column {
	return Column{name: name}
}

func (c Column) expr() {}

type Value struct {
	val interface{}
}

func (v Value) expr() {}

func (c Column) Eq(arg interface{}) Predicate {
	return Predicate{
		left:  c,
		op:    opEq,
		right: Value{val: arg},
	}
}

func (c Column) GT(arg any) Predicate {
	return Predicate{
		left: c,
		op:   opGT,
		right: Value{
			val: arg,
		},
	}
}

func (c Column) EQ(arg any) Predicate {
	return Predicate{
		left: c,
		op:   opEq,
		right: Value{
			val: arg,
		},
	}
}

func (c Column) LT(arg any) Predicate {
	return Predicate{
		left: c,
		op:   opLT,
		right: Value{
			val: arg,
		},
	}
}

func Not(p Predicate) Predicate {
	return Predicate{
		op:    opNot,
		right: p,
	}
}

func (p Predicate) And(r Predicate) Predicate {
	return Predicate{
		left:  p,
		op:    opAnd,
		right: r,
	}
}

func (p Predicate) Or(r Predicate) Predicate {
	return Predicate{
		left:  p,
		op:    opOr,
		right: r,
	}
}

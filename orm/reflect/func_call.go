package reflect

import "reflect"

func IterateFunc(entity any) (map[string]FuncInfo, error) {
	return nil, nil
}

type FuncInfo struct {
	Name        string
	InputTypes  []reflect.Type
	OutputTypes []reflect.Type
	Result      []any
}

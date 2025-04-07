package goop

import (
	"fmt"
	"strings"
)

type ValType string

const (
	ValTypeString   ValType = "string"
	ValTypeNumber   ValType = "number"
	ValTypeBool     ValType = "bool"
	ValTypeList     ValType = "list"
	ValTypeFunction ValType = "function"
)

type Val struct {
	Type ValType
	Val  any
}

func NewStringVal(s string) *Val {
	return &Val{
		Type: ValTypeString,
		Val:  s,
	}
}

func NewFuncVal(f func(args []*Val) (*Val, error)) *Val {
	return &Val{
		Type: ValTypeFunction,
		Val:  f,
	}
}

func (v *Val) String() string {
	switch v.Type {
	case ValTypeString:
		val := v.Val.(string)
		val = strings.ReplaceAll(val, "\\", "\\\\")
		val = strings.ReplaceAll(val, "\n", "\\n")
		val = strings.ReplaceAll(val, "\t", "\\t")
		val = strings.ReplaceAll(val, "\"", "\\\"")
		return fmt.Sprintf("\"%s\"", val)
	case ValTypeNumber:
		val := v.Val.(float64)
		if val == float64(int(val)) {
			return fmt.Sprintf("%d", int(val))
		}
		return fmt.Sprintf("%f", val)
	case ValTypeBool:
		val := v.Val.(bool)
		if val {
			return "#t"
		}
		return "#f"
	default:
		return fmt.Sprintf("%v", v.Val)
	}
}

package object

////////////////////////////////////////////////////////////////////////////////
// DEPENDENCIES
////////////////////////////////////////////////////////////////////////////////

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"strings"

	"github.com/klaytonkowalski/example-interpreter/ast"
)

////////////////////////////////////////////////////////////////////////////////
// VARIABLES
////////////////////////////////////////////////////////////////////////////////

const (
	ObjectInteger        = "Integer"
	ObjectBoolean        = "Boolean"
	ObjectNull           = "Null"
	ObjectReturn         = "Return"
	ObjectError          = "Error"
	ObjectFunction       = "Function"
	ObjectString         = "String"
	ObjectNativeFunction = "Native Function"
	ObjectArray          = "Array"
	ObjectHash           = "Hash"
)

////////////////////////////////////////////////////////////////////////////////
// INTERFACES
////////////////////////////////////////////////////////////////////////////////

type Object interface {
	GetType() string
	GetDebugString() string
}

type Hashable interface {
	GetHashKey() HashKey
}

////////////////////////////////////////////////////////////////////////////////
// STRUCTURES
////////////////////////////////////////////////////////////////////////////////

type Integer struct {
	Value int64
}

type Boolean struct {
	Value bool
}

type Null struct{}

type Return struct {
	Value Object
}

type Error struct {
	Message string
}

type Function struct {
	Parameters  []*ast.Identifier
	Body        *ast.BlockStatement
	Environment *Environment
}

type String struct {
	Value string
}

type Native struct {
	Function NativeFn
}

type Array struct {
	Elements []Object
}

type HashKey struct {
	Type  string
	Value uint64
}

type HashPair struct {
	Key   Object
	Value Object
}

type Hash struct {
	Pairs map[HashKey]HashPair
}

////////////////////////////////////////////////////////////////////////////////
// METHODS
////////////////////////////////////////////////////////////////////////////////

func (i *Integer) GetType() string {
	return ObjectInteger
}

func (i *Integer) GetDebugString() string {
	return fmt.Sprintf("%d", i.Value)
}

func (i *Integer) GetHashKey() HashKey {
	return HashKey{Type: i.GetType(), Value: uint64(i.Value)}
}

func (b *Boolean) GetType() string {
	return ObjectBoolean
}

func (b *Boolean) GetDebugString() string {
	return fmt.Sprintf("%t", b.Value)
}

func (b *Boolean) GetHashKey() HashKey {
	var value uint64
	if b.Value {
		value = 1
	} else {
		value = 0
	}
	return HashKey{Type: b.GetType(), Value: value}
}

func (n *Null) GetType() string {
	return ObjectNull
}

func (n *Null) GetDebugString() string {
	return "null"
}

func (r *Return) GetType() string {
	return ObjectReturn
}

func (r *Return) GetDebugString() string {
	return r.Value.GetDebugString()
}

func (e *Error) GetType() string {
	return ObjectError
}

func (e *Error) GetDebugString() string {
	return "Error: " + e.Message
}

func (f *Function) GetType() string {
	return ObjectFunction
}

func (f *Function) GetDebugString() string {
	var out bytes.Buffer
	params := []string{}
	for _, param := range f.Parameters {
		params = append(params, param.GetDebugString())
	}
	out.WriteString("fn(")
	out.WriteString(strings.Join(params, ","))
	out.WriteString(") {\n")
	out.WriteString(f.Body.GetDebugString())
	out.WriteString("\n}")
	return out.String()
}

func (s *String) GetType() string {
	return ObjectString
}

func (s *String) GetDebugString() string {
	return s.Value
}

func (s *String) GetHashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))
	return HashKey{Type: s.GetType(), Value: h.Sum64()}
}

func (n *Native) GetType() string {
	return ObjectNativeFunction
}

func (n *Native) GetDebugString() string {
	return "native function"
}

func (a *Array) GetType() string {
	return ObjectArray
}

func (a *Array) GetDebugString() string {
	var out bytes.Buffer
	elems := []string{}
	for _, elem := range a.Elements {
		elems = append(elems, elem.GetDebugString())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elems, ","))
	out.WriteString("]")
	return out.String()
}

func (h *Hash) GetType() string {
	return ObjectHash
}

func (h *Hash) GetDebugString() string {
	var out bytes.Buffer
	pairs := []string{}
	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.GetDebugString(), pair.Value.GetDebugString()))
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ","))
	out.WriteString("}")
	return out.String()
}

////////////////////////////////////////////////////////////////////////////////
// FUNCTIONS
////////////////////////////////////////////////////////////////////////////////

type NativeFn func(args ...Object) Object

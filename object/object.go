package object

////////////////////////////////////////////////////////////////////////////////
// DEPENDENCIES
////////////////////////////////////////////////////////////////////////////////

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/klaytonkowalski/example-interpreter/ast"
)

////////////////////////////////////////////////////////////////////////////////
// VARIABLES
////////////////////////////////////////////////////////////////////////////////

const (
	ObjectInteger  = "Integer"
	ObjectBoolean  = "Boolean"
	ObjectNull     = "Null"
	ObjectReturn   = "Return"
	ObjectError    = "Error"
	ObjectFunction = "Function"
)

////////////////////////////////////////////////////////////////////////////////
// INTERFACES
////////////////////////////////////////////////////////////////////////////////

type Object interface {
	GetType() string
	GetDebugString() string
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

////////////////////////////////////////////////////////////////////////////////
// METHODS
////////////////////////////////////////////////////////////////////////////////

func (i *Integer) GetType() string {
	return ObjectInteger
}

func (i *Integer) GetDebugString() string {
	return fmt.Sprintf("%d", i.Value)
}

func (b *Boolean) GetType() string {
	return ObjectBoolean
}

func (b *Boolean) GetDebugString() string {
	return fmt.Sprintf("%t", b.Value)
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

////////////////////////////////////////////////////////////////////////////////
// FUNCTIONS
////////////////////////////////////////////////////////////////////////////////

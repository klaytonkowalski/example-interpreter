package object

////////////////////////////////////////////////////////////////////////////////
// DEPENDENCIES
////////////////////////////////////////////////////////////////////////////////

import "fmt"

////////////////////////////////////////////////////////////////////////////////
// VARIABLES
////////////////////////////////////////////////////////////////////////////////

const (
	ObjectInteger = "Integer"
	ObjectBoolean = "Boolean"
	ObjectNull    = "Null"
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

////////////////////////////////////////////////////////////////////////////////
// FUNCTIONS
////////////////////////////////////////////////////////////////////////////////

package object

////////////////////////////////////////////////////////////////////////////////
// STRUCTURES
////////////////////////////////////////////////////////////////////////////////

type Environment struct {
	store  map[string]Object
	parent *Environment
}

////////////////////////////////////////////////////////////////////////////////
// METHODS
////////////////////////////////////////////////////////////////////////////////

func (e *Environment) GetObject(key string) (Object, bool) {
	obj, ok := e.store[key]
	if !ok && e.parent != nil {
		obj, ok = e.parent.GetObject(key)
	}
	return obj, ok
}

func (e *Environment) SetObject(key string, obj Object) Object {
	e.store[key] = obj
	return obj
}

////////////////////////////////////////////////////////////////////////////////
// FUNCTIONS
////////////////////////////////////////////////////////////////////////////////

func CreateEnvironment() *Environment {
	m := make(map[string]Object)
	return &Environment{store: m, parent: nil}
}

func CreateClosureEnvironment(parent *Environment) *Environment {
	env := CreateEnvironment()
	env.parent = parent
	return env
}
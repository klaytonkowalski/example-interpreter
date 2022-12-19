package evaluator

import "github.com/klaytonkowalski/example-interpreter/object"

var natives = map[string]*object.Native{
	"len": {
		Function: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return createError("Wrong number of arguments to len(); got %d, expected %d.", len(args), 1)
			}
			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			default:
				return createError("Argument type to len() not supported; got %s, expected %s.", args[0].GetType(), object.ObjectString)
			}
		},
	},
	"first": {
		Function: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return createError("Wrong number of arguments to first(); got %d, expected %d.", len(args), 1)
			}
			if args[0].GetType() != object.ObjectArray {
				return createError("Argument type to first() not supported; got %s, expected %s", args[0].GetType(), object.ObjectArray)
			}
			array := args[0].(*object.Array)
			if len(array.Elements) > 0 {
				return array.Elements[0]
			}
			return Null
		},
	},
	"last": {
		Function: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return createError("Wrong number of arguments to last(); got %d, expected %d.", len(args), 1)
			}
			if args[0].GetType() != object.ObjectArray {
				return createError("Argument type to last() not supported; got %s, expected %s", args[0].GetType(), object.ObjectArray)
			}
			array := args[0].(*object.Array)
			if len(array.Elements) > 0 {
				return array.Elements[len(array.Elements)-1]
			}
			return Null
		},
	},
	"rest": {
		Function: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return createError("Wrong number of arguments to rest(); got %d, expected %d.", len(args), 1)
			}
			if args[0].GetType() != object.ObjectArray {
				return createError("Argument type to rest() not supported; got %s, expected %s", args[0].GetType(), object.ObjectArray)
			}
			array := args[0].(*object.Array)
			length := len(array.Elements)
			if length > 0 {
				newElements := make([]object.Object, length-1, length-1)
				copy(newElements, array.Elements[1:length])
				return &object.Array{Elements: newElements}
			}
			return Null
		},
	},
	"push": {
		Function: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return createError("Wrong number of arguments to push(); got %d, expected %d.", len(args), 2)
			}
			if args[0].GetType() != object.ObjectArray {
				return createError("Argument type to push() not supported; got %s, expected %s", args[0].GetType(), object.ObjectArray)
			}
			array := args[0].(*object.Array)
			length := len(array.Elements)
			newElements := make([]object.Object, length+1, length+1)
			copy(newElements, array.Elements)
			newElements[length] = args[1]
			return &object.Array{Elements: newElements}
		},
	},
}

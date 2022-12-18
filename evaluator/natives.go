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
			default:
				return createError("Argument type to len() not supported; got %s, expected %s.", args[0].GetType(), object.ObjectString)
			}
		},
	},
}

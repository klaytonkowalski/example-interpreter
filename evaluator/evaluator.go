package evaluator

////////////////////////////////////////////////////////////////////////////////
// DEPENDENCIES
////////////////////////////////////////////////////////////////////////////////

import (
	"fmt"

	"github.com/klaytonkowalski/example-interpreter/ast"
	"github.com/klaytonkowalski/example-interpreter/object"
)

////////////////////////////////////////////////////////////////////////////////
// VARIABLES
////////////////////////////////////////////////////////////////////////////////

var (
	True  = &object.Boolean{Value: true}
	False = &object.Boolean{Value: false}
	Null  = &object.Null{}
)

////////////////////////////////////////////////////////////////////////////////
// METHODS
////////////////////////////////////////////////////////////////////////////////

func Evaluate(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evaluateProgram(node, env)
	case *ast.ExpressionStatement:
		return Evaluate(node.Expression, env)
	case *ast.LetStatement:
		value := Evaluate(node.Expression, env)
		if isError(value) {
			return value
		}
		env.SetObject(node.Identifier.Value, value)
	case *ast.ReturnStatement:
		value := Evaluate(node.Expression, env)
		if isError(value) {
			return value
		}
		return &object.Return{Value: value}
	case *ast.PrefixExpression:
		rhsObject := Evaluate(node.RHSExpression, env)
		if isError(rhsObject) {
			return rhsObject
		}
		return evaluatePrefixExpression(node.Operator, rhsObject)
	case *ast.InfixExpression:
		lhsObject := Evaluate(node.LHSExpression, env)
		if isError(lhsObject) {
			return lhsObject
		}
		rhsObject := Evaluate(node.RHSExpression, env)
		if isError(rhsObject) {
			return rhsObject
		}
		return evaluateInfixExpression(node.Operator, lhsObject, rhsObject)
	case *ast.BlockStatement:
		return evaluateBlockStatement(node, env)
	case *ast.IfExpression:
		return evaluateIfExpression(node, env)
	case *ast.Integer:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return convertBoolToBoolean(node.Value)
	case *ast.Identifier:
		return evaluateIdentifier(node, env)
	case *ast.Function:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Body: body, Environment: env}
	case *ast.CallExpression:
		function := Evaluate(node.Function, env)
		if isError(function) {
			return function
		}
		args := evaluateExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(function, args)
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// FUNCTIONS
////////////////////////////////////////////////////////////////////////////////

func evaluateProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object
	for _, statement := range program.Statements {
		result = Evaluate(statement, env)
		switch result := result.(type) {
		case *object.Return:
			return result.Value
		case *object.Error:
			return result
		}
	}
	return result
}

func evaluateBlockStatement(bs *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object
	for _, statement := range bs.Statements {
		result = Evaluate(statement, env)
		if result != nil {
			if result.GetType() == object.ObjectReturn || result.GetType() == object.ObjectError {
				return result
			}
		}
	}
	return result
}

func evaluatePrefixExpression(operator string, rhsObject object.Object) object.Object {
	switch operator {
	case "!":
		return evaluateBangExpression(rhsObject)
	case "-":
		return evaluateMinusExpression(rhsObject)
	default:
		return createError("Unknown operator: %s%s", operator, rhsObject.GetType())
	}
}

func evaluateBangExpression(rhsObject object.Object) object.Object {
	switch rhsObject {
	case True:
		return False
	case False:
		return True
	case Null:
		return True
	default:
		return False
	}
}

func evaluateMinusExpression(rhsObject object.Object) object.Object {
	if rhsObject.GetType() != object.ObjectInteger {
		return createError("Wrong expression type: -%s", rhsObject.GetType())
	}
	value := rhsObject.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evaluateInfixExpression(operator string, lhsObject, rhsObject object.Object) object.Object {
	switch {
	case lhsObject.GetType() == object.ObjectInteger && rhsObject.GetType() == object.ObjectInteger:
		return evaluateIntegerExpression(operator, lhsObject, rhsObject)
	case operator == "==":
		return convertBoolToBoolean(lhsObject == rhsObject)
	case operator == "!=":
		return convertBoolToBoolean(lhsObject != rhsObject)
	case lhsObject.GetType() != rhsObject.GetType():
		return createError("Type mismatch: %s %s %s", lhsObject.GetType(), operator, rhsObject.GetType())
	default:
		return createError("Unknown operator: %s %s %s", lhsObject.GetType(), operator, rhsObject.GetType())
	}
}

func evaluateIntegerExpression(operator string, lhsObject, rhsObject object.Object) object.Object {
	lhsValue := lhsObject.(*object.Integer).Value
	rhsValue := rhsObject.(*object.Integer).Value
	switch operator {
	case "+":
		return &object.Integer{Value: lhsValue + rhsValue}
	case "-":
		return &object.Integer{Value: lhsValue - rhsValue}
	case "*":
		return &object.Integer{Value: lhsValue * rhsValue}
	case "/":
		return &object.Integer{Value: lhsValue / rhsValue}
	case "<":
		return convertBoolToBoolean(lhsValue < rhsValue)
	case ">":
		return convertBoolToBoolean(lhsValue > rhsValue)
	case "==":
		return convertBoolToBoolean(lhsValue == rhsValue)
	case "!=":
		return convertBoolToBoolean(lhsValue != rhsValue)
	default:
		return createError("Unknown operator: %s %s %s", lhsObject.GetType(), operator, rhsObject.GetType())
	}
}

func evaluateIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Evaluate(ie.Condition, env)
	if isError(condition) {
		return condition
	}
	if isTruthy(condition) {
		return Evaluate(ie.Then, env)
	}
	if ie.Else != nil {
		return Evaluate(ie.Else, env)
	}
	return Null
}

func evaluateIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	value, ok := env.GetObject(node.Value)
	if !ok {
		return createError("Identifier not found: %s", node.Value)
	}
	return value
}

func evaluateExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object
	for _, exp := range exps {
		evaluated := Evaluate(exp, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}
	return result
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	function, ok := fn.(*object.Function)
	if !ok {
		return createError("Not a function: %s", fn.GetType())
	}
	extendedEnv := extendFunctionEnvironment(function, args)
	evaluated := Evaluate(function.Body, extendedEnv)
	return unwrapReturnValue(evaluated)
}

func extendFunctionEnvironment(fn *object.Function, args []object.Object) *object.Environment {
	env := object.CreateClosureEnvironment(fn.Environment)
	for i, param := range fn.Parameters {
		env.SetObject(param.Value, args[i])
	}
	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.Return); ok {
		return returnValue
	}
	return obj
}

func convertBoolToBoolean(boolean bool) object.Object {
	if boolean {
		return True
	}
	return False
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case Null:
		return false
	case True:
		return true
	case False:
		return false
	default:
		return true
	}
}

func createError(message string, args ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(message, args...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.GetType() == object.ObjectError
	}
	return false
}

package evaluator

////////////////////////////////////////////////////////////////////////////////
// DEPENDENCIES
////////////////////////////////////////////////////////////////////////////////

import (
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

func Evaluate(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evaluateProgram(node)
	case *ast.ExpressionStatement:
		return Evaluate(node.Expression)
	case *ast.ReturnStatement:
		value := Evaluate(node.Expression)
		return &object.Return{Value: value}
	case *ast.PrefixExpression:
		rhsObject := Evaluate(node.RHSExpression)
		return evaluatePrefixExpression(node.Operator, rhsObject)
	case *ast.InfixExpression:
		lhsObject := Evaluate(node.LHSExpression)
		rhsObject := Evaluate(node.RHSExpression)
		return evaluateInfixExpression(node.Operator, lhsObject, rhsObject)
	case *ast.BlockStatement:
		return evaluateBlockStatement(node)
	case *ast.IfExpression:
		return evaluateIfExpression(node)
	case *ast.Integer:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return convertBoolToBoolean(node.Value)
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// FUNCTIONS
////////////////////////////////////////////////////////////////////////////////

func evaluateProgram(program *ast.Program) object.Object {
	var result object.Object
	for _, statement := range program.Statements {
		result = Evaluate(statement)
		if obj, ok := result.(*object.Return); ok {
			return obj.Value
		}
	}
	return result
}

func evaluateBlockStatement(bs *ast.BlockStatement) object.Object {
	var result object.Object
	for _, statement := range bs.Statements {
		result = Evaluate(statement)
		if result != nil && result.GetType() == object.ObjectReturn {
			return result
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
		return Null
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
		return Null
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
	default:
		return Null
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
		return Null
	}
}

func evaluateIfExpression(ie *ast.IfExpression) object.Object {
	condition := Evaluate(ie.Condition)
	if isTruthy(condition) {
		return Evaluate(ie.Then)
	}
	if ie.Else != nil {
		return Evaluate(ie.Else)
	}
	return Null
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

package orchestrator_test

import (
	"calc-LMS-orchestrator/internal/orchestrator"
	"math"
	"testing"
)

func TestAddExpression_SingleDigit(t *testing.T) {
	id, err := orchestrator.AddExpression("2+3")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	expr, found := orchestrator.GetExpression(id)
	if !found {
		t.Errorf("Expression not found")
	}
	if expr.Status != "pending" && expr.Status != "completed" {
		t.Errorf("Unexpected expression status: %s", expr.Status)
	}
}

func TestAddExpression_MultiDigit(t *testing.T) {
	_, err := orchestrator.AddExpression("12+3")
	if err == nil {
		t.Errorf("Expected error for multi-digit number, got nil")
	}
}

func TestAddExpression_Parentheses(t *testing.T) {
	id, err := orchestrator.AddExpression("(2+3)*4")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	expr, found := orchestrator.GetExpression(id)
	if !found {
		t.Errorf("Expression not found")
	}
	if expr.Status != "pending" && expr.Status != "completed" {
		t.Errorf("Unexpected expression status: %s", expr.Status)
	}
	if expr.Status == "pending" && !math.IsNaN(expr.Result) {
		t.Errorf("Expected NaN result for pending expression, got %v", expr.Result)
	}
}

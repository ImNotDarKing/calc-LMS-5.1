package orchestrator

import (
	"context"
	"math"
	"testing"

	"github.com/ImNotDarKing/calc-LMS-5.1/internal/db"
)

func resetState() {
	exprMutex.Lock()
	defer exprMutex.Unlock()
	expressions = make(map[int]*Expression)
	tasks = make(map[int]*Task)
	readyTasks = make([]*Task, 0)
	exprCounter = 0
	taskCounter = 0
}

func TestParseExpression_Simple(t *testing.T) {
	node, err := parseExpression("2+3*4 - (1+1)")
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	if node.Op != "-" {
		t.Errorf("expected root op '-', got %q", node.Op)
	}
	if node.Left.Op != "+" && node.Left.Op != "*" {
		t.Errorf("unexpected left node op %q", node.Left.Op)
	}
}

func TestBuildTasks_And_Complete(t *testing.T) {
	resetState()

	var ids []int
	val, rootID, err := buildTasks(&Node{
		Op: "-",
		Left: &Node{Op: "", Value: 5},
		Right: &Node{Op: "", Value: 2},
	}, 42, &ids)
	if err != nil {
		t.Fatalf("buildTasks error: %v", err)
	}
	if !math.IsNaN(val) {
		t.Error("expected NaN for non-leaf buildTasks return")
	}
	if len(ids) != 1 {
		t.Errorf("expected 1 task, got %d", len(ids))
	}
	taskID := rootID
	if err := CompleteTask(taskID, 3); err != nil {
		t.Fatalf("CompleteTask error: %v", err)
	}
	exprMutex.Lock()
	defer exprMutex.Unlock()
	if expressions[42].Result != 3 {
		t.Errorf("expected expr.Result=3, got %v", expressions[42].Result)
	}
}

func TestAddExpression_And_Queue(t *testing.T) {
	resetState()
	if err := db.InitDB(context.Background(), ":memory:"); err != nil {
		t.Fatalf("InitDB: %v", err)
	}
	if err := db.CreateTables(context.Background()); err != nil {
		t.Fatalf("CreateTables: %v", err)
	}

	id, err := AddExpression("7+8", 100)
	if err != nil {
		t.Fatalf("AddExpression: %v", err)
	}
	if id == 0 {
		t.Errorf("expected nonzero expr ID")
	}
	exprMutex.Lock()
	defer exprMutex.Unlock()
	if len(readyTasks) != 1 {
		t.Errorf("expected 1 ready task, got %d", len(readyTasks))
	}
	task := readyTasks[0]
	if task.Arg1 != 7 || task.Arg2 != 8 || task.Operator != "+" {
		t.Errorf("unexpected ready task content: %#v", task)
	}
}

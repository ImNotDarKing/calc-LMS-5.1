package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ImNotDarKing/calc-LMS-5.1/internal/db"
	"github.com/ImNotDarKing/calc-LMS-5.1/internal/server"
)

func TestFullWorkflow(t *testing.T) {
	if err := db.InitDB(context.Background(), ":memory:"); err != nil {
		t.Fatalf("InitDB: %v", err)
	}
	if err := db.CreateTables(context.Background()); err != nil {
		t.Fatalf("CreateTables: %v", err)
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		server.StartServer()
	}))
	defer ts.Close()

	client := ts.Client()

	registerPayload := map[string]string{"login": "u1", "password": "p1"}
	b, _ := json.Marshal(registerPayload)
	resp, err := client.Post(ts.URL+"/api/v1/register", "application/json", bytes.NewReader(b))
	if err != nil {
		t.Fatalf("register request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("register expected 200, got %d", resp.StatusCode)
	}
	var regResp struct{ Token string }
	if err := json.NewDecoder(resp.Body).Decode(&regResp); err != nil {
		t.Fatalf("decode register: %v", err)
	}
	token := regResp.Token
	if token == "" {
		t.Fatal("empty token on register")
	}

	calcReq := map[string]string{"expression": "1+2"}
	b, _ = json.Marshal(calcReq)
	req, _ := http.NewRequest("POST", ts.URL+"/api/v1/calculate", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("calculate request: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		body, _ := ioutil.ReadAll(resp.Body)
		t.Fatalf("calculate expected 201, got %d: %s", resp.StatusCode, body)
	}
	resp.Body.Close()

	req, _ = http.NewRequest("GET", ts.URL+"/internal/task", nil)
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("get task: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("get task expected 200, got %d", resp.StatusCode)
	}
	var taskResp struct {
		Task struct {
			ID     int     `json:"id"`
			Arg1   float64 `json:"arg1"`
			Arg2   float64 `json:"arg2"`
			OperationTime int `json:"operation_time"`
		} `json:"task"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&taskResp); err != nil {
		t.Fatalf("decode task: %v", err)
	}
	resp.Body.Close()

	time.Sleep(time.Duration(taskResp.Task.OperationTime) * time.Millisecond)
	result := taskResp.Task.Arg1 + taskResp.Task.Arg2
	postBody, _ := json.Marshal(map[string]interface{}{
		"id":     taskResp.Task.ID,
		"result": result,
	})
	resp, err = client.Post(ts.URL+"/internal/task", "application/json", bytes.NewReader(postBody))
	if err != nil {
		t.Fatalf("post task result: %v", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("post task result expected 200, got %d", resp.StatusCode)
	}

	req, _ = http.NewRequest("GET", ts.URL+"/api/v1/expressions", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("list expressions: %v", err)
	}
	defer resp.Body.Close()
	var listResp struct {
		Expressions []struct {
			ID         int     `json:"id"`
			Expression string  `json:"expression"`
			Result     float64 `json:"result"`
		} `json:"expressions"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&listResp); err != nil {
		t.Fatalf("decode list: %v", err)
	}
	if len(listResp.Expressions) != 1 {
		t.Fatalf("expected 1 expr, got %d", len(listResp.Expressions))
	}
	if listResp.Expressions[0].Result != 3 {
		t.Errorf("expected result=3, got %v", listResp.Expressions[0].Result)
	}
}

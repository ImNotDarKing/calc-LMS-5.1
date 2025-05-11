package api_v1_test

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/cookiejar"
    "net/http/httptest"
    "testing"

    "github.com/ImNotDarKing/calc-LMS-5.1/internal/api/v1"
)

func setupMux() *http.ServeMux {
    mux := http.NewServeMux()
    mux.HandleFunc("/api/v1/register", v1.RegisterHandler)
    mux.HandleFunc("/api/v1/login",    v1.LoginHandler)
    mux.HandleFunc("/api/v1/calculate",    v1.SubmitExpression)
    mux.HandleFunc("/api/v1/expressions/", v1.GetExpression)
    mux.HandleFunc("/api/v1/expressions",  v1.ListExpressions)
    mux.HandleFunc("/internal/task", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodGet {
            v1.GetTask(w, r)
        } else if r.Method == http.MethodPost {
            v1.PostTaskResult(w, r)
        } else {
            http.Error(w, `{"error": "Method Not Allowed"}`, http.StatusMethodNotAllowed)
        }
    })
    return mux
}

func newAuthClient(tsURL string, t *testing.T) *http.Client {
    jar, err := cookiejar.New(nil)
    if err != nil {
        t.Fatalf("cannot create cookie jar: %v", err)
    }
    client := &http.Client{Jar: jar}

    regBody, _ := json.Marshal(map[string]string{
        "username": "testuser",
        "password": "secret123",
    })
    resp, err := client.Post(tsURL+"/api/v1/register", "application/json", bytes.NewReader(regBody))
    if err != nil {
        t.Fatalf("register failed: %v", err)
    }
    resp.Body.Close()
    if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
        t.Fatalf("expected 200 or 201 on register, got %d", resp.StatusCode)
    }

    loginBody, _ := json.Marshal(map[string]string{
        "username": "testuser",
        "password": "secret123",
    })
    resp, err = client.Post(tsURL+"/api/v1/login", "application/json", bytes.NewReader(loginBody))
    if err != nil {
        t.Fatalf("login failed: %v", err)
    }
    resp.Body.Close()
    if resp.StatusCode != http.StatusOK {
        t.Fatalf("expected 200 on login, got %d", resp.StatusCode)
    }

    return client
}

func TestCalculateAndListExpressions(t *testing.T) {
    mux := setupMux()
    ts := httptest.NewServer(mux)
    defer ts.Close()

    client := newAuthClient(ts.URL, t)

    expr := map[string]string{"expression": "5*6"}
    body, _ := json.Marshal(expr)
    resp, err := client.Post(ts.URL+"/api/v1/calculate", "application/json", bytes.NewReader(body))
    if err != nil {
        t.Fatalf("POST /calculate error: %v", err)
    }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusOK {
        t.Fatalf("expected 200 OK, got %d", resp.StatusCode)
    }

    var res struct{ ID int }
    if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
        t.Fatalf("cannot parse JSON: %v", err)
    }
    if res.ID == 0 {
        t.Fatalf("expected non-zero ID, got %d", res.ID)
    }

    listResp, err := client.Get(ts.URL + "/api/v1/expressions")
    if err != nil {
        t.Fatalf("GET /expressions error: %v", err)
    }
    defer listResp.Body.Close()
    if listResp.StatusCode != http.StatusOK {
        t.Fatalf("expected 200 OK for list, got %d", listResp.StatusCode)
    }

    var list struct {
        Expressions []struct{ ID int `json:"id"` } `json:"expressions"`
    }
    if err := json.NewDecoder(listResp.Body).Decode(&list); err != nil {
        t.Fatalf("cannot parse list JSON: %v", err)
    }

    found := false
    for _, e := range list.Expressions {
        if e.ID == res.ID {
            found = true
            break
        }
    }
    if !found {
        t.Fatalf("ID %d not found in expressions list", res.ID)
    }
}

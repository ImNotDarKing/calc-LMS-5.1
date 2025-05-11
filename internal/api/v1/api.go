package v1

import (
	"encoding/json"
	"fmt"
	"io"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/ImNotDarKing/calc-LMS-5.1/internal/orchestrator"
	"github.com/ImNotDarKing/calc-LMS-5.1/internal/db"
)


func authenticate(r *http.Request) (int64, int, error) {
    auth := r.Header.Get("Authorization")
    if !strings.HasPrefix(auth, "Bearer ") {
        return 0, http.StatusUnauthorized, errors.New("missing token")
    }
    tokenStr := strings.TrimPrefix(auth, "Bearer ")

    claims := &jwt.RegisteredClaims{}
    tok, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
        return jwtKey, nil
    })
    if err != nil {
        var ve *jwt.ValidationError
        if errors.As(err, &ve) && ve.Errors&jwt.ValidationErrorExpired != 0 {
            return 0, http.StatusUnauthorized, fmt.Errorf("token expired")
        }
        return 0, http.StatusUnauthorized, fmt.Errorf("invalid token")
    }
    if !tok.Valid {
        return 0, http.StatusUnauthorized, fmt.Errorf("invalid token")
    }

    ok, err := db.IsJWTTokenValid(r.Context(), tokenStr)
    if err != nil || !ok {
        return 0, http.StatusUnauthorized, fmt.Errorf("invalid token")
    }

    uid, err := strconv.ParseInt(claims.Subject, 10, 64)
    if err != nil {
        return 0, http.StatusUnauthorized, fmt.Errorf("invalid subject")
    }
    return uid, http.StatusOK, nil
}

type ExpressionRequest struct {
	Expression string `json:"expression"`
}

type ExpressionResponse struct {
	ID int `json:"id"`
}

func SubmitExpression(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, `{"error":"Method Not Allowed"}`, http.StatusMethodNotAllowed)
        return
    }

    uid, status, authErr := authenticate(r)
    if authErr != nil {
        http.Error(w, fmt.Sprintf(`{"error":"%s"}`, authErr.Error()), status)
        return
    }

    var req ExpressionRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.Expression) == "" {
        http.Error(w, `{"error":"Invalid expression"}`, http.StatusUnprocessableEntity)
        return
    }

    id, err := orchestrator.AddExpression(req.Expression, uid)
    if err != nil {
        http.Error(w, fmt.Sprintf(`{"error":"Failed to add expression: %s"}`, err), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(ExpressionResponse{ID: id})
}

type ExpressionInfo struct {
    ID     int     `json:"id"`
    Expr   string  `json:"expression"`
    Result float64 `json:"result"`
}

func ListExpressions(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, `{"error":"Method Not Allowed"}`, http.StatusMethodNotAllowed)
        return
    }

    uid, status, authErr := authenticate(r)
    if authErr != nil {
        http.Error(w, fmt.Sprintf(`{"error":"%s"}`, authErr.Error()), status)
        return
    }

    exprs, err := db.SelectExpressionsByUser(r.Context(), uid)
    if err != nil {
        http.Error(w, `{"error":"Internal Server Error"}`, http.StatusInternalServerError)
        return
    }

    infos := make([]ExpressionInfo, len(exprs))
    for i, e := range exprs {
        infos[i] = ExpressionInfo{
            ID:     int(e.ID),
            Expr:   e.ExprText,  
            Result: e.Result,
        }
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{"expressions": infos})
}

func GetExpression(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, `{"error":"Method Not Allowed"}`, http.StatusMethodNotAllowed)
        return
    }

    uid, status, authErr := authenticate(r)
    if authErr != nil {
        http.Error(w, fmt.Sprintf(`{"error":"%s"}`, authErr.Error()), status)
        return
    }

    parts := strings.Split(r.URL.Path, "/")
    id, err := strconv.ParseInt(parts[len(parts)-1], 10, 64)
    if err != nil {
        http.Error(w, `{"error":"Invalid expression ID"}`, http.StatusBadRequest)
        return
    }

    expRec, err := db.GetExpressionByID(r.Context(), id)
    if err != nil {
        http.Error(w, `{"error":"Expression not found"}`, http.StatusNotFound)
        return
    }
    if expRec.UserID != uid {
        http.Error(w, `{"error":"Forbidden"}`, http.StatusForbidden)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "expression": map[string]interface{}{
            "id":         expRec.ID,
            "expression": expRec.ExprText,
            "result":     expRec.Result,
        },
    })
}

type TaskResponse struct {
	Task *orchestrator.Task `json:"task"`
}

func GetTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error": "Method Not Allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	task, err := orchestrator.GetReadyTask()
	if err != nil {
		http.Error(w, `{"error": "No task available"}`, http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(TaskResponse{Task: task})
}

type TaskResultRequest struct {
	ID     int     `json:"id"`
	Result float64 `json:"result"`
}

func PostTaskResult(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "Method Not Allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, `{"error": "Invalid request"}`, http.StatusBadRequest)
		return
	}
	var req TaskResultRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, `{"error": "Invalid JSON"}`, http.StatusUnprocessableEntity)
		return
	}
	err = orchestrator.CompleteTask(req.ID, req.Result)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, `{"error": "Task not found"}`, http.StatusNotFound)
		} else {
			http.Error(w, `{"error": "Failed to complete task"}`, http.StatusUnprocessableEntity)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

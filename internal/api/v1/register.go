package v1

import (
    "encoding/json"
    "net/http"
    "strconv"
    "time"

    "github.com/ImNotDarKing/calc-LMS-5.1/internal/db"
    "golang.org/x/crypto/bcrypt"
    "github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte("your-secret-key")

type registerReq struct {
    Login    string `json:"login"`
    Password string `json:"password"`
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, `{"error":"Method Not Allowed"}`, http.StatusMethodNotAllowed)
        return
    }

    var req registerReq
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, `{"error":"Bad Request"}`, http.StatusBadRequest)
        return
    }
    defer r.Body.Close()

    hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        http.Error(w, `{"error":"Internal Server Error"}`, http.StatusInternalServerError)
        return
    }

    user := &db.User{
        Login:        req.Login,
        PasswordHash: string(hash),
    }
    id, err := db.InsertUser(r.Context(), user)
    if err != nil {
        http.Error(w, `{"error":"Conflict"}`, http.StatusConflict)
        return
    }

    exp := time.Now().Add(1 * time.Hour)
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
        Subject:   strconv.FormatInt(id, 10),
        ExpiresAt: jwt.NewNumericDate(exp),
        IssuedAt:  jwt.NewNumericDate(time.Now()),
    })
    tstr, err := token.SignedString(jwtKey)
    if err != nil {
        http.Error(w, `{"error":"Internal Server Error"}`, http.StatusInternalServerError)
        return
    }

    jtok := &db.JWTToken{
        UserID:    id,
        Token:     tstr,
        ExpiresAt: exp,
    }
    if _, err := db.InsertJWTToken(r.Context(), jtok); err != nil {
        http.Error(w, `{"error":"Internal Server Error"}`, http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"token": tstr})
}
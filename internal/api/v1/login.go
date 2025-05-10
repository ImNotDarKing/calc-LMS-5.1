package v1

import (
    "net/http"
    "time"
    "strconv"
    "github.com/ImNotDarKing/calc-LMS-5.1/internal/db"

    "github.com/gin-gonic/gin"
    "golang.org/x/crypto/bcrypt"
    "github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte("your-secret-key") 

type loginReq struct {
    Login    string `json:"login"  binding:"required"`
    Password string `json:"password" binding:"required"`
}

func LoginHandler(c *gin.Context) {
    var r loginReq
    if c.BindJSON(&r) != nil {
        c.Status(http.StatusBadRequest)
        return
    }
    user, err := db.GetUserByLogin(c.Request.Context(), r.Login)
    if err != nil {
        c.Status(http.StatusUnauthorized)
        return
    }
    if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(r.Password)) != nil {
        c.Status(http.StatusUnauthorized)
        return
    }
    exp := time.Now().Add(24 * time.Hour)
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
        Subject:   strconv.FormatInt(user.ID, 10),
        ExpiresAt: jwt.NewNumericDate(exp),
        IssuedAt:  jwt.NewNumericDate(time.Now()),
    })
    tstr, err := token.SignedString(jwtKey)
    if err != nil {
        c.Status(http.StatusInternalServerError)
        return
    }
    jtok := &db.JWTToken{
        UserID:    user.ID,
        Token:     tstr,
        ExpiresAt: exp,
    }
    db.InsertJWTToken(c.Request.Context(), jtok)
    c.JSON(http.StatusOK, gin.H{"token": tstr})
}
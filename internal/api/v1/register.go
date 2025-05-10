package v1

import (
    "net/http"
    "github.com/ImNotDarKing/calc-LMS-5.1/internal/db"
    "github.com/gin-gonic/gin"
    "golang.org/x/crypto/bcrypt"
)

type registerReq struct {
    Login    string `json:"login"  binding:"required"`
    Password string `json:"password" binding:"required"`
}

func RegisterHandler(c *gin.Context) {
    var r registerReq
    if c.BindJSON(&r) != nil {
        c.Status(http.StatusBadRequest)
        return
    }
    hash, err := bcrypt.GenerateFromPassword([]byte(r.Password), bcrypt.DefaultCost)
    if err != nil {
        c.Status(http.StatusInternalServerError)
        return
    }
    user := &db.User{
        Login:        r.Login,
        PasswordHash: string(hash),
    }
    if _, err := db.InsertUser(c.Request.Context(), user); err != nil {
        c.Status(http.StatusConflict) 
        return
    }
    c.Status(http.StatusOK)
}

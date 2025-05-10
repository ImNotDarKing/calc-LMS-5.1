package v1

import (
    "net/http"
    "strconv"
    "github.com/ImNotDarKing/calc-LMS-5.1/internal/db"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v4"
)

func ComputeHandler(c *gin.Context) {
    hdr := c.GetHeader("Authorization")
    if len(hdr) < 7 || hdr[:7] != "Bearer " {
        c.Status(http.StatusUnauthorized)
        return
    }
    tokStr := hdr[7:]
    tok, err := jwt.ParseWithClaims(tokStr, &jwt.RegisteredClaims{}, func(t *jwt.Token) (any, error) {
        return jwtKey, nil
    })
    if err != nil || !tok.Valid {
        c.Status(http.StatusUnauthorized)
        return
    }
    claims := tok.Claims.(*jwt.RegisteredClaims)
    uid, _ := strconv.ParseInt(claims.Subject, 10, 64)
    ok, _ := db.IsJWTTokenValid(c.Request.Context(), tokStr)
    if !ok {
        c.Status(http.StatusUnauthorized)
        return
    }
    var rq struct{ Expr string `json:"expr" binding:"required"` }
    if c.BindJSON(&rq) != nil {
        c.Status(http.StatusBadRequest)
        return
    }
    // res := Evaluate(rq.Expr)
    exp := &db.Expression{
        UserID:   uid,
        ExprText: rq.Expr,
        Result:   res,
    }
    db.InsertExpression(c.Request.Context(), exp)
    c.JSON(http.StatusOK, gin.H{"result": res})
}

package db

import (
	"time"
	"context"
	"database/sql"
	
	_ "github.com/glebarez/sqlite"
)

type (
    User struct {
        ID           int64     `db:"id"`
        Login        string    `db:"login"`
        PasswordHash string    `db:"password_hash"`
        CreatedAt    time.Time `db:"created_at"`
    }

    Expression struct {
        ID       int64   `db:"id"`
        UserID   int64   `db:"user_id"`
        ExprText string  `db:"expr_text"`
        Result   float64 `db:"result"`
    }

    JWTToken struct {
        ID        int64     `db:"id"`
        UserID    int64     `db:"user_id"`
        Token     string    `db:"token"`
        IssuedAt  time.Time `db:"issued_at"`
        ExpiresAt time.Time `db:"expires_at"`
    }
)

var Conn *sql.DB

func InitDB(ctx context.Context, file string) error {
    var err error
    Conn, err = sql.Open("sqlite", file)
    if err != nil {
        return err
    }
    return Conn.PingContext(ctx)
}

func CreateTables(ctx context.Context) error {
    users := `
      CREATE TABLE IF NOT EXISTS users (
        id             INTEGER PRIMARY KEY AUTOINCREMENT,
        login          TEXT    NOT NULL UNIQUE,
        password_hash  TEXT    NOT NULL,
        created_at     DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
      );
    `
    exprs := `
      CREATE TABLE IF NOT EXISTS expressions (
        id         INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id    INTEGER NOT NULL REFERENCES users(id),
        expr_text  TEXT    NOT NULL,
        result     REAL    NOT NULL
      );
    `
    tokens := `
      CREATE TABLE IF NOT EXISTS jwt_tokens (
        id          INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id     INTEGER NOT NULL REFERENCES users(id),
        token       TEXT    NOT NULL UNIQUE,
        issued_at   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        expires_at  DATETIME NOT NULL
      );
    `
    if _, err := Conn.ExecContext(ctx, users); err != nil {
        return err
    }
    if _, err := Conn.ExecContext(ctx, exprs); err != nil {
        return err
    }
    if _, err := Conn.ExecContext(ctx, tokens); err != nil {
        return err
    }
    return nil
}

func InsertUser(ctx context.Context, u *User) (int64, error) {
    res, err := Conn.ExecContext(ctx,
        `INSERT INTO users(login,password_hash) VALUES(?,?)`,
        u.Login, u.PasswordHash,
    )
    if err != nil {
        return 0, err
    }
    return res.LastInsertId()
}

func GetUserByLogin(ctx context.Context, login string) (*User, error) {
    u := &User{}
    err := Conn.QueryRowContext(ctx,
        `SELECT id,login,password_hash,created_at FROM users WHERE login = ?`,
        login,
    ).Scan(&u.ID, &u.Login, &u.PasswordHash, &u.CreatedAt)
    if err != nil {
        return nil, err
    }
    return u, nil
}

func InsertExpression(ctx context.Context, e *Expression) (int64, error) {
    res, err := Conn.ExecContext(ctx,
        `INSERT INTO expressions(user_id,expr_text,result) VALUES(?,?,?)`,
        e.UserID, e.ExprText, e.Result,
    )
    if err != nil {
        return 0, err
    }
    return res.LastInsertId()
}

func SelectExpressionsByUser(ctx context.Context, userID int64) ([]Expression, error) {
    rows, err := Conn.QueryContext(ctx,
        `SELECT id,user_id,expr_text,result FROM expressions WHERE user_id = ?`,
        userID,
    )
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var out []Expression
    for rows.Next() {
        var e Expression
        if err := rows.Scan(&e.ID, &e.UserID, &e.ExprText, &e.Result); err != nil {
            return nil, err
        }
        out = append(out, e)
    }
    return out, nil
}

func InsertJWTToken(ctx context.Context, t *JWTToken) (int64, error) {
    res, err := Conn.ExecContext(ctx,
        `INSERT INTO jwt_tokens(user_id,token,expires_at) VALUES(?,?,?)`,
        t.UserID, t.Token, t.ExpiresAt,
    )
    if err != nil {
        return 0, err
    }
    return res.LastInsertId()
}

func IsJWTTokenValid(ctx context.Context, tok string) (bool, error) {
    var expiresAt time.Time
    err := Conn.QueryRowContext(ctx,
        `SELECT expires_at FROM jwt_tokens WHERE token = ?`,
        tok,
    ).Scan(&expiresAt)
    if err == sql.ErrNoRows {
        return false, nil
    }
    if err != nil {
        return false, err
    }
    return expiresAt.After(time.Now()), nil
}

func GetExpressionByID(ctx context.Context, id int64) (*Expression, error) {
    row := Conn.QueryRowContext(ctx,
        `SELECT id, user_id, expr_text, result FROM expressions WHERE id = ?`,
        id,
    )
    var e Expression
    if err := row.Scan(&e.ID, &e.UserID, &e.ExprText, &e.Result); err != nil {
        return nil, err
    }
    return &e, nil
}

func UpdateExpressionResult(ctx context.Context, id int64, result float64) error {
    _, err := Conn.ExecContext(ctx,
        `UPDATE expressions SET result = ? WHERE id = ?`,
        result, id,
    )
    return err
}
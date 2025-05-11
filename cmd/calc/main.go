package main

import (
	"context"
    "log"
    "github.com/ImNotDarKing/calc-LMS-5.1/internal/db"
	"github.com/ImNotDarKing/calc-LMS-5.1/internal/server"
)

func main() {
	ctx := context.Background()
    if err := db.InitDB(ctx, "calc.db"); err != nil {
        log.Fatalf("db.InitDB: %v", err)
    }
    if err := db.CreateTables(ctx); err != nil {
        log.Fatalf("db.CreateTables: %v", err)
    }
	server.StartServer()
}

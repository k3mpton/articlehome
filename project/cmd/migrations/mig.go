package main

import (
	"database/sql"
	"flag"
	dbconnection "legi/newspapers/project/utils/DbConnection"
	"log"
	"strings"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose"
)

var (
	MigMoving = flag.String("m", "up", "Moving migration, up or down")
)

func main() {
	flag.Parse()
	pool := dbconnection.NewConnection()
	defer pool.Close()

	mig := strings.ToLower(*MigMoving)
	if mig != "up" && mig != "down" {
		log.Fatalln("не удалось получить движение миграции")
	}

	db := sql.OpenDB(stdlib.GetPoolConnector(pool))
	defer db.Close()

	dir := "./migrations"
	switch mig {
	case "up":
		if err := goose.Up(db, dir); err != nil {
			log.Fatalf("failed up: %v", err)
		}
	default:
		if err := goose.Down(db, dir); err != nil {
			log.Fatalf("failed down: %v", err)
		}
	}
}

package config

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDatabase() {

	connStr := "postgresql://neondb_owner:npg_Q2bFa4wXgidk@ep-silent-tree-ad2ts1lh-pooler.c-2.us-east-1.aws.neon.tech/neondb?sslmode=require&channel_binding=require"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Database connected successfully")

	DB = db
}

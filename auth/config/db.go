package config

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
)

var DB *pgx.Conn

func ConnectDB() {
	conn, err := pgx.Connect(context.Background(), "postgresql://auth-service-db_owner:npg_zxt5liv2YoSn@ep-red-paper-a5qyrzu4-pooler.us-east-2.aws.neon.tech/auth-service-db?sslmode=require")
	if err != nil {
		panic(err)
	}

	DB = conn

	_, err = DB.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		email VARCHAR(255) NOT NULL,
		password VARCHAR(255) NOT NULL,
		first_name VARCHAR(255) NOT NULL,
		last_name VARCHAR(255) NOT NULL,
		emergency_contact VARCHAR(255) NOT NULL,
		date_of_birth DATE NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`)

	if err != nil {
		panic(err)
	}

	log.Println("Connected to database")
}

func CloseDB() {
	DB.Close(context.Background())
	log.Println("Closed database connection")
}

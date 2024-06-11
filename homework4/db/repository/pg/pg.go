package pg

import (
	"database/sql"
	"db/model"
	"fmt"
	"github.com/google/uuid"
	"log"

	_ "github.com/lib/pq" // PostgreSQL 驅動
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "mysecretpassword"
	dbname   = "postgres"
)

var pgDB *sql.DB

func GetDB() *sql.DB {
	if pgDB != nil {
		return pgDB
	}
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	pgDB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	err = pgDB.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected to PostgreSQL!")
	return pgDB
}

func Close() {
	err := pgDB.Close()
	if err != nil {
		log.Fatal(err)
		return
	}
	pgDB = nil
}

func CreateTable(db *sql.DB) error {
	query := `
        CREATE TABLE IF NOT EXISTS users (
            id VARCHAR(255) PRIMARY KEY,
            name VARCHAR(255) NOT NULL,
            email VARCHAR(255) UNIQUE NOT NULL
        );
    `
	_, err := db.Exec(query)
	if err != nil {
		log.Println(err)
		return err
	}
	fmt.Println("Table created successfully")
	return nil
}

func InsertRecord(db *sql.DB, name, email string) error {
	query := `
        INSERT INTO users (name, email, id)
        VALUES ($1, $2, $3);
    `
	_, err := db.Exec(query, name, email, uuid.New().String())
	if err != nil {
		log.Println(err)
		return err
	}
	fmt.Println("Record inserted successfully")
	return nil
}

func QueryRecords(db *sql.DB) ([]*model.User, error) {
	query := `
        SELECT id, name, email FROM users;
    `
	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(rows)

	users := make([]*model.User, 0)
	fmt.Println("Records:")
	for rows.Next() {
		var id, name, email string
		err = rows.Scan(&id, &name, &email)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		fmt.Println(id, name, email)
		users = append(users, &model.User{
			ID:    id,
			Name:  name,
			Email: email,
		})
	}
	return users, nil
}

func UpdateRecord(db *sql.DB, id, email string) error {
	query := `
        UPDATE users
        SET email = $1
        WHERE id = $2;
    `
	_, err := db.Exec(query, email, id)
	if err != nil {
		log.Println(err)
		return err
	}
	fmt.Println("Record updated successfully")
	return nil
}

func DeleteRecord(db *sql.DB, id string) error {
	query := `
        DELETE FROM users
        WHERE id = $1;
    `
	_, err := db.Exec(query, id)
	if err != nil {
		log.Println(err)
		return err
	}
	fmt.Println("Record deleted successfully")
	return nil
}

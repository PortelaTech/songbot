package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"os"
	"strings"
)


var dbInfo = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbname, sslmode)

//Collecting data from bot
func CollectData(username string, chatid int64, message string, answer []string) error {
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return err
	}
	defer db.Close()
	answ := strings.Join(answer, ", ")
	data := `INSERT INTO users(username, chat_id, message, answer) VALUES($1, $2, $3, $4);`
	if _, err = db.Exec(data, `@`+username, chatid, message, answ); err != nil {
		return err
	}
	return nil
}

//Creating users table in database
func CreateTable() error {
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return err
	}
	defer db.Close()
	if _, err = db.Exec(`CREATE TABLE users(ID SERIAL PRIMARY KEY, TIMESTAMP TIMESTAMP DEFAULT CURRENT_TIMESTAMP, USERNAME TEXT, CHAT_ID INT, MESSAGE TEXT, ANSWER TEXT);`); err != nil {
		return err
	}
	return nil
}

//Getting number of users who using bot
func GetNumberOfUsers() (int64, error) {
	var count int64
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return 0, err
	}
	defer db.Close()
	row := db.QueryRow("SELECT COUNT(DISTINCT username) FROM users;")
	err = row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

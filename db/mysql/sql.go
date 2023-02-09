package sqlconnection

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"
)

var close chan bool

func createConnection(dsn string) (*sql.DB, error) {
	fmt.Println("entering createConnection")
	close <- false
	// create connection
	db, err := sql.Open("mysql", dsn)
	fmt.Println("connection created")
	if err != nil {
		return &sql.DB{}, fmt.Errorf("error while creating connection: %v", err)
	}
	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Second * 5)
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(0)
	fmt.Println("connection set")

	go func() {
		fmt.Println("waiting for close signal")
		select {
		case <-close:
			db.Close()
			log.Println("Closed connection after writing data successfully")
		default:
			db.Close()
			log.Println("Connection timed out, failed to write data")
		}
		fmt.Println("exiting go routine")
	}()

	fmt.Println("exiting createConnection")
	return db, nil
}

func pushUserToUsers(user struct {
	Name     string
	Email    string
	Password string
	Phone    string
}) error {
	db, err := createConnection("")
	fmt.Println("connection created")
	if err != nil {
		return fmt.Errorf("error while creating connection: %v", err)
	}

	// Insert the password record
	passwordResult, err := db.Exec("INSERT INTO passwords (password, salt) VALUES (?, ?)", user.Password, "user_salt")
	fmt.Println("password inserted")
	close <- true
	if err != nil {
		return fmt.Errorf("error while inserting password: %v", err)
	}

	// Retrieve the last inserted password ID
	passwordID, err := passwordResult.LastInsertId()
	fmt.Println("password id retrieved")
	if err != nil {
		return fmt.Errorf("error while retrieving last inserted password ID: %v", err)
	}

	// Insert the user record
	_, err = db.Exec("INSERT INTO users (name, email, phone, password_id) VALUES (?, ?, ?, ?)", user.Name, user.Email, user.Phone, passwordID)
	fmt.Println("user inserted")
	close <- true
	if err != nil {
		return fmt.Errorf("error while inserting user data: %v", err)
	}
	fmt.Println("exiting database connection")
	return nil
}

func SearchUserInUsers(user struct {
	Name     string
	Email    string
	Password string
	Phone    string
}) error {
	fmt.Println("entered SearchUserInUsers")
	db, err := createConnection("")
	fmt.Println("connection created")
	if err != nil {
		return fmt.Errorf("error while creating connection: %v", err)
	}

	// Search the user record
	rows, err := db.Query("SELECT * FROM users WHERE email = ? OR phone = ?", user.Email, user.Phone)
	fmt.Println("search query executed")
	close <- true
	if err != nil {
		return fmt.Errorf("error while checking if user exists: %v", err)
	}
	defer rows.Close()

	// Iterate over the rows
	if rows.Next() {
		fmt.Println("user already exists")
		return errors.New("user already exists")
	} else {
		err := pushUserToUsers(user)
		fmt.Println("user added to db")
		if err != nil {
			return fmt.Errorf("error while adding user to db: %v", err)
		}
	}
	fmt.Println("exiting database connection")
	return nil
}

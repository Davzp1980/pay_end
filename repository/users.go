package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"pay/internal"

	"golang.org/x/crypto/bcrypt"
)

func CreateAdmin(db *sql.DB, name, hashedPassword string) (string, error) {
	isAdmin := true
	var user internal.User

	err := db.QueryRow("INSERT INTO users (name, password, is_admin) VALUES ($1,$2,$3) Returning id", name, hashedPassword, isAdmin).Scan(&user.Name)
	if name == "" {
		return "", errors.New("invalid user's name")
	}
	if err != nil {

		return "", fmt.Errorf("user %s already exists", name)

	}
	return fmt.Sprintf("Admin %s created", user.Name), nil
}

func CreateUser(db *sql.DB, name, hashedPassword string) (string, error) {
	isAdmin := false

	if name == "" {
		return "", errors.New("invalid user's name")
	}

	_, err := db.Query("INSERT INTO users (name, password, is_admin) VALUES ($1,$2,$3)", name, hashedPassword, isAdmin)
	if err != nil {
		return "", fmt.Errorf("user %s already exists", name)

	}
	return fmt.Sprintf("User %s created", name), nil
}

//Поля в postman:
// "name"

func BlockUser(db *sql.DB, name string) (string, error) {
	var user internal.User

	err := db.QueryRow("SELECT name FROM users WHERE name=$1", name).Scan(&user.Name)
	if err != nil {

		return "", err
	}
	expectedName := user.Name
	input_Name := name

	if input_Name == expectedName {
		_, err := db.Exec("UPDATE users SET blocked=$1 WHERE name=$2", true, name)
		if err != nil {
			return "", err
		}
	} else {
		return "", fmt.Errorf("user %s does not exists", name)
	}

	return fmt.Sprintf("User %s blocked", name), nil
}

func UnBlockUser(db *sql.DB, name string) (string, error) {
	var user internal.User

	err := db.QueryRow("SELECT name FROM users WHERE name=$1", name).Scan(&user.Name)
	if err != nil {

		return "", err
	}
	expectedName := user.Name
	input_Name := name

	if input_Name == expectedName {
		_, err := db.Exec("UPDATE users SET blocked=$1 WHERE name=$2", false, name)
		if err != nil {
			return "", err
		}
	} else {
		return "", fmt.Errorf("user %s does not exists", name)
	}

	return fmt.Sprintf("User %s Unblocked", name), nil
}

func ChangeUserPassword(db *sql.DB, name, hashedPassword string) (string, error) {
	_, err := db.Exec("UPDATE users SET password=$1 WHERE name=$2;", hashedPassword, name)
	if err != nil {
		log.Println("Change password error")
		return "", fmt.Errorf("password change error")
	}

	return fmt.Sprintf("Password %s changed", name), nil
}

func GetUser(db *sql.DB, username, password string) (string, error) {
	var user internal.User
	err := db.QueryRow("SELECT * FROM users WHERE name=$1", username).Scan(
		&user.ID, &user.Name, &user.PasswordHash, &user.IsAdmin, &user.Blocked)
	if err != nil {
		log.Println("User does not exists", err)
		return "", err
	}

	if !CheckPassword(password, user.PasswordHash) || username != user.Name {
		log.Println("Wrong password or user name")
		return "", errors.New("wrong password or user name")
	}
	if user.Blocked {
		log.Printf("User %s is blocked", username)
		return "", errors.New("user is blocked")
	}
	return username, nil
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

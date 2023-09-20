package service

import (
	"database/sql"
	"log"
	"math/rand"
	"pay/internal"

	"strconv"
)

func GenerateIban(name string) string {
	i := strconv.Itoa(rand.Intn(1000000000))
	iban := i + name

	return iban
}

func CheckIban(db *sql.DB, inputIban string) (bool, error) {
	var account internal.Account

	err := db.QueryRow("SELECT iban, blocked FROM accounts WHERE iban=$1", inputIban).Scan(&account.Iban, &account.Blocked)
	if err != nil {

		log.Println(err)
		return false, err
	}
	expectedIban := account.Iban

	if inputIban == expectedIban {
		return true, nil
	} else {
		log.Println("Account", inputIban, "Does not exists")
		return false, err

	}

}

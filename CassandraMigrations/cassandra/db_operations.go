package cassandra

import (
	"Alex_workshop/CassandraMigrations/cassandra/entities"
	"log"
)

var (
	//database manager instance
	dbManager    = GetInstance()
	name         string
	age          int
	address      string
	emailId      string
	password     string
	balance      int
	insertStatus bool
	deleteStatus bool
)

const (
	//Query statements for Account transactions
	InsertAccountQuery     = "INSERT INTO accounts (name,age,address,emailId,password,balance) VALUES (?,?,?,?,?,?)"
	DeleteAccountQuery     = "DELETE FROM accounts WHERE emailId = ?"
	DeleteAllAccountsQuery = "TRUNCATE accounts"
	FindAccountQuery       = "SELECT emailid FROM accounts WHERE emailId = ?LIMIT 1 ALLOW FILTERING"
	FindPasswordQuery      = "SELECT password FROM accounts WHERE emailId = ?LIMIT 1 ALLOW FILTERING"
	FindBalanceQuery      = "SELECT balance FROM accounts WHERE emailId = ?LIMIT 1 ALLOW FILTERING"
	FindAllAccountsQuery   = "SELECT * FROM accounts"
	UpdateBalanceQuery     = "UPDATE accounts SET balance = ? WHERE emailId = ?"
)

/*
 These functions are responsible for handling the database operations
 Function names will describe the operation it performs
*/

// Database operations for inserting a Account
func InsertAccount(user entities.User) bool {
	query := dbManager.Session.Query(InsertAccountQuery, user.GetName(), user.GetAge(), user.GetAddress(), user.GetEmailID(), user.GetPassword(), user.GetBalance())
	if err := query.Exec(); err != nil {
		log.Println("Account entry Failed error: %v", err)
		insertStatus = false
	} else {
		log.Println("Account entry successful")
		insertStatus = true
	}
	return insertStatus
}

// Database operations for deleting a Account where emailId is provided
func DeleteAccountWithEmailID(email string) bool {
	query := dbManager.Session.Query(DeleteAccountQuery, email)
	if err := query.Exec(); err != nil {
		log.Println("Deleting Account Failed error: %v", err)
		deleteStatus = false
	} else {

		log.Println("Deleting Account successful")
		deleteStatus = true
	}
	return deleteStatus
}

// Database operations for deleting all Accounts stored in database
func DeleteAllAccounts() {
	query := dbManager.Session.Query(DeleteAllAccountsQuery)
	if err := query.Exec(); err != nil {
		log.Println("Failed to delete all Accounts from database", err)
	} else {
		log.Println("Accounts deletion successful")

	}
}

// Database operations for updating Account balance where Email Id is provided
func UpdateBalance(bal int, email string) {
	query := dbManager.Session.Query(UpdateBalanceQuery, bal, email)
	if err := query.Exec(); err != nil {
		log.Println("Failed to updateBalance", err)
	} else {
		log.Println("Balance updation successful")
	}
}

// Database operations for finding Account where Email Id is provided
func FindAccountWithEmailId(email string) bool {
	query := dbManager.Session.Query(FindAccountQuery, email)
	if err := query.Scan(&emailId); err != nil {
			return false
		} else {
			return true
		}
	}

// Database operations for fetching password for a Account in database
func GetPassword(email string) string {
	query := dbManager.Session.Query(FindPasswordQuery, email)
	if err := query.Scan(&password); err != nil {
		log.Println("password retreival failed", err)
		return ""
	}
	return password
}

// Database operations for fetching Balance for a Account in database
func GetBalance(email string) int {
	query := dbManager.Session.Query(FindBalanceQuery, email)
	if err := query.Scan(&balance); err != nil {
		log.Println("Balance retreival failed", err)
		return 0
	}
	return balance
}

// Database operations for finding all Account entries in database
func FindAllAccounts() {
	query := dbManager.Session.Query(FindAllAccountsQuery)
	iterator := query.Iter()
	for iterator.Scan(&name, &age, &address, &emailId, &password, &balance) {
		log.Println(name, age, address, emailId, password, balance)
	}
}

package entities

type User struct {
	name     string
	age      int
	address  string
	emailId  string
	password string
	balance  int
}

func (u *User) SetName(name string) {
	u.name = name
}

func (u *User) SetAge(age int) {
	u.age = age
}

func (u *User) SetAddress(address string) {
	u.address = address
}

func (u *User) SetEmailID(emailId string) {
	u.emailId = emailId
}

func (u *User) SetPassword(password string) {
	u.password = password
}

func (u *User) SetBalance(balance int) {
	u.balance = balance
}

func (u User) GetName() string {
	return u.name
}

func (u User) GetAge() int {
	return u.age
}

func (u User) GetAddress() string {
	return u.address
}

func (u User) GetEmailID() string {
	return u.emailId
}

func (u User) GetPassword() string {
	return u.password
}

func (u User) GetBalance() int {
	return u.balance
}

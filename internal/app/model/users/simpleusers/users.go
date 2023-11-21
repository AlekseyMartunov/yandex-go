// Package simpleusers uses when db dose dot exists
package simpleusers

// Users interface show witch methods struct need
type Users interface {
	SaveNewUser() (int, error)
}

// User store count of users in app
type User struct {
	id int
}

// NewUser create new User instance
func NewUser() *User {
	return &User{
		id: 0,
	}
}

// SaveNewUser update count of users
func (u *User) SaveNewUser() (int, error) {
	u.id++
	return u.id, nil
}

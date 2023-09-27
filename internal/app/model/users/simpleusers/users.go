package simpleusers

type Users interface {
	SaveNewUser() (int, error)
}

type User struct {
	id int
}

func NewUser() *User {
	return &User{
		id: 0,
	}
}

func (u *User) SaveNewUser() (int, error) {
	u.id++
	return -1, nil
}

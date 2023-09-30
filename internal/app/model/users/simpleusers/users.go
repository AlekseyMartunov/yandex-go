package simpleusers

type Users interface {
	GetFreeID() (int, error)
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

func (u *User) GetFreeID() (int, error) {
	return u.id + 1, nil
}

func (u *User) SaveNewUser() (int, error) {
	u.id++
	return -1, nil
}

package entities

import "golang.org/x/crypto/bcrypt"

type User struct {
	ID            int64  `json:"-"`
	Login         string `json:"login"`
	Password      string `json:"password"`
	CryptPassword []byte `json:"-"`
}

func (u *User) GenerateCryptPassword() error {
	cryptPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.CryptPassword = cryptPassword

	return nil
}

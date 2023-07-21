package entities

type User struct {
	ID            int64  `json:"-"`
	Login         string `json:"login"`
	Password      string `json:"password"`
	CryptPassword []byte `json:"-"`
}

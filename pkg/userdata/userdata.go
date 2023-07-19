package userdata

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/Serhii1Epam/enhanceHttpServer/pkg/appdb"
	"github.com/Serhii1Epam/enhanceHttpServer/pkg/hasher"
)

type UserData struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

type IParse interface {
	ParseBody() *UserData
}

type JsonBytes []byte
type PlainTextBytes []byte

func NewUser() *UserData {
	return &UserData{}
}

func (u UserData) Print() {
	fmt.Printf("User: [%s], Password: [%s]\n", u.User, u.Password)
}

func Parse(i IParse) *UserData {
	return i.ParseBody()
}

func (in JsonBytes) ParseBody() *UserData {
	inUser := NewUser()
	json.Unmarshal(in, &inUser)
	return inUser
}

func (in PlainTextBytes) ParseBody() *UserData {
	var data []byte = in
	inUser := NewUser()
	str := strings.Fields(string(data))
	inUser.User = str[0]
	inUser.Password = str[1]
	return inUser
}

func (u *UserData) Create(db *appdb.Database) error {
	if db == nil {
		return errors.New("Database isn't initialized")
	}
	hsr := hasher.NewHasher(u.Password)
	hsr.HashPassword()
	if err := db.Insert(*hsr, u.User); err != nil {
		return err
	}
	u.Print()
	db.Print()
	return nil
}

func (u *UserData) Login(db *appdb.Database) error {
	if db == nil {
		return errors.New("Database isn't initialized")
	}
	hsr := hasher.NewHasher(u.Password)
	if !hsr.CheckPasswordHash(db.Select(u.User)) {
		customMsg := fmt.Sprintf("Invalid User:[%v], from DB [%s]", u, db.Select(u.User))
		return errors.New(customMsg)
	}
	return nil
}

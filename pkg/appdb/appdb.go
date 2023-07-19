package appdb

import (
	"errors"
	"fmt"

	"github.com/Serhii1Epam/enhanceHttpServer/pkg/hasher"
)

type Database struct {
	UserTable map[string]string
}

func NewDatabase() *Database {
	return &Database{UserTable: make(map[string]string)}
}

func (d *Database) Insert(h hasher.HashingData, user string) error {
	// Check and validate User and hasher
	if user == "" {
		return errors.New("User can't be nil")
	}

	if &h == nil || h.Hash == "" {
		return errors.New("Hash is absent")
	}

	d.UserTable[user] = h.Hash
	return nil
}

func (d Database) Select(user string) string {
	return d.UserTable[user]
}

func (d Database) Print() {
	for k, v := range d.UserTable {
		fmt.Println("[", k, "] value is [", v, "]")
	}
}

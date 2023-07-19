package appdb_test

import (
	"strings"
	"testing"

	"github.com/Serhii1Epam/enhanceHttpServer/pkg/appdb"
	"github.com/Serhii1Epam/enhanceHttpServer/pkg/hasher"
)

func TestDataBase_Create(t *testing.T) {
	db := appdb.NewDatabase()
	if db == nil {
		t.Errorf("Cant create DB!")
	}
}

func TestDataBase_MyInsert(t *testing.T) {
	db := appdb.NewDatabase()
	if db == nil {
		t.Errorf("Cant create DB!")
	}
	want := "eddef9e8e578c2a560c3187c4152c8b6f3f90c1dcf8c88b386ac1a9a96079c2c"

	//Filled hasher test
	hasher1 := hasher.NewHasher("TestPass")
	hasher1.HashPassword()
	if err := db.Insert(*hasher1, "TestUser"); err != nil {
		t.Errorf("Insert failed! %s [%v]", err, hasher1)
	}

	if strings.Compare(db.UserTable["TestUser"], want) != 0 {
		t.Errorf("Data from DB is differ,\nDB[%s],\nwant[%s]\n", db.UserTable["Serg"], want)
	}

	//Empty hasher, user tests
	var hasher2 = &hasher.HashingData{}
	if err := db.Insert(*hasher2, "TestUser"); err == nil {
		t.Errorf("Empty hasher test failed!")
	}

	if err := db.Insert(*hasher2, ""); err == nil {
		t.Errorf("Empty user test failed!")
	}

}

func TestDataBase_MySelect(t *testing.T) {
	db := appdb.NewDatabase()
	if db == nil {
		t.Errorf("Cant create DB!")
	}
	want := "eddef9e8e578c2a560c3187c4152c8b6f3f90c1dcf8c88b386ac1a9a96079c2c"
	hasher := hasher.NewHasher("TestPass")
	hasher.HashPassword()
	if err := db.Insert(*hasher, "TestUser"); err != nil {
		t.Errorf("Insert failed! %s [%v]", err, hasher)
	}

	if strings.Compare(db.Select("TestUser"), want) != 0 {
		t.Errorf("Select failed! %s", db.UserTable["TestUser"])
	}
}

func TestDatabase_Print(t *testing.T) {
	type fields struct {
		UserTable map[string]string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := appdb.Database{
				UserTable: tt.fields.UserTable,
			}
			d.Print()
		})
	}
}

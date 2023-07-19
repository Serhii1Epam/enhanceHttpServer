package userdata_test

import (
	"strings"
	"testing"

	"github.com/Serhii1Epam/enhanceHttpServer/pkg/appdb"
	"github.com/Serhii1Epam/enhanceHttpServer/pkg/userdata"
)

func TestParse(t *testing.T) {
	jb := userdata.JsonBytes("{\"user\":\"TestUser\", \"password\":\"TestPass\"}")
	ptb := userdata.PlainTextBytes("TestUser TestPass")
	userFromJson := jb.ParseBody()
	userFromPlainText := ptb.ParseBody()

	if strings.Compare(userFromJson.User, "TestUser") != 0 ||
		strings.Compare(userFromJson.Password, "TestPass") != 0 {
		t.Errorf("Json parsing error!")
	}

	if strings.Compare(userFromPlainText.User, "TestUser") != 0 ||
		strings.Compare(userFromPlainText.Password, "TestPass") != 0 {
		t.Errorf("Plain text parsing error!")
	}

}

func TestUserData_Create(t *testing.T) {
	var str userdata.IParse
	var db *appdb.Database = nil

	str = userdata.JsonBytes("{\"user\":\"TestUser\", \"password\":\"TestPass\"}")
	usr := userdata.Parse(str)

	if strings.Compare(usr.User, "TestUser") != 0 ||
		strings.Compare(usr.Password, "TestPass") != 0 {
		t.Errorf("Interface Parse() parsing error!")
	}

	//Not initialized DB test
	if err := usr.Create(db); err == nil {
		t.Errorf("Db initialization error.")
	}

	db = appdb.NewDatabase()
	if err := usr.Create(db); err != nil {
		t.Errorf("User creation error. %s", err)
	}
}

func TestUserData_Login(t *testing.T) {
	ptb := userdata.PlainTextBytes("TestUser TestPass")
	usr := ptb.ParseBody()
	var db *appdb.Database = nil

	//Not initialized DB test
	if err := usr.Login(db); err == nil {
		t.Errorf("Db initialization error.")
	}

	db = appdb.NewDatabase()
	if err := usr.Create(db); err != nil {
		t.Errorf("Db connection error. %s", err)
	}

	if err := usr.Login(db); err != nil {
		t.Errorf("Login user error. %s", err)
	}
}

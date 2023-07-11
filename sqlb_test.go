package sqlb

import (
	"database/sql/driver"
	"fmt"

	_ "github.com/17e10/go-sqlb/dialect/mysql"
	_ "github.com/go-sql-driver/mysql"
)

type person struct {
	Id         uint64 `json:"id"`
	FamilyName string `json:"family_name"`
	GivenName  string `json:"given_name"`
	Age        int    `json:"age"`
}

var (
	olivia  = person{1, "Williams", "Olivia", 54}
	kenny   = person{2, "Loggins", "Kenny", 75}
	persons = []person{olivia, kenny}

	_ = persons
)

type noint int

type testValuer string

func (s testValuer) Value() (driver.Value, error) {
	return string(s), nil
}

type errValuer string

func (s errValuer) Value() (driver.Value, error) {
	return "", fmt.Errorf(string(s))
}

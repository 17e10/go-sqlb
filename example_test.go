package sqlb_test

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/17e10/go-sqlb"
)

func ExampleM() {
	var extra []sqlb.Sqler

	// 追加条件1 オーストラリア出身
	extra = append(extra, sqlb.T("birthplace = @", "Australia"))

	// 追加条件2 A から始まる名前
	extra = append(extra, sqlb.T("(@ <= given_name AND given_name < @)", "A", "B"))

	query := sqlb.M(`
		SELECT
			given_name, family_name
		FROM
			person
		WHERE
			gen = @gen
		AND $extra
	`, map[string]any{
		"@gen":   3,
		"$extra": sqlb.And(extra...),
	})

	s, _ := sqlb.Stringify(query)
	fmt.Println(sqlb.Compact(s))
	// Output: SELECT given_name, family_name FROM person WHERE gen = 3 AND birthplace = 'Australia' AND ('A' <= given_name AND given_name < 'B')
}

func ExampleT() {
	type person struct {
		Id         uint64
		GivenName  string
		FamilyName string
		Age        int
	}
	persons := []person{
		{1, "Olivia", "Williams", 54},
		{2, "Kenny", "Loggins", 75},
	}

	query := sqlb.T(
		"INSERT INTO person (#) VALUES @",
		sqlb.Columns((*person)(nil), "id"),
		sqlb.GroupValues(persons, "id"),
	)

	s, _ := sqlb.Stringify(query)
	fmt.Println(sqlb.Compact(s))
	// Output: INSERT INTO person (`given_name`, `family_name`, `age`) VALUES ('Olivia', 'Williams', 54), ('Kenny', 'Loggins', 75)
}

var (
	ctx  context.Context
	conn *sql.Conn
)

func ExampleScan() {
	var person struct {
		ID         uint64
		GivenName  string
		FamilyName string
	}

	// クエリを実行する
	// SELECT id, given_name, family_name FROM person ORDER BY id
	rows, err := sqlb.Query(conn, ctx, sqlb.T(
		"SELECT # FROM person ORDER BY id",
		sqlb.Columns(&person),
	))
	if err != nil {
		return
	}
	defer rows.Close()

	// Rows を読む
	for rows.Next() {
		if err = sqlb.Scan(rows, &person); err != nil {
			return
		}
		fmt.Printf("%#v\n", &person)
	}
}

# go-sqlb

[![GoDev][godev-image]][godev-url]

go-sqlb パッケージは SQL 操作のユーティリティを提供します.

## Usage

```go
import "github.com/17e10/go-sqlb"

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
sqlb.Compact(s) == `SELECT given_name, family_name FROM person WHERE gen = 3 AND birthplace = 'Australia' AND ('A' <= given_name AND given_name < 'B')`
```

## License

This software is released under the MIT License, see LICENSE.

## Author

17e10

[godev-image]: https://pkg.go.dev/badge/github.com/17e10/go-sqlb
[godev-url]: https://pkg.go.dev/github.com/17e10/go-sqlb

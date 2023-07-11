// sqlb パッケージは SQL 操作のユーティリティを提供します.
//
// sqlb は ORM (Object-Relational Mapping) のような大規模な仕組みではなく、
// Go 標準パッケージの database/sql を容易に扱うためのユーティリティ群です.
//
// MySQL で使用する場合 次の import を追加してください.
//
//	import _ "github.com/17e10/b/sqlb/dialect/mysql"
//
// # Sqler
//
// SQL を動的に組み立てる仕組みに Sqler インターフェイスを導入しています.
//
//	type Sqler interface {
//		Sql(w Writer) error
//	}
//
// Sqler は SQL 文字列が必要になったときに Writer に SQL を出力していくことで
// 高速かつメモリ効率よく SQL を生成します.
//
// # SQL 記述
//
// 動的な SQL を簡潔に記述できる T 関数, M 関数を提供しています.
// これらの関数は SQL テンプレートに値や識別子, SQL を展開をできます.
// 構文は非常にシンプルながらとても柔軟かつ強力です.
//
// 値の展開: @
//
// １つの値を展開:
//
//	nil			T("@", nil)							NULL
//	整数			T("@", 123)							123
//	浮動小数点	T("@", 123.45)						123.45
//	文字列		T("@", "abc")						'abc'
//	ブール値		T("@", false)						FALSE
//	バイト列		T("@", []byte{0x41, 0x42, 0x43})	X'414243'
//	日付時刻		T("@", time.Time{})					'2001-01-02 13:14:15.678901'
//
// 特別な展開:
//
//	値リスト				T("@", []any{"a", "b"})						'a', 'b'
//	グループリスト			T("@", [][]any{{"a", "b"}, {"c", "d"}})		('a', 'b'), ('c', 'd')
//	Key-Value ペアリスト	T("@", []Kv{{"k1", "v1"}, {"k2", "v2"}})	`k1` = 'v1', `k2` = 'v2'
//
// 識別子の展開: #
//
// 1つの識別子を展開:
//
//	識別子		T("#", "a")					`a`
//	フィールドリスト	T("#", []string{"a", "b"})	`a`, `b`
//
// SQL の展開: $
//
//	SQL		T("WHERE $", sqler)		WHERE `a` = b
//
// 擬似イコール構文: == @, !== @
//
//	=				T("== @", "a")				= 'a'
//	!=				T("!== @", "a")				!= 'a'
//	IN (...)		T("== @", []any{"a", "b"})	IN ('a', 'b')
//	NOT IN (...)	T("!== @", []any{"a", "b"})	NOT IN ('a', 'b')
//	IS NULL			T("== @", nil)				IS NULL
//	IS NOT NULL		T("!== @", nil)				IS NOT NULL
package sqlb

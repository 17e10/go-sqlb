package sqlb

import "github.com/17e10/go-rob"

func isSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\r' || c == '\n'
}

// Compact は SQL から重複するスペースを削除します.
func Compact(s string) string {
	var l, p, q int
	l = len(s)
	if l == 0 {
		return ""
	}

	b, d := rob.Bytes(s), make([]byte, 0, l)
	first := true
	for {
		// 単語の先頭を探す
		for ; p < l; p++ {
			if !isSpace(b[p]) {
				break
			}
		}
		if p == l {
			break
		}
		// 単語の終わりを探す
		for q = p + 1; q < l; q++ {
			if isSpace(b[q]) {
				break
			}
		}
		if !first {
			d = append(d, ' ')
		} else {
			first = false
		}
		// 単語を追加する
		d = append(d, b[p:q]...)
		p = q
	}
	return rob.String(d)
}

package sqlb

func joinSqler(sep string, v []Sqler) SqlerFunc {
	return func(w Writer) (err error) {
		lv := len(v)
		if lv == 0 {
			return nil
		}
		if err = v[0].Sql(w); err != nil {
			return err
		}
		for i := 1; i < lv; i++ {
			w.WriteString(sep)
			if err = v[i].Sql(w); err != nil {
				return err
			}
		}
		return nil
	}
}

// And は条件式を表す Sqler を AND で繋げます.
func And(v ...Sqler) Sqler {
	return joinSqler(" AND ", v)
}

// Or は条件式を表す Sqler を OR で繋げます.
func Or(v ...Sqler) Sqler {
	return joinSqler(" OR ", v)
}

// Bracket は条件式を表す Sqler を括弧で括ります.
func Bracket(sqler Sqler) Sqler {
	fn := func(w Writer) error {
		w.WriteByte('(')
		if err := sqler.Sql(w); err != nil {
			return err
		}
		w.WriteByte(')')
		return nil
	}
	return SqlerFunc(fn)
}

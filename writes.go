package sqlb

import (
	"database/sql/driver"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"time"
)

// writeIdent は識別子を dialect を通じて文字列にします.
func writeIdent(w Writer, s string) error {
	if err := dialect().WriteIdent(w, s); err != nil {
		return err
	}
	return nil
}

// writeValue は値を dialect を通じて文字列にします.
//
// v が受け取れるのは bool, int, float, string などの基底型や []byte, time.Time です.
// v が driver.Valuer インターフェイスを実装していればそれを利用します.
func writeValue(w Writer, v any) (err error) {
	d := dialect()

	// Valuer を反映する
	if valuer, ok := v.(driver.Valuer); ok {
		v, err = valuer.Value()
		if err != nil {
			return fmt.Errorf("%T.Value() failed: %w", v, err)
		}
	}

	switch x := v.(type) {
	case nil:
		err = d.WriteNull(w)
	case []byte:
		err = d.WriteBytes(w, x)
	case bool:
		err = d.WriteBool(w, x)
	case time.Time:
		err = d.WriteTime(w, x)
	case string:
		err = d.WriteString(w, x)
	case int, int8, int16, int32, int64:
		i64 := reflect.ValueOf(v).Int()
		err = d.WriteInt64(w, i64)
	case uint, uint8, uint16, uint32, uint64:
		u64 := reflect.ValueOf(v).Uint()
		if u64 <= math.MaxInt64 {
			err = d.WriteInt64(w, int64(u64))
		} else {
			err = d.WriteString(w, strconv.FormatUint(u64, 10))
		}
	case float32, float64:
		f64 := reflect.ValueOf(v).Float()
		err = d.WriteFloat64(w, f64)
	default:
		err = fmt.Errorf("got type %T: %w", v, errNoValueType)
	}
	return err
}

package assert

import (
	"reflect"
	"testing"
)

func NoError(t testing.TB, err error) {
	if err != nil {
		t.Helper()
		t.Fatalf("%+v", err)
	}
}

func Error(t testing.TB, err error) {
	if err == nil {
		t.Helper()
		t.Fatal("expected an error")
	}
}

func Equal(t testing.TB, a, b interface{}) {
	if a != b && literalConvert(a) != literalConvert(b) {
		t.Helper()
		t.Fatalf("%#v != %#v", a, b)
	}
}

func DeepEqual(t testing.TB, a, b interface{}) {
	if !reflect.DeepEqual(a, b) {
		t.Helper()
		t.Fatalf("%#v != %#v", a, b)
	}
}

func That(t testing.TB, v bool) {
	if !v {
		t.Helper()
		t.Fatal("expected condition failed")
	}
}

func Nil(t testing.TB, a interface{}) {

	if a == nil {
		return
	}

	rv := reflect.ValueOf(a)
	if !canNil(rv) {
		t.Helper()
		t.Fatalf("%#v cannot be nil", a)
	}
	if !rv.IsNil() {
		t.Helper()
		t.Fatalf("%#v != nil", a)
	}
}

func NotNil(t testing.TB, a interface{}) {
	if a == nil {
		t.Helper()
		t.Fatal("expected not nil")
	}

	rv := reflect.ValueOf(a)
	if !canNil(rv) {
		return
	}
	if rv.IsNil() {
		t.Helper()
		t.Fatalf("%#v == nil", a)
	}
}

func canNil(rv reflect.Value) bool {
	switch rv.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return true
	}
	return false
}

func literalConvert(val interface{}) interface{} {
	switch val := reflect.ValueOf(val); val.Kind() {
	case reflect.Bool:
		return val.Bool()

	case reflect.String:
		return val.Convert(reflect.TypeOf("")).Interface()

	case reflect.Float32, reflect.Float64:
		return val.Float()

	case reflect.Complex64, reflect.Complex128:
		return val.Complex()

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if asInt := val.Int(); asInt < 0 {
			return asInt
		}
		return val.Convert(reflect.TypeOf(uint64(0))).Uint()

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return val.Uint()

	default:
		return val
	}
}

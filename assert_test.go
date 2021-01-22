package assert

import (
	"errors"
	"reflect"
	"testing"
)

func TestNoError(t *testing.T) {
	check(t, false, func(t testing.TB) { NoError(t, nil) })
	check(t, true, func(t testing.TB) { NoError(t, errors.New("some error")) })
}

func TestError(t *testing.T) {
	check(t, true, func(t testing.TB) { Error(t, nil) })
	check(t, false, func(t testing.TB) { Error(t, errors.New("some error")) })
}

func TestEqual(t *testing.T) {
	check(t, false, func(t testing.TB) { Equal(t, nil, nil) })
	check(t, true, func(t testing.TB) { Equal(t, nil, 0) })
	check(t, true, func(t testing.TB) { Equal(t, 0, nil) })

	numericTypes := []reflect.Type{
		reflect.TypeOf(int(0)), reflect.TypeOf(uint(0)),
		reflect.TypeOf(int8(0)), reflect.TypeOf(uint8(0)),
		reflect.TypeOf(int16(0)), reflect.TypeOf(uint16(0)),
		reflect.TypeOf(int32(0)), reflect.TypeOf(uint32(0)),
		reflect.TypeOf(int64(0)), reflect.TypeOf(uint64(0)),
	}

	for _, t1 := range numericTypes {
		for _, t2 := range numericTypes {
			v1 := reflect.New(t1).Elem().Interface()
			v2 := reflect.New(t2).Elem().Interface()
			check(t, false, func(t testing.TB) { Equal(t, v1, v2) })
		}
	}

	type (
		b    bool
		s    string
		f32  float32
		f64  float64
		c64  complex64
		c128 complex128
	)

	check(t, false, func(t testing.TB) { Equal(t, b(true), true) })
	check(t, true, func(t testing.TB) { Equal(t, b(false), true) })

	check(t, false, func(t testing.TB) { Equal(t, s("hi"), "hi") })
	check(t, true, func(t testing.TB) { Equal(t, s("hi"), "he") })

	check(t, false, func(t testing.TB) { Equal(t, f32(1.0), 1.0) })
	check(t, false, func(t testing.TB) { Equal(t, f64(1.0), 1.0) })
	check(t, false, func(t testing.TB) { Equal(t, f32(1.0), f64(1.0)) })
	check(t, true, func(t testing.TB) { Equal(t, f32(1.0), 0.0) })
	check(t, true, func(t testing.TB) { Equal(t, f64(1.0), 0.0) })
	check(t, true, func(t testing.TB) { Equal(t, f32(1.0), f64(0.0)) })

	check(t, false, func(t testing.TB) { Equal(t, c64(1.0+1.0i), 1.0+1.0i) })
	check(t, false, func(t testing.TB) { Equal(t, c128(1.0+1.0i), 1.0+1.0i) })
	check(t, false, func(t testing.TB) { Equal(t, c64(1.0+1.0i), c128(1.0+1.0i)) })
	check(t, true, func(t testing.TB) { Equal(t, c64(1.0+1.0i), 1.0) })
	check(t, true, func(t testing.TB) { Equal(t, c128(1.0+1.0i), 1.0) })
	check(t, true, func(t testing.TB) { Equal(t, c64(1.0+1.0i), c128(1.0)) })

	check(t, false, func(t testing.TB) { Equal(t, []byte("hi"), []byte("hi")) })
	check(t, true, func(t testing.TB) { Equal(t, []byte("hi"), []byte("he")) })
}

func TestDeepEqual(t *testing.T) {
	check(t, false, func(t testing.TB) { DeepEqual(t, []byte("hi"), []byte("hi")) })
	check(t, true, func(t testing.TB) { DeepEqual(t, []byte("hi"), []byte("he")) })

	check(t, true, func(t testing.TB) { DeepEqual(t, int(1), uint(1)) })
}

func TestThat(t *testing.T) {
	check(t, false, func(t testing.TB) { That(t, true) })
	check(t, true, func(t testing.TB) { That(t, false) })
}

func TestTrue(t *testing.T) {
	check(t, false, func(t testing.TB) { True(t, true) })
	check(t, true, func(t testing.TB) { True(t, false) })
}

func TestFalse(t *testing.T) {
	check(t, false, func(t testing.TB) { False(t, false) })
	check(t, true, func(t testing.TB) { False(t, true) })
}

func TestNil(t *testing.T) {
	check(t, false, func(t testing.TB) { Nil(t, nil) })
	check(t, false, func(t testing.TB) { Nil(t, (*int)(nil)) })
	check(t, true, func(t testing.TB) { Nil(t, new(int)) })
	check(t, true, func(t testing.TB) { Nil(t, 1) })
}

func TestNotNil(t *testing.T) {
	check(t, true, func(t testing.TB) { NotNil(t, nil) })
	check(t, true, func(t testing.TB) { NotNil(t, (*int)(nil)) })
	check(t, false, func(t testing.TB) { NotNil(t, new(int)) })
	check(t, false, func(t testing.TB) { NotNil(t, 1) })
}

//
// helpers
//

var sentinel = new(byte)

type recordingTB struct {
	testing.TB // must embed in order to implement
	failed     bool
}

func (tb *recordingTB) Fatal(args ...interface{}) {
	tb.failed = true
	panic(sentinel)
}

func (tb *recordingTB) Fatalf(format string, args ...interface{}) {
	tb.failed = true
	panic(sentinel)
}

func (tb *recordingTB) Helper() {}

func check(t testing.TB, failed bool, fn func(t testing.TB)) {
	rec := new(recordingTB)
	func() {
		defer func() {
			if rec := recover(); rec != nil && rec != sentinel {
				panic(rec)
			}
		}()
		fn(rec)
	}()
	if rec.failed != failed {
		t.Helper()
		t.Fatal("failed does not match")
	}
}

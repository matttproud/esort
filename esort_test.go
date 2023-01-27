package esort

import (
	"errors"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"golang.org/x/exp/slices"
)

type Data struct {
	Bool    bool
	Int8    int8
	Int16   int16
	Int32   int32
	Int64   int64
	Uint8   uint8
	Uint16  uint16
	Uint32  uint32
	Uint64  uint64
	Float32 float32
	Float64 float64
	Byte    byte
	Rune    rune
	Uint    uint
	Int     int
	Pointer uintptr
	String  string
	Bytes   []byte
}

func reverse[S ~[]E, E any](s S) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func TestIntrinsic(t *testing.T) {
	for _, test := range []struct {
		name    string
		s       *Sorter[Data]
		in, out []Data
	}{
		{
			name: "bool asc",
			s:    New[Data]().ByBool(func(d Data) bool { return d.Bool }, Asc),
			in:   []Data{{Bool: true}, {Bool: false}},
			out:  []Data{{Bool: false}, {Bool: true}},
		},
		{
			name: "bool desc",
			s:    New[Data]().ByBool(func(d Data) bool { return d.Bool }, Desc),
			in:   []Data{{Bool: false}, {Bool: true}},
			out:  []Data{{Bool: true}, {Bool: false}},
		},
		{
			name: "int8 asc",
			s:    New[Data]().ByInt8(func(d Data) int8 { return d.Int8 }, Asc),
			in:   []Data{{Int8: 1}, {Int8: 0}},
			out:  []Data{{Int8: 0}, {Int8: 1}},
		},
		{
			name: "int8 desc",
			s:    New[Data]().ByInt8(func(d Data) int8 { return d.Int8 }, Desc),
			in:   []Data{{Int8: 0}, {Int8: 1}},
			out:  []Data{{Int8: 1}, {Int8: 0}},
		},
		{
			name: "int16 asc",
			s:    New[Data]().ByInt16(func(d Data) int16 { return d.Int16 }, Asc),
			in:   []Data{{Int16: 1}, {Int16: 0}},
			out:  []Data{{Int16: 0}, {Int16: 1}},
		},
		{
			name: "int16 desc",
			s:    New[Data]().ByInt16(func(d Data) int16 { return d.Int16 }, Desc),
			in:   []Data{{Int16: 0}, {Int16: 1}},
			out:  []Data{{Int16: 1}, {Int16: 0}},
		},
		{
			name: "int32 asc",
			s:    New[Data]().ByInt32(func(d Data) int32 { return d.Int32 }, Asc),
			in:   []Data{{Int32: 1}, {Int32: 0}},
			out:  []Data{{Int32: 0}, {Int32: 1}},
		},
		{
			name: "int32 desc",
			s:    New[Data]().ByInt32(func(d Data) int32 { return d.Int32 }, Desc),
			in:   []Data{{Int32: 0}, {Int32: 1}},
			out:  []Data{{Int32: 1}, {Int32: 0}},
		},
		{
			name: "int64 asc",
			s:    New[Data]().ByInt64(func(d Data) int64 { return d.Int64 }, Asc),
			in:   []Data{{Int64: 1}, {Int64: 0}},
			out:  []Data{{Int64: 0}, {Int64: 1}},
		},
		{
			name: "int64 desc",
			s:    New[Data]().ByInt64(func(d Data) int64 { return d.Int64 }, Desc),
			in:   []Data{{Int64: 0}, {Int64: 1}},
			out:  []Data{{Int64: 1}, {Int64: 0}},
		},
		{
			name: "uint8 asc",
			s:    New[Data]().ByUint8(func(d Data) uint8 { return d.Uint8 }, Asc),
			in:   []Data{{Uint8: 1}, {Uint8: 0}},
			out:  []Data{{Uint8: 0}, {Uint8: 1}},
		},
		{
			name: "uint8 desc",
			s:    New[Data]().ByUint8(func(d Data) uint8 { return d.Uint8 }, Desc),
			in:   []Data{{Uint8: 0}, {Uint8: 1}},
			out:  []Data{{Uint8: 1}, {Uint8: 0}},
		},
		{
			name: "uint16 asc",
			s:    New[Data]().ByUint16(func(d Data) uint16 { return d.Uint16 }, Asc),
			in:   []Data{{Uint16: 1}, {Uint16: 0}},
			out:  []Data{{Uint16: 0}, {Uint16: 1}},
		},
		{
			name: "uint16 desc",
			s:    New[Data]().ByUint16(func(d Data) uint16 { return d.Uint16 }, Desc),
			in:   []Data{{Uint16: 0}, {Uint16: 1}},
			out:  []Data{{Uint16: 1}, {Uint16: 0}},
		},
		{
			name: "uint32 asc",
			s:    New[Data]().ByUint32(func(d Data) uint32 { return d.Uint32 }, Asc),
			in:   []Data{{Uint32: 1}, {Uint32: 0}},
			out:  []Data{{Uint32: 0}, {Uint32: 1}},
		},
		{
			name: "uint32 desc",
			s:    New[Data]().ByUint32(func(d Data) uint32 { return d.Uint32 }, Desc),
			in:   []Data{{Uint32: 0}, {Uint32: 1}},
			out:  []Data{{Uint32: 1}, {Uint32: 0}},
		},
		{
			name: "uint64 asc",
			s:    New[Data]().ByUint64(func(d Data) uint64 { return d.Uint64 }, Asc),
			in:   []Data{{Uint64: 1}, {Uint64: 0}},
			out:  []Data{{Uint64: 0}, {Uint64: 1}},
		},
		{
			name: "uint64 desc",
			s:    New[Data]().ByUint64(func(d Data) uint64 { return d.Uint64 }, Desc),
			in:   []Data{{Uint64: 0}, {Uint64: 1}},
			out:  []Data{{Uint64: 1}, {Uint64: 0}},
		},
		{
			name: "float32 asc",
			s:    New[Data]().ByFloat32(func(d Data) float32 { return d.Float32 }, Asc),
			in:   []Data{{Float32: 1}, {Float32: 0}},
			out:  []Data{{Float32: 0}, {Float32: 1}},
		},
		{
			name: "float32 desc",
			s:    New[Data]().ByFloat32(func(d Data) float32 { return d.Float32 }, Desc),
			in:   []Data{{Float32: 0}, {Float32: 1}},
			out:  []Data{{Float32: 1}, {Float32: 0}},
		},
		{
			name: "float64 asc",
			s:    New[Data]().ByFloat64(func(d Data) float64 { return d.Float64 }, Asc),
			in:   []Data{{Float64: 1}, {Float64: 0}},
			out:  []Data{{Float64: 0}, {Float64: 1}},
		},
		{
			name: "float64 desc",
			s:    New[Data]().ByFloat64(func(d Data) float64 { return d.Float64 }, Desc),
			in:   []Data{{Float64: 0}, {Float64: 1}},
			out:  []Data{{Float64: 1}, {Float64: 0}},
		},
		{
			name: "byte asc",
			s:    New[Data]().ByByte(func(d Data) byte { return d.Byte }, Asc),
			in:   []Data{{Byte: 1}, {Byte: 0}},
			out:  []Data{{Byte: 0}, {Byte: 1}},
		},
		{
			name: "byte desc",
			s:    New[Data]().ByByte(func(d Data) byte { return d.Byte }, Desc),
			in:   []Data{{Byte: 0}, {Byte: 1}},
			out:  []Data{{Byte: 1}, {Byte: 0}},
		},
		{
			name: "rune asc",
			s:    New[Data]().ByRune(func(d Data) rune { return d.Rune }, Asc),
			in:   []Data{{Rune: 1}, {Rune: 0}},
			out:  []Data{{Rune: 0}, {Rune: 1}},
		},
		{
			name: "rune desc",
			s:    New[Data]().ByRune(func(d Data) rune { return d.Rune }, Desc),
			in:   []Data{{Rune: 0}, {Rune: 1}},
			out:  []Data{{Rune: 1}, {Rune: 0}},
		},
		{
			name: "uint asc",
			s:    New[Data]().ByUint(func(d Data) uint { return d.Uint }, Asc),
			in:   []Data{{Uint: 1}, {Uint: 0}},
			out:  []Data{{Uint: 0}, {Uint: 1}},
		},
		{
			name: "uint desc",
			s:    New[Data]().ByUint(func(d Data) uint { return d.Uint }, Desc),
			in:   []Data{{Uint: 0}, {Uint: 1}},
			out:  []Data{{Uint: 1}, {Uint: 0}},
		},
		{
			name: "int asc",
			s:    New[Data]().ByInt(func(d Data) int { return d.Int }, Asc),
			in:   []Data{{Int: 1}, {Int: 0}},
			out:  []Data{{Int: 0}, {Int: 1}},
		},
		{
			name: "int desc",
			s:    New[Data]().ByInt(func(d Data) int { return d.Int }, Desc),
			in:   []Data{{Int: 0}, {Int: 1}},
			out:  []Data{{Int: 1}, {Int: 0}},
		},
		{
			name: "pointer asc",
			s:    New[Data]().ByPointer(func(d Data) uintptr { return d.Pointer }, Asc),
			in:   []Data{{Pointer: 1}, {Pointer: 0}},
			out:  []Data{{Pointer: 0}, {Pointer: 1}},
		},
		{
			name: "pointer desc",
			s:    New[Data]().ByPointer(func(d Data) uintptr { return d.Pointer }, Desc),
			in:   []Data{{Pointer: 0}, {Pointer: 1}},
			out:  []Data{{Pointer: 1}, {Pointer: 0}},
		},
		{
			name: "string asc",
			s:    New[Data]().ByString(func(d Data) string { return d.String }, Asc),
			in:   []Data{{String: "b"}, {String: "a"}},
			out:  []Data{{String: "a"}, {String: "b"}},
		},
		{
			name: "string desc",
			s:    New[Data]().ByString(func(d Data) string { return d.String }, Desc),
			in:   []Data{{String: "a"}, {String: "b"}},
			out:  []Data{{String: "b"}, {String: "a"}},
		},
		{
			name: "bytes asc",
			s:    New[Data]().ByBytes(func(d Data) []byte { return d.Bytes }, Asc),
			in:   []Data{{Bytes: []byte{1}}, {Bytes: []byte{0}}},
			out:  []Data{{Bytes: []byte{0}}, {Bytes: []byte{1}}},
		},
		{
			name: "bytes desc",
			s:    New[Data]().ByBytes(func(d Data) []byte { return d.Bytes }, Desc),
			in:   []Data{{Bytes: []byte{0}}, {Bytes: []byte{1}}},
			out:  []Data{{Bytes: []byte{1}}, {Bytes: []byte{0}}},
		},
		{
			name: "func asc",
			s:    New[Data]().ByFunc(func(l, r Data) bool { return l.Int < r.Int }, Asc),
			in:   []Data{{Int: 1}, {Int: 0}},
			out:  []Data{{Int: 0}, {Int: 1}},
		},
		{
			name: "func desc",
			s:    New[Data]().ByFunc(func(l, r Data) bool { return l.Int < r.Int }, Desc),
			in:   []Data{{Int: 0}, {Int: 1}},
			out:  []Data{{Int: 1}, {Int: 0}},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Run("normal", func(t *testing.T) {
				out := make([]Data, len(test.in))
				copy(out, test.in)
				slices.SortFunc(out, test.s.Less)
				if diff := cmp.Diff(test.out, out); diff != "" {
					t.Errorf("slices.SortFunc(%v) = %v, want %v\n\ndiff (-want, +got):\n%v", test.in, out, test.out, diff)
				}
			})
			t.Run("inverse", func(t *testing.T) {
				out := make([]Data, len(test.in))
				copy(out, test.in)
				reverse(out)
				slices.SortFunc(out, test.s.Less)
				if diff := cmp.Diff(test.out, out); diff != "" {
					t.Errorf("slices.SortFunc(%v) = %v, want %v\n\ndiff (-want, +got):\n%v", test.in, out, test.out, diff)
				}
			})
		})
	}
}

func TestCompound(t *testing.T) {
	for _, test := range []struct {
		name    string
		s       *Sorter[Data]
		in, out []Data
	}{
		{
			name: "int ASC uint DESC",
			s: New[Data]().
				ByInt(func(d Data) int { return d.Int }, Asc).
				ByUint(func(d Data) uint { return d.Uint }, Desc),
			in: []Data{
				{Int: 1, Uint: 0},
				{Int: 1, Uint: 1},
				{Int: 0, Uint: 0},
				{Int: 0, Uint: 1},
			},
			out: []Data{
				{Int: 0, Uint: 1},
				{Int: 0, Uint: 0},
				{Int: 1, Uint: 1},
				{Int: 1, Uint: 0},
			},
		},
		{
			name: "int DESC uint ASC",
			s: New[Data]().
				ByInt(func(d Data) int { return d.Int }, Desc).
				ByUint(func(d Data) uint { return d.Uint }, Asc),
			in: []Data{
				{Int: 0, Uint: 1},
				{Int: 0, Uint: 0},
				{Int: 1, Uint: 1},
				{Int: 1, Uint: 0},
			},
			out: []Data{
				{Int: 1, Uint: 0},
				{Int: 1, Uint: 1},
				{Int: 0, Uint: 0},
				{Int: 0, Uint: 1},
			}},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Run("normal", func(t *testing.T) {
				out := make([]Data, len(test.in))
				copy(out, test.in)
				slices.SortFunc(out, test.s.Less)
				if diff := cmp.Diff(test.out, out); diff != "" {
					t.Errorf("slices.SortFunc(%v) = %v, want %v\n\ndiff (-want, +got):\n%v", test.in, out, test.out, diff)
				}
			})
			t.Run("inverse", func(t *testing.T) {
				out := make([]Data, len(test.in))
				copy(out, test.in)
				reverse(out)
				slices.SortFunc(out, test.s.Less)
				if diff := cmp.Diff(test.out, out); diff != "" {
					t.Errorf("slices.SortFunc(%v) = %v, want %v\n\ndiff (-want, +got):\n%v", test.in, out, test.out, diff)
				}
			})
		})
	}
}

func TestEmpty(t *testing.T) {
	var err error
	defer func() {
		err = recover().(error)
	}()
	data := []Data{{}, {}}
	sorter := New[Data]()
	slices.SortFunc(data, sorter.Less)
	if got, want := err, errNoProgram; !errors.Is(got, want) {
		t.Errorf("after empty sorter sort panic = %v, want %v", got, want)
	}
}

var benchData = []Data{
	{Int: 0, Uint: 3},
	{Int: 1, Uint: 1},
	{Int: 0, Uint: 0},
	{Int: 0, Uint: 1},
	{Int: 0, Uint: 2},
	{Int: 0, Uint: 3},
	{Int: 2, Uint: 1},
	{Int: 3, Uint: 0},
	{Int: 1, Uint: 0},
}

func Benchmark(b *testing.B) {
	for _, i := range []int{10, 100, 1000} {
		b.Run(fmt.Sprint(i), func(b *testing.B) {
			bench := make([][]Data, 0, b.N)
			var data []Data
			for j := 0; j < i; j++ {
				data = append(data, benchData...)
			}
			for j := 0; j < b.N; j++ {
				bench = append(bench, data)
			}
			sorter := New[Data]().
				ByInt(func(d Data) int { return d.Int }, Desc).
				ByUint(func(d Data) uint { return d.Uint }, Asc)
			b.ResetTimer()
			b.ReportAllocs()
			for j := 0; j < b.N; j++ {
				slices.SortFunc(bench[j], sorter.Less)
			}
		})
	}
}

func BenchmarkBest(b *testing.B) {
	for _, i := range []int{10, 100, 1000} {
		b.Run(fmt.Sprint(i), func(b *testing.B) {
			bench := make([][]Data, 0, b.N)
			var data []Data
			for j := 0; j < i; j++ {
				data = append(data, benchData...)
			}
			for j := 0; j < b.N; j++ {
				bench = append(bench, data)
			}
			sorter := func(l, r Data) bool {
				if l.Int != r.Int {
					return l.Int > r.Int
				}
				return l.Uint < r.Uint
			}

			b.ResetTimer()
			b.ReportAllocs()
			for j := 0; j < b.N; j++ {
				slices.SortFunc(bench[j], sorter)
			}
		})
	}
}

func BenchmarkNaive(b *testing.B) {
	for _, i := range []int{10, 100, 1000} {
		b.Run(fmt.Sprint(i), func(b *testing.B) {
			bench := make([][]Data, 0, b.N)
			var data []Data
			for j := 0; j < i; j++ {
				data = append(data, benchData...)
			}
			for j := 0; j < b.N; j++ {
				bench = append(bench, data)
			}
			sorter := func(l, r Data) bool {
				if l.Int > r.Int {
					return true
				} else if r.Int > l.Int {
					return false
				}
				return l.Uint < r.Uint
			}

			b.ResetTimer()
			b.ReportAllocs()
			for j := 0; j < b.N; j++ {
				slices.SortFunc(bench[j], sorter)
			}
		})
	}
}

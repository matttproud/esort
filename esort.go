// Package esort provides mechanisms for sorting user-defined types according
// to compound criteria extensibly.  It is mutually compatible with package sort
// from the standard library.
//
// Consider a case of sorting a user-defined type of a Person:
//
//	type Person struct {
//		GivenName string
//		ID int
//	}
//
// Suppose GivenName is to be sorted in descending order and ties of the
// GivenName by ID in ascending order:
//
//	sorter := esort.New().
//		ByString(func(p Person) string { return p.GivenName }, esort.Desc).
//		ByInt(func(p Person) int { return p.ID }, esort.Asc)
//
// The data can then be sorted:
//
//	var data []Person
//
//	sort.Slice(data, func(i, j int) bool { return sorter.Less(data[i], data[j]) })
//
// The same mechanism can be used directly with generics:
//
//	slices.SortFunc(data, sorter.Less)
//
// # Sorting Instructions
//
// The methods prefixed with By copy the current sorting rule set and add a
// new instruction to the copy.  Each rule set is evaluated according to the
// order in which the instructions were added to the [Sorter], so earlier
// sorting rules carry higher sorting weight than the previous.
//
// A Sorter without rules is invalid and panics upon use.
//
// Due to the copy on write semantics of the instructions API, rule sets can be
// templatized:
//
//	base := esort.New().
//		ByString(func(p Person) string { return p.GivenName }, esort.Desc).
//
//	idAsc := base.ByInt(func(p Person) int { return p.ID }, esort.Asc)
//	idDesc := base.ByInt(func(p Person) int { return p.ID }, esort.Desc)
//
// # Efficiency
//
// Prefer using the simple, scalar By methods for comparing individual fields
// versus encoding their logic in the [Sorter.ByFunc] method.  These individual
// methods are implemented efficiently.
//
// Prefer:
//
//	sorter := esort.New().
//		ByString(func(p Person) string { return p.GivenName }, esort.Desc).
//		ByInt(func(p Person) int { return p.ID }, esort.Asc)
//
// Over:
//
//	sorter := esort.New().
//		ByFunc(func(l, r Person) bool {
//			if l.GivenName != r.GivenName {
//				return l.GivenName < r.GivenName
//			}
//		}, esort.Desc).
//		ByInt(func(p Person) int { return p.ID }, esort.Asc)
//
// Prefer defining the sorters as top-level variables if they have a static
// sorting basis.  They will allocate zero memory during later program
// runtime — provided of course no user-defined functions registered with
// ByFunc do.
//
// # User-Defined Sorting Functions
//
// The ByFunc method enables types to be sorted by user-defined functions.
// The ByFunc method follows the same logic as the other instructions; it is
// appended to the rule set and evaluated accordingly.
//
// If the function provided to ByFunc compares only one sorting basis, it may
// be implemented in a naive fashion whereby it only compares one side of the
// basis.  Consider a data-transmission object (DTO)
//
//	func lessDTO(l, r *dto.SomeData) bool {
//		return l.GetField() < r.GetField()
//	}
//
//	sorter := esort.New().
//		ByFunc(lessDTO, esort.Asc)
//
// The function lessDTO needn't consider the inverse case:
//
//	func lessDTO(l, r *dto.SomeData) bool {
//		if l.GetField() < r.GetField() {
//			return true
//		}
//		return r.GetField() < l.GetField()
//	}
//
// The [Sorter] checks inverse cases automatically.
//
// # Ergonomics of Value Accessors
//
// The examples above are implemented using anonymous functions, similar to this:
//
//	func(p Person) string { return p.GivenName }
//
// You are free to use free-standing, top-level functions:
//
//	func name(p Person) string { return p.GivenName }
//
//	sorter := esort.New().
//		ByFunc(name, esort.Asc)
//
// And [method expressions] on the types, though consider whether [getters] are
// even appropriate for the API first:
//
//	func (p Person) GivenName string { return p.givenName }
//
//	sorter := esort.New().
//		ByFunc(p.GivenName, esort.Asc)
//
// [method expressions]: https://go.dev/ref/spec#Method_expressions
// [getters]: https://google.github.io/styleguide/go/decisions.html#getters
package esort

// Things to improve in the design:
//
// 1. Eliminate the need to panic at all on the case of empty instructions.

import (
	"bytes"
	"errors"

	"golang.org/x/exp/constraints"
)

// inst is a sorting operation instruction.
type inst[T any] struct {
	Func func(l, r T) bool
	Dir  Dir
}

// Sorter is the representation of a compound sorting program.  A Sorter is
// safe for concurrent use by multiple goroutines.
type Sorter[T any] struct {
	// prog internally emulates the builder pattern.  While this pattern is
	// unconventional for Go, its use here is safe and ergonomic, because there
	// is zero need to consider error handling in this API.
	prog []inst[T]
}

// Dir represents the direction for the sort.
type Dir int

const (
	// Asc sorts the results in ascending order.
	Asc = Dir(iota)
	// Desc sorts the results in descending order.
	Desc
)

// addInst copies the existing sorting program and adds a new instruction to
// the copy.  The copy semantic is used to keep each Sorter safe for use in
// multiple goroutines and to enable basic templatization of the programs to be
// performed à la the builder pattern.
func (s *Sorter[T]) addInst(o inst[T]) *Sorter[T] {
	return &Sorter[T]{
		prog: append(append([]inst[T](nil), s.prog...), o),
	}
}

// ByBool sorts the data by a given boolean value.
func (s *Sorter[T]) ByBool(f func(T) bool, d Dir) *Sorter[T] {
	fn := func(l, r T) bool {
		return !f(l) && f(r)
	}
	return s.addInst(inst[T]{fn, d})
}

// lessFunc sorts any ordered data.
func lessFunc[T any, V constraints.Ordered](f func(T) V) func(l, r T) bool {
	return func(l, r T) bool {
		return f(l) < f(r)
	}
}

// ByInt8 sorts the data by a given int8 value.
func (s *Sorter[T]) ByInt8(f func(T) int8, d Dir) *Sorter[T] {
	fn := lessFunc(f)
	return s.addInst(inst[T]{fn, d})
}

// ByInt16 sorts the data by a given int16 value.
func (s *Sorter[T]) ByInt16(f func(T) int16, d Dir) *Sorter[T] {
	fn := lessFunc(f)
	return s.addInst(inst[T]{fn, d})
}

// ByInt32 sorts the data by a given int32 value.
func (s *Sorter[T]) ByInt32(f func(T) int32, d Dir) *Sorter[T] {
	fn := lessFunc(f)
	return s.addInst(inst[T]{fn, d})
}

// ByInt64 sorts the data by a given int64 value.
func (s *Sorter[T]) ByInt64(f func(T) int64, d Dir) *Sorter[T] {
	fn := lessFunc(f)
	return s.addInst(inst[T]{fn, d})
}

// ByUint8 sorts the data by a given uint8 value.
func (s *Sorter[T]) ByUint8(f func(T) uint8, d Dir) *Sorter[T] {
	fn := lessFunc(f)
	return s.addInst(inst[T]{fn, d})
}

// ByUint16 sorts the data by a given uint16 value.
func (s *Sorter[T]) ByUint16(f func(T) uint16, d Dir) *Sorter[T] {
	fn := lessFunc(f)
	return s.addInst(inst[T]{fn, d})
}

// ByUint32 sorts the data by a given uint32 value.
func (s *Sorter[T]) ByUint32(f func(T) uint32, d Dir) *Sorter[T] {
	fn := lessFunc(f)
	return s.addInst(inst[T]{fn, d})
}

// ByUint64 sorts the data by a given uint64 value.
func (s *Sorter[T]) ByUint64(f func(T) uint64, d Dir) *Sorter[T] {
	fn := lessFunc(f)
	return s.addInst(inst[T]{fn, d})
}

// ByFloat32 sorts the data by a given float32 value.
func (s *Sorter[T]) ByFloat32(f func(T) float32, d Dir) *Sorter[T] {
	fn := lessFunc(f)
	return s.addInst(inst[T]{fn, d})
}

// ByFloat64 sorts the data by a given float64 value.
func (s *Sorter[T]) ByFloat64(f func(T) float64, d Dir) *Sorter[T] {
	fn := lessFunc(f)
	return s.addInst(inst[T]{fn, d})
}

// ByByte sorts the data by a given byte value.
func (s *Sorter[T]) ByByte(f func(T) byte, d Dir) *Sorter[T] {
	fn := lessFunc(f)
	return s.addInst(inst[T]{fn, d})
}

// ByRune sorts the data by a given rune value.
func (s *Sorter[T]) ByRune(f func(T) rune, d Dir) *Sorter[T] {
	fn := lessFunc(f)
	return s.addInst(inst[T]{fn, d})
}

// ByUint sorts the data by a given uint value.
func (s *Sorter[T]) ByUint(f func(T) uint, d Dir) *Sorter[T] {
	fn := lessFunc(f)
	return s.addInst(inst[T]{fn, d})
}

// ByInt sorts the data by a given int value.
func (s *Sorter[T]) ByInt(f func(T) int, d Dir) *Sorter[T] {
	fn := lessFunc(f)
	return s.addInst(inst[T]{fn, d})
}

// ByPointer sorts the data by a given uintptr value.
func (s *Sorter[T]) ByPointer(f func(T) uintptr, d Dir) *Sorter[T] {
	fn := lessFunc(f)
	return s.addInst(inst[T]{fn, d})
}

// ByString sorts the data by a given string value.
func (s *Sorter[T]) ByString(f func(T) string, d Dir) *Sorter[T] {
	fn := lessFunc(f)
	return s.addInst(inst[T]{fn, d})
}

// ByBytes sorts the data by a given byte slice value.
func (s *Sorter[T]) ByBytes(f func(T) []byte, d Dir) *Sorter[T] {
	fn := func(l, r T) bool {
		return bytes.Compare(f(l), f(r)) < 0
	}
	return s.addInst(inst[T]{fn, d})
}

// SortFunc sorts the data according to an arbitrary function.
//
// SortFunc mimics a [sort.Interface.Less] function.  Functions that sort by
// one criterion may return a simple l before r check.  Consider:
//
//	sorter := esort.New().
//		ByFunc(func(l, r T) bool {
//			return l.GivenName < r.GivenName
//		})
//
// Functions that perform compound comparisons must perform symmetric checks
// of r before l in order to produce sensible results:
//
//	sorter := esort.New().
//		ByFunc(func(l, r T) bool {
//			if l.GivenName != r.GivenName {
//				return l.GivenName < r.GivenName
//			}
//			return l.ID < r.ID
//		})
//
// Implementations of SortFunc must not sort in a way that conflicts with a
// pre-existing instruction in a [Sorter]; otherwise inconsistent results may
// be produced.
//
// Using the native By-prefixed functions to generate sorting rule sets is
// preferable to using this API.
type SortFunc[T any] func(l, r T) bool

// ByFunc sorts the data according to an arbitrary given [SortFunc].
// The SortFunc must not the underlying data by that any pre-existing
// intruction does.
func (s *Sorter[T]) ByFunc(f SortFunc[T], d Dir) *Sorter[T] {
	return s.addInst(inst[T]{f, d})
}

// errNoProgram indicates that the sorter has no recorded instructions, meaning
// it can't do anything.
var errNoProgram = errors.New("esort: no sorting instructions provided")

// Less is a sort ordering function that fulfills the contract expected by
// [sort.Interface.Less] and related APIs.
func (s *Sorter[T]) Less(l, r T) bool {
	for i, f := range s.prog {
		if f.Dir == Asc {
			r, l = l, r
		}
		switch i {
		case len(s.prog) - 1:
			return f.Func(l, r)
		default:
			if f.Func(r, l) {
				return true
			} else if f.Func(l, r) {
				return false
			}
			continue
		}
	}
	panic(errNoProgram)
}

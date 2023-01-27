# esort

> Package esort provides mechanisms for sorting user-defined types according
> to compound criteria extensibly.  It is mutually compatible with package sort
> from the standard library.  The package name comes from the combination of
> the words "extensible" and "sort" and shortened as "esort."

## Goals

* Type Safe
* Zero Allocation (in common workflows if sorting basis is known a priori)
* Zero Reflection
* Safe and Accurate
* Reasonable Degree of Ergonomics

## Motivation

This is an experimental library.

Stupid as it sounds: during code review out in the field, I have seen more
half-baked `(sort.Interface).Less` implementations than I am comfortable with,
particularly where compound criteria are evaluated.  I wondered whether the
situation could be solved through an API.

I'm not suggesting anyone should use this.  There are potentially other
solutions to improve ergonomics and performance.  For the time being, it is a
minimally-viable prototype.

Its API documentation is available at https://pkg.go.dev/github.com/matttproud/esort.
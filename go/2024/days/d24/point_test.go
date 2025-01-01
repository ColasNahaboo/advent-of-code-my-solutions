package main

import (
	"testing"
	"reflect"
	"strings"
)

// setup tools

func t_b(s string) *Board[string] {
	b:= ParseBoard[string](strings.Split(s, " "), t_ParseCellLetter)
	return b
}
func t_ParseCellLetter(x, y int, r rune) string { // chars, with _ = empty string
	if r == '_' { return "" 
	} else { return string(r)
	}
}

// tests to do 
func Test(t *testing.T) {
	b := "abcd efgh ijkl mnop"

	t_IR("IR1", t, t_b(b), 2, 1, t_b("abcd efgh ____ ijkl mnop"))
	t_IR("IR2", t, t_b(b), 1, 2, t_b("abcd ____ ____ efgh ijkl mnop"))
	t_IR("IR3", t, t_b(b), 0, 1, t_b("____ abcd efgh ijkl mnop"))
	t_IR("IR4", t, t_b(b), 3, 1, t_b("abcd efgh ijkl ____ mnop"))
	t_IR("IR5", t, t_b(b), 4, 1, t_b("abcd efgh ijkl mnop ____"))

	t_DR("DR1", t, t_b(b), 2, 1, t_b("abcd efgh mnop"))
	t_DR("DR2", t, t_b(b), 1, 2, t_b("abcd mnop"))
	t_DR("DR3", t, t_b(b), 0, 1, t_b("efgh ijkl mnop"))
	t_DR("DR4", t, t_b(b), 2, 1, t_b("abcd efgh mnop"))
	t_DR("DR5", t, t_b(b), 3, 1, t_b("abcd efgh ijkl"))

	t_IC("IC1", t, t_b(b), 2, 1, t_b("ab_cd ef_gh ij_kl mn_op"))
	t_IC("IC2", t, t_b(b), 1, 2, t_b("a__bcd e__fgh i__jkl m__nop"))
	t_IC("IC3", t, t_b(b), 0, 1, t_b("_abcd _efgh _ijkl _mnop"))
	t_IC("IC4", t, t_b(b), 3, 1, t_b("abc_d efg_h ijk_l mno_p"))
	t_IC("IC5", t, t_b(b), 4, 1, t_b("abcd_ efgh_ ijkl_ mnop_"))

	t_DC("DC1", t, t_b(b), 2, 1, t_b("abd efh ijl mnp"))
	t_DC("DC2", t, t_b(b), 1, 2, t_b("ad eh il mp"))
	t_DC("DC3", t, t_b(b), 0, 1, t_b("bcd fgh jkl nop"))
	t_DC("DC4", t, t_b(b), 0, 3, t_b("d h l p"))
	t_DC("DC5", t, t_b(b), 3, 1, t_b("abc efg ijk mno"))
}

// functions to test

func t_IR(label string, t *testing.T, b *Board[string], i, n int, e *Board[string]) {
	// setup here
	t.Run(label, func(t *testing.T) {
		b.InsertRowsBefore(i, n)
		if !reflect.DeepEqual(*b, *e) {
			t.Errorf("exp '%v'\ngot '%v'", *e, *b)
		}})}

func t_DR(label string, t *testing.T, b *Board[string], i, n int, e *Board[string]) {
	// setup here
	t.Run(label, func(t *testing.T) {
		b.DeleteRowsAt(i, n)
		if !reflect.DeepEqual(*b, *e) {
			t.Errorf("exp '%v'\ngot '%v'", *e, *b)
		}})}

func t_IC(label string, t *testing.T, b *Board[string], i, n int, e *Board[string]) {
	// setup here
	t.Run(label, func(t *testing.T) {
		b.InsertColsBefore(i, n)
		if !reflect.DeepEqual(*b, *e) {
			t.Errorf("\nexp '%v'\ngot '%v'", *e, *b)
		}})}

func t_DC(label string, t *testing.T, b *Board[string], i, n int, e *Board[string]) {
	// setup here
	t.Run(label, func(t *testing.T) {
		b.DeleteColsAt(i, n)
		if !reflect.DeepEqual(*b, *e) {
			t.Errorf("expe '%v'\n got '%v'", *e, *b)
		}})}

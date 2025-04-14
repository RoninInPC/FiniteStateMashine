package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func C[anything any](a anything) {
	println(a)
}

func TestDFA1(t *testing.T) {
	dfa := DefaultDFA[string, string]().
		AddTransition("0", "a", "1").
		AddTransition("1", "d", "1").
		SetStart("0").
		SetEnd([]string{"1"}).Build()
	answer := dfa.Correct("a", "d", "d", "d")
	assert.Equal(t, answer, true, nil)
}
func TestDFA2(t *testing.T) {
	dfa := DefaultDFA[string, string]().
		AddTransition("0", "a", "1").
		AddTransition("1", "d", "1").
		SetStart("0").
		SetEnd([]string{"1"}).Build()
	answer := dfa.Correct("a", "d", "d", "a")
	assert.Equal(t, answer, false, nil)
}

type FUNCTEST func()

func Out() {
	println("Out")
}
func In() {
	println("In")
}

func TestDFA3(t *testing.T) {
	dfa := DefaultDFA[FUNCTEST, string]().
		AddTransition(Out, "a", In).
		AddTransition(In, "d", In).
		AddTransition(In, "a", Out).
		SetStart(Out).
		SetEnd([]FUNCTEST{In}).Build()
	dfa.GetDFACurrentCondition()()
	dfa.ToNext("a")
	dfa.GetDFACurrentCondition()()
	dfa.ToNext("a")
	dfa.GetDFACurrentCondition()()
	dfa.ToNext("a")
	dfa.GetDFACurrentCondition()()
	dfa.ToNext("d")
	dfa.GetDFACurrentCondition()()
}

type MyType int

var (
	CONST    MyType = 0
	NO_CONST MyType = 1
)

func (m MyType) String() string {
	if m == CONST {
		return "CONST"
	}
	if m == NO_CONST {
		return "NO_CONST"
	}
	return "NONE"
}

func TestDFA4(t *testing.T) {
	dfa := DefaultDFA[MyType, string]().
		AddTransition(CONST, "a", NO_CONST).
		AddTransition(NO_CONST, "d", NO_CONST).
		AddTransition(NO_CONST, "a", CONST).
		SetStart(CONST).
		SetEnd([]MyType{NO_CONST}).Build()
	println(dfa.GetDFACurrentCondition().String())
	dfa.ToNext("a")
	println(dfa.GetDFACurrentCondition().String())
	dfa.ToNext("a")
	println(dfa.GetDFACurrentCondition().String())
	dfa.ToNext("a")
	println(dfa.GetDFACurrentCondition().String())
	dfa.ToNext("d")
	println(dfa.GetDFACurrentCondition().String())
}

package gost

import (
	"testing"
)

type Fn func()

type Step struct {
	TK *testing.T
	F  []Fn
}

func (s *Step) It(desc string, fn ...Fn) *Step {
	var fns []Fn
	for _, it := range fn {
		fns = append(fns, it)
	}

	s.F = fns

	return s
}

func (s *Step) Run() {
	for _, it := range s.F {
		it()
	}
}

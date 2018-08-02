package states

import (
	"errors"
	"fmt"

	"github.com/amyadzuki/amygolib/onfail"
	"github.com/amyadzuki/amygolib/str"
)

var ErrTooManyNames =
	errors.New("State: Bad arguments to Run/RunOnce: must be () or (string)")

type State struct {
	Data interface{}

	fnCloseRequested func() bool
	fns              map[string]func(*State)
	editing          string
	sCurrent, sNext  string
	state            string
}

func New(fn func() bool) *State {
	return new(State).Init(fn)
}

func (s *State) Init(fn func() bool) *State {
	s.fnCloseRequested = fn
	s.fns = make(map[string]func(*State))
	return s
}

func (s *State) OnEnter(args ...interface{}) *State {
	s.fns[s.editing + "{"] = s.reg(args...)
	return s
}

func (s *State) OnLeave(args ...interface{}) *State {
	s.fns[s.editing + "}"] = s.reg(args...)
	return s
}

func (s *State) Register(args ...interface{}) *State {
	s.fns[s.editing] = s.reg(args...)
	return s
}

func (s *State) Run(name ...string) *State {
	switch len(name) {
	case 0:
	case 1:
		s.sNext = str.Simp(name[0])
	default:
		panic(ErrTooManyNames)
	}
	for !s.fnCloseRequested() && len(s.sNext) > 0 {
		s.runOnce()
	}
}

func (s *State) RunOnce(name ...string) *State {
	switch len(name) {
	case 0:
	case 1:
		s.sNext = str.Simp(name[0])
	default:
		panic(ErrTooManyNames)
	}
	if !s.fnCloseRequested() && len(s.sNext) > 0 {
		return s.runOnce()
	}
}

func (s *State) SetData(data interface{}) *State {
	s.Data = data
	return s
}

func (s *State) SetNext(state string, onFail ...onfail.OnFail) *State {
	if _, ok := s.fns[str.Simp(state)]; ok {
		s.sNext = state
	} else {
		onfail.Fail("Unregistered state: \"" + state + "\"", s, nil, onFail...)
	}
	return s
}

func (s *State) reg(args ...interface{}) func(*State) {
	switch len(args) {
	case 1:
		cb, ok := args[0].(func(*State))
		if ok {
			return cb
		}
	case 2:
		name, nok := args[0].(string)
		cb, cok := args[1].(func(*State))
		if nok && cok {
			s.editing = str.Simp(name)
			return cb
		}
	}
	panic(badBuilderArgs(args...))
}

func (s *State) runOnce() *State {
	s.sCurrent = s.sNext
	main, mok := s.fns[s.state]
	if mok {
		enter, eok := s.fns[s.state + "{"]
		leave, lok := s.fns[s.state + "}"]
		if eok {
			enter(s)
		}
		for !s.fnCloseRequested() && s.sNext == s.sCurrent {
			main(s)
		}
		if lok {
			leave(s)
		}
	}
	return s
}

func badBuilderArgs(args ...interface{}) error {
	msg := "State: Bad arguments to builder method: must be (func(*State)) or (string, func(*State))\nHave: ("
	for aid, arg := range args {
		if aid != 0 {
			msg += ", "
		}
		msg += fmt.Sprintf("%T", arg)
	}
	msg += ")"
	return errors.New(msg)
}

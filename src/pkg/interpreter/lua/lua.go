package lua

import (
	"errors"
	"golua"
	"interpreter"
)

type LuaInterpreter struct {
	cbs []interpreter.EmitCallback
}

var (
	EInvalidCode = errors.New("Could not compile code")
)

func New() *LuaInterpreter {
	return &LuaInterpreter {
		cbs: make([]interpreter.EmitCallback, 0),
	}
}

func (l *LuaInterpreter) Eval(code string) (interpreter.Interpretation, error) {
	s := golua.NewState()
	s.Register("emit", func(s *golua.State) int {
		data := s.ToString(-1)
		s.Pop(-1)
		for _, cb := range l.cbs {
			cb([]byte(data))
		}
		return 0
	})
	s.LoadString(code)
	if s.IsNil(-1) {
		return nil, EInvalidCode
	}
	return func(data, state []byte) []byte {
		// Copy the function to the top of the stack
		s.PushValue(-1)
		s.PushString(string(data))
		s.PushString(string(state))
		s.Call(2, 1)
		return []byte(s.ToString(-1))
	}, nil
}

func (l *LuaInterpreter) RegisterEmitCallback(cb interpreter.EmitCallback) {
	l.cbs = append(l.cbs, cb)
}

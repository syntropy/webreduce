package nop

import (
	"wr/interpreter"
)

type NopInterpreter struct {
	cbs []interpreter.EmitCallback
}

func New() *NopInterpreter {
	return &NopInterpreter{
		cbs: make([]interpreter.EmitCallback, 0),
	}
}

func (n *NopInterpreter) Eval(code string) (interpreter.Interpretation, error) {
	return func(data, state []byte) []byte {
		for _, cb := range n.cbs {
			cb(data)
		}
		return state
	}, nil
}

func (n *NopInterpreter) RegisterEmitCallback(cb interpreter.EmitCallback) {
	n.cbs = append(n.cbs, cb)
}

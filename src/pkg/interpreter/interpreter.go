package interpreter

// Interpreter defines the interface for differenct backend
// languages.
type Interpreter interface {
	// Eval takes a string of code and returns an interpretation
	// which executes the given code.
	Eval(code string) Interpretation

	// The passed function is supposed to be called,
	// whenever the interpreted code calls emit.
	RegisterEmitCallback(cb EmitCallback)
}

// Interpretation is a closure around the interpreted code
// Data is the date to be processed, state is the current
// state of the behaviour. The return value is the new state.
type Interpretation func(data, state []byte) []byte

// The callback is called whenever the interpreted code calls
// emit and the data to be emitted is passed along.
type EmitCallback func(data []byte)

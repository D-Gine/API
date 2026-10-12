package scripting

import "context"

type Engine interface {
	// executes function from a script and returns a Go value converted result
	Eval(ctx context.Context, script string, function string, args []any, kwargs map[string]any) (any, error)
}

// extracts a string result and accepts nil
func AsOptionalString(v any) (*string, error) {
	if v == nil {
		return nil, nil
	}
	s, ok := v.(string)
	if !ok {
		return nil, ErrUnexpectedType("string", v)
	}
	if s == "" {
		return nil, nil
	}
	return &s, nil
}

// ScriptTypeError indicates a script returned an unexpected value type
type ScriptTypeError struct {
	Expected string
	Got      any
}

func (e ScriptTypeError) Error() string {
	return "script returned unexpected type: expected " + e.Expected
}

func ErrUnexpectedType(expected string, got any) error {
	return ScriptTypeError{Expected: expected, Got: got}
}

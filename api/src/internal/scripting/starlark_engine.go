package scripting

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	starlarkjson "go.starlark.net/lib/json"
	"go.starlark.net/starlark"
	"go.starlark.net/syntax"
)

// StarlarkEngine runs user-provided Starlark scripts (a small, safe, Python-like
// language) from Go code. It holds no state, so a single instance can be shared
// and used concurrently.
//
// Basic usage:
//
//	engine := scripting.NewStarlarkEngine()
//	result, err := engine.Eval(ctx, script, "main", []any{1, 2}, nil)
type StarlarkEngine struct{}

// NewStarlarkEngine creates a ready-to-use StarlarkEngine.
func NewStarlarkEngine() *StarlarkEngine {
	return &StarlarkEngine{}
}

// Eval loads a Starlark script, calls one of its functions and returns the
// result converted to plain Go values.
//
// Parameters:
//   - ctx:      context checked before and after the call; if it is cancelled,
//     Eval returns ctx.Err(). Note it does not interrupt a script that is
//     already running.
//   - script:   the full Starlark source code.
//   - function: name of the top-level function in the script to call.
//   - args:     positional arguments, passed in order (may be nil).
//   - kwargs:   named arguments, passed by name (may be nil).
//
// Argument values are converted with goToStarlark, and the returned value with
// starlarkToGo (see those functions for the supported types).
//
// Steps performed:
//  1. Execute the script once to define its functions.
//  2. Look up `function` among the globals and check that it is callable.
//  3. Convert args/kwargs from Go to Starlark.
//  4. Call the function.
//  5. Convert the result back to Go.
//
// Errors are returned (wrapped) when the script has a syntax/runtime error, the
// function does not exist or is not callable, a value cannot be converted, or
// the context is cancelled.
//
// Example:
//
//	script := `
//	def greet(name, punctuation="!"):
//	    return {"message": "Hello " + name + punctuation}
//	`
//	res, err := engine.Eval(ctx, script, "greet",
//		[]any{"Bob"},
//		map[string]any{"punctuation": "?"},
//	)
//	if err != nil {
//		return err
//	}
//	// res == map[string]any{"message": "Hello Bob?"}
//	msg := res.(map[string]any)["message"].(string)
func (e *StarlarkEngine) Eval(ctx context.Context, script string, function string, args []any, kwargs map[string]any) (any, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	globals, err := starlark.ExecFileOptions(&syntax.FileOptions{}, &starlark.Thread{Name: function}, "script.star", script, nil)
	if err != nil {
		return nil, fmt.Errorf("starlark exec error: %w", err)
	}

	fnValue, ok := globals[function]
	if !ok {
		return nil, fmt.Errorf("missing %s function in script", function)
	}
	fn, ok := fnValue.(starlark.Callable)
	if !ok {
		return nil, fmt.Errorf("%s is not callable", function)
	}

	positionals := make(starlark.Tuple, 0, len(args))
	for _, arg := range args {
		sv, err := goToStarlark(arg)
		if err != nil {
			return nil, fmt.Errorf("args conversion error: %w", err)
		}
		positionals = append(positionals, sv)
	}

	var named []starlark.Tuple
	if len(kwargs) > 0 {
		named = make([]starlark.Tuple, 0, len(kwargs))
		for k, v := range kwargs {
			sv, err := goToStarlark(v)
			if err != nil {
				return nil, fmt.Errorf("kwargs conversion error for %q: %w", k, err)
			}
			named = append(named, starlark.Tuple{starlark.String(k), sv})
		}
	}

	thread := &starlark.Thread{Name: function + ".call"}
	result, err := starlark.Call(thread, fn, positionals, named)
	if err != nil {
		return nil, fmt.Errorf("starlark call error: %w", err)
	}

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	return starlarkToGo(result)
}

// starlarkToGo converts a Starlark value into a plain Go value.
//
// Mapping:
//
//	None              -> nil
//	string            -> string
//	bool              -> bool
//	int               -> int64 (error if it does not fit in int64)
//	float             -> float64
//	list, tuple       -> []any (elements converted recursively)
//	dict              -> map[string]any (keys MUST be strings)
//
// Any other type (set, function, ...) returns an error. When reading a result
// in Go, use type assertions, e.g. `n := res.(int64)`.
func starlarkToGo(v starlark.Value) (any, error) {
	if v == starlark.None {
		return nil, nil
	}
	if s, ok := starlark.AsString(v); ok {
		return s, nil
	}
	switch x := v.(type) {
	case starlark.Bool:
		return bool(x), nil
	case starlark.Int:
		i, ok := x.Int64()
		if !ok {
			return nil, fmt.Errorf("integer out of int64 range")
		}
		return i, nil
	case starlark.Float:
		return float64(x), nil
	case *starlark.List:
		out := make([]any, 0, x.Len())
		for i := 0; i < x.Len(); i++ {
			item, err := starlarkToGo(x.Index(i))
			if err != nil {
				return nil, err
			}
			out = append(out, item)
		}
		return out, nil
	case starlark.Tuple:
		out := make([]any, 0, len(x))
		for _, item := range x {
			gv, err := starlarkToGo(item)
			if err != nil {
				return nil, err
			}
			out = append(out, gv)
		}
		return out, nil
	case *starlark.Dict:
		out := map[string]any{}
		for _, item := range x.Items() {
			k, ok := starlark.AsString(item[0])
			if !ok {
				return nil, fmt.Errorf("dict key must be string, got %s", item[0].Type())
			}
			gv, err := starlarkToGo(item[1])
			if err != nil {
				return nil, err
			}
			out[k] = gv
		}
		return out, nil
	default:
		return nil, fmt.Errorf("unsupported starlark return type: %s", v.Type())
	}
}

// goToStarlark converts a Go value into a Starlark value so it can be passed
// to a script.
//
// Mapping:
//
//		nil                          -> None
//		string, bool                 -> string, bool
//		int*, uint*                  -> int
//		float32, float64             -> float
//		map[string]any, []any        -> dict, list (converted recursively)
//	 map[string]json.RawMessage   -> dict raw json
//		any map with string keys     -> dict   (via reflection)
//		any slice or array           -> list   (via reflection)
//
// Structs, pointers and maps with non-string keys are NOT supported and return
// an error; convert them to a map[string]any first.
func goToStarlark(v any) (starlark.Value, error) {
	if v == nil {
		return starlark.None, nil
	}

	switch x := v.(type) {
	case string:
		return starlark.String(x), nil
	case bool:
		return starlark.Bool(x), nil
	case int:
		return starlark.MakeInt(x), nil
	case int8:
		return starlark.MakeInt64(int64(x)), nil
	case int16:
		return starlark.MakeInt64(int64(x)), nil
	case int32:
		return starlark.MakeInt64(int64(x)), nil
	case int64:
		return starlark.MakeInt64(x), nil
	case uint:
		return starlark.MakeUint(uint(x)), nil
	case uint8:
		return starlark.MakeUint64(uint64(x)), nil
	case uint16:
		return starlark.MakeUint64(uint64(x)), nil
	case uint32:
		return starlark.MakeUint64(uint64(x)), nil
	case uint64:
		return starlark.MakeUint64(x), nil
	case float32:
		return starlark.Float(x), nil
	case float64:
		return starlark.Float(x), nil
	case map[string]any:
		d := starlark.NewDict(len(x))
		for k, vv := range x {
			sv, err := goToStarlark(vv)
			if err != nil {
				return nil, err
			}
			if err := d.SetKey(starlark.String(k), sv); err != nil {
				return nil, err
			}
		}
		return d, nil
	case []any:
		items := make([]starlark.Value, 0, len(x))
		for _, vv := range x {
			sv, err := goToStarlark(vv)
			if err != nil {
				return nil, err
			}
			items = append(items, sv)
		}
		return starlark.NewList(items), nil
	case map[string]json.RawMessage:
		jsonBytes, err := json.Marshal(x)
		if err != nil {
			return nil, fmt.Errorf("Error marshaling map: %s", err)
		}
		jsonString := string(jsonBytes)

		thread := &starlark.Thread{Name: "main"}
		decodeFn, ok := starlarkjson.Module.Members["decode"].(starlark.Callable)
		if !ok {
			return nil, errors.New("failed to find starlark json decode function")
		}
		starlarkVal, err := starlark.Call(thread, decodeFn, starlark.Tuple{starlark.String(jsonString)}, nil)
		return starlarkVal, err
	case json.RawMessage:
		jsonString := x.String()
		thread := &starlark.Thread{Name: "main"}
		decodeFn, ok := starlarkjson.Module.Members["decode"].(starlark.Callable)
		if !ok {
			return nil, errors.New("failed to find starlark json decode function")
		}
		starlarkVal, err := starlark.Call(thread, decodeFn, starlark.Tuple{starlark.String(jsonString)}, nil)
		return starlarkVal, err

	}

	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Map:
		if rv.Type().Key().Kind() != reflect.String {
			return nil, fmt.Errorf("map keys must be strings, got %s", rv.Type().Key().String())
		}
		d := starlark.NewDict(rv.Len())
		iter := rv.MapRange()
		for iter.Next() {
			k := iter.Key().String()
			sv, err := goToStarlark(iter.Value().Interface())
			if err != nil {
				return nil, err
			}
			if err := d.SetKey(starlark.String(k), sv); err != nil {
				return nil, err
			}
		}
		return d, nil
	case reflect.Slice, reflect.Array:
		items := make([]starlark.Value, 0, rv.Len())
		for i := 0; i < rv.Len(); i++ {
			sv, err := goToStarlark(rv.Index(i).Interface())
			if err != nil {
				return nil, err
			}
			items = append(items, sv)
		}
		return starlark.NewList(items), nil
	}

	return nil, fmt.Errorf("unsupported type: %T", v)
}

package scripting

import (
	"context"
	"fmt"
	"reflect"

	"go.starlark.net/starlark"
	"go.starlark.net/syntax"
)

type StarlarkEngine struct{}

func NewStarlarkEngine() *StarlarkEngine {
	return &StarlarkEngine{}
}

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

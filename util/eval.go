package util

import (
	"github.com/dop251/goja"
)

func EvalJS(jsCode string) map[string]interface{} {
	vm := goja.New()
	value, err := vm.RunString(jsCode)
	if err != nil {
		return nil
	}

	result := make(map[string]interface{})
	if err := vm.ExportTo(value, &result); err != nil {
		return nil
	}

	return result
}

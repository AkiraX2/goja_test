package main

import (
	"fmt"
	"os"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/console"
	"github.com/dop251/goja_nodejs/require"
)

func newVMWithRequire() (*goja.Runtime, *require.RequireModule) {
	registry := new(require.Registry)
	vm := goja.New()
	req := registry.Enable(vm)
	console.Enable(vm)

	vm.SetFieldNameMapper(goja.UncapFieldNameMapper())
	// https://github.com/dop251/goja#mapping-struct-field-and-method-names
	// use this if we need optionally uncapitalises
	// vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))
	return vm, req

}

func newVMWithAssert() (*goja.Runtime, *require.RequireModule) {
	vm, req := newVMWithRequire()
	// req.Require("./scripts/assert.js")
	vm.RunString(`
		var assert = require('./scripts/assert');
		var { 
			assertTrue,
			assertEqual
		} = assert;
	`)
	return vm, req
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func exportJsError(err error) interface{} {
	if jserr, ok := err.(*goja.Exception); ok {
		fmt.Printf(jserr.Value().Export().(string))
		return jserr.Value().Export()
	}
	return nil
}

// TODO : Compile and RunProgram

func RunString(vm *goja.Runtime, src string) (goja.Value, interface{}) {
	if vm == nil {
		vm, _ = newVMWithAssert()
	}
	res, err := vm.RunString(src)
	checkError(err)
	jserr := exportJsError(err)
	return res, jserr
}

func RunScriptFromFile(vm *goja.Runtime, path string) (goja.Value, interface{}) {
	if vm == nil {
		vm, _ = newVMWithAssert()
	}

	if script, err := os.ReadFile(path); err != nil {
		panic(err)

	} else {
		res, err := vm.RunScript(path, string(script))
		checkError(err)
		jserr := exportJsError(err)
		return res, jserr
	}

}

package main

import (
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/dop251/goja"
)

func TestRunString(t *testing.T) {
	// TODO: test req.Require

	/// general functions
	// passing nil
	t.Run("passing nil or vm", func(t *testing.T) {
		assert.NotPanics(t, func() {
			RunString(nil, `1+1`)
		})

	})

	// passing vm
	t.Run("passing vm", func(t *testing.T) {

		vm := goja.New()
		assert.NotPanics(t, func() {
			RunString(vm, `1+1`)
		})

		assert.NotPanics(t, func() {
			RunString(vm, `1+1`)
			RunString(vm, `2+3`)
			RunString(vm, `1+3`)
		})
	})

	// empty script
	t.Run("empty", func(t *testing.T) {
		assert.NotPanics(t, func() {
			RunString(nil, ``)
		})

		assert.NotPanics(t, func() {
			RunString(nil, `

			
			
			`)
		})

	})

	// comments
	t.Run("comments", func(t *testing.T) {
		assert.NotPanics(t, func() {
			RunString(nil, `
				// this is a comment
			`)
		})

		assert.NotPanics(t, func() {
			RunString(nil, `
				/* this is a comment */
			`)
		})

		assert.NotPanics(t, func() {
			RunString(nil, `
				/**
				* 
				* this is a multiline comment
				* 
				*/

				/* this is a multiline comment */
			`)
		})

		assert.NotPanics(t, func() {
			RunString(nil, `
				/// this is a multiline comment
				/// this is a multiline comment
			`)
		})

	})

	/// fails
	// fail to run non js code
	t.Run("fail on non js code", func(t *testing.T) {
		assert.Panics(t, func() {
			RunString(nil, `
				// python
				print("hello world")
			`)
		})

	})

	// can detect js syntax error
	t.Run("fail on func not defined", func(t *testing.T) {
		assert.Panics(t, func() {
			RunString(nil, `
				1++1	// syntax error
			 `)
		})
		assert.Panics(t, func() {
			RunString(nil, `
				var res = mul(2, 3);	// mul is not defined
			 `)
		})

	})

	// can detect js runtime error
	t.Run("fail on runtime error", func(t *testing.T) {
		assert.Panics(t, func() {
			RunString(nil, `
				var foo = undefined;
				foo.bar;				// <-error
			`)
		})
	})

	t.Run("fail on throw", func(t *testing.T) {
		assert.Panics(t, func() {
			RunString(nil, `
				throw("Test");
			`)
		})

		assert.Panics(t, func() {
			RunString(nil, `
				throw new Error("Test");
			`)
		})
	})

	// js assert
	t.Run("js assert", func(t *testing.T) {
		assert.NotPanics(t, func() {
			RunString(nil, `
 				assertTrue(1 == 1);

				assertEqual(1, 1);
			`)
		})

		assert.Panics(t, func() {
			RunString(nil, `
 				assertTrue(1 == 2); 
			`)
		})
	})

	// js long running
	t.Run("js long running", func(t *testing.T) {
		// FIXME : detect long running
		assert.Panics(t, func() {
			RunString(nil, `
				var i = 0;
				while (i < 1000000) {
					i++;
				}
				throw "Done";	// for watch
			`)
		})
	})

	// await until promise callback done
	t.Run("wait until promise callback done", func(t *testing.T) {
		// NOTE: this is not allowed
		assert.NotPanics(t, func() {
			RunString(nil, `
				 var promise = new Promise(function(resolve, reject) {
					// ...
					resolve("done");
				});
				promise.then(function(value) {
					// ...
					throw "Done";	// for watch
				});

			`)
		})
	})

}

func TestRunScriptFromFile(t *testing.T) {

	// console.log
	t.Run("empty", func(t *testing.T) {
		assert.NotPanics(t, func() {
			RunScriptFromFile(nil, "./scripts/empty.js")
		})
	})

	// non existent
	t.Run("non existent", func(t *testing.T) {
		assert.Panics(t, func() {
			RunScriptFromFile(nil, "non_existent.js")
		})
	})

	// module.exports
	t.Run("existing", func(t *testing.T) {
		assert.Panics(t, func() {
			// this is not allowed
			// file include exports
			RunScriptFromFile(nil, "./scripts/m.js")
		})
	})

}

func TestRunJs(t *testing.T) {

	// primitives
	t.Run("primitives", func(t *testing.T) {
		assert.NotPanics(t, func() {
			RunString(nil, `
				var foo = 0;
				var bar = "";
				var baz = true;
				var id = Symbol("id");
			`)
		})

		assert.NotPanics(t, func() {
			RunString(nil, `
				var foo = null;
				var bar = undefined;
				var baz = {};
				var qux = [];
			`)
		})

		// bigint es2020
		assert.Panics(t, func() {
			RunString(nil, `
				var foo = 1n;
				var baz = 9007199254740991n;
			`)
		})
	})

	// func
	t.Run("func", func(t *testing.T) {
		assert.NotPanics(t, func() {
			RunString(nil, `
			var mul = function(a, b) {
				return a * b;
			}
			var res = mul(2, 3);
		`)
		})
	})

	// Number
	t.Run("Number", func(t *testing.T) {
		assert.NotPanics(t, func() {
			RunString(nil, `
				var num = 1;
				var max = Number.MAX_VALUE;
				var min = Number.MIN_VALUE;
				var nan = Number.NaN;

				assertTrue(max > 0);
				assertTrue(min > 0);
				assertTrue(isNaN(nan));
			`)

			// es6
			RunString(nil, `
				assertTrue(Number.isInteger(1));
				assertTrue(!Number.isInteger(1.5));
				assertTrue(Number.isSafeInteger(Number.MAX_SAFE_INTEGER));
				assertTrue(Number.isSafeInteger(Number.MIN_SAFE_INTEGER));
				assertTrue(!Number.isSafeInteger(12345678901234567890));

				assertTrue(isFinite(1));
				assertTrue(!isFinite(Infinity));	
				assertTrue(isNaN("xxx"));
		
			`)

		})

		assert.Panics(t, func() {
			// num sep _ es2021
			RunString(nil, `
				var num = 1_000_000_000;
				assertTrue(num == 1000000000);
			`)
		})
	})

	// String
	t.Run("String", func(t *testing.T) {
		assert.NotPanics(t, func() {
			RunString(nil, `
				var str = "hello world";
				var len = str.length;
				var sub = str.substring(0, 5);
				var idx = str.indexOf("world");
				var up = str.toUpperCase();
				var low = str.toLowerCase();
				var trim = str.trim();
				var split = str.split(" ");
				var rep = str.replace("world", "earth");

				assertTrue(len == 11);
				assertTrue(sub == "hello");
				assertTrue(idx == 6);
				assertTrue(up == "HELLO WORLD");
				assertTrue(low == "hello world");
				assertTrue(trim == "hello world");
				assertTrue(split[0] == "hello");
				assertTrue(split[1] == "world");
				assertTrue(rep == "hello earth");
			`)

			// includes es6
			RunString(nil, `
				var str = "hello world";
				var res = str.includes("world");
				assertTrue(res == true);
			`)

			// startsWith endsWith es6
			RunString(nil, `
				var str = "hello world";
 				assertTrue(str.startsWith("hello") == true);
 				assertTrue(str.endsWith("world") == true);
			`)

			// padStart padEnd es2017
			RunString(nil, `
				var str = "hello";
				// assertTrue(str.padStart(10, "0") == "00000hello");
				assertTrue(str.padEnd(10, 0) == "hello00000");
			`)

			// trimStart trimEnd es2019
			RunString(nil, `
				var str = "  hello world  ";
				assertTrue(str.trimStart() == "hello world  ");
				assertTrue(str.trimEnd() == "  hello world");
			`)

			// matchAll es2020
			RunString(nil, `
				var str = "hello world";
				var res = str.matchAll(/l/g);

				assertEqual(res.next().value.index, 2);
				assertEqual(res.next().value.index, 3);
				assertEqual(res.next().value.index, 9);
				assertEqual(res.next().done, true);
			`)

			// replaceAll es2021
			RunString(nil, `
				var str = "hello world";
				var res = str.replaceAll("l", "x");

				assertTrue(res == "hexxo worxd");
			`)

		})
	})

	/// nodejs
	/// https://github.com/dop251/goja_nodejs
	/// TODO: commonjs
	/// https://github.com/tliron/commonjs-goja

	/// module
	// exports
	t.Run("module.exports", func(t *testing.T) {
		vm, req := newVMWithRequire()
		_ = req
		assert.Panics(t, func() {
			// this is not allowed
			// module is not defined
			RunString(vm, `
				module.exports = function(a, b) {
					return a * b;
				}
			`)
		})

		// module.exports
		assert.NotPanics(t, func() {
			RunString(nil, `
				var m = require("./scripts/m.js");
				var res = m.test();

				assertTrue(res == "test");
			`)
		})

		// exports
		assert.NotPanics(t, func() {
			RunString(nil, `
				var m = require("./scripts/m.js");
				var res = m.test();

				assertTrue(res == "test");
			`)
		})
	})

	// require
	t.Run("require", func(t *testing.T) {
		vm, _ := newVMWithAssert()
		// req.Require("./scripts/m.js")
		assert.NotPanics(t, func() {
			RunString(vm, `
				 var m = require("./scripts/m.js");
				 var res = m.test();

				 assertTrue(res == "test");
			`)
		})
	})

	// import
	t.Run("import", func(t *testing.T) {
		vm, _ := newVMWithAssert()
		// req.Require("./scripts/m.js")
		assert.Panics(t, func() {
			RunString(vm, `
				import { test } from "./scripts/m.js";
				var res = test();

				assertTrue(res == "test");
			`)
		})
	})

	// require and export
	t.Run("require and export", func(t *testing.T) {
		vm, _ := newVMWithRequire()
		// req.Require("./scripts/assert.js")
		// req.Require("./scripts/m.js")
		// req.Require("./scripts/require_and_export.js")
		assert.NotPanics(t, func() {
			RunString(vm, `
				var m = require("./scripts/require_and_export");
  				var res = m.test();

				m.assertTrue(res == "test");
			`)
		})
	})

	// console
	t.Run("console", func(t *testing.T) {

		assert.NotPanics(t, func() {
			RunString(nil, `
				console.log("hello world");
			`)
		})
	})

	// process
	t.Run("process", func(t *testing.T) {

		assert.Panics(t, func() {
			RunString(nil, `
			const process = require('node:process');
			process.on('exit', (code) => {

			});
			`)

		})

		// TODO: enable process module
		// github.com/dop251/goja_nodejs/process
	})
	// fs
	t.Run("fs", func(t *testing.T) {
		assert.Panics(t, func() {
			// Invalid module
			RunString(nil, `
				var fs = require("fs");
				var data = fs.readFileSync("./scripts/data.txt", "utf8");
				assertTrue(data == "hello world");
			`)
		})
	})

	/// Specs

	// try catch
	t.Run("try catch", func(t *testing.T) {
		assert.NotPanics(t, func() {
			RunString(nil, `
				try {
					throw new Error("Test");
				} catch (e) {
					assertTrue(e.message == "Test");
				}
			`)
			// finally
			RunString(nil, `
				var func = function() {
					let res = 0;
					try {
						// ...
						throw new Error("Test");
					} catch (e) {
						assertTrue(e.message == "Test");
					} finally {
						res = -1;
					}
					return res;
				}

				assertTrue(func() == -1);
			`)
		})
	})

	// typeof

	t.Run("typeof", func(t *testing.T) {
		assert.NotPanics(t, func() {
			RunString(nil, `
				assertTrue(typeof "John" == "string");
				assertTrue(typeof 3.14 == "number");
				assertTrue(typeof true == "boolean");
				assertTrue(typeof undefined == "undefined");
				assertTrue(typeof null == "object");
				assertTrue(typeof {} == "object");
				assertTrue(typeof [] == "object");
				assertTrue(typeof function(){} == "function");
			`)
		})
	})

	// SetTimeout
	t.Run("setTimeout", func(t *testing.T) {
		assert.Panics(t, func() {
			RunString(nil, `
				setTimeout(function() {
					console.log("hello world");
				}, 1000);
			`)
		})
	})

	// SetInterval
	t.Run("setInterval", func(t *testing.T) {
		assert.Panics(t, func() {
			RunString(nil, `
				setInterval(function() {
					console.log("hello world");
				}, 1000);
			`)
		})
	})

	// Date
	t.Run("Date", func(t *testing.T) {
		assert.NotPanics(t, func() {
			RunString(nil, `
				var date = new Date();
				var year = date.getFullYear();
				var month = date.getMonth();
				var day = date.getDate();
				var hour = date.getHours();
				var min = date.getMinutes();

				var now = Date.now(); 
				
				assertTrue(year > 2020);
				assertTrue(month >= 0);
				assertTrue(day >= 1);
				assertTrue(hour >= 0);
				assertTrue(min >= 0);
				assertTrue(now > 0);
			`)
		})

		// TODO
	})

	// JSON
	// https://github.com/dop251/goja?tab=readme-ov-file#json
	// only utf8 allowed
	// utf16 not allowed
	t.Run("JSON", func(t *testing.T) {
		assert.NotPanics(t, func() {
			RunString(nil, `
				var obj = {
					"foo": "bar",
					"baz": 1
				};
				var str = JSON.stringify(obj);
				var obj2 = JSON.parse(str);
				assertTrue(obj2.foo == "bar");
				assertTrue(obj2.baz == 1);
			`)
		})
	})

	// Math
	t.Run("Math", func(t *testing.T) {
		assert.NotPanics(t, func() {
			RunString(nil, `
 				assertEqual(Math.max(1, 2, 3), 	3);
				assertEqual(Math.min(1, 2, 3),  1);
				assertEqual(Math.round(1.5),    2);
				assertEqual(Math.floor(1.5),    1);
				assertEqual(Math.ceil(1.5),     2);
				assertEqual(Math.abs(-1),       1);
				assertEqual(Math.pow(2, 3),     8);
				assertEqual(Math.sqrt(4),       2);
				assertEqual(Math.cbrt(8),       2);
				assertEqual(Math.sin(Math.PI/2), 1);
				assertEqual(Math.log(Math.E),    1);
				assertEqual(Math.log10(100),     2);
				assertEqual(Math.log2(8),        3);
				assertEqual(Math.exp(1),         Math.E);
				assertEqual(Math.random() >= 0,  true);
			`)

			// es6
			RunString(nil, `
				assertEqual(Math.clz32(0), 32);
				assertEqual(Math.imul(2, 3), 6);
				assertEqual(Math.sign(-10), -1);
				assertEqual(Math.trunc(1.5), 1);
				assertEqual(Math.log1p(Math.E - 1), 1);
				assertEqual(Math.expm1(Math.log(Math.E)), Math.E - 1);
				assertEqual(Math.hypot(3, 4), 5);
			`)

		})
	})

	// RegExp
	t.Run("RegExp", func(t *testing.T) {
		assert.NotPanics(t, func() {
			RunString(nil, `
				var re = /hello/;
				var str = "hello world";
				var res = re.test(str);
				assertTrue(res == true);

				var re2 = new RegExp("hello");
				var res2 = re2.test(str);
				assertTrue(res2 == true);

				var re = /hello/;
				var str = "hello world";
				var res = re.exec(str);
				assertTrue(res[0] == "hello");
			`)
		})

		assert.NotPanics(t, func() {
			RunString(nil, `
				// complicated regex
				var re3 = /[a-z]+\d{2,4}/i;
				var str3 = "abc123";
				var res3 = re3.test(str3);
				assertTrue(res3 == true);
				
				// more complicated regex
				var regex = /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,}$/;
				
				assertEqual(regex.test("Password1@"), true);
				assertEqual(regex.test("password1@"), false);
				assertEqual(regex.test("Password1"), false);
				assertEqual(regex.test("Pwd@"), false);
			`)
		})

		// TODO: revise

	})

	// prototype
	t.Run("prototype", func(t *testing.T) {
		assert.NotPanics(t, func() {
			RunString(nil, `
				var obj = {
					"foo": "bar",
					"baz": 1
				}
				var proto = Object.getPrototypeOf(obj);
				assertTrue(proto == Object.prototype);

				proto.test = function() {
					return "test";
				}

				assertTrue(obj.test() == "test");
			`)
		})
	})

	/// es5

	// strict mode
	t.Run("strict mode", func(t *testing.T) {
		assert.NotPanics(t, func() {
			RunString(nil, `
				"use strict";
				var foo = 1;
			`)

		})

		assert.Panics(t, func() {
			RunString(nil, `
				"use strict";
				foo = 1;
			`)
		})
	})

	// Array
	t.Run("Array", func(t *testing.T) {

		// from ES6
		assert.NotPanics(t, func() {
			RunString(nil, `
				var arr = Array.from("hello");
				assertTrue(arr[0] == "h");
				assertTrue(arr[1] == "e");
			`)
		})

		// keys ES6
		assert.NotPanics(t, func() {
			RunString(nil, `
				var arr = [1, 2, 3, 4, 5];
				var keys = arr.keys();

				assertTrue(keys.next().value == 0);
				assertTrue(keys.next().value == 1);
			`)
		})

		// entries ES6
		assert.NotPanics(t, func() {
			RunString(nil, `
				var arr = [-1, -2, -3, -4, -5];
				var pairs = arr.entries();
		
				assertTrue(pairs.next().value[0] == 0);
				assertTrue(pairs.next().value[1] == -2);
			`)
		})

		// find
		assert.NotPanics(t, func() {
			RunString(nil, `
				var arr = [1, 2, 3, 4, 5];
				var res = arr.find(function(x) {
					return x % 2 == 0;
				});

				assertTrue(res == 2);
			`)
		})

		// findIndex
		assert.NotPanics(t, func() {
			RunString(nil, `
				var arr = [1, 2, 3, 4, 5];
				var res = arr.findIndex(function(x) {
					return x % 2 == 0;
				});

				assertTrue(res == 1);
			`)
		})

		// includes es2016
		assert.NotPanics(t, func() {
			RunString(nil, `	
				var arr = [1, 2, 3, 4, 5];

				assertTrue(arr.includes(2));
				assertTrue(!arr.includes(0));
			`)
		})

		// map
		assert.NotPanics(t, func() {
			RunString(nil, `
				var arr = [1, 2, 3, 4, 5];
				var res = arr.map(function(x) {
					return x * 2;
				});

				assertTrue(res[0] == 2);
				assertTrue(res[1] == 4);
			`)
		})

		// filter
		assert.NotPanics(t, func() {
			RunString(nil, `
				var arr = [1, 2, 3, 4, 5];
				var res = arr.filter(function(x) {
					return x % 2 == 0;
				});

				assertTrue(res[0] == 2);
				assertTrue(res[1] == 4);
			`)
		})

		// reduce
		assert.NotPanics(t, func() {
			RunString(nil, `
				var arr = [1, 2, 3, 4, 5];
				var res = arr.reduce(function(acc, x) {
					return acc + x;
				});

				assertTrue(res == 15);
			`)
		})

		// some every
		assert.NotPanics(t, func() {
			RunString(nil, `
				var arr = [1, 2, 3, 4, 5];
				var res = arr.some(function(x) {
					return x % 2 == 0;
				});

				assertTrue(res == true);
			`)

			RunString(nil, `
				var arr = [1, 2, 3, 4, 5];
				var res = arr.every(function(x) {
					return x % 2 == 0;
				});

				assertTrue(res == false);
			`)
		})

		// forEach
		assert.NotPanics(t, func() {
			RunString(nil, `
				var arr = [1, 2, 3];
				var sum = 0;
				arr.forEach(function(x) {
					sum += x;
				});
		
				assertTrue(sum == 6);
			`)
		})

		// sort and reverse
		assert.NotPanics(t, func() {
			RunString(nil, `
				var arr = [3, 2, 1];
				arr.sort();
				assertTrue(arr[0] == 1);

				arr.reverse();
				assertTrue(arr[0] == 3);

			`)

		})

		// flat es2019
		assert.NotPanics(t, func() {
			RunString(nil, `
				var arr = [1, 2, [3, 4]];
				var arr2 = arr.flat();
				assertTrue(arr2[2] == 3);
			`)
		})

		// at es2022
		assert.NotPanics(t, func() {
			RunString(nil, `
				var arr = [1, 2, 3, 4, 5];
				var res = arr.at(2);
				assertTrue(res == 3);
			`)
		})

		// findLast findLastIndex es2023
		assert.NotPanics(t, func() {
			RunString(nil, `
				var arr = [1, 2, 3, 4, 5];
				var res = arr.findLast(function(x) {
					return x % 2 == 0;
				});

				assertTrue(res == 4);

				var res2 = arr.findLastIndex(function(x) {
					return x % 2 == 0;
				});

				assertTrue(res2 == 3);
			`)
		})

		// toSorted toReversed ES2023
		assert.Panics(t, func() {
			RunString(nil, `
				var arr = [3, 2, 1];
				var arr2 = arr.toSorted();
				assertTrue(arr2[0] == 1);

				var arr3 = arr.toReversed();
				assertTrue(arr3[0] == 3);
			`)
		})

		// with ES2023
		assert.Panics(t, func() {
			RunString(nil, `
				var arr = [1, 2, 3, 4, 5];
				arr = arr.with(0, -1);
				assertTrue(arr[0] == -1);
			`)
		})

		// toSpliced ES2023
		assert.Panics(t, func() {
			RunString(nil, `
				var arr = [1, 2, 3, 4, 5];
				var arr2 = arr.toSpliced(2, 1);
				assertTrue(arr2[2] == 4);
			`)
		})

	})

	// Object
	t.Run("Object", func(t *testing.T) {
		// entries es2017
		assert.NotPanics(t, func() {
			RunString(nil, `
				var obj = {
					"foo": "bar",
					"baz": 1
				}
				var pairs = Object.entries(obj)

				assertTrue(pairs[0][0] == "foo");
				assertTrue(pairs[1][1] == 1);
			`)
		})

		// values es 2017
		assert.NotPanics(t, func() {
			RunString(nil, `
				var obj = {
					"foo": "bar",
					"baz": 1
				}
				var values = Object.values(obj)

				assertTrue(values[0] == "bar");
			`)
		})

		// keys
		assert.NotPanics(t, func() {
			RunString(nil, `
				var obj = {
					"foo": "bar",
					"baz": 1
				}
				var keys = Object.keys(obj)

				assertTrue(keys[0] == "foo");
			`)
		})

		// assign
		assert.NotPanics(t, func() {
			RunString(nil, `
				var obj = {
					"foo": "bar",
					"baz": 1
				}
				var obj2 = {
					"qux": 2
				}
				Object.assign(obj, obj2);
				assertTrue(obj.qux == 2);
			`)
		})

		// defineProperty es5
		assert.NotPanics(t, func() {
			RunString(nil, `
				var obj = {};
				Object.defineProperty(obj, "foo", {
					value: "bar",
					writable: false
				});
				assertTrue(obj.foo == "bar");
			`)
		})

		// fromEntries es2019
		assert.NotPanics(t, func() {
			RunString(nil, `
				var arr = [["foo", "bar"], ["baz", 1]];
				var obj = Object.fromEntries(arr);
				assertTrue(obj.foo == "bar");
			`)
		})

		// hasOwn es2022
		assert.Panics(t, func() {
			RunString(nil, `
				var obj = {
					"foo": "bar",
					"baz": 1
				}
				assertTrue(obj.hasOwnProperty("foo") == true);
				assertTrue(obj.hasOwn("foo") == true);
			`)
		})

	})

	// WeakMap
	t.Run("WeakMap", func(t *testing.T) {
		// TODO
		assert.NotPanics(t, func() {
			RunString(nil, `
				var m = new WeakMap();
				var key = {};
				var value = {/* a very large object */};
				m.set(key, value);
				value = undefined;
				m = undefined; // The value does NOT become garbage-collectable at this point
				key = undefined; // Now it does
				// m.delete(key); // This would work too
			`)
		})
	})

	// WeakSet
	t.Run("WeakSet", func(t *testing.T) {
		assert.NotPanics(t, func() {
			RunString(nil, `
				var ws = new WeakSet(); 
				// TODO
			`)
		})
	})

	// WeakRef
	t.Run("WeakRef", func(t *testing.T) {
		assert.Panics(t, func() {
			RunString(nil, `
				var wr = new WeakRef({});
			`)
		})
	})

	///
	/// es6

	// const, let
	t.Run("const, let", func(t *testing.T) {
		assert.NotPanics(t, func() {
			RunString(nil, `
				const foo = 1;
				let bar = 2;
				assertTrue(foo == 1);
				assertTrue(bar == 2);
			`)
		})
	})

	// class
	t.Run("class", func(t *testing.T) {
		assert.NotPanics(t, func() {
			RunString(nil, `
				class Rectangle {
					constructor(height, width) {
						this.height = height;
						this.width = width;
					}

					get area() {
						return this.calcArea();
					}

					calcArea() {
						return this.height * this.width;
					}

					field = 0;	// es2022

					#name = "private";	// es2022

				}
				var square = new Rectangle(2, 2);
				assertTrue(square.height == 2);
				assertTrue(square.area == 4);
				assertTrue(square.field == 0);
				// assertTrue(square.#name == "private");
			`)
		})
	})

	// destructuring
	t.Run("destructuring", func(t *testing.T) {
		assert.NotPanics(t, func() {
			RunString(nil, `
				var obj = {a: 1, b: 2, c: 3};
				var {a, b, c} = obj;
				assertTrue(a == 1);
				assertTrue(b == 2);
				assertTrue(c == 3);
			`)
		})
	})

	// default parameter
	t.Run("default parameter", func(t *testing.T) {
		assert.NotPanics(t, func() {
			RunString(nil, `
				function foo(x = 1) {
					return x;
				}

				assertTrue(foo() == 1);
				assertTrue(foo(2) == 2);
			`)
		})
	})

	// rest parameter
	t.Run("rest parameter", func(t *testing.T) {
		assert.NotPanics(t, func() {
			RunString(nil, `
			function arr(...args) {
				return args;
			  }
			  
			  var x = arr(1, 2, 3);
			  assertTrue(x[0] == 1);
			  assertTrue(x[2] == 3);
			`)
		})
	})

	// spread operator
	t.Run("spread operator", func(t *testing.T) {
		assert.NotPanics(t, func() {
			RunString(nil, `
				var arr = [1, 2, 3];
				var arr2 = [4, 5, 6];

				var arr3 = [...arr, ...arr2];

				assertTrue(arr3[0] == 1);
				assertTrue(arr3[3] == 4);
			`)
		})
	})

	// for of
	t.Run("for of", func(t *testing.T) {
		assert.NotPanics(t, func() {
			RunString(nil, `
				var arr = [1, 2, 3];
				var sum = 0;
				for (var x of arr) {
					sum += x;
				}
				
				assertTrue(sum == 6);
			`)
		})

		// obj
		assert.NotPanics(t, func() {
			RunString(nil, `
				var obj = {
					"foo": 1,
					"bar": 2,
					"baz": 3
				}
				var sum = 0;
				for (var x of Object.values(obj)) {
					sum += x;
				}
				
				assertTrue(sum == 6);
			`)
		})
	})

	// Map
	t.Run("Map", func(t *testing.T) {
		assert.NotPanics(t, func() {
			RunString(nil, `
				var map = new Map([
					["foo", "bar"],
				]);
				map.set("baz", 1);
				assertTrue(map.get("foo") == "bar");
				assertTrue(map.get("baz") == 1);
			`)
		})
	})

	// Set
	t.Run("Set", func(t *testing.T) {
		assert.NotPanics(t, func() {
			RunString(nil, `
				var set = new Set([1, 2, 3]);
				set.add(1);

				assertTrue(set.has(1));
				assertTrue(set.has(2));
			`)
		})
	})

	// Promise
	// FIXME: await callback done
	t.Run("Promise", func(t *testing.T) {
		assert.NotPanics(t, func() {
			RunString(nil, `
				var promise = new Promise(function(resolve, reject) {
					// ...
					resolve("done");
				});
				promise.then(function(value) {
					assertTrue(value == "done");
				});
			`)

			RunString(nil, `
				var promise = new Promise(function(resolve, reject) {
					// ...
					reject("error");
				});

				promise.catch(function(value) {
					assertTrue(value == "error");
				});
			`)

			// Promise.all
			RunString(nil, `
				var promise1 = Promise.resolve(3);
				var promise2 = 42;
				
				var res = Promise.all([promise1, promise2]);
				res.then(function(value) {
					assertTrue(value[0] == 3);
					assertTrue(value[1] == 42);
				});
			`)

			// Promise.allSettled es2020
			RunString(nil, `
				var promise1 = Promise.resolve(3);
				var promise2 = Promise.reject("error");

				var res = Promise.allSettled([promise1, promise2]);
				res.then(function(value) {
					assertTrue(value[0].status == "fulfilled");
					assertTrue(value[1].status == "rejected");
				});
			`)

			// finally es2018
			RunString(nil, `
				var promise = new Promise(function(resolve, reject) {
					// ...
					resolve("done");
				});
				promise.finally(function(value) {
					assertTrue(value == "done");
				});
			`)

			// Promise.any es2021
			RunString(nil, `
				var promise1 = Promise.reject("error");
				var promise2 = Promise.resolve("done");

				var res = Promise.any([promise1, promise2]);
				res.then(function(value) {
					assertTrue(value == "done");
				});
			`)
		})
	})

	// async, await es2017
	// FIXME : await callback done
	t.Run("async, await", func(t *testing.T) {
		assert.NotPanics(t, func() {
			RunString(nil, `
				async function foo() {
					return "done";
				}
				async function bar() {
					return await foo();
				}
				bar().then(function(value) {
					assertTrue(value == "done");
				});
			`)
		})
	})

	// Generators
	t.Run("Generators", func(t *testing.T) {
		assert.NotPanics(t, func() {
			RunString(nil, `
				function* foo() {
					yield 1;
					yield 2;
					yield 3;
				}
				var gen = foo();
				assertTrue(gen.next().value == 1);
				assertTrue(gen.next().value == 2);
			`)
		})
	})

	// template string
	t.Run("template string", func(t *testing.T) {
		assert.NotPanics(t, func() {
			RunScriptFromFile(nil, "./scripts/template_string.js")
		})
	})

	// arrow function
	t.Run("=>", func(t *testing.T) {
		RunString(nil, `
			var mul = (a, b) => a * b;
			var res = mul(2, 3);

			assertTrue(res == 6);
		`)
	})

	// Operator
	t.Run("Operator", func(t *testing.T) {

		assert.NotPanics(t, func() {
			// ** **= es2016
			RunString(nil, `
				var foo = 2;
				foo **= 3;
				assertTrue(foo == 8);
			`)

			// ?? es2020
			RunString(nil, `
			var foo = null;
			var bar = foo ?? "default";
			assertTrue(bar == "default");
			`)

			// ?. es2020
			RunString(nil, `
			var obj = {
				"foo": {
					"bar": "baz"
				}
			}
			var res = obj.foo?.bar;
			assertTrue(res == "baz");
			`)

		})

		assert.Panics(t, func() {
			// &&= ||= es2020
			RunString(nil, `
				var foo = 1;
				foo &&= 2;
				assertTrue(foo == 2);

				var bar = 1;
				bar ||= 2;
				assertTrue(bar == 1);
			`)

		})

		assert.Panics(t, func() {
			// ??=
			RunString(nil, `
				var foo = null;
				foo ??= "default";
				assertTrue(foo == "default");
			`)
		})

	})

}

// Test Accessibility & Security
func TestAccessibility(t *testing.T) {
	// LocalFile
	assert.Panics(t, func() {
		RunString(nil, `
			var fs = require("fs");
		`)
	})

	// Network
	assert.Panics(t, func() {
		RunString(nil, `
			var net = require("net");
		`)
	})
	assert.Panics(t, func() {
		RunString(nil, `
			var http = require("http");
		`)
	})

	// System
	assert.Panics(t, func() {
		RunString(nil, `
			var process = require("process");
		`)
	})

	// TODO

}

func TestException(t *testing.T) {
	// TODO: error handling
	t.Run("error handling", func(t *testing.T) {
		assert.NotPanics(t, func() {
			defer func() {
				if r := recover(); r != nil {
					if jserr, ok := r.(*goja.Exception); ok {
						assert.Equal(t, jserr.Value().ToString().String(), "Test Error")
						assert.Equal(t, jserr, nil)

						re := regexp.MustCompile(`\d+:\d+`)

						match := re.FindString(jserr.Error())
						assert.Equal(t, match, "1:1")

					} else {
						panic("Not a goja.Exception")
					}
				}
			}()

			RunString(nil, `throw "Test Error";`)
		})

		assert.NotPanics(t, func() {
			defer func() {
				if r := recover(); r != nil {
					if jserr, ok := r.(*goja.Exception); ok {
						assert.Equal(t, jserr.Value().ToString().String(), "Error: Test Error")
						assert.Equal(t, jserr.Error(), jserr.Unwrap())

						re := regexp.MustCompile(`\d+:\d+`)

						match := re.FindString(jserr.Error())
						assert.Equal(t, match, "1:7")

					} else {
						panic("Not a goja.Exception")
					}
				}
			}()

			RunString(nil, `throw new Error("Test Error");`)
		})

	})
}

func TestInputAndOutput(t *testing.T) {

	assert.NotPanics(t, func() {

		parmas := map[string]interface{}{
			// "obj": map[string]interface{}{},
		}
		vm, _ := newVMWithAssert()
		vm.Set("parmas", parmas)
		RunString(vm, `
		var obj = {
			test: true,
			foo: -10
			};
		parmas.obj = obj; // obj gets Export()'ed, i.e. copied to a new map[string]interface{} and then this map is set as params["obj"]
		obj.test = false; // note, params.obj.test is still true
	`)
		assert.Equal(t, parmas["obj"].(map[string]interface{})["test"], true)
		assert.Equal(t, parmas["obj"].(map[string]interface{})["foo"], int64(-10))

	})

	type S struct {
		Field int
	}

	assert.NotPanics(t, func() {
		vm, _ := newVMWithAssert()
		obj := S{Field: -1}
		vm.Set("obj", &obj) // note here we pass a pointer to the object
		res, _ := RunString(vm, `
		
			assertEqual(obj.Field, undefined);
			assertEqual(obj.field, -1);
			
			
			obj.field = 10;
		`)

		assert.Equal(t, res.Export(), int64(10))
		assert.Equal(t, obj.Field, 10)
		assert.Equal(t, obj, S{Field: 10})

	})

	assert.NotPanics(t, func() {
		vm, _ := newVMWithAssert()
		obj := S{Field: -1}
		vm.Set("obj", obj) // note here
		res, _ := RunString(vm, `
		
			assertEqual(obj.Field, undefined);
			assertEqual(obj.field, -1);
			
			
			obj.field = 10;
		`)

		assert.Equal(t, res.Export(), int64(10))
		assert.Equal(t, obj.Field, -1)
		assert.Equal(t, obj, S{Field: -1})

	})
	// TODO: customize $ mocking

	// TODO: Built in methods and variables

}

// goroutine safe
// No. An instance of goja.Runtime can only be used by a single goroutine at a time.
// You can create as many instances of Runtime as you like but it's not possible to pass object values between runtimes.
// https://github.com/dop251/goja?tab=readme-ov-file#is-it-goroutine-safe

// TODO: performance
// Non-primary objectives

func TestCallJsFunc(t *testing.T) {

	assert.NotPanics(t, func() {

		const SCRIPT = `
		function sum(a, b) {
			return +a + b;
		}
		`

		vm := goja.New()
		_, err := vm.RunString(SCRIPT)
		if err != nil {
			panic(err)
		}
		sum, ok := goja.AssertFunction(vm.Get("sum"))
		if !ok {
			panic("Not a function")
		}

		res, err := sum(goja.Undefined(), vm.ToValue(40), vm.ToValue(2))
		if err != nil {
			panic(err)
		}

		assert.Equal(t, res.Export(), int64(42))

	})

	assert.NotPanics(t, func() {

		const SCRIPT = `
			function sum(a, b) {
				return +a + b;
			}
			`

		vm := goja.New()
		_, err := vm.RunString(SCRIPT)
		if err != nil {
			panic(err)
		}

		var sum func(int, int) int
		err = vm.ExportTo(vm.Get("sum"), &sum)
		if err != nil {
			panic(err)
		}

		assert.Equal(t, sum(40, 2), 42)
	})
}

func TestInterrupt(t *testing.T) {

	vm, _ := newVMWithAssert()
	time.AfterFunc(200*time.Millisecond, func() {
		vm.Interrupt("Test INT")
	})

	_, err := vm.RunString(`
			var i = 0;
			for (;;) {
				i++;
			}	
		
		`)

	if err == nil {
		t.Fatal("Err is nil")
	}
	// err is of type *InterruptedError and its Value() method returns whatever has been passed to vm.Interrupt()
	if err.(*goja.InterruptedError).Value() != "Test INT" {
		t.Fatal("Err type is wrong")
	}
}

// TODO: real world example

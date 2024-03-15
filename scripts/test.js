var assertTrue = require('./assert').assertTrue;
var assertEqual = require('./assert').assertEqual;

(() => {

    var obj = {
        a: 1,
        b: 2,
        c: 3
    }

    var pairs = Object.entries(obj);
    console.log(pairs)
})();


(() => {
    var arr = [1, 2, 3, 4, 5];

    assertTrue(arr.includes(2));
    assertTrue(!arr.includes(0));
})();

(() => {

    var arr = [-1, -2, -3, -4, -5];
    var pairs = arr.entries();
    

        assertTrue(pairs.next().value[0] == 0);
        assertTrue(pairs.next().value[1] == -2);
})();




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


console.log(Math.tan(Math.PI / 4))

var {sayHi} = require('./m.exports');
sayHi("John");


var str = "hello";
assertTrue(str.padStart(10, "0") == "00000hello");
assertTrue(str.padEnd(10, 0) == "hello00000");


// matchAll
var str = "hello world";
var res = str.matchAll(/l/g);

 assertEqual(res.next().value.index, 2);
 assertEqual(res.next().value.index, 3);
 assertEqual(res.next().value.index, 9);
 assertEqual(res.next().done, true);
 

var promise = new Promise(function(resolve, reject) {
    // ...
    resolve("done");
});
promise.then(function(value) {
    // long running code
    for(var i = 0; i < 1000000000; i++) {
        // ...
    }

    assertTrue(value == "done1");

});
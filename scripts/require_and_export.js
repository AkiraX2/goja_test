
var assert = require('./assert');
var test = require('./m').test;
var { 
    assertTrue,
    assertEqual
} = assert;
 
assertEqual(1, 1);

module.exports = {
    assertEqual: assertEqual,
    assertTrue: assertTrue,
    test:test
}
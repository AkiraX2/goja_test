const { assertEqual } = require("./assert");


function _errorAssert() {
  assertEqual(1, 2);
}

function errorAssert(){
  throw "Test Error"
    // _errorAssert(); // TODO: why is this not being caught?
}

module.exports = {
    errorAssert
}
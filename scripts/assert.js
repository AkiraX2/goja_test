function assertEqual(a, b) {
    if(a !== b) {
        var error =  new Error("Assertion failed", a, b);
        const callerStack = error.stack.split('\n')[2].trim();
        console.error('Assertion failed:', a, b, callerStack);
    }
}

function assertTrue(a) {
    if(!a) {
        var error = new Error("Assertion failed")
        const callerStack = error.stack.split('\n')[2].trim();
        console.error('Assertion failed:', callerStack);
        error.message = 'Assertion failed: ' + callerStack;
        // error.stack = error.stack.split('\n').slice(1).join('\n');
        // error.stack = error.stack.split('\n').slice(2).join('\n');

        throw error;
    }
}


module.exports = {
    assertEqual: assertEqual,
    assertTrue: assertTrue
}
function sayHi(user) {
    console.log(`Hello, ${user}!`);
}

// export function test() {
//     return "test";    
// }

function test() {
    return "test";    
}

module.exports = {
    sayHi: sayHi,
    test: test
}
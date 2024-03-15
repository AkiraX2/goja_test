function sayHi(user) {
    console.log(`Hello, ${user}!`);
}

// export function test() {
//     return "test";    
// }

function test() {
    return "test";    
}

exports.sayHi = sayHi;
exports.test = test;
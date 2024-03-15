var name = "world";
var str = `hello ${name}`;
assertTrue(str == "hello world");


var multiline = `hello
world`;

assertTrue(multiline == "hello\nworld");
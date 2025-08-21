// JavaScript test file for syntax highlighting
function greetUser(name, age) {
    // This is a comment
    const message = "Hello, " + name + "!";
    let isAdult = age >= 18;
    
    if (isAdult) {
        console.log(message);
        console.log("You are an adult.");
    } else {
        console.log("You are a minor.");
    }
    
    return {
        greeting: message,
        adult: isAdult,
        age: age
    };
}

class Person {
    constructor(name, age) {
        this.name = name;
        this.age = age;
    }
    
    greet() {
        return greetUser(this.name, this.age);
    }
}

const user = new Person("Alice", 25);
console.log(user.greet());
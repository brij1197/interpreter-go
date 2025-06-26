# Lox Interpreter

A Go implementation of the Lox programming language interpreter from Robert Nystrom's "Crafting Interpreters" book.

## Try It Online
**Web Interface**: [https://interpreter-go-production.up.railway.app](https://interpreter-go-production.up.railway.app)

Try the Lox interpreter directly in your browser with syntax highlighting, examples, and file upload/download capabilities.

## Branches

- **`master`** - Terminal/CLI version for local use
- **`dev`** - Web interface version with online playground

## Usage

### Web Version (this branch)
Visit the live demo above or run locally:
```bash
go run app/*.go web
# Open http://localhost:8080
```

### Terminal Version (master branch)
```bash
git checkout master

# Run a Lox file
./your_program.sh run script.lox

# Start REPL mode
./your_program.sh prompt

# Other commands
./your_program.sh tokenize script.lox
./your_program.sh parse script.lox
./your_program.sh evaluate script.lox
```

## Features

- **Variables and Scoping**: Local and global variable declarations with lexical scoping
- **Data Types**: Numbers, strings, booleans, and nil
- **Expressions**: Arithmetic, comparison, logical, and assignment operations
- **Control Flow**: If/else statements, while and for loops
- **Functions**: First-class functions with closures and recursion
- **Classes**: Object-oriented programming with inheritance
- **Methods**: Instance methods with `this` binding
- **Inheritance**: Class inheritance with `super` keyword support
- **Web Interface**: Browser-based code editor with examples and file operations

## Language Examples

### Variables and Basic Operations
```lox
var name = "World";
var number = 42;
var isTrue = true;

print "Hello, " + name + "!";
print "Age: " + number;
print "Student: " + isTrue;
```

### Control Flow
```lox
var x = 10;
if (x > 5) {
    print "x is greater than 5";
} else {
    print "x is 5 or less";
}

for (var i = 0; i < 3; i = i + 1) {
    print i;
}
```

### Functions
```lox
fun greet(name) {
    return "Hello, " + name + "!";
}

fun fibonacci(n) {
    if (n <= 1) return n;
    return fibonacci(n - 2) + fibonacci(n - 1);
}

print greet("Alice");
print fibonacci(10);
```

### Classes and Objects
```lox
class Person {
    init(name, age) {
        this.name = name;
        this.age = age;
    }
    
    sayHello() {
        print "Hi, I'm " + this.name;
    }
    
    getAge() {
        return this.age;
    }
}

var person = Person("Bob", 25);
person.sayHello();
print person.getAge();
```

### Inheritance and Super
```lox
class Animal {
    init(name) {
        this.name = name;
    }
    
    speak() {
        print this.name + " makes a sound";
    }
}

class Dog < Animal {
    speak() {
        print this.name + " barks";
        super.speak();
    }
}

var dog = Dog("Rex");
dog.speak();
// Output:
// Rex barks
// Rex makes a sound
```

### Closures
```lox
fun makeCounter() {
    var count = 0;
    fun counter() {
        count = count + 1;
        return count;
    }
    return counter;
}

var counter = makeCounter();
print counter(); // 1
print counter(); // 2
print counter(); // 3
```

## Web Interface Features

- **Code Editor**: Syntax-aware text editor for Lox code
- **Live Examples**: Pre-built examples for all language features
- **File Operations**: Upload `.lox` files and download output
- **Real-time Execution**: Run code and see results instantly
- **Error Handling**: Clear error messages for debugging
- **Terminal Styling**: Dark theme with terminal-like appearance

## Implementation Details

The interpreter follows a tree-walking approach with these main components:

- **Scanner**: Tokenizes source code
- **Parser**: Builds Abstract Syntax Tree (AST) 
- **Resolver**: Performs static analysis and variable resolution
- **Interpreter**: Executes the AST with environment-based variable storage
- **Web Server**: Go HTTP server serving the web interface and API

Key features include lexical scoping, first-class functions with closures, and a complete object system with inheritance.

## Development

- **Web version** auto-deploys from `dev` branch via GitHub Actions to Railway
- **Terminal version** available on `master` branch for local development
- **API Endpoint**: `/api/interpret` accepts POST requests with Lox code
- **Static Files**: Web interface served from `/web` directory
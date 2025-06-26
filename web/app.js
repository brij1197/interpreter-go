const examples = {
    hello: `print "Hello, World!";`,
    variables: `var name = "Alice";
var age = 30;
var isStudent = false;

print "Name: " + name;
print "Age: " + age;
print "Student: " + isStudent;`,
    functions: `fun greet(name) {
    return "Hello, " + name + "!";
}

fun fibonacci(n) {
    if (n <= 1) return n;
    return fibonacci(n - 2) + fibonacci(n - 1);
}

print greet("World");
print "Fibonacci(10): " + fibonacci(10);`,
    classes: `class Person {
    init(name, age) {
        this.name = name;
        this.age = age;
    }
    
    greet() {
        print "Hi, I'm " + this.name + " and I'm " + this.age + " years old.";
    }
}

var person = Person("Bob", 25);
person.greet();`,
    inheritance: `class Animal {
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
dog.speak();`
};

async function runCode(){
    const code = document.getElementById('code').value;
    const output = document.getElementById('output');

    output.textContent = 'Running...';

    try{
        const response = await fetch('/api/interpret',{
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({code})
        })

        const result = await response.json();

        if (result.error) {
            output.innerHTML = `<span class="error">Error: ${result.error}</span>`;
        } else{
            output.textContent = result.output || 'Program completed successfully.';
        }
    } catch(error){
        output.innerHTML = `<span class="error">Network error: ${error.message}</span>`;
    }
}

function loadExample(example){
    if (example && examples[example]){
        document.getElementById('code').value = examples[example];
    }
}

function clearEditor() {
    document.getElementById('code').value = '';
    document.getElementById('output').textContent = 'Ready to run Lox code...';
}

function loadFile(event){
    const file = event.target.files[0];
    if (file){
        const reader = new FileReader();
        reader.onload = function(e){
            document.getElementById('code').value = e.target.result;
        };
        reader.readAsText(file);
    }
}

function downloadOutput(){
    const output = document.getElementById('output').textContent;
    const blob = new Blob([output], {type: 'text/plain'});
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = 'lox_output.txt'
    a.click()
    URL.revokeObjectURL(url);
}
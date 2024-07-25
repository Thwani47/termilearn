# Introduction to Go

Golang, also known as Go, is an open source programming language designed for building scalable, secure, and reliable software.

Go is a statically-typed languge, which means that the type of a variable is known at compile time. This allows for better performance and more secure code.

Go has built-in features that support concurrency, which allows for effecient use of multiple processors and the ability to manage multiple tasks at once. It also has a garbage collector that automatically manages memory allocation and deallocation, making it easier to write safe and stable code.

Go is a very versatile programmnig language that can be used to build a wide range of applications, including web applications, network programming, system administrations, and more.

We'll start with a simple Hello World program to understand how Go programs are structured.

A Hello World program in Go is written as follows

```go
package main

import "fmt"

func main(){
    fmt.Println("Hello, World!")
}
```

Let's explain different sections of the code above:

```go
package main
```

Every Go program must start with the package declaration. Go programs are organized in packages. This allows us to reuse our code.

```go
import "fmt"
```

As we mentioned above that Go programs are organized into packages, here we import the `fmt` package. The `import` keyword allows us to use types and functions defined in other packages. The `fmt` allows us to work with formatted Input/Output functions (i.e, allow us to read input and print output from our programs).

```go
func main(){
    fmt.Println("Hello, World!")
}
```

The keyword `func` is used to define a function in Go. Under the `main` package, a function `main` is required. This is where our program starts. Inside the `main` function, we call the `Println` function defined in the `fmt` package. This method allows us to print output to the standard output. We use this method to print the text "Hello World" to the standard output.

# Summary

- [Introduction](#Introduction)
- [IDE and Plug-in](#IDE-and-Plug-in)
- [Code organization](#Code-organization)
  - [Overview](#Overview)
  - [Workspace](#Workspace)
  - [Gopath Environment Variable](#Gopath-Environment-Variable)
  - [Import package](#Import-package)
  - [Syntax](#Syntax)
    - [Declare a constant](#Declare-a-constant)
    - [Declare a variable](#Declare-a-variable)
    - [Basic var types](#Basic-var-types)
    - [Declare a function](#Declare-a-variable)
    - [Declare a method](#Declare-a-method)
    - [_If_ statements](#_If_-statements)
    - [_For_ loop](#_For_-loop)
    - [Range](#Range)
    - [Switch - cases](#Switch---cases)
    - [Defer](#Defer)
    - [Arrays and Slices](#Arrays-and-Slices)
    - [Maps](#Maps)
    - [Pointer receivers](#Pointer-receivers)
    - [Structs](#Structs)
    - [Interfaces](#Interfaces)

# Introduction

The following document will synthezise the best practice in code organization and formatting for Golang

# IDE and Plug-in

if you're using VSCode you should install the Go Linter and Formating plug in _ms-vscode.go_
if you're using Atom you should install _Go-Plus_ that provides enhanced Go support
If you have a JetBrain Licence you can download [JetBrain-GoLand](https://www.jetbrains.com/go/promo/?gclid=EAIaIQobChMIp9zutb7k5QIVBIbVCh1mJgn7EAAYASAAEgI3tvD_BwE)

# Code organization

## Overview

- Go programmers typically keep **all** their Go code in a single workspace.
- A workspace contains many version control repositories (managed by Git, for example).
- Each repository contains one or more packages.
- Each package consists of one or more Go source files in a single directory.
- The path to a package's directory determines its import path.

## Workspace

A workspace is a directory hierarchy with two directories at its root:

- src contains Go source files, and
- bin contains executable commands.

The Go tool builds and installs binaries to the _bin_ directory.

The _src_ subdirectory typically contains multiple version control repositories (such as for Git or Mercurial) that track the development of one or more source packages.

You can see the full exemple in the workspace set up in the chapter [Installation](https://github.com/apps-cirruseo/golang-tour/tree/master/1_installation)

One of the best practice is to keep **all** your Go dependencies in the same workspace, without using symbolic links

## Gopath Environment Variable

The GOPATH environment variable specifies the location of your workspace. It defaults to a directory named go inside your home directory, so \$HOME/go on Unix based systems

If you would like to work in a different location, you will need to set GOPATH to the path to that directory

```bash
go env -w GOPATH=$HOME/[Newpath]
```

For convenience, add the workspace's bin subdirectory to your PATH:

```bash
$ export PATH=$PATH:$(go env GOPATH)/bin
```

The scripts in the rest of this document use $GOPATH instead of $(go env GOPATH) for brevity. To make the scripts run as written if you have not set GOPATH, you can substitute \$HOME/go in those commands or else run:

```bash
$ export GOPATH=$(go env GOPATH)
```

## Import package

An import path is a string that uniquely identifies a package. A package's import path corresponds to its location inside a workspace or in a remote repository.

The packages from the standard library are given short import paths such as "fmt" and "net/http". For your own packages, you must choose your own path

```golang
import ("$HOME/scr/name")
```

The first statement in a Go source file must be

```golang
package name
```

where _name_ is the package's default name for imports ( make the difference between package name and file name )
Executable commands must always use _package main_.

## Syntax

For further information on syntax, visit [Effective Go](https://golang.org/doc/effective_go.html)

### Declare a constant

Go constant are fixed values declared at compile time and can only be numbers, characters (runes), strings or booleans

They are declared as :

```golang
const Pi = 3.14
```

### Declare a variable

Variables can be initialized just like constants but the initializer can be a general expression computed at run time.

```golang
var (
    home   = os.Getenv("HOME")
    user   = os.Getenv("USER")
    gopath = os.Getenv("GOPATH")
)
```

you can declare several variables at the same time
Or you can even _short declare_ one or several variables :
(this short declaration only works **inside** of a function)

```golang
k := 3 // Println(k) will return 3
lang, test, size := "golang", true, 10 // Println (lang, test, size) will return "golang", true, 10
```

#### Basic var types

You can have various variable types :

```golang
bool
string
int (8,16,32,64)
uint (8,16,32,64)
float32
float64
complex64
complex128
byte // alias for uint8
rune // alias for int32
     // represents a Unicode code point
```

Variables that are not initalized are given the default Zero value (0, false, "")

### Declare a function

A function can take zero or more arguments and is declared like this:

```golang
func add(x int, y int) int {
	return x + y
}
```

In this example, _add_ takes two parameters of type _int_.
You can also write _x, y int_
**Notice that the type comes after the variable name.**

The last _int_ variable declared is the type expected tu be returned by the function (addition of 2 integers expects to return a third integer as a result)

### Declare a method

Go does not have classes. However, you can define methods on types.
A method is a function with a special receiver argument.
The receiver appears in its own argument list between the func keyword and the method name.

```golang
type Vertex struct {
	X, Y float64
}

func (v Vertex) Abs() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}
```

### _If_ statements

in _if_ statements, the expression need not be surrounded by parentheses ( ) but the braces { } are required.

```golang
func sqrt(x float64) string {
	if x < 0 {
		return sqrt(-x) + "i"
	}
	return fmt.Sprint(math.Sqrt(x))
}
```

If statements also exist in short versions :

```golang
if v := math.Pow(x, n); v < lim
```

Variables declared inside an _if_ short statement are also available inside any of the _else_ blocks.

### _For_ loop

Go has only one looping construct, the _for_ loop.
The basic for loop has three components separated by semicolons:

- the init statement: executed before the first iteration
- the condition expression: evaluated before every iteration
- the post statement: executed at the end of every iteration

Like for _if_ statements, the expression need not be surrounded by parentheses ( ) but the braces { } are required.

You can declare your _for_ loop such as :

```golang
func main() {
	sum := 0
	for i := 0; i < 10; i++ {
		sum += i
	}
	fmt.Println(sum)
}
```

_for_ is also used in Go as _while_ loops
You just have to drop the semicolons :

```golang
func main() {
	sum := 1
	for sum < 1000 {
		sum += sum
	}
	fmt.Println(sum)
}
```

### Range

Another variation of the _for_ loop if the _range_ function
You can declare a range of values, in which your loop is going to do its work

```golang
var pow = []int{1, 2, 4, 8, 16, 32, 64, 128}

func main() {
	for i, v := range pow {
		fmt.Printf("2**%d = %d\n", i, v)
	}
}
```

In this sample, the _pow_ variable is an array containing the first power values of 2 (2<sup>0</sup>, 2<sup>1</sup>, etc.)
The _for_ loop will work using the values in this array, and the variable _i_, index of the _for_ loop, to print the values in the form of : 2<sup>i</sup> = int[v]
_%d_ is a placeholder for the values _i_ or _v_ declared as Printf arguments

### Switch - cases

A switch statement is a shorter way to write a sequence of if - else statements. It runs the first case whose value is equal to the condition expression.
Go's _switch_ only run the matching case, not the ones that follows, so the _break_ statement is not needed in Go
Also, Go's cases, can be constants, variables of any types

```golang
func main() {
	fmt.Print("Go runs on ")
	switch os := runtime.GOOS; os {
	case "darwin":
		fmt.Println("OS X.")
	case "linux":
		fmt.Println("Linux.")
	default:
		// freebsd, openbsd,
		// plan9, windows...
		fmt.Printf("%s.\n", os)
    }
```

### Defer

in Go you can use the _defer_ function , which delays the execution of a function until the surrounding functions (before or after) returns values

```golang
func main() {
	defer fmt.Println("world")

	fmt.Println("hello")
}
// this function will wait for print "hello", before printing "world"
```

Deferred function calls are pushed onto a stack. When a function returns, its deferred calls are executed in last-in-first-out order.

### Arrays and Slices

In Go you can use arrays declared as array[n]T is an array of n values of type T.
The first value of an array is array[0]

```golang
var a [2]string
//declares the variable a as an array of 2 strings
```

An array's length is part of its type, so arrays cannot be resized.
Go provides a convenient way of working with arrays : Slices
A slice is a dynamically-sized, flexible view into the elements of an array, so slices are more commonly used

A slice is formed by specifying two indices, a low and high bound, separated by a colon.
This selects a half-open range which includes the first element, but excludes the last one.

```golang
primes := [6]int{2, 3, 5, 7, 11, 13}
var s []int = primes[1:4]
// the array "s" is a slice of the array "primes" which contains values from range 1 to 3
```

Slices are like references to arrays, a slice does not store any data, it just describes a section of an underlying array.

### Maps

A _map_ maps a key to a value.
The _make_ function returns a map of the given type, initialized and ready for use.
Maps are Go’s built-in associative data type, using a Key and a Value (sometimes called hashes or dicts in other languages).

```golang
func main() {

    m := make(map[string]int)

    m["k1"] = 7
    m["k2"] = 13

    fmt.Println("map:", m)

    v1 := m["k1"]
    fmt.Println("v1: ", v1)

    fmt.Println("len:", len(m))

    delete(m, "k2")
    fmt.Println("map:", m)

    _, prs := m["k2"]
    fmt.Println("prs:", prs)

    n := map[string]int{"foo": 1, "bar": 2}
    fmt.Println("map:", n)
}
```

### Pointer receivers

A pointer is used to link a value to its memory adress

```golang
func main() {
	a := 5
	b := &a

	fmt.println (a, b)
}
```

A Pointer receiver is a method that uses a pointer to modify the value of another method

```golang
type Vertex struct {
	X, Y float64
}

func (v Vertex) Abs() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func (v *Vertex) Scale(f float64) {
	v.X = v.X * f
	v.Y = v.Y * f

}

func main() {
	v := Vertex{3, 4} //this line takes both values X and Y, multiplies them by themselves,(3*3 and 4*4) adds up the results (9+16) and applies square root to the result (square root of 25)
	v.Scale(10) //this line multiplies Values X and Y by another factor f (3*10 and 4*10)
	fmt.Println(v.Abs())
}
```

### Structs

Go’s _structs_ are typed collections of fields. They’re useful for grouping data together to form records.

```golang

type person struct {
    name string
    age  int
}

func NewPerson(name string) *person {

    p := person{name: name}
    p.age = 42
    return &p
}

func main() {

    fmt.Println(person{name: "Alice", age: 30})

    fmt.Println(NewPerson("Jon"))

    s := person{name: "Sean", age: 50}
    fmt.Println(s.name)

    sp := &s
    fmt.Println(sp.age)

    sp.age = 51
    fmt.Println(sp.age)
}
```

### Interfaces

An interface type is defined as a set of method signatures.
A value of interface type can hold any value that implements those methods.

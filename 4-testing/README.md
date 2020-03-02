# Testing

Let's write a library, use it for the hello program and test the results
We will re-use the same program " Hello World " we made earlier

## Library Creation

Again, the first step is to choose a package path (we'll use src/stringutil) and create the package directory:

```bash
$ mkdir -p $GOPATH/src/stringutil
```

Next, create a file named _reverse.go_ in that directory with the following contents.

```golang
// Package stringutil contains utility functions for working with strings.
package stringutil

// Reverse returns its argument string reversed rune-wise left to right.
func Reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}
```

Now, test that the package compiles with go build:

```bash
$ go build $GOPATH/src/stringutil
```

This won't produce an output file. Instead it saves the compiled package in the local build cache .

## Import library

Next, update _hello.go_ with the following sample :

```golang
package main

import (
	"fmt"
	"stringutil"
)

func main() {
	fmt.Println(reverse("Golang"))
}
func reverse(chaine string)string{
	return stringutil.Reverse(chaine)
}
```

Install the hello program:

```bash
$ go install $GOPATH/src/hello
```

Running the new version of the program, you should see a new, reversed message:

```bash
$ hello
gnaloG
```

( don't forget to re-install your program _hello.go_ everytime you make a change )

## Testing you app

Go has a lightweight test framework composed of the go test command and the testing package.

Create a file named reverse_test.go in src/stringutil that contains functions named TestXXX with signature func (t \*testing.T). (\_test.go extension is mandatory for the file to be considered by the test framework as a test file)

The test framework runs each such function; if the function calls a failure function such as t.Error or t.Fail, the test is considered to have failed.

Add the following sample to the test file

```Golang
package stringutil

import "testing"

func TestReverse(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"Golang", "gnaloG"}, // This line tests a simple word
		{"Hello, !", "! ,olleH"}, // This line tests a word and a special character
		{"", ""}, // This line test a blank string
 	}
	for _, c := range cases {
		got := Reverse(c.in)
		if got != c.want {
			t.Errorf("Reverse(%q) == %q, want %q", c.in, got, c.want)
		}
	}
}
```

Then run the test with go test:

```bash
$ go test $GOPATH/src/stringutil
ok stringutil
```

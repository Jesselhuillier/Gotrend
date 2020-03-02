# Your first program

To compile and run a simple program, first choose a package path and create a corresponding package directory inside your workspace:

```bash
$ mkdir - p $GOPATH/src/hello
```

Next, create a file named hello.go inside that directory, containing the following Go code.

```golang
package main

import "fmt"

func main() {
	fmt.Println("Hello, world.")
}
```

Now you can build and install that program with the go tool:

```bash
$ go install hello
```

Note that you can run this command from anywhere on your system. The go tool finds the source code by looking for the _hello_ package inside the workspace specified by GOPATH.

You can now run the test app by using :

```bash
$ $GOPATH/bin/hello
Hello, world.
```

or

```bash
$ hello
Hello, world.
```

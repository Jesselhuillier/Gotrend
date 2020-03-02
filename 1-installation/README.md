# Download

Download GO packages at [Official Website](https://golang.org/dl/)
Choose mac os version
Download the package file, open it, and follow the prompts to install the Go tools.
The package installs the Go distribution to /usr/local/go.

The package should put the /usr/local/go/bin directory in your PATH environment variable. You may need to restart any open Terminal sessions for the change to take effect.

# Test installation

Check that Go is installed correctly by setting up a workspace and building a simple program, as follows.

Create your workspace directory, \$HOME/go. (If you'd like to use a different directory, you will need to set the GOPATH environment variable.)

Next, make the directory src/{githubName}/{projectName}/hello inside your workspace, and in that directory create a file named hello.go that looks like:

```Golang
package main

import "fmt"

func main() {
	fmt.Printf("hello, world\n")
}
```

Then build it with the go tool:

```bash
$ cd $HOME/go/src/hello
$ go build
```

The command above will build an executable named hello in the directory alongside your source code. Execute it to see the greeting:

```bash
$ ./hello
hello, world
```

If you see the "hello, world" message then your Go installation is working.

You can run go install to install the binary into your workspace's bin directory or go clean -i to remove it.

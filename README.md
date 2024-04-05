# gyparse

YAML parser made in Go.

## What is it?

gyparse is a simple YAML parser made in Go. It is a simple and easy to use library that can be used to parse YAML files into objects that you can use in your golang projects.

## How to install it?

You can install gyparse using the `go get` command:

```bash
go get github.com/tomassantunes/gyparse@latest
```

## How to use it?

To use gyparse, you need to import the library into your project and then you can use the `Parse` function to parse a YAML file into an object. Here is an example:

```go
package main

import (
    "os"
    "fmt"

	"github.com/tomassantunes/gyparse"
)

func main() {
    input, err := os.ReadFile("./example.yml")
	if err != nil {
		fmt.Println(err)
	}

	obj, err := Parse(string(input))
	if err != nil {
	    fmt.Println(err)
    }

    fmt.Println(obj)
}
```

## Next steps

-   Add support for more complex YAML files

# jq Go Binding

This is a Go binding for the `jq` JSON processor. It allows you to execute `jq` commands from your Go application.

## Features

*   Easy integration of `jq` functionality into your Go code.
*   Supports command-line arguments and file input.
*   Uses `jq`'s filtering capabilities for powerful JSON manipulation.

## Installation

```bash
go get github.com/thekhanj/jq
```

## Usage

```go
package main

import (
	"fmt"
	"log"
	"github.com/thekhanj/jq"
)

func main() {
	jq, err := jq.NewJq(
		jq.WithFilterString(".[0].name"), // Example filter: extract the 'name' field from the first element
		jq.WithFile(bytes.NewReader([]byte(`{"name": "John", "age": 30}`))), // Example file input
	)
	if err != nil {
		log.Fatal(err)
	}

	result, err := jq.Start()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(result))
}
```

## API

### `Jq` Struct

*   `filter io.Reader`: The filter expression to be passed to `jq`.
*   `filePaths []string`: List of file paths to process.
*   `files []io.Reader`: List of readers to process.
*   `options []interface{}`: List of command-line options.
*   `ranOnce bool`: Flag to prevent multiple calls to `Start()`.
*   `tempDir string`: Temporary directory used for creating FIFO files.

### `JqOption` Type

A function type that takes a `*Jq` pointer and allows you to configure the `Jq` instance.

### `NewJq(opts ...JqOption) (*Jq, error)`

Creates a new `Jq` instance with the provided options.

### `Start() ([]byte, error)`

Executes the `jq` command and returns the output as a byte array and an error.

### `JqOption` Functions

Convenience functions for setting various `jq` options:

*   `WithFilterString(filter string) JqOption`: Sets the filter expression.
*   `WithFile(file io.Reader) JqOption`:  Sets an input file.
*   `WithFilePath(path string) JqOption`: Sets a file path.
*   `WithFileData(data []byte) JqOption`: Sets data from a byte array.
*   `WithTempDir(dir string) JqOption`: Sets the temporary directory.

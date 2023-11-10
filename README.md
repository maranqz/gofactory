# Factory linter

[![CI](https://github.com/maranqz/go-factory-lint/actions/workflows/ci.yml/badge.svg)](https://github.com/maranqz/go-factory-lint/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/maranqz/go-factory-lint)](https://goreportcard.com/report/github.com/maranqz/go-factory-lint?dummy=unused)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)

The linter checks that the Structes are created by the Factory, and not directly.
The checking helps to provide invariants without exclusion.
Validation helps you avoid losing the creation of invalid object.


## Use

Installation

    go install github.com/maranqz/go-factory-lint/cmd/go-factory-lint@latest

### Options

`-b`, `--blockedPgs` - list of packages, where the structures should be created by factories. By default, all structures in all packages should be created by factories. 

## Example

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
package main

import (
	"fmt"

	"bad"
)

func main() {
	// Use factory for bad.User
	u := &bad.User{
		ID: -1,
	}

	fmt.Println(u.ID) // -1
	fmt.Println(u.CreatedAt) // time.Time{}
}

```

```go
package bad

import "time"

type User struct {
	ID        int64
	CreatedAt time.Time
}

var sequenceID = int64(0)

func NextID() int64 {
	sequenceID++

	return sequenceID
}


```

</td><td>

```go
package main

import (
	"fmt"

	"good"
)

func main() {
	u := good.NewUser()
	
	fmt.Println(u.ID)        // auto increment
	fmt.Println(u.CreatedAt) // time.Now()
}

```

```go
package user

import "time"

type User struct {
	ID        int64
	CreatedAt time.Time
}

func NewUser() *User {
	return &User{
		ID: nextID(),
		CreatedAt: time.Now(),
	}
}

var sequenceID = int64(0)

func nextID() int64 {
	sequenceID++

	return sequenceID
}

```

</td></tr>
</tbody></table>
# Logging Error Context

lec is a Go library that provides a context management system with logging and error handling capabilities.

## Installation

```zsh
go get github.com/gosuit/lec
```

## Features

• Logging using the sl.Logger.

• Storage of key-value pairs with ability to log values.

• Error management with the ability to add and retrieve errors.

• Implementation of context.Context

## Usage

### Logging

```golang
package main

import (
    "github.com/gosuit/sl"
    "github.com/gosuit/lec"
)

func main() {
    logCfg := // Init *sl.Config

    log := sl.New(logCfg)
    ctx := lec.New(log)

    ctx.AddValue("key1", "value1", false)
    ctx.AddValue("key2", "value2", true)

    some(ctx)
}

func some(ctx lec.Context) {
    log := ctx.Logger()

    log.Info("msg") // result has attribute key2=value2
}
```

### Base context

```golang
package main

import (
    "context"

    "github.com/gosuit/sl"
    "github.com/gosuit/lec"
)

func main() {
    logCfg := // Init *sl.Config

    log := sl.New(logCfg))
    baseCtx := context.Background()

    ctx := lec.NewWithCtx(baseCtx, logg)

    //You can access standard context methods:
    deadline, ok := ctx.Deadline()
    doneChan := ctx.Done()
    err := ctx.Err()
    value := ctx.Value("someKey")
}
```

### Adding Values

```golang
package main

import (
    "github.com/gosuit/sl"
    "github.com/gosuit/lec"
)

func main() {
    logCfg := // Init *sl.Config

    log := sl.New(logCfg)
    ctx := lec.New(log)

    ctx.AddValue("key1", "value1", true) // 'true' indicates the value should be logged

    value := ctx.GetValue("key1")
    if value != nil {
        fmt.Println(value.Val) // Output: value1
    }

    // To get all stored values:

    allValues := ctx.GetValues()
}
```

### Error Handling

```golang
package main

import (
    "github.com/gosuit/sl"
    "github.com/gosuit/lec"
)

func main() {
    logCfg := // Init *sl.Config

    log := sl.New(logCfg)
    ctx := lec.New(log)

    some(ctx)

    if ctx.HasErr() {
        fmt.Println(ctx.GetErr()) // Output: error msg
    }
}

func some(ctx lec.Context) {
    err := doSomething()
    if err != nil {
        ctx.AddErr(err)
    }
}

func doSomething() error {
    return error.New("error msg")
}
```

## Contributing

Contributions are welcome! Please feel free to submit a pull request or open an issue for any enhancements or bug fixes.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
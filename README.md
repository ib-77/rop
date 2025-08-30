# Railway Oriented Programming (ROP)

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Current version: v0.5.3

## Overview

Railway Oriented Programming is a functional programming pattern for handling errors in a clean and composable way. This Go library implements ROP principles using generics, allowing you to build robust error handling pipelines.

## Core Concepts

Railway Oriented Programming visualizes program flow as a railway track:

- **Success Track**: When operations succeed, they continue along the main track
- **Failure Track**: When operations fail, they switch to a parallel error track
- **Composition**: Functions can be chained together, with errors automatically propagated

## Features

- **Type-safe error handling** using Go generics
- **Composable operations** that can be chained together
- **Three result states**: Success, Fail, and Cancel
- **Context support** for cancellation and timeouts
- **Retry mechanisms** with various strategies (fixed, linear, exponential)
- **Parallel processing** with fan-out/fan-in patterns
- **Comprehensive testing** for all components

## Core Components

- **Result[T]**: The central type representing either success or failure
- **Solo**: Operations on single values
- **Mass**: Operations on streams of values
- **Bridge**: Operations across multiple streams
- **Fan**: Fan-out/fan-in operations for parallel processing
- **Group**: Grouping operations for complex workflows

## Basic Usage

```go
// Validate a value
result := solo.Validate(value, func(a int) bool {
    return a < 10
}, "Value must be less than 10")

// Chain operations
result = solo.AndThen(result, func(v int) rop.Result[string] {
    return rop.Success(strconv.Itoa(v))
})

// Handle the result
if result.IsSuccess() {
    fmt.Println("Success:", result.Result())
} else {
    fmt.Println("Error:", result.Err())
}
```

## Advanced Features

### Retry Operations

```go
// Create a retry strategy
strategy := rop.NewExponentialRetryStrategy(5, 2.0, time.Second, nil)
ctx := rop.WithRetry(context.Background(), strategy)

// Use with retry-aware functions
result := solo.RetryableOperation(ctx, input)
```

### Parallel Processing

```go
// Process items in parallel with fan-out/fan-in
results := fan.Process(ctx, inputs, processor, 10)
```

## Installation

```
go get github.com/ib-77/rop
```

## License

MIT License - see [LICENSE.md](LICENSE.md) for details.

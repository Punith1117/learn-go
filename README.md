# Learn Go

A practical Go learning repository focused on the concepts needed to build a concurrent backend.

## Packages & Modules

The `packages-modules` directory contains the example application
for learning how Go packages and modules work.

```text
packages-modules/
├── event/
│   └── event.go
└── main.go
```

### Compile

From the root of the repository, run:

```bash
go build -o packages-modules-app ./packages-modules
```

This compiles the `packages-modules` package and
creates an executable named `packages-modules-app` in the root directory.

### Run

After compiling, run the executable:

```bash
./packages-modules-app
```

## Concurrency Quick Reference

| Concept          | Purpose                                |
| ---------------- | -------------------------------------- |
| Goroutine        | Run work concurrently                  |
| Channel          | Communicate between goroutines         |
| `sync.WaitGroup` | Wait for goroutines to finish          |
| `sync.Mutex`     | Protect shared mutable state           |
| Semaphore        | Limit concurrent operations            |
| Race detector    | Detect unsafe concurrent memory access |

## Testing

Test files must end with _test.go: `inventory_test.go`

Test functions start with Test and receive *testing.T:

```
func TestSomething(t *testing.T) {
    // test
}
```

Run tests:

```bash
go test .
```

Run all tests:

```bash
go test ./...
```

Run tests with race detection:

```bash
go test -race .
```

Force a fresh test run when needed:

```bash
go test -race -count=1 .
```

`-count=1` disables the test cache for that run.
It is mainly useful when debugging concurrent or flaky behavior.

# Examples

Examples demonstrating the public API of `erax`.

```
examples/
├── alien
├── json
├── meta
├── new
├── style
└── wrap
```

## [new](examples/new/main.go)

Creating root errors.

Functions:

- `erax.New`
- `erax.Format`

Run:

```bash
go run ./examples/new/main.go
```

---

## [wrap](examples/wrap/main.go)

Wrapping errors and building error trees.

Functions:

- `erax.Wrap`
- `erax.WrapWithErrors`
- `erax.Format`

Run:

```bash
go run ./examples/wrap/main.go
```

---

## [meta](examples/meta/main.go)

Working with structured metadata.

Functions:

- `erax.WithMeta`
- `erax.AddMeta`
- `erax.GetMeta`
- `erax.F`
- `erax.Format`

Run:

```bash
go run ./examples/meta/main.go
```

---

## [json](examples/json/main.go)

Serializing and restoring errors.

Functions:

- `erax.FormatToJSONString`
- `erax.FormatToJSONMap`
- `erax.FromJSONMap`

Run:

```bash
go run ./examples/json/main.go
```

---

## [alien](examples/alien/main.go)

Interoperability with non-erax errors and the Go standard library.

Functions:

- `erax.Wrap`
- `erax.Cast`
- `errors.Is`
- `errors.As`
- `erax.Format`

Also demonstrates compatibility with:

- `fmt.Errorf`
- `github.com/pkg/errors`

Run:

```bash
go run ./examples/alien/main.go
```

---

## [style](examples/style/main.go)

Customizing terminal output.

Functions:

- `erax.SetBranchColor`
- `erax.SetErrorColor`
- `erax.SetKeyColor`
- `erax.SetValueColor`

Run:

```bash
go run ./examples/style/main.go
```

# Benchmarks

- **Create**: ~3.4x faster, 51% less memory, 39% fewer allocations
- **Serialize**: ~5x faster, 76% less memory, allocations reduced from 42 to 1
- **Create and Serialize**: ~3.6x faster, 64% less memory, 71% fewer allocations

|                      | Before v0.4.0                           | After v0.4.0                            |
|----------------------|-----------------------------------------|-----------------------------------------|
| Create               | `2002 ns/op`	`2136 B/op` `36 allocs/op` | `584 ns/op` `1040 B/op` `22 allocs/op`  |
| Serialize            | `3959 ns/op` `2377 B/op` `42 allocs/op` | `796.7 ns/op` `576 B/op` `1 allocs/op`  |
| Create and Serialize | `5483 ns/op` `4514 B/op` `78 allocs/op` | `1525 ns/op` `1617 B/op` `23 allocs/op` |

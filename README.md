<p align="center">
  <img src="/docs/logo.png" alt="heapcraft logo" width="75%"/>
</p>

**heapcraft** is a Go library offering a suite of heap data structures which include binary, d‑ary, pairing, radix, skew, and leftist heaps.

Available heap types include:

**D-ary Heaps:**
- `DaryHeap` / `SyncDaryHeap`

**Radix Heaps:**
- `RadixHeap` / `SyncRadixHeap`

**Tree-Based Heaps:**
- `SimplePairingHeap` / `SyncSimplePairingHeap`
- `PairingHeap` / `SyncPairingHeap`
- `SimpleSkewHeap` / `SyncSimpleSkewHeap`
- `SkewHeap` / `SyncSkewHeap`
- `SimpleLeftistHeap` / `SyncSimpleLeftistHeap`
- `LeftistHeap` / `SyncLeftistHeap`

---

## ✨ **Features**

| Category            | Details                                                                                    |
| ------------------- | ------------------------------------------------------------------------------------------ |
| **Heap Variants**   | `Binary`, `D‑ary`, `Pairing`, `Radix`, `Skew`, `Leftist`                                  |
| **Implementation Types** | **Simple/Full** for `Pairing`, `Skew`, and `Leftist` heaps; **Single** for `D‑ary`, and `Radix` heaps |
| **Thread Safety**   | Both non-thread-safe and thread-safe versions available (e.g., `DaryHeap` and `SyncDaryHeap`) |
| **Generics**        | Go 1.18+ type parameters—store any custom type                              |
| **Node Tracking**   | Full implementations maintain a map for O(1) lookup and update operations                |
| **Memory Pooling**  | Optional object pooling for improved performance                 |
| **Examples**        | Examples for each heap type in the `examples/` directory                   |

---

## 🚀 **Getting Started**

### Install

```bash
go get github.com/galactixx/heapcraft@latest
```

Then import it in your code:

```go
import "github.com/galactixx/heapcraft"
```

## 🔍 **API**

### Implementation Types

**D-ary Heaps** (`DaryHeap` / `SyncDaryHeap`) provide array-based heap operations:
- `Push(value, priority)` - Add elements
- `Pop()` / `PopValue()` / `PopPriority()` - Remove elements
- `Peek()` / `PeekValue()` / `PeekPriority()` - View without removing
- `Length()`, `IsEmpty()`, `Clear()`, `Clone()`
- `Update(index, value, priority)` - Update element at index
- `Remove(index)` - Remove element at index
- `PopPush(value, priority)` - Pop and push in one operation
- `PushPop(value, priority)` - Push and pop in one operation
- `Register(fn)` / `Deregister(id)` - Callback registration for swaps

**Radix Heaps** (`RadixHeap` / `SyncRadixHeap`) provide monotonic priority queue operations:
- `Push(value, priority)` - Add elements (must be >= last popped priority)
- `Pop()` / `PopValue()` / `PopPriority()` - Remove elements
- `Peek()` / `PeekValue()` / `PeekPriority()` - View without removing
- `Length()`, `IsEmpty()`, `Clear()`, `Clone()`
- `Rebalance()` - Manually trigger bucket rebalancing
- `Merge(other)` - Merge with another radix heap

**Simple Tree-Based Heaps** (`SimplePairingHeap` / `SyncSimplePairingHeap`, `SimpleSkewHeap` / `SyncSimpleSkewHeap`, `SimpleLeftistHeap` / `SyncSimpleLeftistHeap`) provide:
- `Push(value, priority)` - Add elements
- `Pop()` / `PopValue()` / `PopPriority()` - Remove elements
- `Peek()` / `PeekValue()` / `PeekPriority()` - View without removing
- `Length()`, `IsEmpty()`, `Clear()`, `Clone()`

**Full Tree-Based Heaps** (`PairingHeap` / `SyncPairingHeap`, `SkewHeap` / `SyncSkewHeap`, `LeftistHeap` / `SyncLeftistHeap`) extend simple heaps with node tracking:
- All simple heap operations
- `Push()` returns a unique node ID
- `UpdateValue(id, newValue)` - Update node value
- `UpdatePriority(id, newPriority)` - Update node priority
- `Get(id)`, `GetValue(id)`, `GetPriority(id)` - Retrieve by ID

## 📚 **Usage**

### Non-Thread-Safe vs Thread-Safe

Each heap type has both non-thread-safe and thread-safe versions:

```go
// Non-thread-safe version (faster, single-threaded use)
heap := heapcraft.NewDaryHeap[int](4, nil, func(a, b int) bool { 
    return a < b 
}, false)

// Thread-safe version (slower, concurrent use)
syncHeap := heapcraft.NewSyncDaryHeap[int](4, nil, func(a, b int) bool { 
    return a < b 
}, false)
```

### D-ary Heaps

```go
// Binary heap (2-ary) or D-ary heap with custom arity
heap := heapcraft.NewDaryHeap[int](4, nil, func(a, b int) bool { 
    return a < b 
}, false)

// Basic and advanced operations
heap.Push(1, 1)
heap.Update(0, 100, 10)
heap.PopPush(42, 5) 
value, _ := heap.PopValue()
```

### Radix Heaps

```go
// Radix heap for integer priorities
heap := heapcraft.NewRadixHeap[int, uint](nil, false)

// Operations (priorities must be >= last popped)
heap.Push(1, 1)
heap.Push(2, 2)
value, _ := heap.PopValue()
heap.Rebalance()
```

### Simple Tree-Based Heaps

```go
// Simple heap (Pairing, Skew, or Leftist)
heap := heapcraft.NewSimplePairingHeap[int](nil, func(a, b int) bool { 
    return a < b 
}, false)

// Basic operations and merging
heap.Push(1, 1)
heap.Push(2, 2)
value, _ := heap.PopValue()
heap.MergeWith(otherHeap)
```

### Full Tree-Based Heaps

```go
// Full heap with node tracking
heap := heapcraft.NewPairingHeap[int](nil, func(a, b int) bool { 
    return a < b 
}, false)

// Node tracking operations
id := heap.Push(42, 10)
heap.UpdateValue(id, 100)
heap.UpdatePriority(id, 1)
value, _ := heap.GetValue(id)
heap.Remove(id)
```

### Memory Pooling

Enable object pooling for better performance:

```go
heap := heapcraft.NewDaryHeap[int](2, nil, func(a, b int) bool { 
    return a < b 
}, true)
```

### Thread Safety

Use thread-safe versions for concurrent access:

```go
// Thread-safe heap for concurrent use
syncHeap := heapcraft.NewSyncDaryHeap[int](nil, func(a, b int) bool { 
    return a < b 
}, false)

// Multiple goroutines can safely call these methods
go func() {
    syncHeap.Push(1, 1)
}()

go func() {
    value, err := syncHeap.PopValue()
}()
```

## 📈 **Performance Benchmarks**

### Environment

| Parameter | Value |
|-----------|-------|
| Date | 2025-06-13 |
| OS | Windows 10  |
| Architecture | amd64 |
| CPU | AMD EPYC 7763 64-Core Processor |
| Go version | 1.24 |

### Micro-benchmark Highlights <sub><sup>*(pooling was not used in running these benchmarks to show raw timing)*</sup></sub>

#### D-ary and Radix Heaps

<div align="center">

| Heap Type | Insertion | | Deletion | | PushPop | | PopPush | |
|-----------|-----------|-----------|----------|----------|----------|----------|----------|----------|
| | Iterations | Time (ns/op) | Iterations | Time (ns/op) | Iterations | Time (ns/op) | Iterations | Time (ns/op) |
| **BinaryHeap** | 19,106,028 | 53.51 | 3,513,837 | 377.3 | 4,237,483 | 394.2 | 3,552,798 | 385.9 |
| **DaryHeap (d=3)** | 34,123,869 | 40.33 | 4,633,384 | 280.1 | 4,165,548 | 375.3 | 4,126,956 | 379.3 |
| **DaryHeap (d=4)** | 37,802,299 | 26.65 | 5,134,050 | 256.0 | 4,271,446 | 392.8 | 4,049,797 | 391.4 |
| **RadixHeap** | 26,960,899 | 42.64 | 2,183,877 | 553.2 | - | - | - | - |

</div>

#### Tree-based Heaps

<div align="center">

| Heap Type | Insertion | | Deletion | | PushPop | | PopPush | |
|-----------|-----------|-----------|----------|----------|----------|----------|----------|----------|
| | Iterations | Time (ns/op) | Iterations | Time (ns/op) | Iterations | Time (ns/op) | Iterations | Time (ns/op) |
| **LeftistHeap** | 1,523,763 | 735.2 | 1,441,719 | 895.7 | - | - | - | - |
| **SimpleLeftistHeap** | 9,759,120 | 119.5 | 2,294,244 | 656.5 | - | - | - | - |
| **PairingHeap** | 1,774,028 | 616.0 | 4,655,505 | 339.3 | - | - | - | - |
| **SimplePairingHeap** | 23,867,677 | 45.23 | 12,821,868 | 124.3 | - | - | - | - |
| **SkewHeap** | 1,000,000 | 1252 | 1,773,727 | 817.2 | - | - | - | - |
| **SimpleSkewHeap** | 4,878,519 | 404.3 | 2,744,472 | 515.4 | - | - | - | - |

</div>

### Full Micro Benchmarks

<div align="center">

```bash
BenchmarkBinaryHeapInsertion-4          	19106028	        53.51 ns/op	      84 B/op	       0 allocs/op
BenchmarkBinaryHeapDeletion-4           	 3513837	       377.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkBinaryPushPop-4                	 4237483	       394.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkBinaryPopPush-4                	 3552798	       385.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkDaryHeap3Insertion-4           	34123869	        40.33 ns/op	      92 B/op	       0 allocs/op
BenchmarkDaryHeap3Deletion-4            	 4633384	       280.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkDaryHeap3PushPop-4             	 4165548	       375.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkDaryHeap3PopPush-4             	 4126956	       379.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkDaryHeap4Insertion-4           	37802299	        26.65 ns/op	      83 B/op	       0 allocs/op
BenchmarkDaryHeap4Deletion-4            	 5134050	       256.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkDaryHeap4PushPop-4             	 4271446	       392.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkDaryHeap4PopPush-4             	 4049797	       391.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkLeftistHeapInsertion-4         	 1523763	       735.2 ns/op	     201 B/op	       3 allocs/op
BenchmarkLeftistHeapDeletion-4          	 1441719	       895.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkSimpleLeftistHeapInsertion-4   	 9759120	       119.5 ns/op	      48 B/op	       1 allocs/op
BenchmarkSimpleLeftistHeapDeletion-4    	 2294244	       656.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkPairingHeapInsertion-4         	 1774028	       616.0 ns/op	     200 B/op	       3 allocs/op
BenchmarkPairingHeapDeletion-4          	 4655505	       339.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkSimplePairingHeapInsertion-4   	23867677	        45.23 ns/op	      32 B/op	       1 allocs/op
BenchmarkSimplePairingHeapDeletion-4    	12821868	       124.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkRadixHeapInsertion-4           	26960899	        42.64 ns/op	      70 B/op	       0 allocs/op
BenchmarkRadixHeapDeletion-4            	 2183877	       553.2 ns/op	     477 B/op	       3 allocs/op
BenchmarkSkewHeapInsertion-4            	 1000000	      1252 ns/op	     239 B/op	       3 allocs/op
BenchmarkSkewHeapDeletion-4             	 1773727	       817.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkSimpleSkewHeapInsertion-4      	 4878519	       404.3 ns/op	      32 B/op	       1 allocs/op
BenchmarkSimpleSkewHeapDeletion-4       	 2744472	       515.4 ns/op	       0 B/op	       0 allocs/op
```

</div>

## 🤝 **License**

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## 📋 **TODO**

- **Interval Heap** - Implementation for double-ended priority queue
- **Weak Heap** - Alternative heap structure with different performance characteristics
- **Enhanced Benchmarking** - Performance testing under contention and with pooling enabled
- **Pooling Benchmarks** - Comparison of performance with and without object pooling
- **Concurrency Testing** - Thread-safe heap performance under various load patterns

---

## 📞 **Contact & Contributing**

Feel free to open an [issue](https://github.com/galactixx/heapcraft/issues) or a pull request. Discussion and feedback are welcome!

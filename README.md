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
- `PairingHeap` / `SyncPairingHeap`
- `FullPairingHeap` / `SyncFullPairingHeap`
- `SkewHeap` / `SyncSkewHeap`
- `FullSkewHeap` / `SyncFullSkewHeap`
- `LeftistHeap` / `SyncLeftistHeap`
- `FullLeftistHeap` / `SyncFullLeftistHeap`

---

## ✨ **Features**

| Category            | Details                                                                                    |
| ------------------- | ------------------------------------------------------------------------------------------ |
| **Heap Variants**   | `Binary`, `D‑ary`, `Pairing`, `Radix`, `Skew`, `Leftist`                                  |
| **Implementation Types** | **Regular/Full** for `Pairing`, `Skew`, and `Leftist` heaps; **Single** for `D‑ary`, and `Radix` heaps |
| **Thread Safety**   | Both non-thread-safe and thread-safe versions available (e.g., `DaryHeap` and `SyncDaryHeap`) |
| **Generics**        | Go 1.18+ type parameters—store any custom type                              |
| **Node Tracking**   | Full implementations maintain a map for O(1) lookup and update operations                |
| **Memory Pooling**  | Optional object pooling                 |
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

**Regular Tree-Based Heaps** (`PairingHeap` / `SyncPairingHeap`, `SkewHeap` / `SyncSkewHeap`, `LeftistHeap` / `SyncLeftistHeap`) provide:
- `Push(value, priority)` - Add elements
- `Pop()` / `PopValue()` / `PopPriority()` - Remove elements
- `Peek()` / `PeekValue()` / `PeekPriority()` - View without removing
- `Length()`, `IsEmpty()`, `Clear()`, `Clone()`

**Full Tree-Based Heaps** (`FullPairingHeap` / `SyncFullPairingHeap`, `FullSkewHeap` / `SyncFullSkewHeap`, `FullLeftistHeap` / `SynFullcLeftistHeap`) extend simple heaps with node tracking:
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

### Regular Tree-Based Heaps

```go
// Regular heap (Pairing, Skew, or Leftist)
heap := heapcraft.NewPairingHeap[int](nil, func(a, b int) bool { 
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
heap := heapcraft.NewFullPairingHeap[int](nil, func(a, b int) bool { 
    return a < b 
}, heapcraft.HeapConfig{UsePool: false})

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
| Date | 2025-06-28 |
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
| **BinaryHeap** | 19,454,936 | 52.71 | 3,510,987 | 374.7 | 3,807,219 | 393.9 | 3,662,504 | 392.8 |
| **DaryHeap (d=3)** | 32,970,384 | 42.15 | 4,658,214 | 284.4 | 3,828,266 | 398.7 | 3,718,894 | 391.9 |
| **DaryHeap (d=4)** | 38,662,906 | 26.42 | 5,091,660 | 259.8 | 3,833,235 | 400.1 | 3,725,020 | 401.0 |
| **RadixHeap** | 26,938,868 | 42.69 | 2,300,319 | 527.4 | - | - | - | - |

</div>

#### Tree-based Heaps

<div align="center">

| Heap Type | Insertion | | Deletion | | PushPop | | PopPush | |
|-----------|-----------|-----------|----------|----------|----------|----------|----------|----------|
| | Iterations | Time (ns/op) | Iterations | Time (ns/op) | Iterations | Time (ns/op) | Iterations | Time (ns/op) |
| **FullLeftistHeap** | 2,314,108 | 596.1 | 1,293,859 | 985.9 | - | - | - | - |
| **LeftistHeap** | 9,481,090 | 130.8 | 1,891,108 | 772.2 | - | - | - | - |
| **FullPairingHeap** | 2,569,038 | 469.5 | 4,440,087 | 338.8 | - | - | - | - |
| **PairingHeap** | 25,685,535 | 44.34 | 11,306,653 | 123.5 | - | - | - | - |
| **FullSkewHeap** | 1,000,000 | 1146 | 1,641,852 | 847.8 | - | - | - | - |
| **SkewHeap** | 5,029,725 | 438.6 | 3,402,890 | 518.3 | - | - | - | - |

</div>

### Full Micro Benchmarks

<div align="center">

```bash
BenchmarkBinaryHeapInsertion-4         	19454936	        52.71 ns/op	      83 B/op	       0 allocs/op
BenchmarkBinaryHeapDeletion-4          	 3510987	       374.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkBinaryPushPop-4               	 3807219	       393.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkBinaryPopPush-4               	 3662504	       392.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkDaryHeap3Insertion-4          	32970384	        42.15 ns/op	      95 B/op	       0 allocs/op
BenchmarkDaryHeap3Deletion-4           	 4658214	       284.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkDaryHeap3PushPop-4            	 3828266	       398.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkDaryHeap3PopPush-4            	 3716894	       391.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkDaryHeap4Insertion-4          	38662906	        26.42 ns/op	      81 B/op	       0 allocs/op
BenchmarkDaryHeap4Deletion-4           	 5091660	       259.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkDaryHeap4PushPop-4            	 3833235	       400.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkDaryHeap4PopPush-4            	 3725020	       401.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkFullLeftistHeap_Insertion-4   	 2314108	       596.1 ns/op	     168 B/op	       2 allocs/op
BenchmarkFullLeftistHeap_Deletion-4    	 1293859	       985.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkLeftistHeap_Insertion-4       	 9481090	       130.8 ns/op	      48 B/op	       1 allocs/op
BenchmarkLeftistHeap_Deletion-4        	 1891108	       772.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkFullPairingHeap_Insertion-4   	 2569038	       469.5 ns/op	     158 B/op	       2 allocs/op
BenchmarkFullPairingHeap_Deletion-4    	 4440087	       338.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkPairingHeap_Insertion-4       	25685535	        44.34 ns/op	      32 B/op	       1 allocs/op
BenchmarkPairingHeap_Deletion-4        	11306653	       123.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkRadixHeapInsertion-4          	26938868	        42.69 ns/op	      70 B/op	       0 allocs/op
BenchmarkRadixHeapDeletion-4           	 2300319	       527.4 ns/op	     475 B/op	       3 allocs/op
BenchmarkFullSkewHeap_Insertion-4      	 1000000	      1146 ns/op	     239 B/op	       3 allocs/op
BenchmarkFullSkewHeap_Deletion-4       	 1641852	       847.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkSkewHeap_Insertion-4          	 5029725	       438.6 ns/op	      32 B/op	       1 allocs/op
BenchmarkSkewHeap_Deletion-4           	 3402890	       518.3 ns/op	       0 B/op	       0 allocs/op
```

</div>

## 🤝 **License**

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## 📋 **TODO**

- Interval Heap
- Weak Heap
- Enhanced Benchmarking
- Pooling Benchmarks
- Concurrency Testing

---

## 📞 **Contact & Contributing**

Feel free to open an [issue](https://github.com/galactixx/heapcraft/issues) or a pull request. Discussion and feedback are welcome!

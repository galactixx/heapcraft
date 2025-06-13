<p align="center">
  <img src="/docs/logo.png" alt="heapcraft logo" width="75%"/>
</p>

**heapcraft** is a high‑performance Go library offering a comprehensive suite of advanced heap data structures—binary heaps, d‑ary heaps, pairing heaps, radix heaps, skew heaps, and leftist heaps—for lightning‑fast priority‑queue operations.

Use it wherever you need efficient scheduling, graph algorithms, event simulation, load balancing, or any task that requires ordered extraction by priority.

### Available Heap Types

<center>

<table width="100%">

| Heap Type | Implementation | Special Features |
|-----------|----------------|------------------|
| **Binary** | `BinaryHeap` | Standard binary heap |
| **D-ary** | `DaryHeap` | Configurable arity (2-ary, 3-ary, etc.) |
| **Radix** | `RadixHeap` | Integer priorities, bucket-based |
| **Pairing** | `SimplePairingHeap` / `PairingHeap` | Constant-time meld, efficient decrease-key |
| **Skew** | `SimpleSkewHeap` / `SkewHeap` | Self-adjusting, amortized O(log n) |
| **Leftist** | `SimpleLeftistHeap` / `LeftistHeap` | Leftist property, efficient merge |

</table>

</center>

---

## ✨ **Features**

| Category            | Details                                                                                    |
| ------------------- | ------------------------------------------------------------------------------------------ |
| **Heap Variants**   | `Binary`, `D‑ary`, `Pairing`, `Radix`, `Skew`, `Leftist`                                  |
| **Implementation Types** | **Simple/Full** for `Pairing`, `Skew`, and `Leftist` heaps; **Single** for `Binary`, `D‑ary`, and `Radix` heaps |
| **Thread Safety**   | All heaps are thread-safe by default using `sync.RWMutex`                                 |
| **Decrease‑Key / Meld** | Native support where algorithmically possible; constant‑time meld on pairing heaps    |
| **Generics**        | Go 1.18+ type parameters—store any comparable or custom type                              |
| **Node Tracking**   | Full implementations maintain a map for O(1) lookup and update operations                |

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

**D-ary Heaps** (`BinaryHeap`, `DaryHeap`) provide array-based heap operations:
- `Push(value, priority)` - Add elements
- `Pop()` / `PopValue()` / `PopPriority()` - Remove elements
- `Peek()` / `PeekValue()` / `PeekPriority()` - View without removing
- `Length()`, `IsEmpty()`, `Clear()`, `Clone()`
- `Update(index, value, priority)` - Update element at index
- `Remove(index)` - Remove element at index
- `PopPush(value, priority)` - Pop and push in one operation
- `PushPop(value, priority)` - Push and pop in one operation
- `Register(fn)` / `Deregister(id)` - Callback registration for swaps

**Radix Heaps** (`RadixHeap`) provide monotonic priority queue operations:
- `Push(value, priority)` - Add elements (must be >= last popped priority)
- `Pop()` / `PopValue()` / `PopPriority()` - Remove elements
- `Peek()` / `PeekValue()` / `PeekPriority()` - View without removing
- `Length()`, `IsEmpty()`, `Clear()`, `Clone()`
- `Rebalance()` - Manually trigger bucket rebalancing
- `Merge(other)` - Merge with another radix heap

**Simple Tree-Based Heaps** (`SimplePairingHeap`, `SimpleSkewHeap`, `SimpleLeftistHeap`) provide:
- `Push(value, priority)` - Add elements
- `Pop()` / `PopValue()` / `PopPriority()` - Remove elements
- `Peek()` / `PeekValue()` / `PeekPriority()` - View without removing
- `Length()`, `IsEmpty()`, `Clear()`, `Clone()`

**Full Tree-Based Heaps** (`PairingHeap`, `SkewHeap`, `LeftistHeap`) extend simple heaps with node tracking:
- All simple heap operations
- `Push()` returns a unique node ID
- `UpdateValue(id, newValue)` - Update node value
- `UpdatePriority(id, newPriority)` - Update node priority
- `Get(id)`, `GetValue(id)`, `GetPriority(id)` - Retrieve by ID

## 📚 **Usage**

### D-ary Heaps

```go
// Binary heap (2-ary) or D-ary heap with custom arity
heap := heapcraft.NewDaryHeap[int](4, nil, func(a, b int) bool { 
    return a < b 
})

// Basic and advanced operations
heap.Push(1, 1)
heap.Update(0, 100, 10)  // Update element at index
heap.PopPush(42, 5)      // Pop and push in one operation
value, _ := heap.PopValue()
```

### Radix Heaps

```go
// Radix heap for integer priorities
heap := heapcraft.NewRadixHeap[int, uint](nil)

// Operations (priorities must be >= last popped)
heap.Push(1, 1)
heap.Push(2, 2)
value, _ := heap.PopValue()
heap.Rebalance()  // Manual bucket rebalancing
```

### Simple Tree-Based Heaps

```go
// Simple heap (Pairing, Skew, or Leftist)
heap := heapcraft.NewSimplePairingHeap[int](nil, func(a, b int) bool { 
    return a < b 
})

// Basic operations and merging
heap.Push(1, 1)
heap.Push(2, 2)
value, _ := heap.PopValue()
heap.MergeWith(otherHeap)  // Merge with another heap
```

### Full Tree-Based Heaps

```go
// Full heap with node tracking
heap := heapcraft.NewPairingHeap[int](nil, func(a, b int) bool { 
    return a < b 
})

// Node tracking operations
id := heap.Push(42, 10)
heap.UpdateValue(id, 100)
heap.UpdatePriority(id, 1)
value, _ := heap.GetValue(id)
heap.Remove(id)
```

### Thread Safety

All heaps are thread-safe by default:

```go
// Safe to use concurrently
heap := heapcraft.NewSimplePairingHeap[int](nil, func(a, b int) bool { 
    return a < b 
})

// Multiple goroutines can safely call these methods
go func() {
    heap.Push(1, 1)
}()

go func() {
    value, err := heap.PopValue()
}()
```

## 📈 **Performance Benchmarks**

### Environment
<center>

<table width="100%">

| Parameter | Value |
|-----------|-------|
| Date | 2025-06-12 |
| OS | Windows 10  |
| Architecture | amd64 |
| CPU | AMD EPYC 7763 64-Core Processor |
| Go version | 1.24 |

</table>

</center>

### Micro-benchmarks

#### D-ary and Radix Heaps
<center>

<table width="100%">

| Heap | Operation | Iterations | ns/op | B/op | allocs/op |
|------|-----------|-----------:|------:|-----:|----------:|
| **BinaryHeap**        | Insertion | 16,758,372 |  73.9ns |    96B |         0 |
|                       | Deletion  |  3,062,550 | 420.5ns |    16B |         1 |
|                       | PushPop   | 26,092,969 |  43.8ns |    16B |         1 |
|                       | PopPush   | 27,053,104 |  43.6ns |    16B |         1 |
| **DaryHeap (d=3)**    | Insertion | 26,351,965 |  58.1ns |    95B |         0 |
|                       | Deletion  |  3,410,320 | 336.0ns |    16B |         1 |
|                       | PushPop   | 26,072,559 |  43.4ns |    16B |         1 |
|                       | PopPush   | 27,454,550 |  43.0ns |    16B |         1 
| **DaryHeap (d=4)**    | Insertion | 27,776,812 |  38.0ns |    90B |         0 |
|                       | Deletion  |  4,467,697 | 291.3ns |    16B |         1 |
|                       | PushPop   | 27,006,224 |  44.1ns |    16B |         1 |
|                       | PopPush   | 25,732,628 |  42.9ns |    16B |         1 |
| **RadixHeap**         | Insertion | 18,605,545 |  61.6ns |    87B |         0 |
|                       | Deletion  |  2,160,903 | 582.1ns |   494B |         4 |

</table>

</center>

#### Tree-based Heaps
<center>

<table width="100%">

| Heap | Operation | Iterations | ns/op | B/op | allocs/op |
|------|-----------|-----------:|------:|-----:|----------:|
| **LeftistHeap**       | Insertion |  8,288,563 | 139.8ns |    48B |         1 |
|                       | Deletion  |  1,919,752 | 753.3ns |     0B |         0 |
| **PairingHeap**       | Insertion | 16,043,745 | 77.11ns |    32B |         1 |
|                       | Deletion  | 13,623,637 | 121.7ns |     0B |         0 |
| **SkewHeap**          | Insertion |  3,978,996 | 479.0ns |    32B |         1 |
|                       | Deletion  |  3,240,573 | 539.5ns |     0B |         0 |

</table>

</center>

## 🤝 **License**

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

## 📞 **Contact & Contributing**

Feel free to open an [issue](https://github.com/galactixx/heapcraft/issues) or a pull request. Discussion and feedback are welcome!
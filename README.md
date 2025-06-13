<p align="center">
  <img src="/docs/logo.png" alt="heapcraft logo" width="75%"/>
</p>

**heapcraft** is a high‑performance Go library offering a comprehensive suite of advanced heap data structures—binary heaps, d‑ary heaps, pairing heaps, radix heaps, skew heaps, and leftist heaps—for lightning‑fast priority‑queue operations.

Use it wherever you need efficient scheduling, graph algorithms, event simulation, load balancing, or any task that requires ordered extraction by priority.

### Available Heap Types

<div align="center">

| Heap Type | Implementation | Special Features |
|-----------|----------------|------------------|
| **Binary** | `BinaryHeap` | Standard binary heap |
| **D-ary** | `DaryHeap` | Configurable arity (2-ary, 3-ary, etc.) |
| **Radix** | `RadixHeap` | Integer priorities, bucket-based |
| **Pairing** | `SimplePairingHeap` / `PairingHeap` | Constant-time meld, efficient decrease-key |
| **Skew** | `SimpleSkewHeap` / `SkewHeap` | Self-adjusting, amortized O(log n) |
| **Leftist** | `SimpleLeftistHeap` / `LeftistHeap` | Leftist property, efficient merge |

</div>

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

| Parameter | Value |
|-----------|-------|
| Date | 2025-06-12 |
| OS | Windows 10  |
| Architecture | amd64 |
| CPU | AMD EPYC 7763 64-Core Processor |
| Go version | 1.24 |

### Micro-benchmarks

#### D-ary and Radix Heaps

<div align="center">

| Heap Type | Insertion | | Deletion | | PushPop | | PopPush | |
|-----------|-----------|-----------|----------|----------|----------|----------|----------|----------|
| | Iterations | Time (ns/op) | Iterations | Time (ns/op) | Iterations | Time (ns/op) | Iterations | Time (ns/op) |
| **BinaryHeap** | 16,758,372 | 73.9 | 3,062,550 | 420.5 | 26,092,969 | 43.8 | 27,053,104 | 43.6 |
| **DaryHeap (d=3)** | 26,351,965 | 58.1 | 3,410,320 | 336.0 | 26,072,559 | 43.4 | 27,454,550 | 43.0 |
| **DaryHeap (d=4)** | 27,776,812 | 38.0 | 4,467,697 | 291.3 | 27,006,224 | 44.1 | 25,732,628 | 42.9 |
| **RadixHeap** | 18,605,545 | 61.6 | 2,160,903 | 582.1 | - | - | - | - |

</div>

#### Tree-based Heaps

<div align="center">

| Heap Type | Insertion | | Deletion | | PushPop | | PopPush | |
|-----------|-----------|-----------|----------|----------|----------|----------|----------|----------|
| | Iterations | Time (ns/op) | Iterations | Time (ns/op) | Iterations | Time (ns/op) | Iterations | Time (ns/op) |
| **LeftistHeap** | 8,288,563 | 139.8 | 1,919,752 | 753.3 | - | - | - | - |
| **PairingHeap** | 16,043,745 | 77.11 | 13,623,637 | 121.7 | - | - | - | - |
| **SkewHeap** | 3,978,996 | 479.0 | 3,240,573 | 539.5 | - | - | - | - |

</div>

## 🤝 **License**

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

## 📞 **Contact & Contributing**

Feel free to open an [issue](https://github.com/galactixx/heapcraft/issues) or a pull request. Discussion and feedback are welcome!
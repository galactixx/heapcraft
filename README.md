<p align="center">
  <img src="/docs/logo.png" alt="heapcraft logo" width="75%"/>
</p>

**heapcraft** is a high‑performance Go library offering a comprehensive suite of advanced heap data structures—binary heaps, d‑ary heaps, pairing heaps, radix heaps, skew heaps, and leftist heaps—for lightning‑fast priority‑queue operations.

Use it wherever you need efficient scheduling, graph algorithms, event simulation, load balancing, or any task that requires ordered extraction by priority.

---

## ✨ **Features**

| Category            | Details                                                                                    |
| ------------------- | ------------------------------------------------------------------------------------------ |
| **Heap Variants**   | `Binary`, `D‑ary`, `Pairing`, `Radix`, `Skew`, `Leftist`                                  |
| **Two Implementations** | **Simple** and **Full** versions available for `Pairing`, `Skew`, and `Leftist` heaps    |
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

## 📚 **Usage**

### Simple vs Full Implementations

**Simple Heaps** (`Simple*Heap`) provide basic heap operations:
- `Insert(value, priority)` - Add elements
- `Pop()` / `PopValue()` / `PopPriority()` - Remove elements
- `Peek()` / `PeekValue()` / `PeekPriority()` - View without removing
- `Length()`, `IsEmpty()`, `Clear()`, `Clone()`

**Full Heaps** (`*Heap`) extend simple heaps with node tracking:
- All simple heap operations
- `Insert()` returns a unique node ID
- `UpdateValue(id, newValue)` - Update node value
- `UpdatePriority(id, newPriority)` - Update node priority
- `Get(id)`, `GetValue(id)`, `GetPriority(id)` - Retrieve by ID
- `Remove(id)` - Remove specific node

### Simple Heaps

```go
// Simple heap (basic operations only)
simpleHeap := heapcraft.NewSimplePairingHeap[int](nil, func(a, b int) bool { 
    return a < b 
})

// Insert elements
simpleHeap.Insert(1, 1)
simpleHeap.Insert(2, 2)

// Pop elements (with highest priority)
value, err := simpleHeap.PopValue()

// Peek without removing
peekValue, err := simpleHeap.PeekValue()
peekPriority, err := simpleHeap.PeekPriority()

// Check size and clear
length := simpleHeap.Length()
isEmpty := simpleHeap.IsEmpty()
simpleHeap.Clear()
```

### Full Heaps

```go
// Full heap (with node tracking)
fullHeap := heapcraft.NewPairingHeap[int](nil, func(a, b int) bool { 
    return a < b 
})

// Insert and get node ID
id1 := fullHeap.Insert(42, 10)
id2 := fullHeap.Insert(15, 5)

// Update node value and priority
fullHeap.UpdateValue(id1, 100)
fullHeap.UpdatePriority(id2, 1)

// Get node by ID
node, err := fullHeap.Get(id1)
value, err := fullHeap.GetValue(id2)
priority, err := fullHeap.GetPriority(id1)

// Remove specific node
removedNode, err := fullHeap.Remove(id2)
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
    heap.Insert(1, 1)
}()

go func() {
    value, err := heap.PopValue()
}()
```

### Available Heap Types

| Heap Type | Simple Implementation | Full Implementation | Special Features |
|-----------|----------------------|-------------------|------------------|
| **Binary** | `BinaryHeap` | - | Standard binary heap |
| **D-ary** | `DaryHeap` | - | Configurable arity (2-ary, 3-ary, etc.) |
| **Radix** | `RadixHeap` | - | Integer priorities, bucket-based |
| **Pairing** | `SimplePairingHeap` | `PairingHeap` | Constant-time meld, efficient decrease-key |
| **Skew** | `SimpleSkewHeap` | `SkewHeap` | Self-adjusting, amortized O(log n) |
| **Leftist** | `SimpleLeftistHeap` | `LeftistHeap` | Leftist property, efficient merge |

### Advanced Features

```go
// D-ary heap with custom arity (4-ary)
daryHeap := heapcraft.NewDaryHeap[int](4, nil, func(a, b int) bool { 
    return a < b 
})

// Radix heap for integer priorities
radixHeap := heapcraft.NewRadixHeap[int, uint](nil)

// Merge two heaps
heap1 := heapcraft.NewSimplePairingHeap[int](nil, func(a, b int) bool { 
    return a < b 
})
heap2 := heapcraft.NewSimplePairingHeap[int](nil, func(a, b int) bool { 
    return a < b 
})
heap1.MergeWith(heap2)

// Pop and push in one operation
value := heap.PopPush(newValue, newPriority)
```

## 🤝 **License**

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

## 📞 **Contact & Contributing**

Feel free to open an [issue](https://github.com/galactixx/heapcraft/issues) or a pull request. Discussion and feedback are welcome!
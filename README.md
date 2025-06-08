<p align="center">
  <img src="/docs/logo.png" alt="heapcraft logo" width="75%"/>
</p>

**heapcraft** is a high‑performance Go library offering a comprehensive suite of advanced heap data structures—d‑ary heaps, pairing heaps, binary heaps, radix heaps, skew heaps, and leftist heaps—for lightning‑fast priority‑queue operations.

Use it wherever you need efficient scheduling, graph algorithms, event simulation, load balancing, or any task that requires ordered extraction by priority.

---

## ✨ **Features**

| Category            | Details                                                                                    |
| ------------------- | ------------------------------------------------------------------------------------------ |
| Heap Variants       | `Binary`, `D‑ary`, `Pairing`, `Radix`, `Skew`, `Leftist`      |
| Thread‑Safe         | Coarse‑grained thread‑safe variants available for all heap types via `Safe*` wrappers. |
| Decrease‑Key / Meld | Native support where algorithmically possible; constant‑time meld on pairing heaps. |
| Generics            | Go 1.18+ type parameters—store any comparable or custom type.                              |

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

#### Standard Heaps

```go
// Create a new heap
heap := heapcraft.NewSimplePairingHeap[int, int](nil, func(a, b int) bool { return a < b })

// Insert elements
heap.Insert(1, 1)
heap.Insert(2, 2)

// Pop elements
value := heap.Pop() // Returns the element with highest priority
```

#### Thread‑Safe Heaps

```go
// Create a new thread‑safe heap
safeHeap := heapcraft.NewSafeSimplePairingHeap[int, int](nil, func(a, b int) bool { return a < b })

// Thread‑safe operations
safeHeap.Insert(1, 1) // Thread‑safe insert
value := safeHeap.Pop() // Thread‑safe pop
```

All heap types (`Binary`, `D‑ary`, `Pairing`, `Radix`, `Skew`, `Leftist`) have corresponding thread‑safe variants prefixed with `Safe*`. These wrappers provide coarse‑grained synchronization using read‑write locks, making them suitable for concurrent access.

---

## 🤝 **License**

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

## 📞 **Contact & Contributing**

Feel free to open an [issue](https://github.com/galactixx/heapcraft/issues) or a pull request.  Discussion and feedback are welcome!
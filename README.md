<p align="center">
  <img src="/docs/logo.png" alt="heapcraft logo" width="75%"/>
</p>

**heapcraft** is a high‑performance Go library offering a comprehensive suite of advanced heap data structures—d‑ary heaps, pairing heaps, binary heaps, radix heaps, skew heaps, and leftist heaps—for lightning‑fast priority‑queue operations.

Use it wherever you need efficient scheduling, graph algorithms (Dijkstra, A\*), event simulation, load balancing, or any task that requires ordered extraction by priority.

---

## ✨ **Features**

| Category            | Details                                                                                    |
| ------------------- | ------------------------------------------------------------------------------------------ |
| Heap Variants       | `Binary`, `D‑ary`, `Pairing`, `Radix`, `Skew`, `Leftist`      |
| Decrease‑Key / Meld | Native support where algorithmically possible; constant‑time meld on pairing heaps. |
| Generics            | Go 1.18+ type parameters—store any comparable or custom type.                              |

---

## 🚀 **Getting Started**

### Install

```bash
go get github.com/galactixx/heapcraft@latest
```

Then import it in your code:

```go
import "github.com/galactixx/heapcraft"
```

---


## 🤝 **License**

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

## 📞 **Contact & Contributing**

Feel free to open an [issue](https://github.com/galactixx/heapcraft/issues) or a pull request.  Discussion and feedback are welcome!
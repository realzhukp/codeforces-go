package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"sort"
)

type stpNode struct {
	lr       [2]*stpNode
	priority uint
	key      int
}

func (o *stpNode) rotate(d int) *stpNode {
	x := o.lr[d^1]
	o.lr[d^1] = x.lr[d]
	x.lr[d] = o
	return x
}

type sTreap struct {
	rd         uint
	root       *stpNode
	comparator func(a, b int) int
}

func (t *sTreap) fastRand() uint {
	t.rd ^= t.rd << 13
	t.rd ^= t.rd >> 17
	t.rd ^= t.rd << 5
	return t.rd
}

func (t *sTreap) _put(o *stpNode, key int) *stpNode {
	if o == nil {
		return &stpNode{priority: t.fastRand(), key: key}
	}
	cmp := t.comparator(key, o.key)
	o.lr[cmp] = t._put(o.lr[cmp], key)
	if o.lr[cmp].priority > o.priority {
		o = o.rotate(cmp ^ 1)
	}
	return o
}

func (t *sTreap) put(key int) { t.root = t._put(t.root, key) }

func (t *sTreap) _delete(o *stpNode, key int) *stpNode {
	if o == nil {
		return nil
	}
	if cmp := t.comparator(key, o.key); cmp >= 0 {
		o.lr[cmp] = t._delete(o.lr[cmp], key)
	} else {
		if o.lr[1] == nil {
			return o.lr[0]
		}
		if o.lr[0] == nil {
			return o.lr[1]
		}
		cmp2 := 0
		if o.lr[0].priority > o.lr[1].priority {
			cmp2 = 1
		}
		o = o.rotate(cmp2)
		o.lr[cmp2] = t._delete(o.lr[cmp2], key)
	}
	return o
}

func (t *sTreap) delete(key int) { t.root = t._delete(t.root, key) }

func (t *sTreap) min() (min *stpNode) {
	for o := t.root; o != nil; o = o.lr[0] {
		min = o
	}
	return
}

type tpNode struct {
	lr       [2]*tpNode
	priority uint
	key      int
	st       *sTreap
}

func (o *tpNode) rotate(d int) *tpNode {
	x := o.lr[d^1]
	o.lr[d^1] = x.lr[d]
	x.lr[d] = o
	return x
}

type treap struct {
	rd         uint
	root       *tpNode
	comparator func(a, b int) int
}

func newTreap() *treap {
	return &treap{
		rd: 1,
		comparator: func(a, b int) int {
			switch {
			case a < b:
				return 0
			case a > b:
				return 1
			default:
				return -1
			}
		},
	}
}

func (t *treap) fastRand() uint {
	t.rd ^= t.rd << 13
	t.rd ^= t.rd >> 17
	t.rd ^= t.rd << 5
	return t.rd
}

func (t *treap) _put(o *tpNode, key, val int) *tpNode {
	if o == nil {
		st := &sTreap{rd: 1, comparator: t.comparator}
		st.put(val)
		return &tpNode{priority: t.fastRand(), key: key, st: st}
	}
	if cmp := t.comparator(key, o.key); cmp >= 0 {
		o.lr[cmp] = t._put(o.lr[cmp], key, val)
		if o.lr[cmp].priority > o.priority {
			o = o.rotate(cmp ^ 1)
		}
	} else {
		o.st.put(val)
	}
	return o
}

func (t *treap) put(key, val int) { t.root = t._put(t.root, key, val) }

func (t *treap) _delete(o *tpNode, key, val int) *tpNode {
	if o == nil {
		return nil
	}
	if cmp := t.comparator(key, o.key); cmp >= 0 {
		o.lr[cmp] = t._delete(o.lr[cmp], key, val)
	} else {
		o.st.delete(val)
		if o.st.root == nil {
			if o.lr[1] == nil {
				return o.lr[0]
			}
			if o.lr[0] == nil {
				return o.lr[1]
			}
			cmp2 := 0
			if o.lr[0].priority > o.lr[1].priority {
				cmp2 = 1
			}
			o = o.rotate(cmp2)
			o.lr[cmp2] = t._delete(o.lr[cmp2], key, val)
		}
	}
	return o
}

func (t *treap) delete(key, val int) { t.root = t._delete(t.root, key, val) }

func (t *treap) ceiling(key int) (ceiling *tpNode) {
	for o := t.root; o != nil; {
		switch cmp := t.comparator(key, o.key); {
		case cmp == 0:
			ceiling = o
			o = o.lr[0]
		case cmp > 0:
			o = o.lr[1]
		default:
			return o
		}
	}
	return
}

// github.com/EndlessCheng/codeforces-go
func Sol840D(reader io.Reader, writer io.Writer) {
	in := bufio.NewScanner(reader)
	in.Split(bufio.ScanWords)
	out := bufio.NewWriter(writer)
	defer out.Flush()
	read := func() (x int) {
		in.Scan()
		for _, b := range in.Bytes() {
			x = x*10 + int(b-'0')
		}
		return
	}

	n, q := read(), read()
	a := make([]int, n+1)
	for i := 1; i <= n; i++ {
		a[i] = read()
	}
	ans := make([]int, q)
	type query struct {
		blockIdx     int
		l, r, k, idx int
	}
	qs := make([]query, q)
	blockSize := int(math.Round(math.Sqrt(float64(n))))
	for i := range qs {
		l, r, k := read(), read()+1, read()
		qs[i] = query{l / blockSize, l, r, k, i}
	}
	sort.Slice(qs, func(i, j int) bool {
		qi, qj := qs[i], qs[j]
		return qi.blockIdx < qj.blockIdx || qi.blockIdx == qj.blockIdx && qi.r < qj.r
	})

	cntMap := map[int]int{}
	t := newTreap()
	update := func(idx, delta int) {
		v := a[idx]
		cntMap[v] += delta
		cnt := cntMap[v]
		t.delete(cnt-delta, v)
		t.put(cnt, v)
	}
	getAns := func(low int) int {
		if o := t.ceiling(low); o != nil {
			return o.st.min().key
		}
		return -1
	}

	l, r := 1, 1
	for _, q := range qs {
		for ; l < q.l; l++ {
			update(l, -1)
		}
		for ; r < q.r; r++ {
			update(r, 1)
		}
		for l > q.l {
			l--
			update(l, 1)
		}
		for r > q.r {
			r--
			update(r, -1)
		}
		ans[q.idx] = getAns((q.r-q.l)/q.k + 1)
	}
	for _, v := range ans {
		fmt.Fprintln(out, v)
	}
}

//func main() {
//	Sol840D(os.Stdin, os.Stdout)
//}

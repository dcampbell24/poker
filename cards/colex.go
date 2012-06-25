package cards

func mulRange(a, b int64) int64 {
	for i := a + 1; i <= b; i++ {
		a *= i
	}
	return a
}

// Optimized for k <= n - k.
func binomial(n, k int64) int64 {
	return mulRange(n-k+1, n) / mulRange(1, k)
}

func newColex() func([]int32) int64 {
	tmp := make([]int32, 0, 52)
	return func(set []int32) int64 {
		sum := int64(0)
		tmp = tmp[:len(set)]
		copy(tmp, set)
		sort(tmp)
		for i, v := range tmp {
			sum += binomial(int64(v)-1, int64(i)+1)
		}
		return sum
	}
}

// Assumes the values of a and b are sorted in ascending order.
func less(a, b []int32) bool {
	if len(a) > len(b) {
		return true
	}
	if len(a) < len(b) {
		return false
	}
	for i, v := range a {
		if v < b[i] {
			return true
		}
		if v > b[i] {
			return false
		}
	}
	return false
}

func sort(a []int32) {
	if len(a) < 2 {
		return
	}
	for i := 1; i < len(a); i++ {
		for j := i; j > 0 && a[j] < a[j-1]; j-- {
			a[j], a[j-1] = a[j-1], a[j]
		}
	}
}

func sort2(a [][]int32) {
	if len(a) < 2 {
		return
	}
	for i := 1; i < len(a); i++ {
		for j := i; j > 0 && less(a[j], a[j-1]); j-- {
			a[j], a[j-1] = a[j-1], a[j]
		}
	}
}

// Rules for a canonical hand:
//	1. The cards are in sorted order
//
//	2. The i-th suit must have at least as many cards as all later suits.  If a
//	   suit isn't present, it counts as having 0 cards.
//
//	3. If two suits have the same number of cards, the ranks in the first suit
//	   must be lower or equal lexicographically (e.g., [1, 3] <= [2, 4]).
//
//	4. Must be a valid hand (no duplicate cards).
func newCanonical() func([]int32) []int32 {
	w := make([][]int32, 4)
	for i := range w {
		w[i] = make([]int32, 0, 13)
	}
	return func(cards []int32) []int32 {
		for i := range w {
			w[i] = w[i][:0]
		}
		for _, v := range cards {
			r := v % 4
			w[r] = append(w[r], (v-1)/4)
		}
		for i := range w {
			sort(w[i])
		}
		sort2(w)
		cs := make([]int32, 0, len(cards))
		for i, v := range w {
			for _, c := range v {
				cs = append(cs, 4*c+int32(i)+1)
			}
		}
		return cs
	}
}

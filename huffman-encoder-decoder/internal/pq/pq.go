package pq

type CharCount struct {
	Chr   rune
	Count int
	Mask  string
	Left  *CharCount
	Right *CharCount
}

type PQ struct {
	Arr []*CharCount
}

func (pq *PQ) Len() int { return len(pq.Arr) }
func (pq *PQ) Less(i, j int) bool {
	return pq.Arr[i].Count < pq.Arr[j].Count
}
func (pq *PQ) Swap(i, j int) {
	pq.Arr[i], pq.Arr[j] = pq.Arr[j], pq.Arr[i]
}
func (pq *PQ) Push(x any) {
	pq.Arr = append(pq.Arr, x.(*CharCount))
}
func (pq *PQ) Pop() any {
	c := pq.Arr[len(pq.Arr)-1]
	pq.Arr = pq.Arr[:len(pq.Arr)-1]
	return c
}

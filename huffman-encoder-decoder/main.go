package main

import (
	"bufio"
	"container/heap"
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

type charCount struct {
	chr   rune
	count int
	mask  string
	left  *charCount
	right *charCount
}

type pq struct {
	arr []*charCount
}

func (pq *pq) Len() int { return len(pq.arr) }
func (pq *pq) Less(i, j int) bool {
	return pq.arr[i].count < pq.arr[j].count
}
func (pq *pq) Swap(i, j int) {
	pq.arr[i], pq.arr[j] = pq.arr[j], pq.arr[i]
}
func (pq *pq) Push(x any) {
	pq.arr = append(pq.arr, x.(*charCount))
}
func (pq *pq) Pop() any {
	c := pq.arr[len(pq.arr)-1]
	pq.arr = pq.arr[:len(pq.arr)-1]
	return c
}

type trie struct {
	chr      rune
	value    rune
	isEnd    bool
	children [2]*trie
}

func NewTrie() *trie {
	return &trie{
		isEnd:    false,
		children: [2]*trie{},
	}
}

func (t *trie) Insert(mask string, value rune) {
	if len(mask) == 0 {
		return
	}

	node := t
	for _, c := range mask {
		if node.children[0] == nil || node.children[1] == nil {
			left := NewTrie()
			right := NewTrie()
			node.children = [2]*trie{left, right}
		}

		if c == '1' {
			node = node.children[1]
		} else {
			node = node.children[0]
		}
	}

	node.isEnd = true
	node.value = value
}

func main() {
	// Get file
	if len(os.Args) < 2 {
		os.Exit(1)
	}

	filename := os.Args[1]
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// DECODE
	reader1 := bufio.NewReader(file)
	charByMask, err := decodeHeader(reader1)
	if err != nil {
		fmt.Println("ERROR DECODING " + err.Error())
		os.Exit(1)
	}

	searchTree := NewTrie()
	for k, v := range charByMask {
		searchTree.Insert(k, v)
	}

	outFile1, err := os.Create("output.txt")
	if err != nil {
		fmt.Printf("Error creating new file " + err.Error())
		os.Exit(1)
	}
	writer1 := bufio.NewWriter(outFile1)

	node := searchTree
	for {
		b, err := reader1.ReadByte()
		if err != nil {
			break
		}

		byteMask := ""
		for range 8 {
			if b&1 == 1 {
				byteMask = "1" + byteMask
			} else {
				byteMask = "0" + byteMask
			}
			b = b >> 1
		}

		for _, m := range byteMask {
			if m == '0' {
				node = node.children[0]
			} else {
				node = node.children[1]
			}

			if node == nil {
				os.Exit(5)
			}

			if node.isEnd {
				if _, err := writer1.WriteRune(node.value); err != nil {
					fmt.Printf("Error writing to file " + err.Error())
					os.Exit(1)
				}

				node = searchTree
			}
		}

	}

	writer1.Flush()
	outFile1.Close()
	file.Close()
	os.Exit(0)

	// Get Frequency
	charCountMap, err := getCharCount(bufio.NewReader(file))
	if err != nil {
		fmt.Println("getcharcount " + err.Error())
		os.Exit(1)
	}

	// Sort
	charCountArr := make([]*charCount, 0)

	for k, v := range charCountMap {
		charCountArr = append(charCountArr, &charCount{
			chr:   k,
			count: v,
		})
	}
	q := pq{arr: charCountArr}
	heap.Init(&q)

	// Prefix table / tree
	for q.Len() > 1 {
		first := heap.Pop(&q).(*charCount)
		second := heap.Pop(&q).(*charCount)

		parentNode := &charCount{
			count: first.count + second.count,
		}

		if first.count < second.count {
			parentNode.left = first
			parentNode.right = second
		} else {
			parentNode.left = second
			parentNode.right = first
		}

		heap.Push(&q, parentNode)
	}

	root := heap.Pop(&q).(*charCount)
	postOrder(root, "")
	maskByChar := make(map[rune]string)
	getMaskByChar(root, maskByChar)

	// Generate new file with masks
	outFile, err := os.OpenFile("out.huff", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("outfile " + err.Error())
		os.Exit(1)
	}

	// Write data
	if _, err := file.Seek(0, 0); err != nil {
		fmt.Println("seek " + err.Error())
		os.Exit(1)
	}

	reader := bufio.NewReader(file)

	// Write Header
	writeHeader(outFile, maskByChar)

	var currByte byte
	var bitCount int
	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Read error")
			os.Exit(1)
		}

		mask := maskByChar[r]

		for _, c := range mask {
			currByte = currByte << 1
			if c == '1' {
				currByte = currByte | 1
			}

			bitCount++
			if bitCount == 8 {
				_, err = outFile.Write([]byte{currByte})
				if err != nil {
					fmt.Println("Write error " + err.Error())
					os.Exit(1)
				}

				currByte = 0
				bitCount = 0
			}
		}
	}

	if bitCount > 0 {
		currByte = currByte << (8 - bitCount)
		_, err := outFile.Write([]byte{currByte})
		if err != nil {
			fmt.Println("Write error " + err.Error())
			os.Exit(1)
		}
	}

	file.Close()
	outFile.Close()
}

func getCharCount(reader *bufio.Reader) (map[rune]int, error) {
	charCountMap := make(map[rune]int)
	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		charCountMap[r]++

	}
	return charCountMap, nil
}

func postOrder(root *charCount, mask string) {
	if root.left != nil {
		postOrder(root.left, mask+"0")
	}

	if root.right != nil {
		postOrder(root.right, mask+"1")
	}

	root.mask = mask
}

func getMaskByChar(root *charCount, maskByChar map[rune]string) {
	if root.left != nil {
		getMaskByChar(root.left, maskByChar)
	}

	if root.right != nil {
		getMaskByChar(root.right, maskByChar)
	}

	if root.left == nil && root.right == nil {
		maskByChar[root.chr] = root.mask
	}
}

func writeHeader(file *os.File, maskByChar map[rune]string) error {
	writer := bufio.NewWriter(file)

	for key, mask := range maskByChar {
		writer.WriteRune(key)

		buf := make([]byte, 8)
		binary.BigEndian.PutUint64(buf, uint64(len(mask)))
		writer.Write(buf)

		var currByte byte
		var bitCount int
		for _, c := range mask {
			currByte = currByte << 1
			if c == '1' {
				currByte = currByte | 1
			}

			bitCount++
			if bitCount == 8 {
				writer.WriteByte(currByte)
				currByte = 0
				bitCount = 0
			}
		}

		if bitCount > 0 {
			writer.WriteByte(currByte)
		}
		writer.WriteByte(byte(','))
	}

	writer.WriteByte(byte('\a'))
	writer.WriteByte(byte('\a'))
	writer.Flush()
	return nil
}

func decodeHeader(reader *bufio.Reader) (map[string]rune, error) {
	charByMask := make(map[string]rune)
	for {
		k, _, err := reader.ReadRune()
		if err != nil {
			return nil, err
		}

		buf := make([]byte, 8)
		_, err = reader.Read(buf)
		if err != nil {
			return nil, err
		}

		maskLen := binary.BigEndian.Uint64(buf)

		mask := ""
		noOfBytes := maskLen
		if maskLen%8 > 0 {
			noOfBytes = (maskLen + (8 - (maskLen % 8)))
		}
		noOfBytes = noOfBytes / 8

		for range noOfBytes {
			b, err := reader.ReadByte()
			if err != nil {
				return nil, err
			}

			byteMask := ""
			for range 8 {
				if maskLen <= 0 {
					break
				}
				if b&1 == 1 {
					byteMask = "1" + byteMask
				} else {
					byteMask = "0" + byteMask
				}
				b = b >> 1
				maskLen--
			}

			mask = mask + byteMask
		}
		charByMask[mask] = k

		token, err := reader.ReadByte()
		if err != nil {
			return nil, err
		}

		if token != ',' {
			return nil, fmt.Errorf("MALFORMED HEADER")
		}

		nextToken, err := reader.Peek(2)
		if err != nil {
			return nil, err
		}

		if token == ',' && nextToken[0] == '\a' && nextToken[1] == '\a' {
			reader.ReadByte()
			reader.ReadByte()
			break
		}
	}

	return charByMask, nil
}

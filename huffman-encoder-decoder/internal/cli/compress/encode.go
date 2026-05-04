package compress

import (
	"bufio"
	"container/heap"
	"fmt"
	"io"
	"os"

	"github.com/madraceee/coding-challenges/huffman-encoder-decoder/internal/header"
	"github.com/madraceee/coding-challenges/huffman-encoder-decoder/internal/pq"
)

func Encode(inputFileName, outputFileName string) error {
	file, err := os.Open(inputFileName)
	if err != nil {
		return err
	}
	defer file.Close()

	charCountMap, err := getCharCount(bufio.NewReader(file))
	if err != nil {
		return fmt.Errorf("Error reading %s - %s", inputFileName, err)
	}

	// Sort
	charCountArr := make([]*pq.CharCount, 0)

	for k, v := range charCountMap {
		charCountArr = append(charCountArr, &pq.CharCount{
			Chr:   k,
			Count: v,
		})
	}
	prefixTree := pq.PQ{Arr: charCountArr}
	heap.Init(&prefixTree)

	// Prefix tree
	for prefixTree.Len() > 1 {
		first := heap.Pop(&prefixTree).(*pq.CharCount)
		second := heap.Pop(&prefixTree).(*pq.CharCount)

		parentNode := &pq.CharCount{
			Count: first.Count + second.Count,
		}

		if first.Count < second.Count {
			parentNode.Left = first
			parentNode.Right = second
		} else {
			parentNode.Left = second
			parentNode.Right = first
		}

		heap.Push(&prefixTree, parentNode)
	}

	root := heap.Pop(&prefixTree).(*pq.CharCount)
	postOrder(root, "")
	maskByChar := make(map[rune]string)
	getMaskByChar(root, maskByChar)

	// Generate new file with masks
	outFile, err := os.OpenFile(outputFileName, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Write data
	if _, err := file.Seek(0, 0); err != nil {
		return err
	}

	// Write Header
	if err := header.WriteHeader(bufio.NewWriter(outFile), maskByChar); err != nil {
		return err
	}

	// Write body
	var currByte byte
	var bitCount int
	reader := bufio.NewReader(file)
	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
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
					return err
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
			return err
		}
	}

	return nil
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

func postOrder(root *pq.CharCount, mask string) {
	if root.Left != nil {
		postOrder(root.Left, mask+"0")
	}

	if root.Right != nil {
		postOrder(root.Right, mask+"1")
	}

	root.Mask = mask
}

func getMaskByChar(root *pq.CharCount, maskByChar map[rune]string) {
	if root.Left != nil {
		getMaskByChar(root.Left, maskByChar)
	}

	if root.Right != nil {
		getMaskByChar(root.Right, maskByChar)
	}

	if root.Left == nil && root.Right == nil {
		maskByChar[root.Chr] = root.Mask
	}
}

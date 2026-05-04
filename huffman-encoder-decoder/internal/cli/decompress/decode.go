package decompress

import (
	"bufio"
	"fmt"
	"os"

	"github.com/madraceee/coding-challenges/huffman-encoder-decoder/internal/header"
	"github.com/madraceee/coding-challenges/huffman-encoder-decoder/internal/trie"
)

func Decode(inputFileName, outputFileName string) error {
	inFile, err := os.Open(inputFileName)
	if err != nil {
		return err
	}
	defer inFile.Close()

	reader := bufio.NewReader(inFile)
	charByMask, err := header.DecodeHeader(reader)
	if err != nil {
		fmt.Println("ERROR DECODING " + err.Error())
		os.Exit(1)
	}

	searchTree := trie.NewTrie()
	for k, v := range charByMask {
		searchTree.Insert(k, v)
	}

	outFile, err := os.Create(outputFileName)
	if err != nil {
		return err
	}
	defer outFile.Close()

	writer := bufio.NewWriter(outFile)
	node := searchTree
	for {
		b, err := reader.ReadByte()
		if err != nil {
			break
		}

		for i := range 8 {
			bit := b >> (7 - i)
			if bit&1 == 1 {
				node = node.Children[1]
			} else {
				node = node.Children[0]
			}

			if node == nil {
				return fmt.Errorf("Malformed Header")
			}

			if node.IsEnd {
				if _, err := writer.WriteRune(node.Value); err != nil {
					return err
				}
				node = searchTree
			}
		}

	}

	if err := writer.Flush(); err != nil {
		return err
	}
	return nil
}

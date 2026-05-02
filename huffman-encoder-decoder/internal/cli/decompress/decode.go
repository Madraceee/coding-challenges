package decompress

import (
	"bufio"
	"fmt"
	"os"

	"github.com/madraceee/coding-challenges/huffman-encoder-decoder/internal/header"
	"github.com/madraceee/coding-challenges/huffman-encoder-decoder/internal/trie"
)

func Decode(inputFileName, outputFileName string) error {
	file, err := os.Open(inputFileName)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
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
		fmt.Printf("Error creating new file " + err.Error())
		os.Exit(1)
	}
	writer := bufio.NewWriter(outFile)

	node := searchTree
	for {
		b, err := reader.ReadByte()
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
				node = node.Children[0]
			} else {
				node = node.Children[1]
			}

			if node == nil {
				os.Exit(5)
			}

			if node.IsEnd {
				if _, err := writer.WriteRune(node.Value); err != nil {
					fmt.Printf("Error writing to file " + err.Error())
					os.Exit(1)
				}

				node = searchTree
			}
		}

	}

	if err := writer.Flush(); err != nil {
		return err
	}
	outFile.Close()
	file.Close()
	return nil
}

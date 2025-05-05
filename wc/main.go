package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"unicode"
)

type FileStats struct {
	ByteCount      int64
	CharacterCount int64
	LineCount      int64
	WordCount      int64
}

func main() {
	var isByte, isLine, isCharacter, isWord bool

	flag.BoolVar(&isByte, "c", false, "Print no. of bytes")
	flag.BoolVar(&isLine, "l", false, "Print no. of lines")
	flag.BoolVar(&isWord, "w", false, "Print no. of words")
	flag.BoolVar(&isCharacter, "m", false, "Print no. of characters")
	flag.Parse()

	if !isByte && !isLine && !isCharacter && !isWord {
		isByte = true
		isLine = true
		isWord = true
	}

	fileName := flag.Arg(0)

	info, err := os.Stdin.Stat()
	if err != nil {
		fmt.Printf("Err: %s\n", err)
		os.Exit(1)
	}

	var file *os.File
	if (info.Mode() & os.ModeCharDevice) == 0 {
		file = os.Stdin
	} else {
		file, err = os.Open(fileName)
		if fileName == "" {
			flag.Usage()
			os.Exit(0)
		}
		if err != nil {
			fmt.Printf("Err: %s\n", err)
			os.Exit(1)
		}
		defer file.Close()
	}

	reader := bufio.NewReader(file)
	fileStats, err := GetStats(reader)

	output := ""

	if isLine {
		output += strconv.FormatInt(fileStats.LineCount, 10) + " "
	}

	if isWord {
		output += strconv.FormatInt(fileStats.WordCount, 10) + "  "
	}

	if isByte {
		output += strconv.FormatInt(fileStats.ByteCount, 10) + "  "
	}

	if isCharacter {
		output += strconv.FormatInt(fileStats.CharacterCount, 10) + "  "
	}

	fmt.Printf("  %s%s\n", output, fileName)
}

func GetStats(reader *bufio.Reader) (*FileStats, error) {
	var byteCount int64
	var characterCount int64
	var lineCount int64
	var wordCount int64

	insideWord := false
	for {
		r, size, err := reader.ReadRune()

		if r != unicode.ReplacementChar {
			characterCount++
			byteCount += int64(size)
		} else if err != nil {
			return nil, err
		} else {
			return nil, errors.New("invalid charactar")
		}

		if size == 0 {
			break
		}

		if unicode.IsSpace(r) {
			if insideWord == true {
				wordCount++
			}
			insideWord = false

			if r == '\n' {
				lineCount++
			}
		} else {
			insideWord = true
		}
	}

	return &FileStats{
		ByteCount:      byteCount,
		CharacterCount: characterCount - 1,
		LineCount:      lineCount,
		WordCount:      wordCount,
	}, nil
}

package header

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"os"
)

func WriteHeader(file *os.File, maskByChar map[rune]string) error {
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

func DecodeHeader(reader *bufio.Reader) (map[string]rune, error) {
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

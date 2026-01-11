package cmd

import (
	"bufio"
	"errors"
	"io"
	"strconv"
	"strings"
)

func GetCutData(reader *bufio.Reader, delimiter string) (*[]string, error) {
	result := make([]string, 0)
	isEof := false

	if CharacterColumnRange != "" {
		delimiter = ""
	}

	for {
		line := ""
		for {
			bytes, isPrefix, err := reader.ReadLine()
			if err == io.EOF {
				isEof = true
				break
			}
			if err != nil {
				return nil, err
			}

			line += string(bytes)

			if isPrefix == false {
				break
			}
		}

		if isEof {
			break
		}

		arr := strings.Split(line, delimiter)
		printRange, err := GetColumnRange(len(arr))
		if err != nil {
			return nil, err
		}

		resultArr := make([]string, 0)

		for _, col := range printRange {
			resultArr = append(resultArr, arr[col])
		}
		result = append(result, strings.Join(resultArr, delimiter))

	}

	return &result, nil
}

func GetColumnRange(maxLen int) ([]int, error) {
	ColumnRange := ""
	if FieldColumnRange != "" {
		ColumnRange = FieldColumnRange
	} else {
		ColumnRange = CharacterColumnRange
	}

	intRange := make([]int, 0)
	// No fields Specified
	if ColumnRange == "" {
		for i := range maxLen {
			intRange = append(intRange, i)
		}
		return intRange, nil
	}

	isList := false
	isRange := false

	for _, c := range ColumnRange {
		if c == ',' {
			isList = true
			break
		} else if c == '-' {
			isRange = true
			break
		}
	}

	fieldsDelimiter := " "
	if isList {
		fieldsDelimiter = ","
	} else if isRange {
		fieldsDelimiter = "-"
	}

	fields := strings.Split(ColumnRange, fieldsDelimiter)

	if isRange {
		if len(fields) > 2 {
			return []int{}, errors.New("invalid field range")
		}
		start := 0
		end := maxLen - 1
		var err error
		if fields[0] != "" {
			start, err = strconv.Atoi(fields[0])
			if err != nil {
				return []int{}, errors.New("invalid field range")
			}
		}
		if fields[1] != "" {
			end, err = strconv.Atoi(fields[1])
			if err != nil {
				return []int{}, errors.New("invalid field range")
			}
		}

		if start > end {
			return []int{}, errors.New("invalid field range")
		}

		for i := start - 1; i <= end; i++ {
			if end > maxLen {
				break
			}
			intRange = append(intRange, i)
		}
	} else {
		for _, field := range fields {
			val, err := strconv.Atoi(field)
			if err != nil {
				return []int{}, errors.New("invalid field range")
			}

			if val == 0 {
				return []int{}, errors.New("invalid field range")
			}
			if val > maxLen {
				continue
			}
			intRange = append(intRange, val-1)
		}
	}

	return intRange, nil
}

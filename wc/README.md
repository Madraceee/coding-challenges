# wc (Word Count)

A simple clone of the unix command to count the number of characters, lines, words, bytes for a given file in Go.

`ccwc [FLAG] [FILE]`

ccwc supports count by characters(-m), lines(-l), words(-w) and bytes(-c).

#### Examples
```
./ccwc test.txt
  7145 58164  342190  test.txt

./ccwc -c test.txt
  342190  test.txt

./ccwc -l test.txt
  7145 test.txt

/ccwc -m test.txt
  339292  test.txt

./ccwc -w test.txt
  58164  test.txt
```

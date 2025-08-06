# cut

A simple clone of cut command available on unix like operating systems written in Go.

#### Usage
`cut OPTION... [FILE]... [FLAG]...`

Flags
```
-c, --characters  Print the specified characters
-d, --delimiter   Delimiter required to perform split each row. (default "\t")
-f, --fields      Print the specified column(1-indexed)
```

#### Example
```
./cut -f1,2 test_data/sample.tsv
f0      f1
0       1
5       6
10      11
15      16
20      21


./cut -f1 -d, test_data/fourchords.csv  | head -n1
Song title
```

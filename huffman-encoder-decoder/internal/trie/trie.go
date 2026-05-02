package trie

type trie struct {
	Value    rune
	IsEnd    bool
	Children [2]*trie
}

func NewTrie() *trie {
	return &trie{
		IsEnd:    false,
		Children: [2]*trie{},
	}
}

func (t *trie) Insert(mask string, value rune) {
	if len(mask) == 0 {
		return
	}

	node := t
	for _, c := range mask {
		if node.Children[0] == nil || node.Children[1] == nil {
			left := NewTrie()
			right := NewTrie()
			node.Children = [2]*trie{left, right}
		}

		if c == '1' {
			node = node.Children[1]
		} else {
			node = node.Children[0]
		}
	}

	node.IsEnd = true
	node.Value = value
}

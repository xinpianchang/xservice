package stringx

// refer: https://github.com/tal-tech/go-zero/blob/master/core/stringx/trie.go

// Trie interface
type Trie interface {
	// Filter filter sentence get masked sentence and get keywords
	Filter(text string) (string, []string)
	// FindKeywords get sentence keywords
	FindKeywords(text string) []string
	// AddWord add keywords
	// Attention: NOT thread safe
	AddWords(text ...string)
}

//
type trieNode struct {
	node
}

var trieMask = '*'

type scope struct {
	start int
	stop  int
}

type node struct {
	children map[rune]*node
	end      bool
}

func (n *node) add(word string) {
	chars := []rune(word)
	if len(chars) == 0 {
		return
	}

	nd := n
	for _, char := range chars {
		if nd.children == nil {
			child := new(node)
			nd.children = map[rune]*node{
				char: child,
			}
			nd = child
		} else if child, ok := nd.children[char]; ok {
			nd = child
		} else {
			child := new(node)
			nd.children[char] = child
			nd = child
		}
	}

	nd.end = true
}

// NewTrie new trie
func NewTrie(words ...string) Trie {
	t := &trieNode{}
	t.AddWords(words...)
	return t
}

// Filter filter sentence get masked sentence and get keywords
func (t *trieNode) Filter(text string) (string, []string) {
	chars := []rune(text)
	if len(chars) == 0 {
		return text, nil
	}

	scopes := t.findScopes(chars)
	keywords := t.collectKeywords(chars, scopes)

	for _, match := range scopes {
		for i := match.start; i < match.stop; i++ {
			chars[i] = trieMask
		}
	}

	return string(chars), keywords
}

// FindKeywords get sentence keywords
func (t *trieNode) FindKeywords(text string) []string {
	chars := []rune(text)
	if len(chars) == 0 {
		return nil
	}

	scopes := t.findScopes(chars)
	return t.collectKeywords(chars, scopes)
}

// AddWord add keywords
func (t *trieNode) AddWords(words ...string) {
	for _, word := range words {
		t.add(word)
	}
}

func (t *trieNode) findScopes(chars []rune) []scope {
	var scopes []scope
	size := len(chars)
	start := -1

	for i := 0; i < size; i++ {
		child, ok := t.children[chars[i]]
		if !ok {
			continue
		}

		if start < 0 {
			start = i
		}
		if child.end {
			scopes = append(scopes, scope{
				start: start,
				stop:  i + 1,
			})
		}

		for j := i + 1; j < size; j++ {
			grandchild, ok := child.children[chars[j]]
			if !ok {
				break
			}

			child = grandchild
			if child.end {
				scopes = append(scopes, scope{
					start: start,
					stop:  j + 1,
				})
			}
		}

		start = -1
	}

	return scopes
}

func (t *trieNode) collectKeywords(chars []rune, scopes []scope) []string {
	set := make(map[string]bool)
	for _, v := range scopes {
		set[string(chars[v.start:v.stop])] = true
	}

	var i int
	keywords := make([]string, len(set))
	for k := range set {
		keywords[i] = k
		i++
	}

	return keywords
}

package stringx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrie(t *testing.T) {
	tests := []struct {
		in       string
		out      string
		keywords []string
	}{
		{
			in:  "Lorem ipsum dolor sit amet, ipsum nominati ocurreret ei per",
			out: "Lorem ***** dolor sit amet, ***** ******** ocurreret ei per",
			keywords: []string{
				"ipsum",
				"nominati",
			},
		},
		{
			in:       "ea timeam aliquip tacimates nec",
			out:      "ea timeam aliquip tacimates nec",
			keywords: []string{},
		},
	}

	trie := NewTrie(
		"ipsum",
		"nominati",
	)

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			out, keywords := trie.Filter(tt.in)
			assert.Equal(t, tt.out, out)
			assert.ElementsMatch(t, tt.keywords, keywords)
		})
	}
}

func BenchmarkTrie1(b *testing.B) {
	b.ReportAllocs()

	trie := NewTrie(
		"ipsum",
		"nominati",
	)

	for i := 0; i < b.N; i++ {
		_, _ = trie.Filter("Lorem ipsum dolor sit amet, ipsum nominati ocurreret ei per")
	}
}

func BenchmarkTrie2(b *testing.B) {
	b.ReportAllocs()

	trie := NewTrie(
		"ipsum",
		"nominati",
	)

	for i := 0; i < b.N; i++ {
		_, _ = trie.Filter("Lorem ipsum dolor sit amet, ipsum nominati ocurreret ei per, in quo tation nonumy, no iusto luptatum gloriatur vel. Per at solet quaestio, admodum feugait splendide ei vis. Mea ad mutat possit. Dicant nonumy animal duo id, no fugit platonem sea. In has zril labitur menandri, his dolorem eleifend et, eu ius wisi solet scribentur.")
	}
}

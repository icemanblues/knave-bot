package shakespeare

import (
	"math/rand"
	"strings"
	"time"
)

// Generator generates sentences following a formula.
type Generator struct {
	prefix  string
	columns [][]string
	postfix string
}

// Sentence the result of this generator's formula
func (g Generator) Sentence() string {
	return g.Generate(" ")
}

// Generate follows the formula using the delim
func (g Generator) Generate(delim string) string {
	builder := strings.Builder{}
	if g.prefix != "" {
		builder.WriteString(g.prefix)
	}

	rand.Seed(time.Now().Unix())
	for i, col := range g.columns {
		if i != 0 || g.prefix != "" {
			builder.WriteString(delim)
		}
		r := rand.Intn(len(col))
		builder.WriteString(col[r])
	}

	if g.postfix != "" {
		if g.prefix != "" || g.columns != nil {
			builder.WriteString(delim)
		}
		builder.WriteString(g.postfix)
	}

	return builder.String()
}

// New constructs a Generator
func New(pre, post string, cols [][]string) *Generator {
	return &Generator{
		prefix:  pre,
		postfix: post,
		columns: cols,
	}
}

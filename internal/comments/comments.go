// Package comments provides utilities for extracting comments from a package.
package comments

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/tools/go/packages"
)

// A DefaultMarkerPrefix that is commonly used by comment markers.
const DefaultMarkerPrefix = "+"

type fl struct {
	Filename string
	Line     int
}

// Comments for a particular package.
type Comments struct {
	groups map[fl]*ast.CommentGroup
	fset   *token.FileSet
}

// In returns all comments in a particular package.
func In(p *packages.Package) Comments {
	groups := map[fl]*ast.CommentGroup{}

	for _, f := range p.Syntax {
		for _, g := range f.Comments {
			p := p.Fset.Position(g.End())
			groups[fl{Filename: p.Filename, Line: p.Line}] = g
		}
	}
	return Comments{groups: groups, fset: p.Fset}
}

func (c Comments) For(o types.Object) string {
	p := c.fset.Position(o.Pos())
	return c.groups[fl{Filename: p.Filename, Line: p.Line - 1}].Text()
}

func (c Comments) Before(o types.Object) string {
	p := c.fset.Position(o.Pos())
	g := c.groups[fl{Filename: p.Filename, Line: p.Line - 1}]

	if g == nil {
		// No comment group ends immediately before this object. Check for one
		// ending two lines back.
		return c.groups[fl{Filename: p.Filename, Line: p.Line - 2}].Text()
	}

	// A comment group ends immediately before this object. Check for another
	// one ending two lines back from where it starts.
	start := c.fset.Position(g.List[0].Slash)
	return c.groups[fl{Filename: start.Filename, Line: start.Line - 2}].Text()
}

type Markers map[string][]string

func ParseMarkers(comment string) Markers {
	return ParseMarkersWithPrefix(DefaultMarkerPrefix, comment)
}

func ParseMarkersWithPrefix(prefix, comment string) Markers {
	m := map[string][]string{}

	for _, line := range strings.Split(comment, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if !strings.HasPrefix(line, prefix) {
			continue
		}
		kv := strings.SplitN(line[len(prefix):], "=", 2)
		k, v := kv[0], ""
		if len(kv) > 1 {
			v = kv[1]
		}
		m[k] = append(m[k], v)
	}

	return m
}

func (m Markers) IsTrue(key string) bool {
	if len(m[key]) != 1 {
		return false
	}
	return strings.EqualFold(m[key][0], "true")
}

package tree

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/afero"
	"path/filepath"
	"sort"
	"strings"
)

const (
	connector       = "├── "
	nextPrefix      = "│   "
	connectorChild  = "└── "
	nextPrefixChild = "    "
)

type Node struct {
	Name     string  `json:"name"`
	Children []*Node `json:"children,omitempty"`
}

type Tree struct {
	fs      afero.Fs
	root    *Node
	path    string
	exclude []string
}

func NewTree(fs afero.Fs, path string, exclude ...string) *Tree {
	return &Tree{
		fs:      fs,
		path:    path,
		root:    &Node{Name: filepath.Base(path)},
		exclude: exclude,
	}
}

func (t *Tree) MakeTree() error {
	if t.fs == nil {
		return errors.New("nil filesystem")
	}
	return t.buildNode(t.path, t.root)
}

func (t *Tree) ToJSON() (string, error) {
	data, err := json.MarshalIndent(t.root, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (t *Tree) ToMarkdown() string {
	var b strings.Builder
	t.writeMarkdown(&b, t.root, "")
	return b.String()
}

func (t *Tree) ToString() string {
	var b strings.Builder
	t.writeTreeFormat(&b, t.root, "", true)
	return b.String()
}

func (t *Tree) writeTreeFormat(b *strings.Builder, node *Node, prefix string, isRoot bool) {
	if isRoot {
		b.WriteString(".\n")
	}

	for i, child := range node.Children {
		isLast := i == len(node.Children)-1
		conn := connector
		next := nextPrefix

		if isLast {
			conn = connectorChild
			next = nextPrefixChild
		}

		_, _ = fmt.Fprintf(b, "%s%s%s\n", prefix, conn, child.Name)

		if len(child.Children) > 0 {
			t.writeTreeFormat(b, child, prefix+next, false)
		}
	}
}

func (t *Tree) writeMarkdown(b *strings.Builder, node *Node, prefix string) {
	_, _ = fmt.Fprintf(b, "%s- %s\n", prefix, node.Name)
	for _, child := range node.Children {
		t.writeMarkdown(b, child, "  "+prefix)
	}
}

func (t *Tree) buildNode(path string, parent *Node) error {
	entries, err := afero.ReadDir(t.fs, path)
	if err != nil {
		return err
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].IsDir() != entries[j].IsDir() {
			return entries[i].IsDir()
		}
		return entries[i].Name() < entries[j].Name()
	})

	for _, entry := range entries {
		if t.shouldExclude(entry.Name()) {
			continue
		}

		child := &Node{Name: entry.Name()}
		parent.Children = append(parent.Children, child)

		if entry.IsDir() {
			err := t.buildNode(filepath.Join(path, entry.Name()), child)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (t *Tree) shouldExclude(name string) bool {
	for _, ex := range t.exclude {
		if ex == name {
			return true
		}
	}
	return false
}

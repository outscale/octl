package tree

import (
	"fmt"
	"io"
)

const (
	emptySpace   = "   "
	middleItem   = "├─ "
	continueItem = "│  "
	lastItem     = "└─ "
)

type Tree interface {
	String() string
	Children() []Tree
}

type line struct {
	prefix string
}

func WriteTo(t Tree, w io.Writer) error {
	if _, err := fmt.Fprintln(w, t); err != nil {
		return err
	}
	return line{}.writeChildren(t, w)
}

func (l line) writeChildren(t Tree, w io.Writer) error {
	children := t.Children()
	for i, c := range children {
		last := i == len(children)-1
		prefix := l.prefix
		next := line{}
		if last {
			prefix += lastItem
			next.prefix = l.prefix + emptySpace
		} else {
			prefix += middleItem
			next.prefix = l.prefix + continueItem
		}
		if _, err := fmt.Fprintln(w, prefix+c.String()); err != nil {
			return err
		}
		if err := next.writeChildren(c, w); err != nil {
			return err
		}
	}
	return nil
}

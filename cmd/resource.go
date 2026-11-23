package cmd

import (
	"context"
	"strings"

	"github.com/outscale/octl/pkg/debug"
	"github.com/outscale/octl/pkg/tree"
	"github.com/samber/lo"
)

type ResourceType string

const Text ResourceType = "text"

type (
	deleteFunc func(ctx context.Context) error

	Resource struct {
		Type ResourceType
		ID   string
		Name string

		children []*Resource

		delete deleteFunc
	}
)

func (r *Resource) String() string {
	b := strings.Builder{}

	switch r.Type {
	case Text:
		_, _ = b.WriteString(r.ID)
	default:
		_, _ = b.WriteString(string(r.Type))
		_, _ = b.WriteString("/")
		_, _ = b.WriteString(r.ID)
		if r.Name != "" {
			_, _ = b.WriteString(" (")
			_, _ = b.WriteString(r.Name)
			_, _ = b.WriteString(")")
		}
	}
	return b.String()
}

func (r *Resource) Children() []tree.Tree {
	return lo.Map(r.children, func(c *Resource, _ int) tree.Tree {
		return tree.Tree(c)
	})
}

func (r *Resource) addToChild(childID string, childType ResourceType, subr *Resource) bool {
	child, found := lo.Find(r.children, func(r *Resource) bool { return r.ID == childID && r.Type == childType })
	if found {
		child.children = append(child.children, subr)
		return true
	}
	for _, c := range r.children {
		if c.addToChild(childID, childType, subr) {
			return true
		}
	}
	debug.Println("parent", childType, childID, "not found for", subr.Type, subr.ID)
	return false
}

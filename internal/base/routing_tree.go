package base

import (
	"errors"
	"fmt"
	"strings"
)

var errPrefixNotFound = errors.New("item with specified prefix not found")

// RoutingTreeSearchFunc is a function that executes inside of APITree.Find()
// It receives a node pointer, value pointer and current URI part being parsed.
// Return values tells if element was found
type RoutingTreeSearchFunc[T any] func(node *RoutingTree[T], value *T, key string) (*RoutingTree[T], bool)

// RoutingTree is a generic suffix tree, which can be used in HTTP and/or AMQP router
// implementations.
type RoutingTree[T any] struct {
	Branches map[string]*RoutingTree[T]
	Item     *T
}

func NewRoutingTree[T any]() RoutingTree[T] {
	return RoutingTree[T]{
		Branches: make(map[string]*RoutingTree[T]),
	}
}

// Add creates a node using URI as a path.
func (u *RoutingTree[T]) Add(uri []string, value T) (err error) {
	node := u
	for i, part := range uri {
		if b, found := node.Branches[part]; found {
			node = b
		} else {
			if node.Branches == nil {
				node.Branches = make(map[string]*RoutingTree[T])
			}
			node.Branches[part] = &RoutingTree[T]{}
			node = node.Branches[part]
		}
		if i == len(uri)-1 {
			if node.Item != nil {
				err = fmt.Errorf("routing tree collision: duplicate mapping for URI '%s'", JoinURLPath(uri))
				return
			}
			node.Item = &value
		}
	}
	return
}

// Add searches for a node with a path specified by URI and a search function provided by user.
func (u RoutingTree[T]) Find(uri []string, lazy bool, searchFunc RoutingTreeSearchFunc[T]) (value T, err error) {
	node := &u
	for _, part := range uri {
		var found bool
		node, found = searchFunc(node, &value, part)
		if !found {
			err = errPrefixNotFound
			return
		} else if lazy && node.Item != nil {
			return
		}
	}
	if node.Item != nil {
		value = *node.Item
		return
	}
	err = errPrefixNotFound
	return
}

func SplitURLPath(path string) []string {
	urlWithoutParameters := strings.Split(path, "?")[0]
	return strings.Split(urlWithoutParameters, "/")
}

func JoinURLPath(path []string) string {
	return strings.Join(path, "/")
}

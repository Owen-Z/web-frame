package gee

import (
	"fmt"
	"strings"
)

type node struct {
	pattern  string  // wait for matching route	  e.g. /p/:lang
	part     string  // part of route  e.g. :lang
	children []*node // child node	e.g. [doc, tutorial, intro]
	isWild   bool    // the node needed to match vaguely	e.g. part contain : or * 时为true
}

func (n *node) String() string {
	return fmt.Sprintf("node{pattern=%s, part=%s, isWild=%t}", n.pattern, n.part, n.isWild)
}

// find the first node match the part
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

func (n *node) matchChildren(part string) []*node {
	// slice need to check out the capacity is enough or not
	// if the capacity is not enough, the cap should be increased
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// insert insert pattern by order
func (n *node) insert(pattern string, parts []string, height int) {
	// reach the bottom of the tree
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	// check if the node is existed
	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

// search the node by the selected path
func (n *node) search(parts []string, height int) *node {
	// check the node is the bottom node or the node has prefixed "*"
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	// find the vague children
	children := n.matchChildren(part)

	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}

	return nil
}

func (n *node) travel(list *([]*node)) {
	if n.pattern != "" {
		*list = append(*list, n)
	}
	for _, child := range n.children {
		child.travel(list)
	}
}

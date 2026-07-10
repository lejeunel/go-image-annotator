package table

import (
	. "maragu.dev/gomponents"
)

type Cell struct {
	Content    Node
	ExtraAttr  string
	ExtraClass string
}

type CellOption func(*Cell)

func WithCellAttr(a string) CellOption {
	return func(c *Cell) {
		c.ExtraAttr = a
	}
}

func WithCellClass(class string) CellOption {
	return func(c *Cell) {
		c.ExtraClass = class
	}
}

func NewCell(n Node, opts ...CellOption) Cell {
	c := &Cell{Content: n}
	for _, opt := range opts {
		opt(c)
	}
	return *c
}

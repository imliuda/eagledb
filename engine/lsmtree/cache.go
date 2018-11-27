package lsmtree

import (
	"github.com/eagledb/eagledb/point"
)

type Cache struct {
}

func NewCache() *Cache {
	return &Cache{}
}

func (c *Cache) Add([]*point.Point) error {
	return nil
}

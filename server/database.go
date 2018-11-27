package server

import (
	"github.com/eagledb/eagledb/engine"
	"github.com/eagledb/eagledb/point"
)

type Database struct {
	name   string
	engine engine.Engine
}

func (db *Database) WritePoints(points []*point.Point) error {
	return db.engine.WritePoints(points)
}

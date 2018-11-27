package engine

import (
	"errors"
	"github.com/eagledb/eagledb/engine/lsmtree"
	"github.com/eagledb/eagledb/point"
	"reflect"
)

var engines = map[string]reflect.Type{}

func NewEngine(name string) (Engine, error) {
	if engine, ok := engines[name]; ok {
		return reflect.New(engine).Interface().(Engine), nil
	}
	return nil, errors.New("no such engine")
}

type Engine interface {
	WritePoints(points []*point.Point) error
}

func init() {
	engines["lsmtree"] = reflect.TypeOf((*lsmtree.LsmTree)(nil))
}

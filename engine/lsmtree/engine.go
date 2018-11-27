package lsmtree

import (
	"github.com/eagledb/eagledb/point"
	"log"
	"time"
)

type LsmTree struct {
	wal       *Wal
	cache     *Cache
	compactor *Compactor
}

func (e *LsmTree) Start() error {
	e.wal = NewWal()
	e.cache = NewCache()
	e.compactor = NewCompactor()

	e.compactor.Schedule()

	return nil
}

func (e *LsmTree) WritePoints(points []*point.Point) error {
	/*
		err := e.wal.Append()
		if err != nil {
			return err
		}

		err = e.Cache.Add(points)
	*/
	log.Println(points)
	return nil
}

func (e *LsmTree) QueryPoints() ([]*point.Point, error) {

	return nil, nil
}

func (e *LsmTree) DeletePoints(mesurement string, tags map[string]string, start, end time.Time) {
}

func (e *LsmTree) CreateSnapShot(id string) {

}

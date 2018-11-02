package lsmtree

import (
	"github.com/eagledb/eagledb/model"
)

type Engine struct {
	wal       *Wal
	cache     *Cache
	compactor *Compactor
	config    *Config
	index     *Index
	store     *Store
}

func (e *Engine) Start() error {
	wal, err := NewWal(e.config.Wal)
	if err != nil {
		return err
	}
	e.wal = wal

	cache, err := NewCache(e.config.Cache)
	if err != nil {
		return nil
	}
	e.cache = cache

	e.compactor = NewCompactor(e.config.Compactor)
	e.compactor.Start()
}

func (e *Engine) AddPoint(points []*model.Point) error {
	err := e.wal.Append()
	if err != nil {
		return err
	}

	err = e.Cache.Add(points)

}

func (e *Engine) QueryPoint() ([]*mode.Point, error) {

}

func (e *Engine) DeletePoint(mesurement string, tags map[string]string, start, end time.Time) {
}

func (e *Engine) CreateSnapShot(id string) {

}

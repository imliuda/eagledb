package model

type Engine interface {
	WritePoint([]*Point) error
}

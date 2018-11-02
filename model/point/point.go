package model

import (
	"sort"
	"time"
)

func escape(string) string {

}

func unescape(string) string {

}

type TimeValue struct {
	time  time.Time
	value interface{}
}

type Tag struct {
	Key   []byte
	Value []byte
}

type Tags []Tag

type Point struct {
	mesurement string
	tags       Tags
	fields     map[string]interface{}
	time       time.Time
}

func (p *Point) Name() []byte {
	return p.name
}

func (p *Point) SetName(name string) {
	p.name = name
}

func (p *Point) Tags() map[string]string {
	return p.tags
}

func (p *Point) AddTag(key, value string) {
	p.tags[key] = value
}

func (p *Point) Fields() map[string]interface{} {
	return p.fields
}

func (p *Point) Time() time.Time {
	return p.time
}

func (p *Point) Key() []byte {
	key := []byte(p.mesurement)
	keys := []string{}
	for k, v := range p.tags {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		key += []byte(k) + []byte(p.tags[k])
	}
	return key
}

func (p *Point) String() string {
	s := p.mesurement
	for k, v := range p.tags {
		s += "," + k + "=" + v
	}
	s += " "
	for k, v := range p.feilds {
		s += k + "=" + v + ","
	}
	s = s[:len(s)-1]
	s += strconv.FormatInt(p.time.UnixNano(), 10)
	return s
}

func NewPoint() {
	return &Point{}
}

func ParsePoint(buffer []byte) ([]*Point, error) {

}

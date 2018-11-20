package point

import (
	"fmt"
	"log"
	"bytes"
	"strconv"
)

type FieldType int

const (
	Integer = iota
	Float
	String
	Boolean
	Null
)

type Tag struct {
	Key   []byte
	Value []byte
}

type Field struct {
	Key []byte
	Type  FieldType
	Value interface{}
}

type Point struct {
	name   []byte
	tags   []Tag
	fields []Field
	time   int64
}

func NewPoint() *Point {
	return &Point{}
}

const (
	PARSE_START = iota
	PARSE_NAME
	PARSE_TAG_KEY
	PARSE_TAG_VALUE
	PARSE_FIELD_KEY_START
	PARSE_FIELD_KEY
	PARSE_FIELD_VALUE
	PARSE_FIELD_VALUE_STRING
	PARSE_FIELD_VALUE_NUMBER
	PARSE_FIELD_VALUE_BOOLEAN_TR
	PARSE_FIELD_VALUE_BOOLEAN_TRU
	PARSE_FIELD_VALUE_BOOLEAN_TRUE
	PARSE_FIELD_VALUE_BOOLEAN_FA
	PARSE_FIELD_VALUE_BOOLEAN_FAL
	PARSE_FIELD_VALUE_BOOLEAN_FALS
	PARSE_FIELD_VALUE_BOOLEAN_FALSE
	PARSE_FIELD_VALUE_NU
	PARSE_FIELD_VALUE_NUL
	PARSE_FIELD_VALUE_NULL
	PARSE_FIELD_VALUE_DONE
	PARSE_TIME_START
	PARSE_TIME
	PARSE_DONE
)

func escape(value []byte, chars []byte) []byte {
	r := make([]byte, 0, len(value))
	for _, c := range value {
		for _, e := range chars {
			if c == e || c == byte('\\'){
				r = append(r, '\\')
			}
		}
		r = append(r, c)
	}
	return r
}

func unescape(value []byte, chars []byte) []byte {
	r := make([]byte, 0, len(value))
	vlen := len(value)
	for i, c := range value {
		if c == byte('\\') {
			if i == vlen - 1 {
				r = append(r, c)
			} else {
				for _, e := range chars {
					if value[i+1] == e || value[i+1] == byte('\\') {
						continue
					} else {
						r = append(r, c)
					}
				}
			}
		} else {
			r = append(r, c)
		}
	}
	return r
}

func escapeStringValue (value []byte) []byte {
	return value
}

func Parses(buffer string) ([]*Point, error) {
	return Parse([]byte(buffer))
}

func Parse(buffer []byte) ([]*Point, error) {
	points := []*Point{}

	var point *Point
	var i int
	var c byte
	var key []byte
	var buflen = len(buffer)
	var start = 0
	var stat = PARSE_START
	for i, c = range buffer {
		log.Println(i, c)
		if stat == PARSE_START {
			if c == byte(' ') || c == byte('\n') {
				continue
			} else {
				point = &Point{}
				start = i
				stat = PARSE_NAME
			}
		} else if stat == PARSE_NAME {
			if c == byte('\n') {
				goto line_error
			} else if c == byte(' ') && buffer[i-1] != byte('\\') {
				name := buffer[start:i]
				name = escape(name, []byte(", "))
				log.Println(name)
				point.SetName(name)
				stat = PARSE_FIELD_KEY_START
			} else if c == byte(',') && buffer[i-1] != byte('\\') {
				name := buffer[start:i]
				name = escape(name, []byte(", "))
				point.SetName(name)
				start = i + 1
				stat = PARSE_TAG_KEY
			} else {
				continue
			}
		} else if stat == PARSE_TAG_KEY {
			if c == byte('\n') {
				goto line_error
			} else if c == byte(' ') && buffer[i-1] != byte('\\') {
				goto space_error
			} else if c == byte('=') && buffer[i-1] != byte('\\') {
				key := buffer[start:i]
				key = escape(key, []byte("= "))
				key = []byte(key)
				start = i + 1
				stat = PARSE_TAG_VALUE
			} else {
				continue
			}
		} else if stat == PARSE_TAG_VALUE {
			if c == byte('\n') {
				goto line_error
			} else if c == byte(' ') && buffer[i-1] != byte('\\') {
				value := buffer[start:i]
				value = escape(value, []byte(",="))
				point.AddTag(key, value)
				stat = PARSE_FIELD_KEY_START
			} else if c == byte(',') && buffer[i-1] != byte('\\') {
				value := buffer[start:i]
				value = escape(value, []byte(",="))
				point.AddTag(key, value)
				stat = PARSE_TAG_KEY
			} else {
				continue
			}
		} else if stat == PARSE_FIELD_KEY_START {
			if c == byte('\n') {
				goto line_error
			} else if c != byte(' ') {
				start = i
				stat = PARSE_FIELD_KEY
			} else {
				continue
			}
		} else if stat == PARSE_FIELD_KEY {
			if c == byte('\n') {
				goto line_error
			} else if c == byte(' ') && buffer[i-1] != byte('\\') {
				goto space_error
			} else if c == byte('=') && buffer[i-1] != byte('\\') {
				key = buffer[start:i]
				key = escape([]byte(key), []byte("= "))
				start = i + 1
				stat = PARSE_FIELD_VALUE
			} else {
				continue
			}
		} else if stat == PARSE_FIELD_VALUE {
			if c == byte('"') {
				stat = PARSE_FIELD_VALUE_STRING
			} else if c >= byte('0') && c <= byte('9') || c == byte('-'){
				stat = PARSE_FIELD_VALUE_NUMBER
			} else if c == byte('t') {
				stat = PARSE_FIELD_VALUE_BOOLEAN_TR
			} else if c == byte('f') {
				stat = PARSE_FIELD_VALUE_BOOLEAN_FA
			} else if c == byte('n') {
				stat = PARSE_FIELD_VALUE_NU
			} else {
				goto value_error
			}
		} else if stat == PARSE_FIELD_VALUE_STRING {
			if c == byte('"') && buffer[i-1] != byte('\\') {
				value := buffer[start:i]
				value = escapeStringValue(value)
				point.AddField(key, value, String)
				stat = PARSE_FIELD_VALUE_DONE
			} else {
				continue
			}
		} else if stat == PARSE_FIELD_VALUE_NUMBER {
			if c >= byte('0') && c <= byte('9') ||
				c == byte('E') || c == byte('e') ||
				c == byte('+') || c == byte('-') {
				continue
			} else if c == byte(' ') || c == byte(','){
				value := buffer[start:i]
				iv, err := strconv.ParseInt(string(value), 10, 64)
				if err == nil {
					point.AddField(key, iv, Integer)
				} else {
					fv, err := strconv.ParseFloat(string(value), 64)
					if err != nil {
						goto value_error
					} else {
						point.AddField(key, fv, Float)
					}
				}
				if c == byte(',') {
					stat = PARSE_FIELD_KEY
				} else if c == byte(' ') {
					stat = PARSE_FIELD_VALUE_DONE
				}
			} else {
				goto value_error
			}
		} else if stat == PARSE_FIELD_VALUE_BOOLEAN_TR {
			if c == byte('r') {
				stat = PARSE_FIELD_VALUE_BOOLEAN_TRU
			} else {
				goto value_error
			}
		} else if stat == PARSE_FIELD_VALUE_BOOLEAN_TRU {
			if c == byte('u') {
				stat = PARSE_FIELD_VALUE_BOOLEAN_TRUE
			} else {
				goto value_error
			}
		} else if stat == PARSE_FIELD_VALUE_BOOLEAN_TRUE {
			if c == byte('e') {
				point.AddField(key, true, Boolean)
				stat = PARSE_FIELD_VALUE_DONE
			} else {
				goto value_error
			}
		} else if stat == PARSE_FIELD_VALUE_BOOLEAN_FA {
			if c == byte('a') {
				stat = PARSE_FIELD_VALUE_BOOLEAN_FAL
			} else {
				goto value_error
			}
		} else if stat == PARSE_FIELD_VALUE_BOOLEAN_FAL {
			if c == byte('l') {
				stat = PARSE_FIELD_VALUE_BOOLEAN_FALS
			} else {
				goto value_error
			}
		} else if stat == PARSE_FIELD_VALUE_BOOLEAN_FALS {
			if c == byte('s') {
				stat = PARSE_FIELD_VALUE_BOOLEAN_FALSE
			} else {
				goto value_error
			}
		} else if stat == PARSE_FIELD_VALUE_BOOLEAN_FALSE {
			if c == byte('e') {
				point.AddField(key, false, Boolean)
				stat = PARSE_FIELD_VALUE_DONE
			} else {
				goto value_error
			}
		} else if stat == PARSE_FIELD_VALUE_DONE {
			if c == ',' {
				stat = PARSE_FIELD_KEY
			} else if c == byte(' ') {
				stat = PARSE_TIME_START
			} else {
				stat = PARSE_TIME
			}
		} else if stat == PARSE_TIME_START {
			if c == byte('\n') {
				goto line_error
			} else if c != byte(' ') {
				start = i
				stat = PARSE_TIME
			} else {
				continue
			}
		} else if stat == PARSE_TIME {
			log.Println("time", i, c)
			if c < '0' && c > '9' || i == buflen - 1{
				stime := buffer[start:i]
				itime, err := strconv.ParseInt(string(stime), 10, 64)
				if err != nil {
					log.Println(stime, err)
					goto time_error
				}
				point.SetTime(itime)
				points = append(points, point)
				stat = PARSE_START
			} else {
				continue
			}
		}
	}

	if stat != PARSE_START {
		goto partial_error
	}

	return points, nil

partial_error:
	return nil, fmt.Errorf("uncompleted point")
line_error:
	return nil, fmt.Errorf("invalid character '\n' at index %d", i)
space_error:
	return nil, fmt.Errorf("invalid character ' ' at index %d", i)
value_error:
	return nil, fmt.Errorf("invalid value")
time_error:
	return nil, fmt.Errorf("invalid time")
}

func (p *Point) Name() []byte {
	return p.name
}

func (p *Point) SetName(name []byte) {
	p.name = name
}

func (p *Point) Tags() []Tag {
	return p.tags
}

func (p *Point) AddTag(key, value []byte) {
	p.tags = append(p.tags, Tag{Key: key, Value: value})
}

func (p *Point) Fields() []Field {
	return p.fields
}

func (p *Point) AddField(key []byte, value interface{}, vtype FieldType) {

}

func (p *Point) SetTime(tm int64) {
	p.time = tm
}

func (p *Point) Time() int64 {
	return p.time
}

func (p *Point) Key() []byte {
	key := p.name
	for _, tag := range p.tags {
		key = append(key, tag.Key...)
		key = append(key, tag.Value...)
	}
	return key
}

func (p *Point) String() string {
	buf := bytes.Buffer{}
	buf.Write(p.name)
	for _, tag := range p.tags {
		buf.WriteRune(',')
		buf.Write(tag.Key)
		buf.WriteRune('=')
		buf.Write(tag.Value)
	}
	buf.WriteRune(' ')
	for _, field := range p.fields {
		buf.Write(field.Key)
		buf.WriteRune('=')
		switch field.Type {
			case Integer:
				buf.WriteString(strconv.FormatInt(field.Value.(int64), 10))
			case Float:
				buf.WriteString(strconv.FormatFloat(field.Value.(float64), 'f', -1, 64))
			case String:
				buf.WriteString(field.Value.(string))
			case Boolean:
				if field.Value.(bool) {
					buf.WriteString("true")
				} else {
					buf.WriteString("false")
				}
			case Null:
				buf.WriteString("null")
			default:
				panic("unknown value type of point")
		}
		buf.WriteRune(',')
	}
	buf.Truncate(buf.Len()-1)
	buf.WriteString(strconv.FormatInt(p.time, 10))
	return buf.String()
}

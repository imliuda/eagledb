package point

import (
	"bytes"
	"fmt"
	"log"
	"sort"
	"strconv"
)

type ValueType int

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
	Key   []byte
	Type  ValueType
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

var (
	WrongValueTypeError = fmt.Errorf("wrong value type")
)

func escape(value []byte, chars []byte) []byte {
	r := make([]byte, 0, len(value))
	for _, c := range value {
		for _, e := range chars {
			if c == e || c == byte('\\') {
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
			if i == vlen-1 {
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

func escapeStringValue(value []byte) []byte {
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
		//log.Println(i, c)
		if stat == PARSE_START {
			if c == ' ' || c == '\n' {
				continue
			} else {
				point = &Point{}
				start = i
				stat = PARSE_NAME
			}
		} else if stat == PARSE_NAME {
			if c == '\n' {
				goto line_error
			} else if c == ' ' && buffer[i-1] != '\\' {
				name := buffer[start:i]
				name = escape(name, []byte(", "))
				log.Println(name)
				point.SetName(name)
				stat = PARSE_FIELD_KEY_START
			} else if c == ',' && buffer[i-1] != '\\' {
				name := buffer[start:i]
				name = escape(name, []byte(", "))
				point.SetName(name)
				start = i + 1
				stat = PARSE_TAG_KEY
			} else {
				continue
			}
		} else if stat == PARSE_TAG_KEY {
			if c == '\n' {
				goto line_error
			} else if c == ' ' && buffer[i-1] != '\\' {
				goto space_error
			} else if c == '=' && buffer[i-1] != '\\' {
				key = buffer[start:i]
				key = escape(key, []byte("= "))
				start = i + 1
				stat = PARSE_TAG_VALUE
			} else {
				continue
			}
		} else if stat == PARSE_TAG_VALUE {
			if c == '\n' {
				goto line_error
			} else if c == ' ' && buffer[i-1] != '\\' {
				value := buffer[start:i]
				value = escape(value, []byte(",="))
				point.AddTag(key, value)
				stat = PARSE_FIELD_KEY_START
			} else if c == ',' && buffer[i-1] != '\\' {
				value := buffer[start:i]
				value = escape(value, []byte(",="))
				point.AddTag(key, value)
				start = i + 1
				stat = PARSE_TAG_KEY
			} else {
				continue
			}
		} else if stat == PARSE_FIELD_KEY_START {
			if c == '\n' {
				goto line_error
			} else if c != ' ' {
				start = i
				stat = PARSE_FIELD_KEY
			} else {
				continue
			}
		} else if stat == PARSE_FIELD_KEY {
			if c == '\n' {
				goto line_error
			} else if c == ' ' && buffer[i-1] != '\\' {
				goto space_error
			} else if c == '=' && buffer[i-1] != '\\' {
				key = buffer[start:i]
				key = escape(key, []byte("= "))
				start = i + 1
				stat = PARSE_FIELD_VALUE
			} else {
				continue
			}
		} else if stat == PARSE_FIELD_VALUE {
			if c == '"' {
				stat = PARSE_FIELD_VALUE_STRING
			} else if c >= '0' && c <= '9' || c == '-' {
				stat = PARSE_FIELD_VALUE_NUMBER
			} else if c == 't' {
				stat = PARSE_FIELD_VALUE_BOOLEAN_TR
			} else if c == 'f' {
				stat = PARSE_FIELD_VALUE_BOOLEAN_FA
			} else if c == 'n' {
				stat = PARSE_FIELD_VALUE_NU
			} else {
				goto value_error
			}
		} else if stat == PARSE_FIELD_VALUE_STRING {
			if c == '"' && buffer[i-1] != '\\' {
				value := buffer[start:i]
				value = escapeStringValue(value)
				point.AddStringField(key, value)
				stat = PARSE_FIELD_VALUE_DONE
			} else {
				continue
			}
		} else if stat == PARSE_FIELD_VALUE_NUMBER {
			log.Println("number", i, c)
			if c >= '0' && c <= '9' ||
				c == 'E' || c == 'e' ||
				c == '+' || c == '-' {
				continue
			} else if c == ' ' || c == ',' {
				value := buffer[start:i]
				iv, err := strconv.ParseInt(string(value), 10, 64)
				if err == nil {
					point.AddIntegerField(key, iv)
				} else {
					fv, err := strconv.ParseFloat(string(value), 64)
					if err != nil {
						goto value_error
					} else {
						point.AddFloatField(key, fv)
					}
				}
				if c == ',' {
					start = i + 1
					stat = PARSE_FIELD_KEY
				} else if c == ' ' {
					stat = PARSE_FIELD_VALUE_DONE
				}
			} else {
				goto value_error
			}
		} else if stat == PARSE_FIELD_VALUE_BOOLEAN_TR {
			if c == 'r' {
				stat = PARSE_FIELD_VALUE_BOOLEAN_TRU
			} else {
				goto value_error
			}
		} else if stat == PARSE_FIELD_VALUE_BOOLEAN_TRU {
			if c == 'u' {
				stat = PARSE_FIELD_VALUE_BOOLEAN_TRUE
			} else {
				goto value_error
			}
		} else if stat == PARSE_FIELD_VALUE_BOOLEAN_TRUE {
			if c == 'e' {
				point.AddBooleanField(key, true)
				stat = PARSE_FIELD_VALUE_DONE
			} else {
				goto value_error
			}
		} else if stat == PARSE_FIELD_VALUE_BOOLEAN_FA {
			if c == 'a' {
				stat = PARSE_FIELD_VALUE_BOOLEAN_FAL
			} else {
				goto value_error
			}
		} else if stat == PARSE_FIELD_VALUE_BOOLEAN_FAL {
			if c == 'l' {
				stat = PARSE_FIELD_VALUE_BOOLEAN_FALS
			} else {
				goto value_error
			}
		} else if stat == PARSE_FIELD_VALUE_BOOLEAN_FALS {
			if c == 's' {
				stat = PARSE_FIELD_VALUE_BOOLEAN_FALSE
			} else {
				goto value_error
			}
		} else if stat == PARSE_FIELD_VALUE_BOOLEAN_FALSE {
			if c == 'e' {
				point.AddBooleanField(key, false)
				stat = PARSE_FIELD_VALUE_DONE
			} else {
				goto value_error
			}
		} else if stat == PARSE_FIELD_VALUE_NU {
			if c == 'u' {
				stat = PARSE_FIELD_VALUE_NUL
			} else {
				goto value_error
			}
		} else if stat == PARSE_FIELD_VALUE_NUL {
			if c == 'l' {
				stat = PARSE_FIELD_VALUE_NULL
			} else {
				goto value_error
			}
		} else if stat == PARSE_FIELD_VALUE_NULL {
			if c == 'l' {
				point.AddNullField(key)
				stat = PARSE_FIELD_VALUE_DONE
			} else {
				goto value_error
			}
		} else if stat == PARSE_FIELD_VALUE_DONE {
			if c == '\n' {
				goto line_error
			} else if c == ',' {
				start = i + 1
				stat = PARSE_FIELD_KEY
			} else if c == ' ' {
				stat = PARSE_TIME_START
			} else {
				start = i
				stat = PARSE_TIME
			}
		} else if stat == PARSE_TIME_START {
			if c == '\n' {
				goto line_error
			} else if c != ' ' {
				start = i
				stat = PARSE_TIME
			} else {
				continue
			}
		} else if stat == PARSE_TIME {
			log.Println("time", i, c)
			if c < '0' && c > '9' || i == buflen-1 {
				var stime []byte
				if i == buflen-1 {
					stime = buffer[start : i+1]
				} else {
					stime = buffer[start:i]
				}
				itime, err := strconv.ParseInt(string(stime), 10, 64)
				if err != nil {
					log.Println(string(stime), err)
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
	tag := Tag{
		Key:   key,
		Value: value,
	}

	p.tags = append(p.tags, tag)

	sort.Slice(p.tags, func(i, j int) bool {
		ki := append(p.tags[i].Key, p.tags[i].Value...)
		kj := append(p.tags[j].Key, p.tags[j].Value...)
		return bytes.Compare(ki, kj) < 0
	})
}

func (p *Point) Fields() []Field {
	return p.fields
}

func (p *Point) AddIntegerField(key []byte, value int64) {
	p.fields = append(p.fields, Field{
		Key:   key,
		Value: value,
		Type:  Integer,
	})
}

func (p *Point) AddFloatField(key []byte, value float64) {
	p.fields = append(p.fields, Field{
		Key:   key,
		Value: value,
		Type:  Float,
	})
}

func (p *Point) AddStringField(key []byte, value []byte) {
	p.fields = append(p.fields, Field{
		Key:   key,
		Value: value,
		Type:  String,
	})
}

func (p *Point) AddBooleanField(key []byte, value bool) {
	p.fields = append(p.fields, Field{
		Key:   key,
		Value: value,
		Type:  Boolean,
	})
}

func (p *Point) AddNullField(key []byte) {
	p.fields = append(p.fields, Field{
		Key:   key,
		Value: nil,
		Type:  Null,
	})
}

func (p *Point) SetTime(tm int64) {
	p.time = tm
}

func (p *Point) Time() int64 {
	return p.time
}

func (p *Point) SeriesKey() []byte {
	key := bytes.Buffer{}
	key.Write(p.name)
	for _, tag := range p.tags {
		key.Write(tag.Key)
		key.Write(tag.Value)
	}
	return key.Bytes()
}

func (p *Point) String() string {
	buf := bytes.Buffer{}
	buf.Write(p.name)
	for _, tag := range p.tags {
		buf.WriteByte(',')
		buf.Write(tag.Key)
		buf.WriteByte('=')
		buf.Write(tag.Value)
	}
	buf.WriteByte(' ')
	for _, field := range p.fields {
		buf.Write(field.Key)
		buf.WriteByte('=')
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
		buf.WriteByte(',')
	}
	buf.Truncate(buf.Len() - 1)
	buf.WriteByte(' ')
	buf.WriteString(strconv.FormatInt(p.time, 10))
	return buf.String()
}

type FieldIterator struct {
	seriesKey []byte
	fields    []Field
	index     int
}

func NewFieldIterator(point *Point) *FieldIterator {
	iter := FieldIterator{}
	iter.seriesKey = point.SeriesKey()
	iter.fields = point.Fields()
	iter.index = -1
	log.Println(string(iter.seriesKey))
	return &iter
}

func (i *FieldIterator) Reset() {
	i.index = -1
}

func (i *FieldIterator) Iterate() bool {
	if i.index == len(i.fields)-1 {
		return false
	} else {
		i.index += 1
		return true
	}
}

func (i *FieldIterator) Key() []byte {
	return i.fields[i.index].Key
}

func (i FieldIterator) Type() ValueType {
	return i.fields[i.index].Type
}

func (i FieldIterator) Value() interface{} {
	return i.fields[i.index].Value
}

func (i *FieldIterator) IntegerValue() (int64, error) {
	if i.fields[i.index].Type != Integer {
		return 0, WrongValueTypeError
	}
	return i.fields[i.index].Value.(int64), nil
}

func (i *FieldIterator) FloatValue() (float64, error) {
	if i.fields[i.index].Type != Float {
		return 0, WrongValueTypeError
	}
	return i.fields[i.index].Value.(float64), nil
}

func (i *FieldIterator) StringValue() ([]byte, error) {
	if i.fields[i.index].Type != String {
		return []byte{}, WrongValueTypeError
	}
	return i.fields[i.index].Value.([]byte), nil
}

func (i *FieldIterator) BooleanValue() (bool, error) {
	if i.fields[i.index].Type != Boolean {
		return false, WrongValueTypeError
	}
	return i.fields[i.index].Value.(bool), nil
}

package point

import (
	"bytes"
	"fmt"
	"sort"
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
	Key   []byte
	Value interface{}
	Type  FieldType
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
	PARSE_DONE_LF
	PARSE_DONE
)

func escape(value []byte, chars []byte) []byte {
	r := make([]byte, 0, len(value))
	for _, c := range value {
		for _, e := range chars {
			if c == e || c == '\\' {
				r = append(r, '\\')
			}
		}
		r = append(r, c)
	}
	return r
}

func unescape(value []byte, chars []byte) ([]byte, error) {
	r := make([]byte, 0, len(value))
	escape := false
	for _, c := range value {
		if c == '\\' && !escape {
			escape = true
		} else {
			if escape {
				found := false
				for _, e := range chars {
					if c == e || c == '\\' {
						found = true
						break
					}
				}
				if !found {
					return nil, fmt.Errorf("unknown escape character %c", c)
				}
				r = append(r, c)
				escape = false
			} else {
				r = append(r, c)
			}
		}
	}
	return r, nil
}

func escapeName(value []byte) []byte {
	return escape(value, []byte(", "))
}

func unescapeName(value []byte) ([]byte, error) {
	return unescape(value, []byte(", "))
}

func escapeTagKey(value []byte) []byte {
	return escape(value, []byte("= "))
}

func unescapeTagKey(value []byte) ([]byte, error) {
	return unescape(value, []byte("= "))
}

func escapeTagValue(value []byte) []byte {
	return escape(value, []byte(", "))
}

func unescapeTagValue(value []byte) ([]byte, error) {
	return unescape(value, []byte(", "))
}

func escapeFieldKey(value []byte) []byte {
	return escape(value, []byte("= "))
}

func unescapeFieldKey(value []byte) ([]byte, error) {
	return unescape(value, []byte("= "))
}

func escapeStringValue(value []byte) []byte {
	buffer := &bytes.Buffer{}
	buffer.Grow(len(value))

	for _, c := range value {
		if c == '"' {
			buffer.WriteString(`\"`)
		} else if c == '\\' {
			buffer.WriteString(`\\`)
		} else if c == '/' {
			buffer.WriteString(`\/`)
		} else if c == '\b' {
			buffer.WriteString(`\b`)
		} else if c == '\f' {
			buffer.WriteString(`\f`)
		} else if c == '\n' {
			buffer.WriteString(`\n`)
		} else if c == '\r' {
			buffer.WriteString(`\r`)
		} else if c == '\t' {
			buffer.WriteString(`\t`)
		} else {
			buffer.WriteByte(c)
		}
	}

	return buffer.Bytes()
}

func unescapeStringValue(value []byte) ([]byte, error) {
	vlen := len(value)
	buffer := &bytes.Buffer{}
	buffer.Grow(vlen)
	for i := 0; i < vlen; i++ {
		if value[i] == '\\' {
			i++
			if i == vlen {
				return nil, fmt.Errorf("end with invalid '\\'")
			}
			if value[i] == '"' {
				buffer.WriteByte('"')
			} else if value[i] == '\\' {
				buffer.WriteByte('\\')
			} else if value[i] == '/' {
				buffer.WriteByte('/')
			} else if value[i] == 'b' {
				buffer.WriteByte('\b')
			} else if value[i] == 'f' {
				buffer.WriteByte('\f')
			} else if value[i] == 'n' {
				buffer.WriteByte('\n')
			} else if value[i] == 'r' {
				buffer.WriteByte('\r')
			} else if value[i] == 't' {
				buffer.WriteByte('\t')
			} else if value[i] == 'u' {
				i++
				if i+4 > vlen {
					return nil, fmt.Errorf("expecting 4 hex digits")
				}
				code, err := strconv.ParseInt(string(value[i:i+4]), 16, 32)
				if err != nil {
					return nil, fmt.Errorf("parse code point error")
				}
				buffer.WriteRune(rune(code))
				i += 3
			}
		} else {
			buffer.WriteByte(value[i])
		}
	}

	return buffer.Bytes(), nil
}

func Parses(buffer string) ([]*Point, error) {
	return Parse([]byte(buffer))
}

func Parse(buffer []byte) ([]*Point, error) {
	points := []*Point{}

	var point *Point
	var err error
	var i int
	var c byte
	var key []byte
	var buflen = len(buffer)
	var start = 0
	var stat = PARSE_START

	for i, c = range buffer {
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
			} else if (c == ',' || c == ' ') && buffer[i-1] != '\\' {
				name := buffer[start:i]
				name, err = unescapeName(name)
				if err != nil {
					goto unescape_error
				}
				point.SetName(name)
				if c == ' ' {
					stat = PARSE_FIELD_KEY_START
				} else {
					start = i + 1
					stat = PARSE_TAG_KEY
				}
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
				key, err = unescapeTagKey(key)
				if err != nil {
					goto unescape_error
				}
				start = i + 1
				stat = PARSE_TAG_VALUE
			} else {
				continue
			}
		} else if stat == PARSE_TAG_VALUE {
			if c == '\n' {
				goto line_error
			} else if (c == ',' || c == ' ') && buffer[i-1] != '\\' {
				value := buffer[start:i]
				value, err = unescapeTagValue(value)
				if err != nil {
					goto unescape_error
				}
				point.AddTag(key, value)
				if c == ' ' {
					stat = PARSE_FIELD_KEY_START
				} else {
					start = i + 1
					stat = PARSE_TAG_KEY
				}
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
				key, err = unescapeFieldKey(key)
				if err != nil {
					goto unescape_error
				}
				start = i + 1
				stat = PARSE_FIELD_VALUE
			} else {
				continue
			}
		} else if stat == PARSE_FIELD_VALUE {
			if c == '"' {
				start = i + 1
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
				value, err = unescapeStringValue(value)
				if err != nil {
					goto unescape_error
				}
				point.AddStringField(key, value)
				stat = PARSE_FIELD_VALUE_DONE
			} else {
				continue
			}
		} else if stat == PARSE_FIELD_VALUE_NUMBER {
			if c >= '0' && c <= '9' || c == '.' ||
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
			if c < '0' || c > '9' || i == buflen-1 {
				var stime []byte
				if i == buflen-1 {
					stime = buffer[start : i+1]
				} else {
					stime = buffer[start:i]
				}
				itime, err := strconv.ParseInt(string(stime), 10, 64)
				if err != nil {
					goto time_error
				}
				point.SetTime(itime)
				points = append(points, point)

				if c == '\n' || i == buflen-1 {
					start = i + 1
					stat = PARSE_START
				} else if c == ' ' {
					stat = PARSE_DONE
				} else {
					goto extra_error
				}
			} else {
				continue
			}
		} else if stat == PARSE_DONE {
			if c == '\n' || i == buflen-1 {
				start = i + 1
				stat = PARSE_START
			} else if c != ' ' {
				goto extra_error
			} else {
				continue
			}
		}
	}

	if stat != PARSE_START && stat != PARSE_DONE {
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
	return nil, fmt.Errorf("invalid value of '%c' near index %d", buffer[i], i)
time_error:
	return nil, fmt.Errorf("invalid time of '%c' near index %d", buffer[i], i)
extra_error:
	return nil, fmt.Errorf("extra character '%c' at %d", buffer[i], i)
unescape_error:
	return nil, fmt.Errorf("unescape error near index %d: %s", i, err)
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
	buf.Write(escapeName(p.name))
	for _, tag := range p.tags {
		buf.WriteByte(',')
		buf.Write(escapeTagKey(tag.Key))
		buf.WriteByte('=')
		buf.Write(escapeTagValue(tag.Value))
	}
	buf.WriteByte(' ')
	for _, field := range p.fields {
		buf.Write(escapeFieldKey(field.Key))
		buf.WriteByte('=')
		switch field.Type {
		case Integer:
			buf.WriteString(strconv.FormatInt(field.Value.(int64), 10))
		case Float:
			buf.WriteString(strconv.FormatFloat(field.Value.(float64), 'f', -1, 64))
		case String:
			buf.WriteByte('"')
			buf.WriteString(string(escapeStringValue(field.Value.([]byte))))
			buf.WriteByte('"')
		case Boolean:
			if field.Value.(bool) {
				buf.WriteString("true")
			} else {
				buf.WriteString("false")
			}
		case Null:
			buf.WriteString("null")
		}
		buf.WriteByte(',')
	}
	buf.Truncate(buf.Len() - 1)
	buf.WriteByte(' ')
	buf.WriteString(strconv.FormatInt(p.time, 10))
	return buf.String()
}

type FieldIterator struct {
	fields []Field
	index  int
}

func NewFieldIterator(point *Point) *FieldIterator {
	iter := FieldIterator{}
	iter.fields = point.Fields()
	iter.index = -1
	return &iter
}

func (i *FieldIterator) Reset() {
	i.index = -1
}

func (i *FieldIterator) Next() bool {
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

func (i FieldIterator) Type() FieldType {
	return i.fields[i.index].Type
}

func (i FieldIterator) Value() interface{} {
	return i.fields[i.index].Value
}

func (i *FieldIterator) IntegerValue() int64 {
	return i.fields[i.index].Value.(int64)
}

func (i *FieldIterator) FloatValue() float64 {
	return i.fields[i.index].Value.(float64)
}

func (i *FieldIterator) StringValue() []byte {
	return i.fields[i.index].Value.([]byte)
}

func (i *FieldIterator) BooleanValue() bool {
	return i.fields[i.index].Value.(bool)
}

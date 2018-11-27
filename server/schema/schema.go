package scheme

import (
	"encoding/binary"
	"os"
)

type Scheme struct {
	dir    string
	tables map[string]*Table
}

func Open(dir string) (*Scheme, error) {

}

func (s *Scheme) loadTable(path) (*Table, error) {

}

func (s *Scheme) CreateTable(name string, cols []*Column) {

}

func (s *Scheme) DropTable(name string) {

}

func (s *Scheme) GetTable(name string) (*Table, error) {

}

func (t *Table) Insert(name string, value []interface{}) error {

}

func (t *Table) Update(name string, value interface{}, conditions bool) error {

}

func (t *Table) Delete(name string, donctitinos bool) error {

}

func (t *Table) Flush(name string) error {

}

type Table struct {
	name      string
	file      *os.File
	fields    []Field
	fieldmap  map[string]FieldType
	datastart int
	buffer    []byte
}

type FieldType int

const (
	Integer = iota
	Float
	String
	Boolean
)

type Field struct {
	Name []byte
	Type FieldType
}

// scheme_count: int32
//
// scheme_name_size: int32
// scheme_name:
// scheme_type: int32
// ...
//
// entry_size: int32
// entry:
// next: +count
// ...
func NewMeta(path string, fields []Field) (*Meta, error) {
	meta := &Meta{}
	fieldMap := map[string]int{}
	for _, field := range fields {
		if field.Type < Integer || field.Type > Boolean {
			return nil, err
		}
		fieldMap[bytes.ToLower(field.Name)] = field.Type
	}
	meta.fieldmap = fieldMap

	if _, err := os.Stat(path); os.IsNotExsit(err) {
		file, err := os.Create(path, 0755)
		if err != nil {
			return nil, err
		}
		meta.file = file

		buffer := bytes.NewBuffer{}
		u32buf := [4]byte{}
		binary.LittleEndian.PutUint32(u32buf, len(fields))
		buffer.Write(u32buf)

		for fname, ftype := range fieldMap {
			binary.LittleEndian.PutUint32(u32buf, len(name))
			buffer.Write(u32buf)
			buffer.Write(fname)
			binnary.LittleEndian.PutUint32(u32buf, ftype)
			buffer.Write(u32buf)
		}
		_, err := meta.file.WriteAt(buffer.Bytes(), 0)
		if err != nil {
			return nil, err
		}
		err := meta.file.Sync()
		if err != nil {
			return nil, err
		}
	} else {
		file, err := os.Open()
		if err != nil {
			return nil, err
		}
		meta.file = file
	}

	buffer, err := ioutil.ReadAll(meta.file)
	if err != nil {
		return nil, err
	}
	meta.buffer = buffer
	index := 0
	diskFields := map[string]int{}
	nfield := binary.LittleEndian.Uint32(buffer[index : index+2])
	index += 2
	for i := 0; i < nfield; i++ {
		u32buf := [4]byte{}
		_, err := meta.file.ReadAt(u32buf, index)
		if err != nil {
			return nil, err
		}
		index += 4
		flen := binary.LittleEndian.Uint32(u32buf)

		name := make([]byte, flen)
		_, err := meta.file.ReadAt(name, index)
		if err != nil {
			return nil, err
		}
		index += flen

		_, err := meta.file.ReadAt(u32buf, index)
		if err != nil {
			return nil, err
		}
		index += 4
		ftype := binary.LittleEndian.Uint32(u32buf)
		diskFields[name] = ftype
	}
	meta.datastart = index

	if len(fieldMap) != len(diskFields) {
		return nil, errors.New("scheme not match")
	}
	for fname, ftype := range diskFields {
		if value, ok := fieldMap[fname]; ok {
			if value != ftype {
				return nil, errors.New("scheme not match")
			}
		} else {
			return nil, errors.New("scheme not match")
		}
	}

	return meta, nil
}

func (m *Meta) GetAll() ([]interface{}, error) {

}

func (m *Meta) Add(entry interface{}) error {

}

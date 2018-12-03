package schema

import (
	"log"
)

const (
	PAGE_SIZE = 8192
)

type ValueType int

const (
	Integer = iota
	Float
	String
	Boolean
)

type Column struct {
	Name string
	Type ValueType
}

type Table struct {
	pages   []*Page
	columns []Column
	dirties map[int]*Page
}

type Page struct {
	buffer []byte
}

func (p *Page) Validate() bool {
}

func (p *Page) NextRecord() []byte {
}

func NewTable(path string, cols []Column) (*Table, error) {

}

func LoadTable(path string) (*Table, error) {

}

func (t *Table) PageCount() int {
	return len(t.pages)
}

func (t *Table) GetPage(num int) *Page {
	return t.pages[num]
}

func (t *Table) Select(conditions, values []interface{}) error {
	pageCount := t.PageCount()
	for i := 0; i < pageCount; i++ {
		page := t.GetPage(i)
		record := page.NextRecord()
		for record != nil {
			record = page.NextRecord()
			for _, col := range t.columns {
				switch col.Type {
				case Integer:
					u64 := binary.LittleEndian.Uint64(record[0:4])
					log.Println(u64)
				}
			}
		}
	}
}

func (t *Table) Insert(values []interface{}) error {
	trxid := GetTrxId()
	pageNum := 0
	page := nil
	for _, value := range values {
		record := t.BuildRecord(value)
		page := buffer.GetPage(t.Name, pageNum)
		if page == nil {
			page := NewPage()
			pageNum++
		}
		if page.FreeSpace() > len(record) {
			cpage := t.CopyPage(i)
			for cpage.FreeSpace() > len(record) {
				cpage.AddRecord(record)
			}
		} else {
			cpage.SetLSN(lsn)
			xlog.AppendRecord()
			t.AddDirtyPage(pageNum, cpage)
		}
	}
	return err
}

func (t *Table) Flush() error {

}

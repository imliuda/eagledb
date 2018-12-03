package schema

const (
	PAGE_SIZE = 8192
)

type Buffer struct {
	pages []*Page

	dirties []*Page
}

type Page struct {
	buffer []byte
	file   os.File
	offset int64
}

func (b *Buffer) LoadTable(name string) error {

}

func (b *Buffer) PageCount(name string) (error, int) {

}

func (b *Buffer) GetPage(index int) (*Page, error) {

}

func (b *Buffer) FlushTable()

package schema

import (
	"sync"
	"sync/atomic"
)

var trxid uint32 = 0

func GetTrxId() uint32 {
	return atomic.AddUint32(trxid)
}

type LSN int64

type XLog struct {
	control *os.File
	lsn     LSN
	current *os.File
	offset  int64
	lock    sync.Mutex
}

type XLogRecord struct {
	TrxState int
}

func (l LSN) String() string {

}

func XLogOpen(dir string) (*XLog, error) {
	xlog := &XLog{}

	_, err := os.Stat(dir)
	if os.IsNotExsit(err) {
		err = os.MkdirAll(dir)
		if err != nil {
			return err
		}
	}

	file, err := os.OpenFile(dir+"/control", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	xlog.control = file

	err = parseControl(file, xlog)
	if err != nil {
		return nil
	}

}

func (x *XLog) Append(record XlogRecord) error {
	if record.TrxState == TRX_BEGIN {
		offset := x.NextPage()
		x.currlog.Write(offset, record.Data())
	}
}

func (x *XLog) Recover() error {

}

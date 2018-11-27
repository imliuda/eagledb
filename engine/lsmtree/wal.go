package lsmtree

import (
	"bytes"
	"os"
	"sync"
	"time"
)

type WalEntry struct {
}

type WalConfig struct {
	SyncPeriod time.Duration
	FileSize   int64
}

type Wal struct {
	dir      string
	files    []os.File
	currfile os.File
	config   *WalConfig
	entries  []*WalEntry
	lock     sync.Mutex
}

func NewWal() *Wal {
	return &Wal{}
}

// 打开一个新的wal文件，并启动同步协程。
func (l *Wal) Open() error {
	_, err := os.Stat(l.dir)
	if err != nil {
		err = os.MkdirAll(l.dir, 0755)
		if err != nil {
			return nil
		}
	}
	_, err = os.Open(l.dir)
	if err != nil {
		return err
	}

	return nil
}

// 数据落盘，并关闭wal文件。
func (l *Wal) Close() error {

	return nil
}

func getBytes(entries []*WalEntry) []byte {
	buffer := bytes.Buffer{}
	for _, entry := range entries {
		buffer.Write(entry.Bytes())
	}
	return buffer.Bytes()
}

// 写人wal文件时，将WalEntry通过管道发送给写入协程。
// 写入协程根据配置的写入模式，马上写入或稍后批量写入。
func (l *Wal) Append(entries []*WalEntry) error {
	if l.config.SyncPeriod.Seconds() == float64(0) {
		buffer := getBytes(entries)
		_, err := l.currfile.Write(buffer)
		if err != nil {
			return err
		}
		err = l.currfile.Sync()
		return err
	} else {
		l.lock.Lock()
		l.entries = append(l.entries, entries...)
		l.lock.Unlock()
		return nil
	}
}

func (l *Wal) Read() ([]*WalEntry, error) {
	return nil, nil
}

func (l *Wal) scheduleSync() error {
	return nil
}

func (e *WalEntry) Bytes() []byte {
	return []byte{}
}

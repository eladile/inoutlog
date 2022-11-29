package timelogger

import (
	"sync"
	"time"

	"inoutlog/readerwriter"
	"inoutlog/records"
)

type Logger interface {
	In(time.Time) error
	Out(time.Time) (int, error)
	GetAll() *records.Records
}

type logger struct {
	writer      readerwriter.Writer
	records     *records.Records
	recordsLock sync.RWMutex
	tariff      int
	extra       int
}

func NewLogger(tariff, extra int, path string) (Logger, error) {
	readerWriter := readerwriter.NewReaderWriter(path, tariff, extra)
	recs, err := readerWriter.Read()
	if err != nil {
		return nil, err
	}
	return &logger{
		writer:  readerWriter,
		records: recs,
		tariff:  tariff,
		extra:   extra,
	}, nil
}

func (l *logger) In(t time.Time) error {
	l.recordsLock.Lock()
	defer l.recordsLock.Unlock()
	err := l.records.In(t)
	if err != nil {
		return err
	}
	return l.writer.Write(l.records)
}

func (l *logger) Out(t time.Time) (int, error) {
	l.recordsLock.Lock()
	defer l.recordsLock.Unlock()
	pay, err := l.records.Out(t)
	if err != nil {
		return 0, err
	}
	return pay, l.writer.Write(l.records)
}

func (l *logger) GetAll() *records.Records {
	l.recordsLock.RLock()
	defer l.recordsLock.RUnlock()

	return l.records
}

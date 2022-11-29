package readerwriter

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"inoutlog/records"
)

type Writer interface {
	Write(*records.Records) error
}

type Reader interface {
	Read() (*records.Records, error)
}

type ReaderWriter interface {
	Reader
	Writer
}

type readerWriter struct {
	path          string
	tariff, extra int
}

func NewReaderWriter(path string, tariff, extra int) ReaderWriter {
	return &readerWriter{
		path:   path,
		tariff: tariff,
		extra:  extra,
	}
}

func (rw readerWriter) Read() (*records.Records, error) {
	f, err := os.Open(rw.path)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		f, err = os.Create(rw.path)
		if err != nil {
			return nil, err
		}
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	records := records.NewRecords(rw.tariff, rw.extra)
	if len(data) != 0 {
		err = json.Unmarshal(data, records)
	}
	return records, err

}

func (rw readerWriter) Write(records *records.Records) error {
	data, err := json.Marshal(records)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(rw.path, data, 0x666)
}

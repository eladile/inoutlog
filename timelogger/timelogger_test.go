package timelogger

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"inoutlog/records"
)

type readerWriterMock struct {
	records *records.Records
}

func (r readerWriterMock) Read() (*records.Records, error) {
	return r.records, nil
}

func (r readerWriterMock) Write(records *records.Records) error {
	r.records = records
	return nil
}

func TestTimeLogger(t *testing.T) {
	recs := records.NewRecords(50, 30)
	timeLogger := logger{
		writer:  readerWriterMock{records: recs},
		records: recs,
		tariff:  50,
		extra:   30,
	}
	base := time.Now()
	out1 := base.Add(4*time.Hour + 30*time.Minute)

	// in must be called first
	_, err := timeLogger.Out(time.Now())
	require.Error(t, err)
	require.Empty(t, timeLogger.GetAll().Records)

	// in - in
	require.NoError(t, timeLogger.In(out1))
	// keep both in
	require.NoError(t, timeLogger.In(base))
	recsSlice := timeLogger.GetAll().Records
	require.Len(t, recsSlice, 2)
	record := recsSlice[1]
	require.True(t, base.Equal(*record.In))
	// then out regards the last in when calculating the result fields
	_, err = timeLogger.Out(out1)
	require.NoError(t, err)
	recsSlice = timeLogger.GetAll().Records
	require.Len(t, recsSlice, 2)
	record = recsSlice[1]
	expected := records.Record{
		In:        &base,
		Out:       &out1,
		TotalTime: 4*time.Hour + 30*time.Minute,
		TotalPay:  4.5*50 + 30,
	}
	require.True(t, expected.In.Equal(*record.In))
	require.True(t, expected.Out.Equal(*record.Out))
	require.Equal(t, expected.TotalTime, record.TotalTime)
	require.Equal(t, expected.TotalPay, record.TotalPay)
	require.Equal(t, expected, record)

	// out after in-out
	_, err = timeLogger.Out(time.Now())
	require.Error(t, err)
	recsSlice = timeLogger.GetAll().Records
	require.Len(t, recsSlice, 2)
	record = recsSlice[1]
	require.Equal(t, expected, record)
}

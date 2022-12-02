package records

import (
	"errors"
	"time"
)

var (
	ErrMissingIn  = errors.New("missing in time")
	ErrMissingOut = errors.New("missing out time")
)

type Records struct {
	Records []Record `json:"records"`
	Tariff  int      `json:"tariff"`
	Extra   int      `json:"extra"`
}

func NewRecords(tariff, extra int) *Records {
	return &Records{
		Records: make([]Record, 0, 10),
		Tariff:  tariff,
		Extra:   extra,
	}
}

func (r *Records) In(t time.Time) error {
	r.Records = append(r.Records, Record{In: &t})
	return nil
}

func (r *Records) Out(t time.Time) (Record, error) {
	if len(r.Records) == 0 {
		return Record{}, ErrMissingIn
	}
	last := r.Records[len(r.Records)-1]
	if last.In != nil && last.Out != nil {
		return Record{}, ErrMissingIn
	}
	last.Out = &t
	err := last.FillPayment(r.Tariff, r.Extra)
	if err != nil {
		return Record{}, err
	}
	r.Records[len(r.Records)-1] = last
	return last, nil
}

type Record struct {
	In        *time.Time    `json:"in,omitempty"`
	Out       *time.Time    `json:"out,omitempty"`
	TotalTime time.Duration `json:"totalTime,omitempty"`
	TotalPay  int           `json:"totalPay,omitempty"`
	Paid      int           `json:"paid,omitempty"`
}

// FillPayment will update the record according to the hours worked.
// tariffHourly is per hour
// extra is added to the total value
func (r *Record) FillPayment(tariffHourly, extra int) error {
	if r.In == nil {
		return ErrMissingIn
	}
	if r.Out == nil {
		return ErrMissingOut
	}
	if r.In.After(*r.Out) {
		r.TotalPay = 0
		r.TotalTime = 0
	}
	r.TotalTime = r.Out.Sub(*r.In)
	hoursWorked := float64(r.TotalTime) / float64(time.Hour)
	r.TotalPay = int(hoursWorked*float64(tariffHourly) + float64(extra))
	return nil
}

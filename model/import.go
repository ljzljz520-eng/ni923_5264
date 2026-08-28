package model

type ImportBatch struct {
	ID       string   `json:"id"`
	Source   string   `json:"source"`
	Rows     int      `json:"rows"`
	Accepted int      `json:"accepted"`
	Rejected int      `json:"rejected"`
	Errors   []string `json:"errors"`
}

func (b ImportBatch) Complete() bool {
	return b.Rows == b.Accepted+b.Rejected
}

func (b ImportBatch) Summary() string {
	if b.Rejected == 0 {
		return "all rows accepted"
	}
	return "accepted rows with validation errors"
}

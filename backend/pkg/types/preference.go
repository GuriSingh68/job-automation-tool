package types

type Preference struct {
	ID        int64     `json:"id" db:"id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	JobTitles []*string `json:"job_titles" db:"preferred_role"`
	Locations []*string `json:"locations" db:"locations"`
	Keywords  []*string `json:"keywords" db:"keywords"`
	JobTypes  []*string `json:"job_types" db:"job_types"` // e.g., full-time, part-time, contract
}

const (
	FullTime   string = "full-time"
	PartTime   string = "part-time"
	Contract   string = "contract"
	Internship string = "internship"
	Co_op      string = "co-op"
)

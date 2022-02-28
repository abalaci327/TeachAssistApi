package teachassist

import "time"

type UserMetadata struct {
	Username      string
	Password      string
	StudentId     string
	SessionToken  string
	SessionExpiry time.Time
}

type Course struct {
	Metadata    CourseMetadata `json:"metadata"`
	Assessments []Assessment   `json:"assessments"`
	MarkWeights *MarkWeights   `json:"mark_weights"`
}

type CourseMetadata struct {
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Id          *string   `json:"id"`
	Block       string    `json:"block"`
	Room        string    `json:"room"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	CurrentMark *float32  `json:"current_mark"`
	MidtermMark *float32  `json:"midterm_mark"`
	FinalMark   *float32  `json:"final_mark"`
}

type Assessment struct {
	Name          string `json:"name"`
	Knowledge     []Mark `json:"knowledge"`
	Thinking      []Mark `json:"thinking"`
	Communication []Mark `json:"communication"`
	Application   []Mark `json:"application"`
	Other         []Mark `json:"other"`
	Culminating   []Mark `json:"culminating"`
}

type MarkWeights struct {
	Knowledge     float32 `json:"knowledge"`
	Thinking      float32 `json:"thinking"`
	Communication float32 `json:"communication"`
	Application   float32 `json:"application"`
	Other         float32 `json:"other"`
	Culminating   float32 `json:"culminating"`
}

type Mark struct {
	Numerator   float32 `json:"numerator"`
	Denominator float32 `json:"denominator"`
	Weighting   float32 `json:"weighting"`
}

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
	Name        string    `json:"name" example:"French"`
	Code        string    `json:"code" example:"FSF4UZ-02"`
	Id          *string   `json:"id" example:"462052"`
	Block       string    `json:"block" example:"3"`
	Room        string    `json:"room" example:"215"`
	StartDate   time.Time `json:"start_date" example:"2022-02-07T00:00:00Z"`
	EndDate     time.Time `json:"end_date" example:"2022-06-30T00:00:00Z"`
	CurrentMark *float32  `json:"current_mark" example:"0.95"`
	MidtermMark *float32  `json:"midterm_mark" example:"0.9"`
	FinalMark   *float32  `json:"final_mark" example:"0.95"`
}

type Assessment struct {
	Name          string `json:"name" example:"Test de compréhension écrite"`
	Knowledge     []Mark `json:"knowledge"`
	Thinking      []Mark `json:"thinking"`
	Communication []Mark `json:"communication"`
	Application   []Mark `json:"application"`
	Other         []Mark `json:"other"`
	Culminating   []Mark `json:"culminating"`
}

type MarkWeights struct {
	Knowledge     float32 `json:"knowledge" example:"14"`
	Thinking      float32 `json:"thinking" example:"14"`
	Communication float32 `json:"communication" example:"14"`
	Application   float32 `json:"application" example:"14"`
	Other         float32 `json:"other" example:"14"`
	Culminating   float32 `json:"culminating" example:"30"`
}

type Mark struct {
	Numerator   float32 `json:"numerator" example:"39"`
	Denominator float32 `json:"denominator" example:"40"`
	Weighting   float32 `json:"weighting" example:"15"`
}

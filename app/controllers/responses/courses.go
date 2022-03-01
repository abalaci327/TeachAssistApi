package responses

import "TeachAssistApi/app/teachassist"

type AllCoursesResponse struct {
	Metadata []teachassist.CourseMetadata `json:"metadata"`
}

type CourseIDResponse struct {
	Assessments []teachassist.Assessment `json:"assessments"`
	Weights     teachassist.MarkWeights  `json:"weights"`
}

package teachassist

import "time"

type UserMetadata struct {
	Username      string    `json:"username"`
	Password      string    `json:"password"`
	StudentId     string    `json:"student_id"`
	SessionToken  string    `json:"session_token"`
	SessionExpiry time.Time `json:"session_expiry"`
}

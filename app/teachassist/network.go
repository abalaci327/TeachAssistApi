package teachassist

import (
	"TeachAssistApi/app"
	"fmt"
	"net/http"
	"strings"
	"time"
)

func loginUser(username, password string) (UserMetadata, error) {
	url := "https://ta.yrdsb.ca/yrdsb/index.php"
	method := "POST"

	payload := strings.NewReader(fmt.Sprintf("password=%v&subject_id=0&submit=Login&username=%v", password, username))

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return UserMetadata{}, app.CreateError(app.NetworkError)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		return UserMetadata{}, app.CreateError(app.NetworkError)
	}
	defer res.Body.Close()

	var studentId, sessionToken string
	var sessionExpiry time.Time
	for _, v := range res.Cookies() {
		switch v.Name {
		case "student_id":
			studentId = v.Value
		case "session_token":
			sessionToken = v.Value
			sessionExpiry = v.Expires
		}
	}

	if studentId == "" || sessionToken == "" || sessionToken == "deleted" {
		return UserMetadata{}, app.CreateError(app.AuthError)
	}

	return UserMetadata{
		Username:      username,
		Password:      password,
		StudentId:     studentId,
		SessionToken:  sessionToken,
		SessionExpiry: sessionExpiry,
	}, nil
}

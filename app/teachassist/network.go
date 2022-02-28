package teachassist

import (
	"TeachAssistApi/app"
	"fmt"
	"io/ioutil"
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

func getAllCourses(user *UserMetadata) (*string, error) {
	url := fmt.Sprintf("https://ta.yrdsb.ca/live/students/listReports.php?student_id=%s", user.StudentId)
	return authenticatedRequest(url, user)
}

func getCourseById(id string, user *UserMetadata) (*string, error) {
	url := fmt.Sprintf("https://ta.yrdsb.ca/live/students/viewReport.php?subject_id=%s&student_id=%s", id, user.StudentId)
	return authenticatedRequest(url, user)
}

func authenticatedRequest(url string, user *UserMetadata) (*string, error) {
	method := "GET"

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, app.CreateError(app.NetworkError)
	}

	req.Header.Add("Cookie", fmt.Sprintf("session_token=%s; student_id=%s", user.SessionToken, user.StudentId))

	res, err := client.Do(req)
	if err != nil {
		return nil, app.CreateError(app.NetworkError)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, app.CreateError(app.NetworkError)
	}

	webData := string(body)
	if webData == "" || strings.Contains(webData, "Login") {
		return nil, app.CreateError(app.AuthError)
	}

	return &webData, nil
}

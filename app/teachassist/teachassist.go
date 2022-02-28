package teachassist

func LoginUser(username, password string) (UserMetadata, error) {
	return loginUser(username, password)
}

func GetAllCourses(user *UserMetadata) (*[]CourseMetadata, error) {
	courses, err := getAllCourses(user)
	if err != nil {
		return nil, err
	}
	return parseCourseMetadata(courses)
}

func GetCourseByID(id string, user *UserMetadata) (*MarkWeights, *[]Assessment, error) {
	courses, err := getCourseById(id, user)
	if err != nil {
		return nil, nil, err
	}
	return parseCourse(courses)
}

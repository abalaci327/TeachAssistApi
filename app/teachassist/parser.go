package teachassist

import (
	"TeachAssistApi/app"
	"github.com/PuerkitoBio/goquery"
	"strconv"
	"strings"
	"time"
)

func parseCourseMetadata(document *string) (metadata *[]CourseMetadata, err error) {
	defer func() {
		if r := recover(); r != nil {
			metadata = nil
			err = app.CreateError(app.ParsingError)
		}
	}()

	metadata = &[]CourseMetadata{}

	doc, er := goquery.NewDocumentFromReader(strings.NewReader(*document))
	if er != nil {
		return nil, app.CreateError(app.ParsingError)
	}

	msgBox := doc.Find(".green_border_message").Eq(1)
	table := msgBox.Find("div").First().Find("table > tbody").First()

	rows := table.Find("tr")
	for i := 1; i < len(rows.Nodes); i++ {
		row := rows.Eq(i)
		cols := row.Find("td")
		if len(cols.Nodes) < 3 {
			continue
		}

		generalInfoString := trimWhitespace(cols.Eq(0).Text())
		dateInfoString := trimWhitespace(cols.Eq(1).Text())
		markInfoString := trimWhitespace(cols.Eq(2).Text())

		generalInfo := strings.Split(generalInfoString, ":")
		dates := strings.Split(dateInfoString, "~")
		markLocked := false
		var markInfo []string
		if strings.Contains(markInfoString, "Please") || markInfoString == "" {
			markLocked = true
			markInfo = strings.Split(markInfoString, "P")
		} else {
			markInfo = strings.Split(markInfoString, "=")
		}

		// General Info
		var courseId *string = nil
		if !markLocked {
			val, _ := cols.Eq(2).Find("a").First().Attr("href")
			split := strings.Split(val, "subject_id=")
			id := trimWhitespace(strings.Split(split[len(split)-1], "&")[0])
			courseId = &id
		}
		courseCode := trimWhitespace(generalInfo[0])
		courseName := trimWhitespace(strings.Split(generalInfo[1], "Block")[0])
		if courseName == "" {
			courseName = courseCode
		}
		moreInfo := strings.Split(generalInfo[2], "-")
		var courseBlock, courseRoom string
		if len(moreInfo) > 1 {
			courseBlock = trimWhitespace(moreInfo[0])
		}
		rm := strings.Split(moreInfo[1], "rm.")
		if len(rm) > 1 {
			courseRoom = trimWhitespace(rm[1])
		}

		// Dates
		startDate, er := time.Parse("2006-01-02", trimWhitespace(dates[0]))
		if er != nil {
			startDate = time.Now()
		}
		endDate, er := time.Parse("2006-01-02", trimWhitespace(dates[1]))
		if er != nil {
			startDate = time.Now()
		}

		// Marks
		var currentMark, midtermMark, finalMark *float32
		if len(markInfo) > 1 {
			if !markLocked {
				curMark := trimPercentage(markInfo[1])
				t := *parseFloat(curMark) / 100
				currentMark = &t
			}
			if strings.Contains(markInfo[0], "MIDTERM") && strings.Contains(markInfo[0], "FINAL") {
				midMarkArr := strings.Split(strings.Split(markInfo[0], "%")[0], " ")
				midMark := trimPercentage(midMarkArr[len(midMarkArr)-1])
				t := *parseFloat(midMark) / 100
				midtermMark = &t

				finMarkArr := strings.Split(markInfo[0], ":")
				finMark := trimPercentage(strings.Split(finMarkArr[len(finMarkArr)-1], " ")[0])
				t = *parseFloat(finMark) / 100
				finalMark = &t
			} else if strings.Contains(markInfo[0], "MIDTERM") {
				midMarkArr := strings.Split(strings.Split(markInfo[0], "%")[0], " ")
				midMark := trimPercentage(midMarkArr[len(midMarkArr)-1])
				t := *parseFloat(midMark) / 100
				midtermMark = &t
			} else if strings.Contains(markInfo[0], "FINAL") {
				finMarkArr := strings.Split(markInfo[0], ":")
				finMark := trimPercentage(strings.Split(finMarkArr[len(finMarkArr)-1], " ")[0])
				t := *parseFloat(finMark) / 100
				finalMark = &t
			}
		}
		print()

		m := CourseMetadata{
			Name:        courseName,
			Code:        courseCode,
			Id:          courseId,
			Block:       courseBlock,
			Room:        courseRoom,
			StartDate:   startDate,
			EndDate:     endDate,
			CurrentMark: currentMark,
			MidtermMark: midtermMark,
			FinalMark:   finalMark,
		}
		*metadata = append(*metadata, m)
	}
	return
}

func parseCourse(document *string) (weights *MarkWeights, assessments *[]Assessment, err error) {
	if strings.Contains(*document, "Student Reports for") {
		return nil, nil, app.CreateError(app.InvalidCourseIdError)
	}

	defer func() {
		if r := recover(); r != nil {
			weights = nil
			assessments = nil
			err = app.CreateError(app.ParsingError)
		}
	}()

	weights = &MarkWeights{}
	assessments = &[]Assessment{}

	doc, er := goquery.NewDocumentFromReader(strings.NewReader(*document))
	if er != nil {
		return nil, nil, app.CreateError(app.ParsingError)
	}

	tables := doc.Find("table").FilterFunction(filterForAttribute("width", "100%"))
	tables = tables.First().AddClass("assessments")
	rows := tables.First().Parent().Find(".assessments > tbody > tr")
	markRows := &goquery.Selection{}
	for i := 1; i < len(rows.Nodes); i += 2 {
		markRows = markRows.AddSelection(rows.Eq(i))
	}

	markRows.Each(func(_ int, selection *goquery.Selection) {
		nameRaw := selection.Find("td").FilterFunction(filterForAttribute("rowspan", "2"))
		name := trimWhitespace(nameRaw.Text())

		knowledgeRaw := selectTableBody(selection, "ffffaa")
		thinkingRaw := selectTableBody(selection, "c0fea4")
		communicationRaw := selectTableBody(selection, "afafff")
		applicationRaw := selectTableBody(selection, "ffd490")
		otherCulminatingRaw := selectTableBody(selection, "#dedede")

		knowledge := parseMarkFromTable(knowledgeRaw.Find("tr > td"))
		thinking := parseMarkFromTable(thinkingRaw.Find("tr > td"))
		communication := parseMarkFromTable(communicationRaw.Find("tr > td"))
		application := parseMarkFromTable(applicationRaw.Find("tr > td"))
		other := parseMarkFromTable(otherCulminatingRaw.Find("tr > td").FilterFunction(filterForAttribute("bgcolor", "#dedede")))
		culminating := parseMarkFromTable(otherCulminatingRaw.Find("tr > td").FilterFunction(filterForAttribute("bgcolor", "#cccccc")))

		*assessments = append(*assessments, Assessment{
			Name:          name,
			Knowledge:     *knowledge,
			Thinking:      *thinking,
			Communication: *communication,
			Application:   *application,
			Other:         *other,
			Culminating:   *culminating,
		})
	})

	weightingsTable := doc.Find(".green_border_message").Last().Find("div > table:nth-child(2)")

	knowledgeRow := weightingsTable.Find("tr").FilterFunction(filterForAttribute("bgcolor", "#ffffaa"))
	thinkingRow := weightingsTable.Find("tr").FilterFunction(filterForAttribute("bgcolor", "#c0fea4"))
	communicationRow := weightingsTable.Find("tr").FilterFunction(filterForAttribute("bgcolor", "#afafff"))
	applicationRow := weightingsTable.Find("tr").FilterFunction(filterForAttribute("bgcolor", "#ffd490"))
	otherRow := weightingsTable.Find("tr").FilterFunction(filterForAttribute("bgcolor", "#eeeeee"))
	culminatingRow := weightingsTable.Find("tr").FilterFunction(filterForAttribute("bgcolor", "#cccccc"))

	knowledgeWeighting := parseCourseWeighting(knowledgeRow)
	thinkingWeighting := parseCourseWeighting(thinkingRow)
	communicationWeighting := parseCourseWeighting(communicationRow)
	applicationWeighting := parseCourseWeighting(applicationRow)
	otherWeighting := parseCourseWeighting(otherRow)
	culminatingWeighting := parseCulminatingCourseWeighting(culminatingRow)
	if knowledgeWeighting == nil || thinkingWeighting == nil || communicationWeighting == nil || applicationWeighting == nil || otherWeighting == nil || culminatingWeighting == nil {
		return nil, nil, app.CreateError(app.ParsingError)
	}

	weights = &MarkWeights{
		Knowledge:     *knowledgeWeighting,
		Thinking:      *thinkingWeighting,
		Communication: *communicationWeighting,
		Application:   *applicationWeighting,
		Other:         *otherWeighting,
		Culminating:   *culminatingWeighting,
	}

	return
}

func trimWhitespace(s string) string {
	return strings.TrimSpace(s)
}

func trimPercentage(s string) string {
	r := strings.NewReplacer("%", "", "\n", "", "\t", "", " ", "")
	return r.Replace(s)
}

func parseFloat(s string) *float32 {
	mark, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return nil
	} else {
		mark32 := float32(mark)
		return &mark32
	}
}

func filterForAttribute(key, value string) func(int, *goquery.Selection) bool {
	return func(i int, selection *goquery.Selection) bool {
		if attr, _ := selection.Attr(key); attr == value {
			return true
		}
		return false
	}
}

func selectTableBody(selection *goquery.Selection, bgColor string) *goquery.Selection {
	return selection.Find("td").FilterFunction(filterForAttribute("bgcolor", bgColor)).Find("table > tbody")
}

func parseMarkFromTable(selection *goquery.Selection) (marks *[]Mark) {
	marks = &[]Mark{}

	selection.Each(func(_ int, s *goquery.Selection) {
		markStr := strings.Split(trimWhitespace(s.Text()), "=")

		parts := strings.Split(markStr[0], "/")
		var numeratorStr, denominatorStr string
		// TODO: Potentially add indicator for the state when it is just " / 5" - i.e. not yet marked by teacher or missing
		if len(parts) == 1 {
			numeratorStr = "0"
			denominatorStr = trimWhitespace(parts[0])
		} else {
			numeratorStr = trimWhitespace(parts[0])
			denominatorStr = trimWhitespace(parts[1])
		}

		weightStr := trimWhitespace(s.Find("font").Text())
		if weightStr == "no weight" {
			weightStr = "0"
		} else {
			weightStr = trimWhitespace(strings.Split(weightStr, "=")[1])
		}

		numerator := parseFloat(numeratorStr)
		denominator := parseFloat(denominatorStr)
		weight := parseFloat(weightStr)

		*marks = append(*marks, Mark{
			Numerator:   *numerator,
			Denominator: *denominator,
			Weighting:   *weight,
		})
	})

	return
}

func parseCourseWeighting(selection *goquery.Selection) *float32 {
	selection = selection.Find("td").FilterFunction(filterForAttribute("align", "right"))
	weightingStr := trimPercentage(selection.Eq(1).Text())
	return parseFloat(weightingStr)
}

func parseCulminatingCourseWeighting(selection *goquery.Selection) *float32 {
	selection = selection.Find("td").FilterFunction(filterForAttribute("align", "right"))
	weightingStr := trimPercentage(selection.Eq(0).Text())
	return parseFloat(weightingStr)
}

package entity

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Lab struct {
	ID   int
	Name string
}

type LabParserRules struct {
}

type LabParser struct {
	Lab
	ParserRules LabParserRules
}

func (l *LabParser) ParseResults(txt string) CheckupResults {
	var checkupRes CheckupResults

	checkupRes.Checkup.Lab.Name = l.Name
	checkupRes.Checkup.Lab.ID = l.ID

	if strings.Contains(txt, "Дата/время взятия материала:") {
		checkupRes.Checkup.Date = time.Now() // TODO: strings.Split(strings.TrimSpace(strings.Split(line, ":")[1]), " ")[0]
	}
	checkupRes.Checkup.Profile = l.ParsePatient(txt)

	lines := strings.Split(txt, "\n")

	// Установить значения для полей Analysis
	for _, line := range lines {
		line = strings.Trim(line, "\"")

		m := l.ParseMarker(line)
		if m.Marker.Name != "" {
			checkupRes.Results = append(checkupRes.Results, m)
		}
	}

	return checkupRes
}

func (l *LabParser) ParsePatient(txt string) Profile {
	p := Profile{}

	re := regexp.MustCompile(`Ф\.И\..+?: (\S+ \S+ \S+)`)
	patient := regexp.MustCompile(`пациент:.+?(\S+ \S+ \S+)`)

	matches := re.FindStringSubmatch(txt)
	if len(matches) == 0 {
		matches = patient.FindStringSubmatch(txt)
	}

	if len(matches) > 0 {
		p.FName = matches[1]
	}

	re = regexp.MustCompile(`Дата рождения: (\d{2}/\d{2}/\d{4})`)
	matches = re.FindStringSubmatch(txt)
	//if len(matches) > 0 {
	//	p.DateOfBirth = matches[1]
	//}

	return p
}

func (l *LabParser) ParseMarker(txt string) MarkerResult {
	//testResultPattern := regexp.MustCompile(`^(\S.+?)\s+(\d+\.\d+)\s+(\S+)\s+(\d+\.\d+[-<>\d.,\s]*)$`)
	testResultPattern := regexp.MustCompile(`^(.*?)\s+([\d.,]+(?:\s?10\^\d+)?)\s*([\p{L}/%^]+)\s*([\d.,]+[-=–][\d.,]+)$`)

	txt = strings.Trim(txt, "|")
	matches := testResultPattern.FindStringSubmatch(txt)
	if len(matches) <= 0 {
		return MarkerResult{}
	}

	val, _ := strconv.ParseFloat(matches[2], 64)
	return MarkerResult{
		Marker: Marker{
			Name: strings.TrimSpace(matches[1]),
		},
		Value: val,
		Unit:  Unit{Name: strings.TrimSpace(matches[3])},
		//ReferenceRange: strings.TrimSpace(matches[4]),
	}
}

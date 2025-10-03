package docupanda

import (
	"labra/internal/entity"
	"time"
)

type Document struct {
	File File `json:"file"`
}

type File struct {
	Contents string `json:"contents"` // base64-encoded content
	Filename string `json:"filename"`
}

type Report struct {
	Laboratory      string         `json:"laboratory"`
	PhoneNumber     string         `json:"phoneNumber"`
	Website         string         `json:"website"`
	Patient         Patient        `json:"patient"`
	TestDetails     TestDetails    `json:"testDetails"`
	AnalysisResults []AnalysisItem `json:"analysisResults"`
}

type Patient struct {
	Name      string `json:"name"`
	BirthDate string `json:"birthDate"`
	Gender    string `json:"gender"`
}

type TestDetails struct {
	TestType    string `json:"testType"`
	SampleDate  string `json:"sampleDate"`
	OrderNumber string `json:"orderNumber"`
}

type AnalysisItem struct {
	TestName       string         `json:"testName"`
	Result         Result         `json:"result,omitempty"`
	ReferenceRange ReferenceRange `json:"referenceRange,omitempty"`
}

type Result struct {
	Amount float64 `json:"amount"`
	Unit   string  `json:"unit"`
}

type ReferenceRange struct {
	From string `json:"from"`
	To   string `json:"to"`
}

func (r Patient) buildEntity() entity.Profile {
	layout := "2014-09-12T11:45:26.371Z"
	date, _ := time.Parse(layout, r.BirthDate)

	return entity.Profile{
		FName:       r.Name,
		LName:       r.Name,
		Gender:      r.Gender,
		DateOfBirth: date,
	}
}

func (a AnalysisItem) buildEntity() entity.MarkerResult {
	return entity.MarkerResult{
		Marker: entity.Marker{
			Name: a.TestName,
			ReferenceRange: entity.Range{
				From: a.ReferenceRange.From,
				To:   a.ReferenceRange.To,
			},
		},
		Value: a.Result.Amount,
		Unit: entity.Unit{
			Name: a.Result.Unit,
		},
		CreatedAt: time.Time{},
	}
}

func (r Report) buildEntity() entity.CheckupResults {
	layout := "2014-09-12T11:45:26.371Z"
	date, _ := time.Parse(layout, r.TestDetails.SampleDate)

	result := entity.CheckupResults{}
	result.Checkup = entity.Checkup{
		Title:   r.TestDetails.TestType,
		Profile: r.Patient.buildEntity(),
		Lab: entity.Lab{
			Name: r.Laboratory,
		},
		Date:    date,
		Comment: r.TestDetails.OrderNumber,
	}

	result.Results = make(entity.MarkerResults, 0, len(r.AnalysisResults))
	for _, res := range r.AnalysisResults {
		result.Results = append(result.Results, res.buildEntity())
	}

	return result
}

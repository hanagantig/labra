package deepseek

import (
	"labra/internal/entity"
)

// Define the struct for the request body
type Request struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

// Define the struct for each message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Define the struct for the response
type DeepSeekResponse struct {
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
}

// Struct for the parsed content inside the "content" field
type ParsedContent struct {
	PatientName  string   `json:"patient_name"`
	PatientBDate string   `json:"patient_bdate"`
	Markers      []Marker `json:"markers"`
}

type Marker struct {
	Name  string  `json:"name"`
	Unit  string  `json:"unit"`
	Value float64 `json:"value"`
}

func (p ParsedContent) BuildEntity() entity.CheckupResults {
	res := entity.CheckupResults{
		Checkup: entity.Checkup{
			Profile: entity.Profile{
				FName: p.PatientName,
				//DateOfBirth: p.PatientBDate,
			},
		},
	}

	res.Results = make(entity.MarkerResults, 0, len(p.Markers))
	for _, marker := range p.Markers {
		res.Results = append(res.Results, entity.MarkerResult{
			Marker: entity.Marker{Name: marker.Name},
			Value:  marker.Value,
			Unit: entity.Unit{
				Name: marker.Unit,
			},
		})
	}

	return res
}

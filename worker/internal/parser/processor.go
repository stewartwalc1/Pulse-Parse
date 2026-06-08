package parser

import (
	"time"
)


type LabReport struct {
	ID        string    `json:"id"`
	PatientID string    `json:"patient_id"`
	RawText   string    `json:"raw_text"`
	CreatedAt time.Time `json:"created_at"`
}


type Biomarker struct {
	Name       string  `json:"name"`        // e.g., "Vitamin D", "Ferritin"
	Value      float64 `json:"value"`       // e.g., 24.5
	Unit       string  `json:"unit"`        // e.g., "ng/mL"
	Reference  string  `json:"reference"`   // e.g., "30.0 - 100.0"
	IsCritical bool    `json:"is_critical"` // Flag if completely out of bounds
}


type AnalysisResult struct {
	ReportID   string      `json:"report_id"`
	Biomarkers []Biomarker `json:"biomarkers"`
	Insights   string      `json:"insights"`    // LLM-generated root-cause wellness trends
	ParsedAt   time.Time   `json:"parsed_at"`
}


func ExtractAndAnalyze(report LabReport) (*AnalysisResult, error) {

	
	mockMarkers := []Biomarker{
		{
			Name:       "Vitamin D, 25-Hydroxy",
			Value:      22.4,
			Unit:       "ng/mL",
			Reference:  "30.0 - 100.0",
			IsCritical: true,
		},
		{
			Name:       "Hemoglobin A1c",
			Value:      5.6,
			Unit:       "%",
			Reference:  "< 5.7",
			IsCritical: false,
		},
	}

	return &AnalysisResult{
		ReportID:   report.ID,
		Biomarkers: mockMarkers,
		Insights:   "Patient exhibits suboptimal Vitamin D levels, which may correlate with reports of fatigue. HbA1c is within normal limits but optimal metabolic health suggests keeping this stable.",
		ParsedAt:   time.Now(),
	}, nil
}
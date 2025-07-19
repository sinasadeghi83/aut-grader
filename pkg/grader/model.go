package grader

import (
	"time"

	m "github.com/sinasadeghi83/aut-grader/pkg/platform/model"
)

type Project struct {
	m.Model
	Name     string    `json:"name"`
	Due      time.Time `json:"due"`
	Sections []Section `json:"sections"`
}

type Section struct {
	m.Model
	Name        string   `json:"name"`
	DependsOn   *Section `json:"depends_on"`
	DependsOnID *uint    `json:"depends_on_id"`
	ProjectID   uint     `json:"project_id"`
}

type Scenario struct {
	m.Model
	Name        string    `json:"name"`
	DependsOn   *Scenario `json:"depends_on"`
	DependsOnID *uint     `json:"depends_on_id"`
	Section     Section   `json:"section"`
	SectionID   uint      `json:"section_id"`
}

type Test struct {
	m.Model
	Name     string    `json:"name"`
	Request  TRequest  `json:"request" gorm:"embedded"`
	Response TResponse `json:"response" gorm:"embedded"`

	DependsOn   *Test `json:"depends_on"`
	DependsOnID *uint `json:"depends_on_id"`

	Scenario   Scenario `json:"scenario"`
	ScenarioID uint     `json:"scenario_id"`
}

type TRequest struct {
	Url     string    `json:"url"`
	Method  string    `json:"method"`
	Headers []THeader `json:"headers"`
	ReqBody string    `json:"body"`
}

type THeader struct {
	m.Model
	Key    string `json:"key" gorm:"hkey"`
	Value  string `json:"value" gorm:"hvalue"`
	Test   Test   `json:"test"`
	TestID uint   `json:"test_id"`
}

type TResponse struct {
	StatusCode uint   `json:"status_code"`
	ResBody    string `json:"body"`
	// Headers    []THeader `json:"headers"` // TODO
}

// GradingStatus defines the status of a grading task.
type GradingStatus string

const (
	StatusProcessing GradingStatus = "processing"
	StatusPassed     GradingStatus = "passed"
	StatusFailed     GradingStatus = "failed"
)

// Grader Result Structs for Database Storage
type ProjectResult struct {
	m.Model
	ProjectID   uint            `json:"project_id"`
	ProjectName string          `json:"project_name"`
	UserID      uint            `json:"user_id"` // Link to the user who initiated the grading
	Status      GradingStatus   `json:"status"`
	Message     string          `json:"message,omitempty"`
	Sections    []SectionResult `json:"sections" gorm:"foreignKey:ProjectResultID"`
}

type SectionResult struct {
	m.Model
	SectionID       uint             `json:"section_id"`
	SectionName     string           `json:"section_name"`
	ProjectResultID uint             `json:"project_result_id"` // Foreign key to ProjectResult
	Status          GradingStatus    `json:"status"`
	Message         string           `json:"message,omitempty"`
	TotalScenarios  uint             `json:"total_scenarios"`
	PassedScenarios uint             `json:"passed_scenarios"`
	Score           float64          `json:"score"` // Percentage of passed scenarios
	Scenarios       []ScenarioResult `json:"scenarios" gorm:"foreignKey:SectionResultID"`
}

type ScenarioResult struct {
	m.Model
	ScenarioID      uint          `json:"scenario_id"`
	ScenarioName    string        `json:"scenario_name"`
	SectionResultID uint          `json:"section_result_id"` // Foreign key to SectionResult
	Status          GradingStatus `json:"status"`
	Message         string        `json:"message,omitempty"`
	Tests           []TestResult  `json:"tests" gorm:"foreignKey:ScenarioResultID"`
}

type TestResult struct {
	m.Model
	TestID           uint          `json:"test_id"`
	TestName         string        `json:"test_name"`
	ScenarioResultID uint          `json:"scenario_result_id"` // Foreign key to ScenarioResult
	Status           GradingStatus `json:"status"`
	Message          string        `json:"message,omitempty"`
	// Store request/response details if needed for auditing/debugging
	ActualStatusCode     uint   `json:"actual_status_code,omitempty"`
	ExpectedStatusCode   uint   `json:"expected_status_code,omitempty"`
	ActualResponseBody   string `json:"actual_response_body,omitempty"`
	ExpectedResponseBody string `json:"expected_response_body,omitempty"`
}

func (Project) TableName() string {
	return "projects"
}

func (Section) TableName() string {
	return "sections"
}

func (Scenario) TableName() string {
	return "scenarios"
}

func (Test) TableName() string {
	return "tests"
}

func (THeader) TableName() string {
	return "theaders"
}

func (ProjectResult) TableName() string {
	return "project_results"
}

func (SectionResult) TableName() string {
	return "section_results"
}

func (ScenarioResult) TableName() string {
	return "scenario_results"
}

func (TestResult) TableName() string {
	return "test_results"
}

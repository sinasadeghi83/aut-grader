package grader

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/go-resty/resty/v2"
	"gorm.io/gorm"
)

var (
	variableRegex     = regexp.MustCompile(`\$<(\w+)>`)
	substitutionRegex = regexp.MustCompile(`\{\{(\w+)\}\}`)
)

type Grader struct {
	BaseUrl   string
	UserID    uint
	variables map[string]interface{}
}

func NewGrader(baseUrl string, userID uint) *Grader {
	return &Grader{
		BaseUrl:   baseUrl,
		UserID:    userID,
		variables: make(map[string]interface{}),
	}
}

// GradeProject is the main entry point for grading a project.
func (gd *Grader) GradeProject(db *gorm.DB, projID uint) (*ProjectResult, error) {
	var proj Project
	if err := db.Find(&proj, projID).Error; err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	projectResult := &ProjectResult{
		ProjectID:   proj.ID,
		ProjectName: proj.Name,
		UserID:      gd.UserID,
		Status:      StatusProcessing,
		Message:     "Processing...",
	}

	if err := db.Create(projectResult).Error; err != nil {
		return nil, fmt.Errorf("failed to create initial project result: %w", err)
	}

	if err := gd.processProject(db, proj, projectResult.ID); err != nil {
		projectResult.Message = fmt.Sprintf("Error processing project: %s", err.Error())
		projectResult.Status = StatusFailed
	}

	gd.updateProjectResultStatus(db, projectResult)

	if err := db.Save(projectResult).Error; err != nil {
		return projectResult, fmt.Errorf("failed to save final project result: %w", err)
	}

	db.Preload("Sections").First(projectResult, projectResult.ID)
	return projectResult, nil
}

// processProject iterates through the sections of a project and processes them.
func (gd *Grader) processProject(db *gorm.DB, proj Project, projectResultID uint) error {
	var secs []Section
	if err := db.Model(&proj).Where("depends_on_id IS NULL").Association("Sections").Find(&secs); err != nil {
		return fmt.Errorf("failed to load sections for project %d: %w", proj.ID, err)
	}

	for i := 0; i < len(secs); i++ {
		sec := secs[i]
		isPass, err := gd.processSingleSection(db, sec, projectResultID)
		if err != nil {
			// Log the error but continue processing other sections
			fmt.Printf("Error processing section %d: %v\n", sec.ID, err)
		} else if isPass {
			var dependeds []Section
			if err := db.Model(&sec).Association("DependsOns").Find(&dependeds); err != nil {
				fmt.Printf("Error processing depended section %d: %v\n", sec.ID, err)
			} else {
				secs = append(secs, dependeds...)
			}
		}
	}
	return nil
}

// processSingleSection handles the grading of a single section.
func (gd *Grader) processSingleSection(db *gorm.DB, sec Section, projectResultID uint) (bool, error) {
	gd.variables = make(map[string]interface{}) // Reset variables for each section
	sectionResult := &SectionResult{
		SectionID:       sec.ID,
		SectionName:     sec.Name,
		ProjectResultID: projectResultID,
		Status:          StatusProcessing,
		Message:         "Processing...",
	}
	if err := db.Create(sectionResult).Error; err != nil {
		return false, fmt.Errorf("failed to create initial section result: %w", err)
	}

	if err := gd.processSection(db, sec, sectionResult.ID); err != nil {
		sectionResult.Message = fmt.Sprintf("Error processing section: %s", err.Error())
		sectionResult.Status = StatusFailed
	}

	gd.updateSectionResultStatus(db, sectionResult)

	return sectionResult.Status == StatusPassed, db.Save(sectionResult).Error
}

// processSection iterates through the scenarios of a section and processes them.
func (gd *Grader) processSection(db *gorm.DB, sec Section, sectionResultID uint) error {
	var scns []Scenario
	if err := db.Model(&sec).Where("depends_on_id IS NULL").Association("Scenarios").Find(&scns); err != nil {
		return fmt.Errorf("failed to load scenarios for section %d: %w", sec.ID, err)
	}

	for i := 0; i < len(scns); i++ {
		scn := scns[i]
		isPass, err := gd.processSingleScenario(db, scn, sectionResultID)
		if err != nil {
			fmt.Printf("Error processing scenario %d: %v\n", scn.ID, err)
		} else if isPass {
			var dependeds []Scenario
			if err := db.Model(&scn).Association("DependsOns").Find(&dependeds); err != nil {
				fmt.Printf("Error processing depended scenario %d: %v\n", scn.ID, err)
			} else {
				scns = append(scns, dependeds...)
			}
		}
	}
	return nil
}

// processSingleScenario handles the grading of a single scenario.
func (gd *Grader) processSingleScenario(db *gorm.DB, scn Scenario, sectionResultID uint) (bool, error) {
	scenarioResult := &ScenarioResult{
		ScenarioID:      scn.ID,
		ScenarioName:    scn.Name,
		SectionResultID: sectionResultID,
		Status:          StatusProcessing,
		Message:         "Processing...",
	}
	if err := db.Create(scenarioResult).Error; err != nil {
		return false, fmt.Errorf("failed to create initial scenario result: %w", err)
	}

	if err := gd.processScenario(db, scn, scenarioResult.ID); err != nil {
		scenarioResult.Message = fmt.Sprintf("Error processing scenario: %s", err.Error())
		scenarioResult.Status = StatusFailed
	}

	gd.updateScenarioResultStatus(db, scenarioResult)

	return scenarioResult.Status == StatusPassed, db.Save(scenarioResult).Error
}

// processScenario iterates through the tests of a scenario and processes them.
func (gd *Grader) processScenario(db *gorm.DB, scn Scenario, scenarioResultID uint) error {
	var tests []Test
	if err := db.Model(&scn).Where("depends_on_id IS NULL").Association("Tests").Find(&tests); err != nil {
		return fmt.Errorf("failed to load tests for scenario %d: %w", scn.ID, err)
	}

	for i := 0; i < len(tests); i++ {
		test := tests[i]
		isPass, err := gd.processSingleTest(db, test, scenarioResultID)
		if err != nil {
			fmt.Printf("Error processing test %d: %v\n", test.ID, err)
		} else if isPass {
			var dependeds []Test
			if err := db.Model(&test).Association("DependsOns").Find(&dependeds); err != nil {
				fmt.Printf("Error processing depended test %d: %v\n", test.ID, err)
			} else {
				tests = append(tests, dependeds...)
			}
		}
	}
	return nil
}

// processSingleTest handles the grading of a single test.
func (gd *Grader) processSingleTest(db *gorm.DB, test Test, scenarioResultID uint) (bool, error) {
	testResult := &TestResult{
		TestID:               test.ID,
		TestName:             test.Name,
		ScenarioResultID:     scenarioResultID,
		ExpectedStatusCode:   test.Response.StatusCode,
		ExpectedResponseBody: test.Response.ResBody,
		Status:               StatusProcessing,
		Message:              "Processing...",
	}
	if err := db.Create(testResult).Error; err != nil {
		return false, fmt.Errorf("failed to create initial test result: %w", err)
	}

	gd.executeTest(db, test, testResult)

	return testResult.Status == StatusPassed, db.Save(testResult).Error
}

// executeTest runs a single test and updates its result.
func (gd *Grader) executeTest(db *gorm.DB, test Test, testResult *TestResult) {
	resp, err := gd.makeRequest(test)
	if err != nil {
		testResult.Status = StatusFailed
		testResult.Message = fmt.Sprintf("request failed: %w", err)
		return
	}

	testResult.ActualStatusCode = uint(resp.StatusCode())
	testResult.ActualResponseBody = resp.String()

	if err := gd.validateResponse(test, resp); err != nil {
		testResult.Status = StatusFailed
		testResult.Message = err.Error()
	} else {
		testResult.Status = StatusPassed
		testResult.Message = "Test passed successfully!"
	}
}

// makeRequest executes the HTTP request for a test.
func (gd *Grader) makeRequest(test Test) (*resty.Response, error) {
	client := resty.New()

	// Substitute variables in the URL, headers, and body
	fullURL := gd.substituteVariables(gd.BaseUrl + test.Request.Url)
	req := client.R()

	for _, header := range test.Request.Headers {
		req.SetHeader(gd.substituteVariables(header.Key), gd.substituteVariables(header.Value))
	}

	if test.Request.ReqBody != "" {
		req.SetBody(gd.substituteVariables(test.Request.ReqBody))
	}

	switch test.Request.Method {
	case "GET":
		return req.Get(fullURL)
	case "POST":
		return req.Post(fullURL)
	case "PUT":
		return req.Put(fullURL)
	case "DELETE":
		return req.Delete(fullURL)
	case "PATCH":
		return req.Patch(fullURL)
	default:
		return nil, fmt.Errorf("unsupported HTTP method: %s", test.Request.Method)
	}
}

// validateResponse checks if the HTTP response matches the expected outcome.
func (gd *Grader) validateResponse(test Test, resp *resty.Response) error {
	if uint(resp.StatusCode()) != test.Response.StatusCode {
		return fmt.Errorf("status code mismatch: expected %d, got %d", test.Response.StatusCode, resp.StatusCode())
	}

	// TODO
	// Process headers for variables
	// for _, header := range test.Response.Headers {
	// 	gd.processHeaderForVariables(header, resp.Header())
	// }

	if test.Response.ResBody != "" {
		var actualBody, expectedBody map[string]interface{}
		if err := json.Unmarshal(resp.Body(), &actualBody); err != nil {
			return fmt.Errorf("failed to unmarshal actual response body: %w", err)
		}
		if err := json.Unmarshal([]byte(test.Response.ResBody), &expectedBody); err != nil {
			return fmt.Errorf("failed to unmarshal expected response body: %w", err)
		}
		r, err := gd.compareJSON(actualBody, expectedBody)
		if !r {
			return err
		}
	}
	return nil
}

// updateProjectResultStatus updates the final status of a project result based on its sections.
func (gd *Grader) updateProjectResultStatus(db *gorm.DB, projectResult *ProjectResult) {
	var failedSections int64
	db.Model(&SectionResult{}).Where("project_result_id = ? AND status = ?", projectResult.ID, StatusFailed).Count(&failedSections)

	if failedSections == 0 {
		projectResult.Status = StatusPassed
		projectResult.Message = "All sections passed."
	} else {
		projectResult.Status = StatusFailed
		projectResult.Message = "Some sections failed."
	}
}

// updateSectionResultStatus updates the final status of a section result based on its scenarios.
func (gd *Grader) updateSectionResultStatus(db *gorm.DB, sectionResult *SectionResult) {
	var totalScenarios, passedScenarios int64
	db.Model(&ScenarioResult{}).Where("section_result_id = ?", sectionResult.ID).Count(&totalScenarios)
	db.Model(&ScenarioResult{}).Where("section_result_id = ? AND status = ?", sectionResult.ID, StatusPassed).Count(&passedScenarios)

	sectionResult.TotalScenarios = uint(totalScenarios)
	sectionResult.PassedScenarios = uint(passedScenarios)

	if totalScenarios > 0 {
		sectionResult.Score = (float64(passedScenarios) / float64(totalScenarios)) * 100
	} else {
		sectionResult.Score = 100 // Or 0, depending on desired behavior for no scenarios
	}

	if passedScenarios == totalScenarios {
		sectionResult.Status = StatusPassed
		sectionResult.Message = fmt.Sprintf("All %d scenarios passed. Score: %.2f%%", totalScenarios, sectionResult.Score)
	} else {
		sectionResult.Status = StatusFailed
		sectionResult.Message = fmt.Sprintf("%d/%d scenarios passed. Score: %.2f%%", passedScenarios, totalScenarios, sectionResult.Score)
	}
}

// updateScenarioResultStatus updates the final status of a scenario result based on its tests.
func (gd *Grader) updateScenarioResultStatus(db *gorm.DB, scenarioResult *ScenarioResult) {
	var totalTests, failedTests int64
	db.Model(&TestResult{}).Where("scenario_result_id = ?", scenarioResult.ID).Count(&totalTests)
	db.Model(&TestResult{}).Where("scenario_result_id = ? AND status = ?", scenarioResult.ID, StatusFailed).Count(&failedTests)

	if failedTests == 0 {
		scenarioResult.Status = StatusPassed
		scenarioResult.Message = "All tests passed."
	} else {
		scenarioResult.Status = StatusFailed
		scenarioResult.Message = fmt.Sprintf("%d/%d tests failed.", failedTests, totalTests)
	}
}

// compareJSON compares two JSON objects (maps) field by field.
func (gd *Grader) compareJSON(actual, expected map[string]interface{}) (bool, error) {
	for key, expectedVal := range expected {
		// Handle variable substitution in expected keys
		substitutedKey := gd.substituteVariables(key)
		actualVal, ok := actual[substitutedKey]
		if !ok {
			return false, fmt.Errorf("key '%s' not found in actual response", substitutedKey)
		}

		// Compare values recursively
		r, err := gd.jsonValueEquals(actualVal, expectedVal)
		if !r {
			return false, fmt.Errorf("error at key '%s': %w", substitutedKey, err)
		}
	}
	return true, nil
}

// jsonValueEquals compares two interface{} values, handling different JSON types recursively.
func (gd *Grader) jsonValueEquals(val1, val2 interface{}) (bool, error) {
	if val1 == nil && val2 == nil {
		return true, nil
	}
	if val1 == nil || val2 == nil {
		return false, fmt.Errorf("unexpected null value")
	}

	if expectedStr, ok := val2.(string); ok {
		matches := variableRegex.FindStringSubmatch(expectedStr)
		if len(matches) > 1 {
			varName := matches[1]
			gd.variables[varName] = val1
			return true, nil
		}
		val2 = gd.substituteVariables(expectedStr)
	}

	switch v1 := val1.(type) {
	case map[string]interface{}:
		v2, ok := val2.(map[string]interface{})
		if !ok {
			return false, fmt.Errorf("type mismatch: expected a JSON object")
		}
		return gd.compareJSON(v1, v2)
	case []interface{}:
		v2, ok := val2.([]interface{})
		if !ok {
			return false, fmt.Errorf("type mismatch: expected a JSON array")
		}
		return gd.compareJSONArray(v1, v2)
	default:
		// For primitive types, use fmt.Sprintf for a robust comparison
		if fmt.Sprintf("%v", val1) != fmt.Sprintf("%v", val2) {
			return false, fmt.Errorf("value mismatch: expected '%v', got '%v'", val2, val1)
		}
		return true, nil
	}
}

// compareJSONArray compares two JSON arrays, checking if all elements in expected are present in actual, regardless of order.
func (gd *Grader) compareJSONArray(actual, expected []interface{}) (bool, error) {
	if len(actual) != len(expected) {
		return false, fmt.Errorf("array length mismatch: expected %d, got %d", len(expected), len(actual))
	}

	actualCopy := make([]interface{}, len(actual))
	copy(actualCopy, actual)

	for _, expectedItem := range expected {
		found := false
		for i, actualItem := range actualCopy {
			if actualItem == nil {
				continue
			}
			if eq, _ := gd.jsonValueEquals(actualItem, expectedItem); eq {
				actualCopy[i] = nil // Mark as used
				found = true
				break
			}
		}
		if !found {
			return false, fmt.Errorf("expected item not found in array: %v", expectedItem)
		}
	}
	return true, nil
}

// substituteVariables replaces placeholders like {{var_name}} with their stored values.
func (gd *Grader) substituteVariables(input string) string {
	return substitutionRegex.ReplaceAllStringFunc(input, func(placeholder string) string {
		varName := substitutionRegex.FindStringSubmatch(placeholder)[1]
		if val, ok := gd.variables[varName]; ok {
			return fmt.Sprintf("%v", val)
		}
		return placeholder // Return the original placeholder if the variable is not found
	})
}

// processHeaderForVariables extracts variables from response headers.
func (gd *Grader) processHeaderForVariables(expectedHeader THeader, actualHeaders map[string][]string) {
	substitutedKey := gd.substituteVariables(expectedHeader.Key)

	// Handle variable capture in value
	matches := variableRegex.FindStringSubmatch(expectedHeader.Value)
	if len(matches) > 1 {
		varName := matches[1]
		if values, ok := actualHeaders[substitutedKey]; ok && len(values) > 0 {
			gd.variables[varName] = values[0]
		}
		return
	}

	// Standard header validation after substitution
	substitutedValue := gd.substituteVariables(expectedHeader.Value)
	if values, ok := actualHeaders[substitutedKey]; ok {
		found := false
		for _, v := range values {
			if v == substitutedValue {
				found = true
				break
			}
		}
		if !found {
			// Optional: Log or handle header mismatch
		}
	} else {
		// Optional: Log or handle missing header
	}
}

package core

import "fmt"

type TestResult struct {
	Success bool
	Message string
}

func (testResult *TestResult) Set(success bool, message string) {
	testResult.Success = success
	testResult.Message = message
}

func ValidateParams(input *ApiResourceInput, requiredParams []string) *TestResult {
	message := "Missing params: "
	missingParams := []string{}
	if len(input.Query) == 0 {
		for _, param := range requiredParams {
			message += fmt.Sprintf(" %v", param)
		}
		return &TestResult{Success: false, Message: message}
	} else {
		for _, param := range requiredParams {
			if input.Query.Get(param) == "" {
				missingParams = append(missingParams, param)
			}
		}
		if len(missingParams) > 0 {
			for _, param := range missingParams {
				message += fmt.Sprintf(" %v", param)
			}
			return &TestResult{Success: false, Message: message}
		} else {
			return &TestResult{Success: true, Message: ""}
		}
	}
}

const SourceIdError = "Missing or Invalid sourceId"
const InvalidConnectionError = "Your connection configuration is invalid."
const UnsetConnectionError = "Your connection configuration is not set."
const UnmarshallingError = "There was a problem unmarshalling the response"

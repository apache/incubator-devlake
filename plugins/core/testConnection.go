package core

type TestResult struct {
	Success bool
	Message string
}

func (testResult *TestResult) Set(success bool, message string) {
	testResult.Success = success
	testResult.Message = message
}

const ReadError = "There was a problem reading source data."
const SourceIdError = "Missing or Invalid sourceId"
const InvalidConnectionError = "Your connection configuration is invalid."
const UnsetConnectionError = "Your connection configuration is not set."
const UnmarshallingError = "There was a problem unmarshalling the response"

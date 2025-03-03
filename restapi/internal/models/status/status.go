package status

const (
	Err     = "Error"
	Success = "Success"
)

// Status of response
type Status struct {
	Status string `json:"status"` // Error or Success
	Error  string `json:"error,omitempty"`
}

// OK returns a success status
func OK() Status {
	return Status{
		Status: Success,
	}
}

// Error returns an error status
func Error(err string) Status {
	return Status{
		Status: Err,
		Error:  err,
	}
}

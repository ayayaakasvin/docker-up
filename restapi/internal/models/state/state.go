package state

const (
	Err     = "Error"
	Success = "Success"
)

// State of response
type State struct {
	Status string `json:"status"` // Error or Success
	Error  string `json:"error,omitempty"`
}

// OK returns a success state
func OK() State {
	return State{
		Status: Success,
	}
}

// Error returns an error state
func Error(err string) State {
	return State{
		Status: Err,
		Error:  err,
	}
}

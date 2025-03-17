package state

const (
	Success = "Success"
	Err 	= "Error"
)

type State struct {
	Status 	string	`json:"status"`
	Error 	string 	`json:"error,omitempty"`
}

func Ok () State {
	return State{
		Status: Success,
	}
}

func Error (errorMsg string) State {
	return State{
		Status: Err,
		Error: errorMsg,
	}
}
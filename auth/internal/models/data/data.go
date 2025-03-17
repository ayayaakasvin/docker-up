package data

type Data map[string]any

func NewDate () Data {
	return make(Data)
}
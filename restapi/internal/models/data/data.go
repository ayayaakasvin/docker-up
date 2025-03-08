package data

type Data map[string]any

func NewData() Data {
	return make(Data)
}
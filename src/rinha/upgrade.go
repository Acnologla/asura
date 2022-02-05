package rinha

type Upgrade struct {
	Name   string    `json:"name"`
	Value  string    `json:"value"`
	Childs []Upgrade `json:"childs"`
}

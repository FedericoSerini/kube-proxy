package models

type Pod struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Status    string `json:"status"`
}

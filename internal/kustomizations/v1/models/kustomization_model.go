package models

type Kustomization struct {
	Namespace             string `json:"namespace"`
	Name                  string `json:"name"`
	Status                string `json:"status"`
	LastAppliedRevision   string `json:"last_applied_revision"`
	LastAttemptedRevision string `json:"last_attemped_revision"`
	Message               string `json:"message"`
}

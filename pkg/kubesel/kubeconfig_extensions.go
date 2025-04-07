package kubesel

const (
	kcextApiVersion           = "dev.eth-p.kubesel/v1"
	kcextManagedByKubeselKind = "ManagedByKubesel"
)

type kcextManagedByKubesel struct {
	SessionOwner `json:"owner"`
}

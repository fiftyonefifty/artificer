package contracts

import (
	"artificer/pkg/api/models"
)

type IClientStore interface {
	GetClient(id string) (found bool, client models.Client)
}

package contracts

import (
	"artificer/pkg/client/models"
)

type IClientStore interface {
	GetClient(id string) (found bool, client models.Client)
}

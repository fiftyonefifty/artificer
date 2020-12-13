package contracts

import (
	"artificer/pkg/client/models"
	"context"
)

type IClientStore interface {
	GetClient(ctx context.Context, id string) (found bool, client models.Client)
}

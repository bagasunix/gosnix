package vehicle

import (
	"context"

	"github.com/gofrs/uuid"

	"github.com/bagasunix/gosnix/internal/domain/base"
)

type Common interface {
	Save(ctx context.Context, v *Vehicle) error
	FindByCustomer(ctx context.Context, customerID uuid.UUID) (result base.SliceResult[[]*Vehicle])
	FindByParam(ctx context.Context, params map[string]interface{}) (result base.SingleResult[*Vehicle])
	Update(ctx context.Context, v *Vehicle) error
	Delete(ctx context.Context, id uuid.UUID) error
}
type Repository interface {
	base.Repository
	Common
}

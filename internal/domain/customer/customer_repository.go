package customer

import (
	"context"

	"github.com/bagasunix/gosnix/internal/domain/base"
)

type Common interface {
	Save(ctx context.Context, m *Customer) error
	SaveTx(ctx context.Context, tx any, m *Customer) error
	Delete(ctx context.Context, id int) error
	Updates(ctx context.Context, id int, m *Customer) error
	FindAll(ctx context.Context, limit int, offset int, search string) (result base.SliceResult[*Customer])
	CountCustomer(ctx context.Context, search string) (int, error)

	FindByParams(ctx context.Context, params map[string]interface{}) (result base.SingleResult[*Customer])
}
type Repository interface {
	base.Repository
	Common
}

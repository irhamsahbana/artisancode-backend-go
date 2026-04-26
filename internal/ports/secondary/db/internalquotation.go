package repository

import (
	"context"

	"codebase-app/internal/entity/coreentity"
)

type InternalQuotationRepository interface {
	GetQuotations(ctx context.Context, filter coreentity.InternalCommerceListFilter) ([]coreentity.InternalQuotation, int, error)
	CreateQuotation(ctx context.Context, data coreentity.InternalQuotation) (*coreentity.InternalQuotation, error)
	GetQuotation(ctx context.Context, tenantID, id string) (*coreentity.InternalQuotation, error)
	ApproveQuotation(ctx context.Context, tenantID, id string) (*coreentity.InternalQuotation, error)
	MarkQuotationConverted(ctx context.Context, tenantID, quotationID, orderID string) (*coreentity.InternalQuotation, error)
}

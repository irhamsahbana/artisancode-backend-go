package mapper

import (
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/entity/restentity"
)

func InternalClientFromCoreToRest(item coreentity.InternalClient) restentity.InternalClientResource {
	return restentity.InternalClientResource{
		ID:          item.ID,
		Name:        item.Name,
		Code:        item.Code,
		OwnerNames:  item.OwnerNames,
		OwnerEmails: item.OwnerEmails,
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	}
}

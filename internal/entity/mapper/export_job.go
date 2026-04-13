package mapper

import (
	"encoding/json"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/entity/restentity"
)

func ExportJobFromCoreToRest(item coreentity.ExportJob) restentity.ExportJob {
	return restentity.ExportJob{
		ID:              item.ID,
		ResourceType:    item.ResourceType,
		ResourceLabel:   item.ResourceLabel,
		Format:          string(item.Format),
		Status:          string(item.Status),
		RequestedByName: item.RequestedByName,
		ErrorMessage:    item.ErrorMessage,
		StartedAt:       item.StartedAt,
		CompletedAt:     item.CompletedAt,
		ExpiresAt:       item.ExpiresAt,
		CreatedAt:       item.CreatedAt,
		DownloadURL:     item.DownloadURL,
	}
}

func ExportJobParamsToJSON(req restentity.CreateExportJobReq) (string, error) {
	payload := map[string]any{
		"q":                req.Q,
		"employee_id":      req.EmployeeID,
		"type":             req.Type,
		"source":           req.Source,
		"status":           req.Status,
		"selfie_status":    req.SelfieStatus,
		"org_unit_id":      req.OrgUnitID,
		"branch_id":        req.BranchID,
		"work_location_id": req.WorkLocationID,
		"exception_type":   req.ExceptionType,
		"attendance_date":  req.AttendanceDay,
		"date_from":        req.DateFrom,
		"date_to":          req.DateTo,
	}

	raw, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	return string(raw), nil
}

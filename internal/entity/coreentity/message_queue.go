package coreentity

type CleanupProcessedMessageQueueReq struct {
	Before string
	Limit  int
}

type CleanupProcessedMessageQueueResp struct {
	Deleted int
}

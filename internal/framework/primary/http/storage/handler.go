package handler

import (
	"codebase-app/internal/middleware"
	corePorts "codebase-app/internal/ports/core"

	"github.com/gofiber/fiber/v3"
)

type storageHandler struct {
	core corePorts.StorageCore
}

func NewStorageHandler(core corePorts.StorageCore) *storageHandler {
	return &storageHandler{
		core: core,
	}
}

func (h *storageHandler) Register(router fiber.Router) {
	protected := router.Group("/", middleware.Auth)
	protected.Post("/upload-url", h.createUploadURL)
	protected.Post("/upload", h.uploadFile)
	protected.Delete("/*", h.deleteFile)
	protected.Get("/", h.listFiles)
	router.Get("/private/*", middleware.ValidateSignedURL, h.getPrivateFile)
}

package handler

import (
	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/middleware"
	corePorts "codebase-app/internal/ports/core"
	"codebase-app/pkg/errmsg"
	"codebase-app/pkg/response"
	"errors"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
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
	protected.Post("/upload", h.uploadFile)
	protected.Delete("/:filename", h.deleteFile)
	protected.Get("/", h.listFiles)
	router.Get("/private/:filename", middleware.ValidateSignedURL, h.getPrivateFile)
}

func (h *storageHandler) uploadFile(c *fiber.Ctx) error {
	var (
		ctx = c.UserContext()
	)

	file, err := c.FormFile("file")
	if err != nil {
		if errors.Is(err, fasthttp.ErrMissingFile) {
			message := errmsg.NewCustomErrors(http.StatusBadRequest).
				Add("file", "file is required").
				SetMessage("file is required")
			return c.Status(message.Code).JSON(response.Error(message))
		}
		log.Ctx(ctx).Error().Err(err).Msg("Failed to parse form file")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	req := &coreentity.UploadFileReq{
		File: file,
	}
	req.Filename = c.FormValue("filename")
	req.IsPublic = parseBool(c.FormValue("is_public"))
	req.GeneratePresignedURL = parseBool(c.FormValue("generate_presigned_url"))
	if folder := c.FormValue("folder"); folder != "" {
		req.Folder = common.S3Folder(folder)
	}

	resp, err := h.core.UploadFile(ctx, req)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to upload file")
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(resp, ""))
}

func (h *storageHandler) deleteFile(c *fiber.Ctx) error {
	var (
		ctx      = c.UserContext()
		filename = c.Params("filename")
	)

	req := &coreentity.DeleteFileReq{
		Filename: filename,
	}

	if err := h.core.DeleteFile(ctx, req); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to delete file")
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(nil, ""))
}

func (h *storageHandler) listFiles(c *fiber.Ctx) error {
	ctx := c.UserContext()
	resp, err := h.core.ListFiles(ctx)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to list files")
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(resp, ""))
}

func (h *storageHandler) getPrivateFile(c *fiber.Ctx) error {
	var (
		ctx      = c.UserContext()
		filename = c.Params("filename")
		tenantID = c.Query("tenant_id")
		folder   = c.Query("folder")
	)

	filter := coreentity.FileFilter{
		TenantID: tenantID,
		Folder:   common.S3Folder(folder),
		Filename: filename,
	}

	url, err := h.core.GetFileURL(ctx, filter)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any("filename", filename).Msg("Failed to get file URL")
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(url, ""))
}

func parseBool(value string) bool {
	if value == "" {
		return false
	}
	result, err := strconv.ParseBool(value)
	if err != nil {
		return false
	}
	return result
}

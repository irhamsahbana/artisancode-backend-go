package handler

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"codebase-app/internal/adapter"
	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/config"
	"codebase-app/pkg/errmsg"
	"codebase-app/pkg/jwthandler"
	pkgvalidator "codebase-app/pkg/validator"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

type storageCoreStub struct {
	createUploadURLFunc func(ctx context.Context, req coreentity.PresignUploadURLReq) (*coreentity.PresignUploadURLResp, error)
	uploadFileFunc      func(ctx context.Context, req *coreentity.UploadFileReq) (*coreentity.UploadFileResp, error)
	deleteFileFunc      func(ctx context.Context, req *coreentity.DeleteFileReq) error
	listFilesFunc       func(ctx context.Context) ([]types.Object, error)
	getFileURLFunc      func(ctx context.Context, filter coreentity.FileFilter) (string, error)
}

func (s *storageCoreStub) UploadFile(ctx context.Context, req *coreentity.UploadFileReq) (*coreentity.UploadFileResp, error) {
	if s.uploadFileFunc == nil {
		return nil, errors.New("unexpected UploadFile call")
	}

	return s.uploadFileFunc(ctx, req)
}

func (s *storageCoreStub) CreateUploadURL(ctx context.Context, req coreentity.PresignUploadURLReq) (*coreentity.PresignUploadURLResp, error) {
	if s.createUploadURLFunc == nil {
		return nil, errors.New("unexpected CreateUploadURL call")
	}

	return s.createUploadURLFunc(ctx, req)
}

func (s *storageCoreStub) DeleteFile(ctx context.Context, req *coreentity.DeleteFileReq) error {
	if s.deleteFileFunc == nil {
		return errors.New("unexpected DeleteFile call")
	}

	return s.deleteFileFunc(ctx, req)
}

func (s *storageCoreStub) ListFiles(ctx context.Context) ([]types.Object, error) {
	if s.listFilesFunc == nil {
		return nil, errors.New("unexpected ListFiles call")
	}

	return s.listFilesFunc(ctx)
}

func (s *storageCoreStub) GetFileURL(ctx context.Context, filter coreentity.FileFilter) (string, error) {
	if s.getFileURLFunc == nil {
		return "", errors.New("unexpected GetFileURL call")
	}

	return s.getFileURLFunc(ctx, filter)
}

func (s *storageCoreStub) CleanupExpiredFiles(ctx context.Context, req coreentity.CleanupExpiredFilesReq) (*coreentity.CleanupExpiredFilesResp, error) {
	return nil, errors.New("unexpected CleanupExpiredFiles call")
}

func TestCreateUploadURL(t *testing.T) {
	setupTestValidator()

	t.Run("returns bad request when body is invalid json", func(t *testing.T) {
		h := NewStorageHandler(&storageCoreStub{})
		app := newTestApp(t, func(app *fiber.App) {
			app.Post("/storage/upload-url", h.createUploadURL)
		})

		req := httptest.NewRequest(http.MethodPost, "/storage/upload-url", bytes.NewBufferString("{"))
		req.Header.Set("Content-Type", fiber.MIMEApplicationJSON)

		resp, body := performJSONRequest(t, app, req)

		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
		require.Equal(t, false, body["success"])
	})

	t.Run("returns bad request when validation fails", func(t *testing.T) {
		h := NewStorageHandler(&storageCoreStub{})
		app := newTestApp(t, func(app *fiber.App) {
			app.Post("/storage/upload-url", h.createUploadURL)
		})

		req := httptest.NewRequest(http.MethodPost, "/storage/upload-url", bytes.NewBufferString(`{}`))
		req.Header.Set("Content-Type", fiber.MIMEApplicationJSON)

		resp, body := performJSONRequest(t, app, req)

		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
		require.Equal(t, false, body["success"])
		errorsBody := body["errors"].(map[string]any)
		require.Contains(t, errorsBody, "filename")
		require.Contains(t, errorsBody, "content_type")
		require.Contains(t, errorsBody, "folder")
	})

	t.Run("returns mapped response on success", func(t *testing.T) {
		var gotReq coreentity.PresignUploadURLReq

		h := NewStorageHandler(&storageCoreStub{
			createUploadURLFunc: func(ctx context.Context, req coreentity.PresignUploadURLReq) (*coreentity.PresignUploadURLResp, error) {
				gotReq = req
				return &coreentity.PresignUploadURLResp{
					FileID:   "file-1",
					Filename: "public/docs/report.pdf",
					URL:      "https://upload.example.com",
					Method:   http.MethodPut,
					Headers:  map[string]string{"Content-Type": "application/pdf"},
				}, nil
			},
		})
		app := newTestApp(t, func(app *fiber.App) {
			app.Post("/storage/upload-url", h.createUploadURL)
		})

		req := httptest.NewRequest(http.MethodPost, "/storage/upload-url", bytes.NewBufferString(`{
			"filename":"report.pdf",
			"original_filename":"Q2 Report.pdf",
			"content_type":"application/pdf",
			"folder":"documents",
			"is_public":true
		}`))
		req.Header.Set("Content-Type", fiber.MIMEApplicationJSON)

		resp, body := performJSONRequest(t, app, req)

		require.Equal(t, http.StatusOK, resp.StatusCode)
		require.Equal(t, "tenant-1", gotReq.TenantID)
		require.Equal(t, "user-1", gotReq.CreatedBy)
		require.Equal(t, common.UserContext{UserID: "user-1", TenantID: "tenant-1"}, gotReq.UserCtx)
		require.Equal(t, "report.pdf", gotReq.Filename)
		require.NotNil(t, gotReq.OriginalFilename)
		require.Equal(t, "Q2 Report.pdf", *gotReq.OriginalFilename)
		require.Equal(t, "application/pdf", gotReq.ContentType)
		require.Equal(t, common.S3Folder("documents"), gotReq.Folder)
		require.True(t, gotReq.IsPublic)

		data := body["data"].(map[string]any)
		require.Equal(t, "file-1", data["file_id"])
		require.Equal(t, "public/docs/report.pdf", data["object_key"])
		require.Equal(t, "https://upload.example.com", data["upload_url"])
		require.Equal(t, http.MethodPut, data["method"])
	})

	t.Run("returns core error response", func(t *testing.T) {
		h := NewStorageHandler(&storageCoreStub{
			createUploadURLFunc: func(ctx context.Context, req coreentity.PresignUploadURLReq) (*coreentity.PresignUploadURLResp, error) {
				return nil, errmsg.NewCustomErrors(http.StatusConflict).
					SetMessage("upload conflict").
					Add("filename", "already exists")
			},
		})
		app := newTestApp(t, func(app *fiber.App) {
			app.Post("/storage/upload-url", h.createUploadURL)
		})

		req := httptest.NewRequest(http.MethodPost, "/storage/upload-url", bytes.NewBufferString(`{
			"filename":"report.pdf",
			"content_type":"application/pdf",
			"folder":"documents"
		}`))
		req.Header.Set("Content-Type", fiber.MIMEApplicationJSON)

		resp, body := performJSONRequest(t, app, req)

		require.Equal(t, http.StatusConflict, resp.StatusCode)
		require.Equal(t, false, body["success"])
		require.Equal(t, "upload conflict", body["message"])
	})
}

func TestUploadFile(t *testing.T) {
	t.Run("returns bad request when file is missing", func(t *testing.T) {
		h := NewStorageHandler(&storageCoreStub{})
		app := newTestApp(t, func(app *fiber.App) {
			app.Post("/storage/upload", h.uploadFile)
		})

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		require.NoError(t, writer.WriteField("filename", "avatar.png"))
		require.NoError(t, writer.Close())

		req := httptest.NewRequest(http.MethodPost, "/storage/upload", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		resp, payload := performJSONRequest(t, app, req)

		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
		require.Equal(t, false, payload["success"])
	})

	t.Run("passes parsed form fields to core", func(t *testing.T) {
		var gotReq *coreentity.UploadFileReq

		h := NewStorageHandler(&storageCoreStub{
			uploadFileFunc: func(ctx context.Context, req *coreentity.UploadFileReq) (*coreentity.UploadFileResp, error) {
				gotReq = req
				return &coreentity.UploadFileResp{
					Filename: "uploads/avatar.png",
					URL:      "https://cdn.example.com/uploads/avatar.png",
				}, nil
			},
		})
		app := newTestApp(t, func(app *fiber.App) {
			app.Post("/storage/upload", h.uploadFile)
		})

		req := newMultipartUploadRequest(
			t,
			"/storage/upload",
			map[string]string{
				"filename":               "avatar.png",
				"is_public":              "true",
				"generate_presigned_url": "1",
				"folder":                 "avatars",
			},
			"file",
			"avatar.png",
			[]byte("image-bytes"),
		)

		resp, body := performJSONRequest(t, app, req)

		require.Equal(t, http.StatusOK, resp.StatusCode)
		require.NotNil(t, gotReq)
		require.Equal(t, "avatar.png", gotReq.Filename)
		require.True(t, gotReq.IsPublic)
		require.True(t, gotReq.GeneratePresignedURL)
		require.Equal(t, common.S3Folder("avatars"), gotReq.Folder)
		require.NotNil(t, gotReq.File)
		require.Equal(t, "avatar.png", gotReq.File.Filename)

		data := body["data"].(map[string]any)
		require.Equal(t, "uploads/avatar.png", data["Filename"])
		require.Equal(t, "https://cdn.example.com/uploads/avatar.png", data["URL"])
	})

	t.Run("returns core error response", func(t *testing.T) {
		h := NewStorageHandler(&storageCoreStub{
			uploadFileFunc: func(ctx context.Context, req *coreentity.UploadFileReq) (*coreentity.UploadFileResp, error) {
				return nil, errmsg.NewCustomErrors(http.StatusUnprocessableEntity).
					SetMessage("upload rejected")
			},
		})
		app := newTestApp(t, func(app *fiber.App) {
			app.Post("/storage/upload", h.uploadFile)
		})

		req := newMultipartUploadRequest(
			t,
			"/storage/upload",
			map[string]string{"filename": "avatar.png"},
			"file",
			"avatar.png",
			[]byte("image-bytes"),
		)

		resp, body := performJSONRequest(t, app, req)

		require.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
		require.Equal(t, false, body["success"])
		require.Equal(t, "upload rejected", body["message"])
	})
}

func TestDeleteFile(t *testing.T) {
	t.Run("passes tenant and wildcard path to core", func(t *testing.T) {
		var gotReq *coreentity.DeleteFileReq

		h := NewStorageHandler(&storageCoreStub{
			deleteFileFunc: func(ctx context.Context, req *coreentity.DeleteFileReq) error {
				gotReq = req
				return nil
			},
		})
		app := newTestApp(t, func(app *fiber.App) {
			app.Delete("/storage/*", h.deleteFile)
		})

		req := httptest.NewRequest(http.MethodDelete, "/storage/private/docs/report.pdf", nil)

		resp, body := performJSONRequest(t, app, req)

		require.Equal(t, http.StatusOK, resp.StatusCode)
		require.Nil(t, body["data"])
		require.Equal(t, &coreentity.DeleteFileReq{
			TenantID: "tenant-1",
			Filename: "private/docs/report.pdf",
		}, gotReq)
	})

	t.Run("returns core error response", func(t *testing.T) {
		h := NewStorageHandler(&storageCoreStub{
			deleteFileFunc: func(ctx context.Context, req *coreentity.DeleteFileReq) error {
				return errmsg.NewCustomErrors(http.StatusNotFound).
					SetMessage("file not found")
			},
		})
		app := newTestApp(t, func(app *fiber.App) {
			app.Delete("/storage/*", h.deleteFile)
		})

		req := httptest.NewRequest(http.MethodDelete, "/storage/private/docs/missing.pdf", nil)

		resp, body := performJSONRequest(t, app, req)

		require.Equal(t, http.StatusNotFound, resp.StatusCode)
		require.Equal(t, false, body["success"])
		require.Equal(t, "file not found", body["message"])
	})
}

func TestListFiles(t *testing.T) {
	t.Run("returns list from core", func(t *testing.T) {
		called := false

		h := NewStorageHandler(&storageCoreStub{
			listFilesFunc: func(ctx context.Context) ([]types.Object, error) {
				called = true
				return []types.Object{{Key: aws.String("uploads/avatar.png")}}, nil
			},
		})
		app := newTestApp(t, func(app *fiber.App) {
			app.Get("/storage", h.listFiles)
		})

		req := httptest.NewRequest(http.MethodGet, "/storage", nil)

		resp, body := performJSONRequest(t, app, req)

		require.Equal(t, http.StatusOK, resp.StatusCode)
		require.True(t, called)
		require.Equal(t, true, body["success"])
		data := body["data"].([]any)
		require.Len(t, data, 1)
	})

	t.Run("returns core error response", func(t *testing.T) {
		h := NewStorageHandler(&storageCoreStub{
			listFilesFunc: func(ctx context.Context) ([]types.Object, error) {
				return nil, errmsg.NewCustomErrors(http.StatusServiceUnavailable).
					SetMessage("storage unavailable")
			},
		})
		app := newTestApp(t, func(app *fiber.App) {
			app.Get("/storage", h.listFiles)
		})

		req := httptest.NewRequest(http.MethodGet, "/storage", nil)

		resp, body := performJSONRequest(t, app, req)

		require.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)
		require.Equal(t, false, body["success"])
		require.Equal(t, "storage unavailable", body["message"])
	})
}

func TestGetPrivateFile(t *testing.T) {
	t.Run("maps wildcard path and query params", func(t *testing.T) {
		var gotFilter coreentity.FileFilter

		h := NewStorageHandler(&storageCoreStub{
			getFileURLFunc: func(ctx context.Context, filter coreentity.FileFilter) (string, error) {
				gotFilter = filter
				return "https://signed.example.com/private/docs/report.pdf", nil
			},
		})
		app := newTestApp(t, func(app *fiber.App) {
			app.Get("/storage/private/*", h.getPrivateFile)
		})

		req := httptest.NewRequest(http.MethodGet, "/storage/private/docs/report.pdf?tenant_id=tenant-1&folder=documents", nil)

		resp, body := performJSONRequest(t, app, req)

		require.Equal(t, http.StatusOK, resp.StatusCode)
		require.Equal(t, coreentity.FileFilter{
			TenantID: "tenant-1",
			Folder:   common.S3Folder("documents"),
			Filename: "private/docs/report.pdf",
		}, gotFilter)
		require.Equal(t, "https://signed.example.com/private/docs/report.pdf", body["data"])
	})

	t.Run("returns core error response", func(t *testing.T) {
		h := NewStorageHandler(&storageCoreStub{
			getFileURLFunc: func(ctx context.Context, filter coreentity.FileFilter) (string, error) {
				return "", errmsg.NewCustomErrors(http.StatusForbidden).
					SetMessage("signature expired")
			},
		})
		app := newTestApp(t, func(app *fiber.App) {
			app.Get("/storage/private/*", h.getPrivateFile)
		})

		req := httptest.NewRequest(http.MethodGet, "/storage/private/docs/report.pdf?tenant_id=tenant-1&folder=documents", nil)

		resp, body := performJSONRequest(t, app, req)

		require.Equal(t, http.StatusForbidden, resp.StatusCode)
		require.Equal(t, false, body["success"])
		require.Equal(t, "signature expired", body["message"])
	})
}

func TestParseBool(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{name: "empty string defaults false", value: "", want: false},
		{name: "true string", value: "true", want: true},
		{name: "numeric truthy string", value: "1", want: true},
		{name: "invalid string defaults false", value: "yes", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, parseBool(tt.value))
		})
	}
}

func TestRegister(t *testing.T) {
	setupTestValidator()
	config.Envs = &config.Config{}
	config.Envs.Guard.JwtPrivateKey = "test-secret"

	h := NewStorageHandler(&storageCoreStub{
		getFileURLFunc: func(ctx context.Context, filter coreentity.FileFilter) (string, error) {
			return "https://signed.example.com/private/docs/report.pdf", nil
		},
	})
	app := fiber.New()
	h.Register(app.Group("/storage"))

	t.Run("protected upload-url route requires auth middleware", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/storage/upload-url", bytes.NewBufferString(`{}`))
		req.Header.Set("Content-Type", fiber.MIMEApplicationJSON)

		resp, body := performJSONRequest(t, app, req)

		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		require.Equal(t, false, body["success"])
		require.Equal(t, "Unauthorized", body["message"])
	})

	t.Run("private route requires valid signed url middleware", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/storage/private/docs/report.pdf", nil)
		req.Header.Set("Authorization", "Bearer "+mustGenerateAuthToken(t))

		resp, body := performJSONRequest(t, app, req)

		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		require.Equal(t, false, body["success"])
		require.Equal(t, "URL not valid", body["message"])
	})

	t.Run("private route reaches handler when signature is valid", func(t *testing.T) {
		expires := time.Now().Add(time.Minute).Unix()
		path := "/storage/private/docs/report.pdf"
		signature := signStorageURL("http://example.com"+path, expires, config.Envs.Guard.JwtPrivateKey)
		url := path + "?tenant_id=tenant-1&folder=documents&expires=" + strconv.FormatInt(expires, 10) + "&signature=" + signature
		req := httptest.NewRequest(http.MethodGet, url, nil)
		req.Host = "example.com"
		req.Header.Set("Authorization", "Bearer "+mustGenerateAuthToken(t))

		resp, body := performJSONRequest(t, app, req)

		require.Equal(t, http.StatusOK, resp.StatusCode)
		require.Equal(t, "https://signed.example.com/private/docs/report.pdf", body["data"])
	})
}

func setupTestValidator() {
	adapter.Adapters = &adapter.Adapter{}
	adapter.Adapters.Sync(adapter.WithValidator(pkgvalidator.NewValidator()))
}

func signStorageURL(basePath string, expires int64, secret string) string {
	data := basePath + strconv.FormatInt(expires, 10)
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func mustGenerateAuthToken(t *testing.T) string {
	t.Helper()

	token, err := jwthandler.GenerateTokenString(jwthandler.CostumClaimsPayload{
		UserID:          "user-1",
		UserName:        "User One",
		TenantID:        "tenant-1",
		TenantName:      "Tenant One",
		Roles:           []string{"owner"},
		TokenExpiration: time.Now().Add(time.Hour),
	})
	require.NoError(t, err)

	return token
}

func newTestApp(t *testing.T, register func(app *fiber.App)) *fiber.App {
	t.Helper()

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		ctx := context.WithValue(c.UserContext(), common.UserContextKeyClaims, common.UserContext{
			UserID:   "user-1",
			TenantID: "tenant-1",
		})
		c.SetUserContext(ctx)
		return c.Next()
	})

	register(app)

	return app
}

func performJSONRequest(t *testing.T, app *fiber.App, req *http.Request) (*http.Response, map[string]any) {
	t.Helper()

	resp, err := app.Test(req, -1)
	require.NoError(t, err)

	bodyBytes, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.NoError(t, resp.Body.Close())

	var body map[string]any
	require.NoError(t, json.Unmarshal(bodyBytes, &body))

	return resp, body
}

func newMultipartUploadRequest(
	t *testing.T,
	target string,
	fields map[string]string,
	fileField string,
	filename string,
	content []byte,
) *http.Request {
	t.Helper()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for key, value := range fields {
		require.NoError(t, writer.WriteField(key, value))
	}

	part, err := writer.CreateFormFile(fileField, filename)
	require.NoError(t, err)

	_, err = part.Write(content)
	require.NoError(t, err)
	require.NoError(t, writer.Close())

	req := httptest.NewRequest(http.MethodPost, target, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	return req
}

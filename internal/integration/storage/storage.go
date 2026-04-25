package integration

import (
	"bytes"
	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/config"
	integrationPorts "codebase-app/internal/ports/integration"
	"codebase-app/pkg"
	"codebase-app/pkg/errmsg"
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/rs/zerolog/log"
)

type storage struct {
	client        *s3.Client
	presignClient *s3.PresignClient
	uploader      *manager.Uploader
	downloader    *manager.Downloader
}

var _ integrationPorts.StorageContract = &storage{}

func NewStorageIntegration(c *s3.Client) *storage {
	return &storage{
		client:        c,
		presignClient: s3.NewPresignClient(c),
		uploader:      manager.NewUploader(c),
		downloader:    manager.NewDownloader(c),
	}
}

func (s *storage) UploadFile(ctx context.Context, req *coreentity.UploadFileReq) (*coreentity.UploadFileResp, error) {
	if req.File == nil {
		return nil, errmsg.NewCustomErrors(400).Add("file", "file is required")
	}

	file, err := req.File.Open()
	if err != nil {
		log.Ctx(ctx).Err(err).Any("payload", req).Msg("error while opening file")
		return nil, err
	}
	defer file.Close()

	req.Filename = pkg.SanitizeFilename(req.File.Filename, true)
	var acl types.ObjectCannedACL

	if !req.IsPublic {
		req.Filename = "private/" + req.Filename
		acl = types.ObjectCannedACLPrivate
	} else {
		req.Filename = "public/" + req.Filename
		acl = types.ObjectCannedACLPublicRead
	}

	result, err := s.uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(config.Envs.Storage.Bucket),
		Key:    aws.String(req.Filename),
		Body:   file,
		ACL:    acl,
	})
	if err != nil {
		log.Ctx(ctx).Err(err).Any("payload", req).Msg("error while uploading file")
		return nil, err
	}

	url := result.Location
	if !req.IsPublic && req.GeneratePresignedURL {
		presignResp, err := s.presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(config.Envs.Storage.Bucket),
			Key:    aws.String(req.Filename),
		})
		if err != nil {
			log.Ctx(ctx).Err(err).Any("payload", req).Msg("error while presigning file")
			return nil, err
		}
		url = presignResp.URL
	}

	return &coreentity.UploadFileResp{
		Filename: req.Filename,
		URL:      url,
	}, nil
}

func (s *storage) UploadBytes(ctx context.Context, req *coreentity.UploadBytesReq) (*coreentity.UploadFileResp, error) {
	if len(req.Body) == 0 {
		return nil, errmsg.NewCustomErrors(400).Add("body", "body is required")
	}

	var acl types.ObjectCannedACL
	if !req.IsPublic {
		acl = types.ObjectCannedACLPrivate
	} else {
		acl = types.ObjectCannedACLPublicRead
	}

	input := &s3.PutObjectInput{
		Bucket: aws.String(config.Envs.Storage.Bucket),
		Key:    aws.String(req.Filename),
		Body:   bytes.NewReader(req.Body),
		ACL:    acl,
	}
	if req.ContentType != "" {
		input.ContentType = aws.String(req.ContentType)
	}

	result, err := s.uploader.Upload(ctx, input)
	if err != nil {
		log.Ctx(ctx).Err(err).Any(common.LogKeyPayload, req).Msg("error while uploading bytes")
		return nil, err
	}

	return &coreentity.UploadFileResp{
		Filename: req.Filename,
		URL:      result.Location,
	}, nil
}

func (s *storage) PresignUploadURL(ctx context.Context, req *coreentity.PresignUploadURLReq) (*coreentity.PresignUploadURLResp, error) {
	if req.Filename == "" {
		return nil, errmsg.NewCustomErrors(400).Add("filename", "filename is required")
	}

	input := &s3.PutObjectInput{
		Bucket: aws.String(config.Envs.Storage.Bucket),
		Key:    aws.String(req.Filename),
	}
	if req.ContentType != "" {
		input.ContentType = aws.String(req.ContentType)
	}

	ps, err := s.presignClient.PresignPutObject(ctx, input, s3.WithPresignExpires(15*time.Minute))
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, req).Msg("error while presigning upload file")
		return nil, err
	}

	return &coreentity.PresignUploadURLResp{
		Filename: req.Filename,
		URL:      ps.URL,
		Method:   "PUT",
		Headers: map[string]string{
			"Content-Type": req.ContentType,
		},
	}, nil
}

func (s *storage) DeleteFile(ctx context.Context, req *coreentity.DeleteFileReq) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(config.Envs.Storage.Bucket),
		Key:    aws.String(req.Filename),
	})
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, req).Msg("error while deleting file")
		return err
	}

	return nil
}

func (s *storage) ListFiles(ctx context.Context) ([]types.Object, error) {
	objects := []types.Object{}
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(config.Envs.Storage.Bucket),
	}

	paginator := s3.NewListObjectsV2Paginator(s.client, input)

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("error while listing files")
			return objects, errmsg.NewCustomErrors(500, errmsg.WithMessage("failed to get page of results"))
		}

		objects = append(objects, page.Contents...)
	}

	return objects, nil
}

func (s *storage) GetFileURL(ctx context.Context, filter coreentity.FileFilter) (string, error) {
	filename := filter.Filename
	if filename == "" {
		return "", errmsg.NewCustomErrors(400).Add("filename", "filename is required")
	}

	expiresIn := filter.PresignExpires
	if expiresIn <= 0 {
		expiresIn = 15 * time.Minute
	}

	ps, err := s.presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(config.Envs.Storage.Bucket),
		Key:    aws.String(filename),
	}, s3.WithPresignExpires(expiresIn))
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, filter).Msg("error while presigning file")
		return "", err
	}

	return ps.URL, nil
}

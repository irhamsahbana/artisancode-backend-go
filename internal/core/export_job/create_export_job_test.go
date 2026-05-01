package core

import (
	"context"
	"errors"
	"testing"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	infraConfig "codebase-app/internal/infrastructure/config"
	integrationMocks "codebase-app/internal/ports/integration/mocks"
	repositoryMocks "codebase-app/internal/ports/secondary/db/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func init() {
	infraConfig.Envs = &infraConfig.Config{}
	infraConfig.Envs.App.Name = "export-job-core-test"
}

func TestCreateExportJob(t *testing.T) {
	ctx := context.Background()
	ownerCtx := common.UserContext{
		UserID:   "user-1",
		TenantID: "tenant-1",
		Roles:    []string{"owner"},
	}
	baseData := coreentity.ExportJobCreate{
		UserCtx:       ownerCtx,
		TenantID:      "tenant-1",
		RequestedBy:   "user-1",
		ResourceType:  "attendance",
		ResourceLabel: "Attendance Report",
		ProcessorKey:  "attendance_logs",
		Format:        coreentity.ExportJobFormatXLSX,
		ParamsJSON:    `{"language":"en"}`,
	}
	createdJob := &coreentity.ExportJob{
		ID:           "export-1",
		TenantID:     "tenant-1",
		RequestedBy:  "user-1",
		ResourceType: "attendance",
		ProcessorKey: "attendance_logs",
		Format:       coreentity.ExportJobFormatXLSX,
		Status:       coreentity.ExportJobStatusPending,
	}

	tests := []struct {
		name      string
		data      coreentity.ExportJobCreate
		setup     func(repo *repositoryMocks.ExportJobRepository, tx *repositoryMocks.Transactor, bus *integrationMocks.MessagePublisher)
		want      *coreentity.ExportJob
		wantError bool
	}{
		{
			name: "success",
			data: baseData,
			setup: func(repo *repositoryMocks.ExportJobRepository, tx *repositoryMocks.Transactor, bus *integrationMocks.MessagePublisher) {
				tx.EXPECT().
					WithinTransaction(mock.Anything, mock.AnythingOfType("func(context.Context) error")).
					RunAndReturn(func(txCtx context.Context, fn func(context.Context) error) error {
						return fn(txCtx)
					})
				repo.EXPECT().
					CreateExportJob(mock.Anything, baseData).
					Return(createdJob, nil)
				bus.EXPECT().
					PublishJSON(mock.Anything, common.MessageTopicExportJobRequested, coreentity.ExportJobRequestedEvent{
						JobID:    createdJob.ID,
						TenantID: createdJob.TenantID,
					}).
					Return(nil)
			},
			want: createdJob,
		},
		{
			name: "validation error when user is not owner or admin",
			data: coreentity.ExportJobCreate{
				UserCtx: common.UserContext{
					UserID:   "user-2",
					TenantID: "tenant-1",
					Roles:    []string{"member"},
				},
				TenantID:     "tenant-1",
				RequestedBy:  "user-2",
				ResourceType: "attendance",
				ProcessorKey: "attendance_logs",
				Format:       coreentity.ExportJobFormatXLSX,
			},
			wantError: true,
		},
		{
			name: "validation error when resource metadata is missing",
			data: coreentity.ExportJobCreate{
				UserCtx:     ownerCtx,
				TenantID:    "tenant-1",
				RequestedBy: "user-1",
				Format:      coreentity.ExportJobFormatXLSX,
			},
			wantError: true,
		},
		{
			name: "dependency error from repository",
			data: baseData,
			setup: func(repo *repositoryMocks.ExportJobRepository, tx *repositoryMocks.Transactor, bus *integrationMocks.MessagePublisher) {
				repoErr := errors.New("create export job failed")
				tx.EXPECT().
					WithinTransaction(mock.Anything, mock.AnythingOfType("func(context.Context) error")).
					RunAndReturn(func(txCtx context.Context, fn func(context.Context) error) error {
						return fn(txCtx)
					})
				repo.EXPECT().
					CreateExportJob(mock.Anything, baseData).
					Return(nil, repoErr)
			},
			wantError: true,
		},
		{
			name: "dependency error from publisher",
			data: baseData,
			setup: func(repo *repositoryMocks.ExportJobRepository, tx *repositoryMocks.Transactor, bus *integrationMocks.MessagePublisher) {
				publishErr := errors.New("publish export job event failed")
				tx.EXPECT().
					WithinTransaction(mock.Anything, mock.AnythingOfType("func(context.Context) error")).
					RunAndReturn(func(txCtx context.Context, fn func(context.Context) error) error {
						return fn(txCtx)
					})
				repo.EXPECT().
					CreateExportJob(mock.Anything, baseData).
					Return(createdJob, nil)
				bus.EXPECT().
					PublishJSON(mock.Anything, common.MessageTopicExportJobRequested, coreentity.ExportJobRequestedEvent{
						JobID:    createdJob.ID,
						TenantID: createdJob.TenantID,
					}).
					Return(publishErr)
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := repositoryMocks.NewExportJobRepository(t)
			tx := repositoryMocks.NewTransactor(t)
			bus := integrationMocks.NewMessagePublisher(t)
			if tt.setup != nil {
				tt.setup(repo, tx, bus)
			}
			core := NewExportJobCore(Config{
				Repo: repo,
				Tx:   tx,
				Bus:  bus,
			})

			got, err := core.CreateExportJob(ctx, tt.data)

			if tt.wantError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestCreateExportJobMissingMessageBus(t *testing.T) {
	core := NewExportJobCore(Config{})

	got, err := core.CreateExportJob(context.Background(), coreentity.ExportJobCreate{
		UserCtx: common.UserContext{
			UserID:   "user-1",
			TenantID: "tenant-1",
			Roles:    []string{"owner"},
		},
		TenantID:     "tenant-1",
		RequestedBy:  "user-1",
		ResourceType: "attendance",
		ProcessorKey: "attendance_logs",
		Format:       coreentity.ExportJobFormatXLSX,
	})

	require.Error(t, err)
	require.Equal(t, (*coreentity.ExportJob)(nil), got)
}

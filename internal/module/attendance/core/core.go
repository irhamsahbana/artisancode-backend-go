package core

import (
	corePorts "codebase-app/internal/ports/core"
	portsRepo "codebase-app/internal/ports/repository"
)

var _ corePorts.AttendanceCore = &attendanceCore{}

type attendanceCore struct {
	repo portsRepo.AttendanceRepository
}

type AttendanceCoreConfig struct {
	Repo portsRepo.AttendanceRepository
}

func NewAttendanceCore(cfg AttendanceCoreConfig) *attendanceCore {
	return &attendanceCore{
		repo: cfg.Repo,
	}
}

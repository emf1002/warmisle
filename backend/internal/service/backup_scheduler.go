// Package service implements the scheduled backup trigger for WarmIsle (暖屿).
// BackupScheduler wraps BackupService and fires ScheduleBackup() at the
// configured daily time via a background goroutine.
package service

import (
	"fmt"
	"log"
	"time"
)

// BackupScheduler manages automatic scheduled backups using a configurable
// daily timer. It wraps BackupService and runs ScheduleBackup() at the
// configured time each day.
type BackupScheduler struct {
	svc     *BackupService
	ticker  *time.Ticker
	stopCh  chan struct{}
	running bool
}

// NewBackupScheduler creates a new BackupScheduler for the given BackupService.
func NewBackupScheduler(svc *BackupService) *BackupScheduler {
	return &BackupScheduler{
		svc:    svc,
		stopCh: make(chan struct{}),
	}
}

// Start begins the scheduling loop. It reads the configured ScheduleTime
// from the cloud drive config, calculates the next occurrence, and spawns
// a background goroutine that triggers scheduled backups daily.
// If ScheduleEnabled is not set in the config, Start logs a message and
// returns without launching the goroutine.
func (s *BackupScheduler) Start() {
	cfg, err := s.svc.GetConfig()
	if err != nil {
		log.Printf("[scheduler] 获取配置失败: %v", err)
		return
	}

	if cfg.ScheduleEnabled == 0 {
		log.Printf("[scheduler] 自动备份未启用")
		return
	}

	if s.running {
		log.Printf("[scheduler] 调度器已在运行中")
		return
	}

	scheduleTime := cfg.ScheduleTime
	if scheduleTime == "" {
		scheduleTime = "03:00"
	}

	nextRun := calcNextRun(scheduleTime)
	log.Printf("[scheduler] 下次备份时间: %s", nextRun.Format("2006-01-02 15:04:05"))

	s.running = true
	s.stopCh = make(chan struct{})

	go func() {
		// Sleep until the first scheduled time
		select {
		case <-time.After(time.Until(nextRun)):
		case <-s.stopCh:
			return
		}

		// Start the daily ticker for subsequent runs
		s.ticker = time.NewTicker(24 * time.Hour)
		defer s.ticker.Stop()

		for {
			select {
			case <-s.ticker.C:
				cfg, err := s.svc.GetConfig()
				if err != nil {
					log.Printf("[scheduler] 获取配置失败，跳过本次备份: %v", err)
					continue
				}
				if cfg.ScheduleEnabled == 0 {
					log.Printf("[scheduler] 自动备份已禁用，跳过本次执行")
					continue
				}
				log.Printf("[scheduler] 开始执行定时备份...")
				s.svc.ScheduleBackup()
		case <-s.stopCh:
				return
			}
		}
	}()
}

// Stop gracefully stops the scheduling loop and releases resources
// (the internal ticker and stop channel).
func (s *BackupScheduler) Stop() {
	if !s.running {
		return
	}
	close(s.stopCh)
	if s.ticker != nil {
		s.ticker.Stop()
		s.ticker = nil
	}
	s.running = false
	log.Printf("[scheduler] 调度器已停止")
}

// Reconfigure stops and restarts the scheduler to pick up configuration
// changes (e.g., schedule enabled/disabled, schedule time).
func (s *BackupScheduler) Reconfigure() {
	log.Printf("[scheduler] 重新加载调度配置...")
	s.Stop()
	s.Start()
}

// calcNextRun parses a "HH:MM" string and returns the next occurrence as a
// time.Time in the local timezone. If the calculated time has already passed
// today, it returns the same time on the following day. Invalid or empty
// inputs silently fall back to 03:00.
func calcNextRun(scheduleTime string) time.Time {
	hour, min := 3, 0
	if _, err := fmt.Sscanf(scheduleTime, "%d:%d", &hour, &min); err != nil {
		hour, min = 3, 0
	}
	if hour < 0 || hour > 23 {
		hour = 3
	}
	if min < 0 || min > 59 {
		min = 0
	}

	now := time.Now()
	nextRun := time.Date(now.Year(), now.Month(), now.Day(), hour, min, 0, 0, time.Local)
	if nextRun.Before(now) {
		nextRun = nextRun.Add(24 * time.Hour)
	}
	return nextRun
}

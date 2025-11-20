package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/larsks/kwctl/pkg/radio/types"
)

// RadioStatus represents the complete radio state
type RadioStatus struct {
	Vfos   [2]VfoStatus `json:"Vfos"`
	PttVfo int          `json:"PttVfo"`
	CtlVfo int          `json:"CtlVfo"`
}

// VfoStatus represents the state of a single VFO
type VfoStatus struct {
	Vfo            types.DisplayVFO `json:"Vfo"`
	ChannelNumber  int              `json:"ChannelNumber"`
	ChannelName    string           `json:"ChannelName"`
	TxPower        string           `json:"TxPower"`
	Mode           string           `json:"Mode"`
	SquelchSetting int              `json:"SquelchSetting"`
	SquelchStatus  int              `json:"SquelchStatus"`
}

// statusUpdate represents a status update or error
type statusUpdate struct {
	status RadioStatus
	err    error
}

// AppModel holds the application state
type AppModel struct {
	status      RadioStatus
	lastUpdate  time.Time
	errorMsg    string
	updateTimer *time.Ticker
	kwctl       *KwCtl
	statusChan  chan statusUpdate
	stopChan    chan struct{}
	logger      *slog.Logger
}

// NewAppModel creates a new application model
func NewAppModel(kwctlCmd string, logger *slog.Logger) *AppModel {
	kwctl, err := NewKwCtl(kwctlCmd, logger)
	if err != nil {
		logger.Error("failed to parse kwctl command", "error", err)
		// Return a model with nil kwctl - errors will be handled in fetchAndSendStatus
		return &AppModel{
			updateTimer: time.NewTicker(1 * time.Second),
			statusChan:  make(chan statusUpdate, 1),
			stopChan:    make(chan struct{}),
			logger:      logger,
		}
	}

	return &AppModel{
		updateTimer: time.NewTicker(1 * time.Second),
		kwctl:       kwctl,
		statusChan:  make(chan statusUpdate, 1), // Buffered to avoid blocking
		stopChan:    make(chan struct{}),
		logger:      logger,
	}
}

// StartPolling starts the background status polling goroutine
func (m *AppModel) StartPolling() {
	go m.pollStatus()
}

// pollStatus runs in a background goroutine and polls radio status
func (m *AppModel) pollStatus() {
	// Do an immediate update
	m.fetchAndSendStatus()

	// Then poll on timer
	for {
		select {
		case <-m.updateTimer.C:
			m.fetchAndSendStatus()
		case <-m.stopChan:
			return
		}
	}
}

// fetchAndSendStatus fetches status and sends it to the channel (non-blocking)
func (m *AppModel) fetchAndSendStatus() {
	// Check if kwctl was initialized successfully
	if m.kwctl == nil {
		select {
		case m.statusChan <- statusUpdate{err: fmt.Errorf("kwctl not initialized")}:
		default:
		}
		return
	}

	// Execute the status command
	output, err := m.kwctl.RunWithOutput("status")
	if err != nil {
		select {
		case m.statusChan <- statusUpdate{err: err}:
		default:
		}
		return
	}

	var status RadioStatus
	if err := json.Unmarshal(output, &status); err != nil {
		select {
		case m.statusChan <- statusUpdate{err: err}:
		default:
		}
		return
	}

	// Send successful status update (non-blocking)
	select {
	case m.statusChan <- statusUpdate{status: status}:
	default:
	}
}

// HandleStatusUpdate processes a status update from the channel
func (m *AppModel) HandleStatusUpdate(update statusUpdate) {
	if update.err != nil {
		m.errorMsg = update.err.Error()
		m.logger.Warn("status update failed", "error", update.err)
		return
	}

	m.status = update.status
	m.lastUpdate = time.Now()
	m.errorMsg = ""
}

// Cleanup stops the update timer and polling goroutine
func (m *AppModel) Cleanup() {
	if m.updateTimer != nil {
		m.updateTimer.Stop()
	}
	close(m.stopChan)
}

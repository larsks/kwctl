package main

import (
	"errors"
	"io"
	"testing"
	"time"

	"go.bug.st/serial"
)

// mockPort implements serial.Port interface for testing
type mockPort struct {
	writeData  []byte
	readData   []byte
	readPos    int
	writeErr   error
	readErr    error
	writeCount int    // Track number of writes
	flushData  []byte // Data to return during flush (first read)
}

func (m *mockPort) Read(p []byte) (n int, err error) {
	if m.readErr != nil {
		return 0, m.readErr
	}

	// During flush phase (after first write, before second write)
	if m.writeCount == 1 && m.readPos < len(m.flushData) {
		// Return flush data byte-by-byte
		n = copy(p, m.flushData[m.readPos:])
		m.readPos += n
		return n, nil
	}

	// If flush phase is complete but no flush data existed
	if m.writeCount == 1 && len(m.flushData) == 0 {
		return 0, io.EOF // Simulate empty buffer
	}

	// After flush phase, return actual command response from readData
	// Calculate offset into readData (readPos includes flushData bytes)
	dataOffset := m.readPos - len(m.flushData)

	if dataOffset >= len(m.readData) {
		return 0, io.EOF
	}
	n = copy(p, m.readData[dataOffset:])
	m.readPos += n
	return n, nil
}

func (m *mockPort) Write(p []byte) (n int, err error) {
	if m.writeErr != nil {
		return 0, m.writeErr
	}
	m.writeData = append(m.writeData, p...)
	m.writeCount++
	return len(p), nil
}

func (m *mockPort) Close() error {
	return nil
}

func (m *mockPort) SetMode(mode *serial.Mode) error {
	return nil
}

func (m *mockPort) SetDTR(dtr bool) error {
	return nil
}

func (m *mockPort) SetRTS(rts bool) error {
	return nil
}

func (m *mockPort) GetModemStatusBits() (*serial.ModemStatusBits, error) {
	return nil, nil
}

func (m *mockPort) SetReadTimeout(t time.Duration) error {
	return nil
}

func (m *mockPort) Drain() error {
	return nil
}

func (m *mockPort) ResetInputBuffer() error {
	return nil
}

func (m *mockPort) ResetOutputBuffer() error {
	return nil
}

func (m *mockPort) Break(d time.Duration) error {
	return nil
}

func TestSendCommand_NoArgs(t *testing.T) {
	mock := &mockPort{
		flushData: []byte("some old data\r"),
		readData:  []byte("ID TM-V71\r"),
	}
	radio := &Radio{
		device: "/dev/null",
		port:   mock,
	}

	result, err := radio.SendCommand("ID")
	if err != nil {
		t.Fatalf("SendCommand failed: %v", err)
	}

	if result != "TM-V71" {
		t.Errorf("expected result 'TM-V71', got '%s'", result)
	}

	// Verify command was sent correctly
	expectedWrite := []byte("\rID\r")
	if string(mock.writeData) != string(expectedWrite) {
		t.Errorf("expected write data '%s', got '%s'", expectedWrite, mock.writeData)
	}
}

func TestSendCommand_WithArgs(t *testing.T) {
	mock := &mockPort{
		flushData: []byte("\r"),
		readData:  []byte("SET OK,123,456\r"),
	}
	radio := &Radio{
		device: "/dev/null",
		port:   mock,
	}

	result, err := radio.SendCommand("SET", "freq", "power", "mode")
	if err != nil {
		t.Fatalf("SendCommand failed: %v", err)
	}

	if result != "OK,123,456" {
		t.Errorf("expected result 'OK,123,456', got '%s'", result)
	}

	// Verify command was sent correctly
	expectedWrite := []byte("\rSET freq,power,mode\r")
	if string(mock.writeData) != string(expectedWrite) {
		t.Errorf("expected write data '%s', got '%s'", expectedWrite, mock.writeData)
	}
}

func TestSendCommand_ResponseWithoutArgs(t *testing.T) {
	mock := &mockPort{
		flushData: []byte("\r"),
		readData:  []byte("CMD\r"),
	}
	radio := &Radio{
		device: "/dev/null",
		port:   mock,
	}

	result, err := radio.SendCommand("CMD")
	if err != nil {
		t.Fatalf("SendCommand failed: %v", err)
	}

	if result != "" {
		t.Errorf("expected empty result, got '%s'", result)
	}
}

func TestSendCommand_WriteError(t *testing.T) {
	mock := &mockPort{
		writeErr: errors.New("write failed"),
	}
	radio := &Radio{
		device: "/dev/null",
		port:   mock,
	}

	_, err := radio.SendCommand("ID")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestSendCommand_ReadError(t *testing.T) {
	mock := &mockPort{
		readErr: errors.New("read failed"),
	}
	radio := &Radio{
		device: "/dev/null",
		port:   mock,
	}

	_, err := radio.SendCommand("ID")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestSendCommand_EmptyBuffer(t *testing.T) {
	// Test that SendCommand doesn't block when the buffer is empty (timeout scenario)
	// No flushData means the buffer flush will get EOF immediately
	mock := &mockPort{
		readData: []byte("ID TM-V71\r"), // Only response data, no stale data
	}
	radio := &Radio{
		device: "/dev/null",
		port:   mock,
	}

	result, err := radio.SendCommand("ID")
	if err != nil {
		t.Fatalf("SendCommand failed: %v", err)
	}

	if result != "TM-V71" {
		t.Errorf("expected result 'TM-V71', got '%s'", result)
	}
}

func TestSendCommand_TimeoutDuringBufferFlush(t *testing.T) {
	// Test that buffer flush handles timeouts gracefully (EOF during flush is OK)
	// No flushData means the flush read will immediately get EOF
	mock := &mockPort{
		readData: []byte("ID TM-V71\r"),
	}
	radio := &Radio{
		device: "/dev/null",
		port:   mock,
	}

	result, err := radio.SendCommand("ID")
	if err != nil {
		t.Fatalf("SendCommand failed: %v", err)
	}

	if result != "TM-V71" {
		t.Errorf("expected result 'TM-V71', got '%s'", result)
	}
}

package event

import "testing"

func TestTelemetryIsDisabled(t *testing.T) {
	Init()
	SetNonInteractive(true)
	SetContinueBySessionID(true)
	SetContinueLastSession(true)
	send("test event", "key", "value")
	Error("test error", "key", "value")
	Flush()

	if GetID() == "" {
		t.Fatal("GetID() returned an empty identifier")
	}
}

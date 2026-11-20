package event

// Telemetry is disabled in the Debian slim build.
// Keep the public event API so the rest of Crush does not need build tags.

func Init() {
	distinctId = getDistinctId()
}

func GetID() string {
	return distinctId
}

func SetNonInteractive(bool) {}

func SetContinueBySessionID(bool) {}

func SetContinueLastSession(bool) {}

// send is intentionally a no-op: this build makes no analytics requests.
func send(string, ...any) {}

func Error(any, ...any) {}

func Flush() {}

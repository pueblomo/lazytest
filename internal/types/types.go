package types

type TestStatus string

const (
	StatusNotStarted TestStatus = "not_run"
	StatusPassed     TestStatus = "passed"
	StatusFailed     TestStatus = "failed"
	StatusSkipped    TestStatus = "skipped"
	StatusRunning    TestStatus = "running"

	IconPassed     = "✓"
	IconFailed     = "✗"
	IconNotStarted = "-"
	IconRunning    = "*"
	IconSkipped    = "⚬"
	IconWatching   = "◉"
)

type TestCase struct {
	Name       string
	Filepath   string
	Output     string
	TestStatus TestStatus
	Watched    watched
	Selected   bool
}

type watched struct {
	IsWatching   bool
	StopWatching func() error
}

func (t *TestCase) FilterValue() string { return t.Name }

func (t *TestCase) TestStatusIcon() string {
	switch t.TestStatus {
	case StatusRunning:
		return IconRunning
	case StatusPassed:
		return IconPassed
	case StatusFailed:
		return IconFailed
	case StatusNotStarted:
		return IconNotStarted
	case StatusSkipped:
		return IconSkipped
	default:
		return "❓ "
	}
}

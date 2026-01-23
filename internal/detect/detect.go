package detect

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pueblomo/lazytest/internal/drivers"
)

func DetectDriver(root string) tea.Cmd {
	return func() tea.Msg {
		for _, driver := range drivers.AllDrivers() {
			detected, err := driver.Detect(root)
			if err != nil {
				continue
			}
			if detected {
				return DriverDetectMsg{
					Driver: driver,
				}
			}
		}

		return DriverDetectMsg{Err: fmt.Errorf("no test driver found")}
	}
}

type DriverDetectMsg struct {
	Driver drivers.Driver
	Err    error
}

package drivers

import (
	"context"

	"github.com/charmbracelet/bubbles/list"
)

type Driver interface {
	Detect(root string) (bool, error)
	Name() string
	DetectTestFiles(ctx context.Context, root string) ([]string, error)
	RunTest(Ctx context.Context, root string, testCase list.Item) error
}

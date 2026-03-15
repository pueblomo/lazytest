package drivers

import (
	"testing"
)

func TestAllDrivers_ReturnsNonEmpty(t *testing.T) {
	drivers := AllDrivers()
	if len(drivers) == 0 {
		t.Error("AllDrivers() should return at least one driver")
	}
}

func TestAllDrivers_ReturnsSlice(t *testing.T) {
	drivers := AllDrivers()
	if drivers == nil {
		t.Error("AllDrivers() should not return nil")
	}
}

func TestAllDrivers_ContainsVitestDriver(t *testing.T) {
	drivers := AllDrivers()

	foundVitest := false
	for _, driver := range drivers {
		if driver.Name() == "vitest" {
			foundVitest = true
			break
		}
	}

	if !foundVitest {
		t.Error("AllDrivers() should contain VitestDriver")
	}
}

func TestAllDrivers_DriversImplementInterface(t *testing.T) {
	drivers := AllDrivers()

	for i, driver := range drivers {
		if driver == nil {
			t.Errorf("Driver at index %d is nil", i)
			continue
		}

		// Test that each driver implements the Driver interface methods
		name := driver.Name()
		if name == "" {
			t.Errorf("Driver at index %d has empty name", i)
		}

		// Verify the driver can be used as a Driver interface
		var _ Driver = driver
	}
}

func TestAllDrivers_ReturnsNewSliceEachTime(t *testing.T) {
	// Call twice and verify they're different slices (not the same reference)
	drivers1 := AllDrivers()
	drivers2 := AllDrivers()

	// Modify first slice
	if len(drivers1) > 0 {
		drivers1[0] = nil
	}

	// Verify second slice is not affected
	if len(drivers2) > 0 && drivers2[0] == nil {
		t.Error("AllDrivers() should return a new slice each time, not a shared reference")
	}
}

func TestAllDrivers_OrderIsConsistent(t *testing.T) {
	drivers1 := AllDrivers()
	drivers2 := AllDrivers()

	if len(drivers1) != len(drivers2) {
		t.Errorf("AllDrivers() returned different lengths: %d vs %d", len(drivers1), len(drivers2))
	}

	for i := 0; i < len(drivers1) && i < len(drivers2); i++ {
		if drivers1[i].Name() != drivers2[i].Name() {
			t.Errorf("AllDrivers() returned different order at index %d: %s vs %s",
				i, drivers1[i].Name(), drivers2[i].Name())
		}
	}
}

func TestAllDrivers_Count(t *testing.T) {
	drivers := AllDrivers()

	// Currently should have exactly 2 drivers (Vitest, Go)
	// This test will need updating when more drivers are added
	expectedCount := 4
	if len(drivers) != expectedCount {
		t.Errorf("AllDrivers() returned %d drivers, expected %d", len(drivers), expectedCount)
	}
}

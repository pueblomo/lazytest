package drivers

func AllDrivers() []Driver {
	return []Driver{
		&VitestDriver{},
	}
}

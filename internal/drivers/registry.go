package drivers

func AllDrivers() []Driver {
	return []Driver{
		&VitestDriver{},
		&GoTestDriver{},
		&MavenDriver{},
	}
}

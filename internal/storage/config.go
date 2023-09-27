package storage

type Driver string

const (
	DriverLocal Driver = "local"
	DriverS3    Driver = "s3"
	DriverGCS   Driver = "gcs"
)

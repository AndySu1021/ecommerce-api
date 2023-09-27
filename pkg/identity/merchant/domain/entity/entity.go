package entity

// Merchant Aggregate Root
type Merchant struct {
	ID          uint64
	Name        string
	Code        string
	Host        string
	EncryptSalt string
}

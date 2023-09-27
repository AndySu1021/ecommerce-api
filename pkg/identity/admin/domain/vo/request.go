package vo

type LoginParams struct {
	Email               string `json:"email" binding:"required"`
	Password            string `json:"password" binding:"required"`
	MerchantID          uint64
	MerchantHost        string
	MerchantEncryptSalt string
}

type LogoutParams struct {
	Email        string
	MerchantHost string
}

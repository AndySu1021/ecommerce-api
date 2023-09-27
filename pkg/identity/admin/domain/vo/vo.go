package vo

type CreateMemberParams struct {
	MerchantID uint64
	Email      string
	Password   string
}

type AdminLoginParams struct {
	MerchantID uint64
	Email      string
	Password   string
}

type ResetPasswordRepoParams struct {
	MerchantID uint64
	Email      string
	Password   string
}

type UpdateMemberPasswordRepoParams struct {
	OldPassword string
	Password    string
	MemberID    int64
}

type Merchant struct {
	ID          uint64
	Host        string
	EncryptSalt string
}

package vo

type CreateMemberParams struct {
	MerchantID uint64
	Email      string
	Password   string
}

type MemberLoginParams struct {
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
	MemberID    uint64
}

type Merchant struct {
	ID          uint64
	Host        string
	EncryptSalt string
}

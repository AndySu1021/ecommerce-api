package event

type MemberForgetPassword struct {
	LogoUrl      string
	MerchantName string
	RealName     string
	Email        string
	Code         string
}

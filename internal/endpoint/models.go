package endpoint

type GetTokenRequest struct {
	Guid string `validate:"required,uuid" json:"guid" example:"a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"`
}
type GetTokenResponse struct {
	Token   string `json:"token"`
	Refresh string `json:"refresh"`
}
type RefreshTokenRequest struct {
	Token   string `json:"token" validate:"required"`
	Refresh string `json:"refresh" validate:"required"`
}
type RefreshTokenResponse struct {
	Token   string `json:"token"`
	Refresh string `json:"refresh"`
}

package endpoint

type GetTokenRequest struct {
	Guid string `validate:"required" json:"guid" example:"a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"`
}
type GetTokenResponse struct {
	Token   string `json:"token"`
	Refresh string `json:"refresh"`
}

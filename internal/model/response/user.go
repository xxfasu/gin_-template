package response

type Login struct {
	AccessToken string `json:"accessToken"`
}

type GetProfile struct {
	UserId   string `json:"userId"`
	Nickname string `json:"nickname" example:"alan"`
}

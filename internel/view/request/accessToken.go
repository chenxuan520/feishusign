package request

type (
	AccessToken struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int64  `json:"expires_in"`
	}

	UserInfo struct {
		OpenId  string `json:"open_id"`
		Name    string `json:"name"`
		UserId  string `json:"user_id"`
		UnionId string `json:"union_id"`
	}
)

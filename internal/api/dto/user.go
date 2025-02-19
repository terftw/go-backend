package dto

type UserUpdate struct {
	Name    *string `json:"name,omitempty"`
	Picture *string `json:"picture,omitempty"`
}

type GoogleUserInfo struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

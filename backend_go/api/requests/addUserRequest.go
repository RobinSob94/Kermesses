package requests

type AddUserRequest struct {
	Type    string `json:"type"`
	UserIds []uint `json:"user_ids"`
}

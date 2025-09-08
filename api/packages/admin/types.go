package mapiadmin

type NewUserReq struct {
	AccessID  string `json:"username"`
	AccessKey string `json:"password"`
}

package mapiadmin

import (
	"github.com/lm-gautam-bhagat/minio-server/constants"
)

type NewUserReq struct {
	AccessID  string               `json:"username"`
	AccessKey string               `json:"password"`
	Policy    constants.PolicyType `json:"policy"`
}

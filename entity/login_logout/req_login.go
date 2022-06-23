package login_logout

type LoginReq struct {
	EquipmentID string `form:"equipmentId"`
	NickName    string `form:"nickName"`
	Age         int    `form:"age"`
	Country     string `form:"country"`
}

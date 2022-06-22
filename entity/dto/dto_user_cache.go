package dto

type UserCache struct {
	IsOnline  int    `db:"is_online"`
	NickName  string `db:"nickname"`
	AvatarUrl string `db:"avatar_url"`
}

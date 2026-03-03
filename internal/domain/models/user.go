package models

type User struct {
	Id 	 			string `json:"id"`
	Email 			string `json:"email"`
	Name 			string `json:"name"`
	PasswordHash 	string `json:"-"`
	AuthProvider 	string `json:"auth_provider"`
	AvatarUrl 		string `json:"avatar_url"`
	CreatedAt 		string `json:"created_at"`
}
package model

type (
	CreateUserRequest struct {
		UserID   string `json:"uid"`
		Name     string `json:"name"`
		Email    string
		Root_dir string `json:"root_dir"`
	}
	CreateUserResponse struct {
		Ok bool `json:"ok"`
	}
	GetUserRequest struct {
		Id string `json:"id"`
	}
	GetUserResponse struct {
		Ok bool `json:"ok"`
	}
)

type User struct {
	UId      string
	Name     string
	Root_dir string
	//Email string
}

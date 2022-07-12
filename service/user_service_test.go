package service

import (
	"hexagonal/caching"
	"hexagonal/repository"
	"testing"
)

func Test_userService_Authentication(t *testing.T) {
	type fields struct {
		cache    caching.AppCache
		userRepo repository.UserRepository
	}
	type args struct {
		payload *AuthenReq
	}
	type Req struct {
		Email    string `json:"email" bson:"email" db:"email"`
		Password string `json:"password" bson:"password"`
	}

	type Res struct {
		Code int
		Msg  string
	}
	tests := []struct {
		Name string
		Req  Req
		Res  Res
	}{
		// TODO: Add test cases.
		{
			Name: "case 1",
		},
		{
			Name: "case 2",
		},
		{
			Name: "case 3",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {

			// fmt.Println(tt.Req.Email)
			// s := userService{
			// 	cache:    tt.cache,
			// 	userRepo: tt.userRepo,
			// }
			// got, err := s.Authentication((*AuthenReq)(&tt.Req))
			// if (err != nil) != tt.wantErr {
			// 	t.Errorf("userService.Authentication() error = %v, wantErr %v", err, tt.wantErr)
			// 	return
			// }
			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("userService.Authentication() = %v, want %v", got, tt.want)
			// }
		})
	}
}

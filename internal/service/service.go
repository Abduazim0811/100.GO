package service

import (
	"100.GO/internal/entity/origin"
	"100.GO/internal/entity/user"
	"100.GO/internal/infrastructura/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserService struct {
	Repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{Repo: repo}
}

func (u *UserService) Createuser(user *user.CreateUser) error {
	return u.Repo.AddUser(user)
}

func (u *UserService) GetuserByEmail(email string) (*user.Login, error) {
	return u.Repo.GetUserByEmail(email)
}

func (u *UserService) Createorigin(req origin.CreateOrigin) error {
	return u.Repo.AddOrigin(req)
}

func (u *UserService) GetoriginById(reqId string) (*origin.GetOrigin, error) {
	return u.Repo.GetByIdOrigin(reqId)
}

func (u *UserService) GetAllorigins() ([]*origin.GetOrigin, error) {
	return u.Repo.GetAllOrigins()
}

func (u *UserService) Updateorigin(id primitive.ObjectID, req origin.CreateOrigin) error {
	return u.Repo.UpdateOrigin(id, req)
}

func (u *UserService) Deleteorigin(id primitive.ObjectID) error {
	return u.Repo.DeleteOrigin(id)
}

func (u *UserService) OriginGetall()([]*origin.CreateOrigin, error){
	return u.Repo.OriginGetAll()
}

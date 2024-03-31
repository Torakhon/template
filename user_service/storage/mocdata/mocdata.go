package mocdata

import (
	"context"
	"github.com/jmoiron/sqlx"
	pb "user_service/genproto/users"
)

type UserRepo struct {
	db *sqlx.DB
}

// MocNewUserRepo ...
func MocNewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (u *UserRepo) CreateUser(ctx context.Context, req *pb.CreateUserReq) (*pb.User, error) {
	return &pb.User{
		Id:        "49a2bac0-a42f-4807-a2ba-c0a42f080785",
		UserName:  "Ali",
		FirstName: "To'raxon",
		LastName:  "Jo'raxonov",
		Email:     "ali@mail.vom",
		Password:  "20030505",
		Role:      "Admin",
		Bio:       "Test bio",
		WebSite:   "Test web site",
		CreatedAt: "2005-05-05",
		UpdatedAt: "2006-06-06",
		DeletedAt: "",
		Posts:     nil,
	}, nil
}

func (u *UserRepo) GetUser(ctx context.Context, req *pb.GetUserReq) (*pb.User, error) {
	return &pb.User{
		Id:        "49a2bac0-a42f-4807-a2ba-c0a42f080785",
		UserName:  "Ali",
		FirstName: "To'raxon",
		LastName:  "Jo'raxonov",
		Email:     "ali@mail.vom",
		Password:  "20030505",
		Role:      "Admin",
		Bio:       "Test bio",
		WebSite:   "Test web site",
		CreatedAt: "2005-05-05",
		UpdatedAt: "2006-06-06",
		DeletedAt: "",
		Posts:     nil,
	}, nil
}

func (u *UserRepo) UpdateUser(ctx context.Context, req *pb.UpdateUserReq) (*pb.User, error) {
	return &pb.User{
		Id:        "49a2bac0-a42f-4807-a2ba-c0a42f080785",
		UserName:  "Ali_0505",
		FirstName: "To'raxon",
		LastName:  "Jo'raxonov",
		Email:     "ali20030505@mail.vom",
		Password:  "20030505",
		Role:      "Admin",
		Bio:       "Test bio",
		WebSite:   "Test web site",
		CreatedAt: "2005-05-05",
		UpdatedAt: "2006-06-06",
		DeletedAt: "",
		Posts:     nil,
	}, nil
}

func (u *UserRepo) DeleteUser(ctx context.Context, req *pb.DeleteUserReq) (*pb.DeleteUserRes, error) {
	return &pb.DeleteUserRes{
		Status: true,
	}, nil
}

func (u *UserRepo) GetAllUsers(ctx context.Context, req *pb.GetAllUsersReq) (*pb.GetAllUsersRes, error) {
	return &pb.GetAllUsersRes{
		Users: []*pb.User{
			{
				Id:        "5e51b8a2-543f-4572-91b8-a2543f157276",
				UserName:  "Ali_007",
				FirstName: "Ali",
				LastName:  "Jo'raxonov",
				Email:     "ali20030505@gmail.com",
				Password:  "Tj3884536",
				Role:      "Admin",
				Bio:       "Mod data bio",
				WebSite:   "Test website",
				CreatedAt: "2000-01-01",
				UpdatedAt: "2001-02-02",
				DeletedAt: "2003-03-03",
				Posts:     nil,
			},
			{
				Id:        "6e61c9b3-654g-4683-92c9-b3554g267277",
				UserName:  "John_123",
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john.doe@mail.com",
				Password:  "002002002",
				Role:      "User",
				Bio:       "New bio",
				WebSite:   "New website",
				CreatedAt: "2005-05-05",
				UpdatedAt: "2006-06-06",
				DeletedAt: "",
				Posts:     nil,
			},
		},
	}, nil
}

func (u *UserRepo) CheckUniques(ctx context.Context, req *pb.CheckUniqReq) (*pb.CheckUniqRes, error) {
	return &pb.CheckUniqRes{Code: 777}, nil
}

func (u *UserRepo) Login(ctx context.Context, req *pb.LoginReq) (*pb.LoginRes, error) {
	return &pb.LoginRes{
		Id:       "5e51b8a2-543f-4572-91b8-a2543f157276",
		Email:    "ali20030505@gmail.com",
		Password: "123123123",
		Role:     "SuperAdmin",
	}, nil
}

func (u *UserRepo) UpdateRole(ctx context.Context, req *pb.UpdateRoleReq) (*pb.UpdateRoleRes, error) {
	return &pb.UpdateRoleRes{
		Stats: true,
	}, nil
}

func (u *UserRepo) UpdateEmail(ctx context.Context, req *pb.UpdateEmailReq) (*pb.UpdateEmailRes, error) {
	return &pb.UpdateEmailRes{Status: true}, nil
}

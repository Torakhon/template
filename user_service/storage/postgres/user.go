package postgres

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"math/rand/v2"
	pb "user_service/genproto/users"
)

type UserRepo struct {
	db *sqlx.DB
}

// NewUserRepo ...
func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{db: db}
}
func (u *UserRepo) CreateUser(ctx context.Context, req *pb.CreateUserReq) (*pb.User, error) {
	var user pb.User
	query := `INSERT INTO users(
								id,
								user_name,
								first_name,
								last_name,
								email,
								password,
								role,
								bio,
								website
										) VALUES 
								($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING
									id,
									user_name,
									first_name,
									last_name,
									email,
									role,
									bio,
									website,
									created_at ,
									updated_at `

	err := u.db.QueryRow(query,
		req.Id,
		req.UserName,
		req.FirstName,
		req.LastName,
		req.Email,
		req.Password,
		req.Role,
		req.Bio,
		req.WebSite).Scan(
		&user.Id,
		&user.UserName,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Role,
		&user.Bio,
		&user.WebSite,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *UserRepo) GetUser(ctx context.Context, req *pb.GetUserReq) (*pb.User, error) {
	var user pb.User
	query := fmt.Sprintf(`SELECT id,
									user_name,
									first_name,
									last_name,
									email,
									role,
									bio,
									website,
									created_at ,
									updated_at FROM users
										WHERE %s = $1`, req.Field)

	err := u.db.QueryRow(query, req.Value).Scan(
		&user.Id,
		&user.UserName,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Role,
		&user.Bio,
		&user.WebSite,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *UserRepo) UpdateUser(ctx context.Context, req *pb.UpdateUserReq) (*pb.User, error) {
	var user pb.User
	query := `UPDATE users SET
							user_name = $1,
							first_name = $2,
							last_name = $3,
							password = $4,
							bio = $5,
							website = $6
							WHERE id = $7 and deleted_at IS NULL 
							RETURNING
							id,
							user_name ,
							first_name ,
							last_name ,
							role,
							bio ,
							website`

	err := u.db.QueryRow(query,
		req.UserName,
		req.FirstName,
		req.LastName,
		req.Password,
		req.Bio,
		req.WebSite,
		req.Id).Scan(
		&user.Id,
		&user.UserName,
		&user.FirstName,
		&user.LastName,
		&user.Role,
		&user.Bio,
		&user.WebSite,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *UserRepo) DeleteUser(ctx context.Context, req *pb.DeleteUserReq) (*pb.DeleteUserRes, error) {
	return &pb.DeleteUserRes{
		Status: true,
	}, nil
}

func (u *UserRepo) GetAllUsers(ctx context.Context, req *pb.GetAllUsersReq) (*pb.GetAllUsersRes, error) {
	var users pb.GetAllUsersRes
	offset := req.Limit * (req.Page - 1)
	query := `SELECT id,
									user_name,
									first_name,
									last_name,
									email,
									role,
									bio,
									website,
									created_at ,
									updated_at FROM users WHERE deleted_at IS NULL LIMIT $1 OFFSET $2`
	rows, err := u.db.Query(query, req.Limit, offset)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var user pb.User
		err := rows.Scan(
			&user.Id,
			&user.UserName,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.Role,
			&user.Bio,
			&user.WebSite,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users.Users = append(users.Users, &user)
	}
	return &users, nil
}

func (u *UserRepo) CheckUniques(ctx context.Context, req *pb.CheckUniqReq) (*pb.CheckUniqRes, error) {
	var count int
	err := u.db.QueryRow(fmt.Sprintf(`SELECT count(1) FROM users WHERE %s = $1`, req.Field), req.Value).Scan(&count)
	if err != nil {
		return &pb.CheckUniqRes{
			Code: 0,
		}, err
	}
	if count != 0 {
		return &pb.CheckUniqRes{
			Code: 0,
		}, err
	}
	num := rand.Int32() % 1000000
	return &pb.CheckUniqRes{
		Code: num,
	}, nil

}

func (u *UserRepo) Login(ctx context.Context, req *pb.LoginReq) (*pb.LoginRes, error) {
	var response pb.LoginRes
	query := `SELECT id, password,role,email FROM users WHERE email = $1`
	err := u.db.QueryRow(query, req.Email).Scan(&response.Id, &response.Password, &response.Role, &response.Email)
	if err != nil {
		return nil, err
	}
	return &pb.LoginRes{
		Password: response.Password,
		Role:     response.Role,
		Id:       response.Id,
		Email:    response.Email,
	}, nil
}

func (u *UserRepo) UpdateRole(ctx context.Context, req *pb.UpdateRoleReq) (*pb.UpdateRoleRes, error) {
	query := `UPDATE users SET
				role = $1
				WHERE id = $2 and deleted_at IS NULL`
	_, err := u.db.Exec(query, req.NewRole, req.Id)
	if err != nil {
		return &pb.UpdateRoleRes{
			Stats: false,
		}, err
	}
	return &pb.UpdateRoleRes{
		Stats: true,
	}, nil
}

func (u *UserRepo) UpdateEmail(ctx context.Context, req *pb.UpdateEmailReq) (*pb.UpdateEmailRes, error) {
	query := `UPDATE users SET
				email = $1
				WHERE id = $2 and deleted_at IS NULL`
	_, err := u.db.Exec(query, req.Email, req.Id)
	if err != nil {
		return &pb.UpdateEmailRes{
			Status: false,
		}, err
	}
	return &pb.UpdateEmailRes{
		Status: true,
	}, nil
}

package mongoDBdatabase

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math/rand"
	"time"
	pb "user_service/genproto/users"
)

type UserRepo struct {
	collection *mongo.Collection
}

// NewUserRepo ...
func NewUserRepo(client *mongo.Client, dbName, collectionName string) *UserRepo {
	return &UserRepo{
		collection: client.Database(dbName).Collection(collectionName),
	}
}

func (u *UserRepo) CreateUser(ctx context.Context, req *pb.CreateUserReq) (*pb.User, error) {
	user := &pb.User{
		Id:        req.Id,
		UserName:  req.UserName,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  req.Password,
		Role:      req.Role,
		Bio:       req.Bio,
		WebSite:   req.WebSite,
	}

	_, err := u.collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserRepo) GetUser(ctx context.Context, req *pb.GetUserReq) (*pb.User, error) {
	var user pb.User

	filter := bson.M{req.Field: req.Value}
	err := u.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *UserRepo) UpdateUser(ctx context.Context, req *pb.UpdateUserReq) (*pb.User, error) {
	filter := bson.M{"id": req.Id}
	update := bson.M{
		"$set": bson.M{
			"userName":  req.UserName,
			"firstName": req.FirstName,
			"lastName":  req.LastName,
			"password":  req.Password,
			"bio":       req.Bio,
			"website":   req.WebSite,
		},
	}

	_, err := u.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return u.GetUser(ctx, &pb.GetUserReq{Field: "id", Value: req.Id})
}

func (u *UserRepo) DeleteUser(ctx context.Context, req *pb.DeleteUserReq) (*pb.DeleteUserRes, error) {
	filter := bson.M{req.Field: req.Value}
	update := bson.M{"$set": bson.M{"deleted_at": primitive.NewDateTimeFromTime(time.Now())}}

	_, err := u.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return &pb.DeleteUserRes{Status: true}, nil
}

func (u *UserRepo) GetAllUsers(ctx context.Context, req *pb.GetAllUsersReq) (*pb.GetAllUsersRes, error) {
	var users []*pb.User

	findOptions := options.Find().SetLimit(int64(req.Limit)).SetSkip(int64(req.Limit * (req.Page - 1)))
	cursor, err := u.collection.Find(ctx, bson.M{"deleted_at": nil}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var user pb.User
		err := cursor.Decode(&user)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return &pb.GetAllUsersRes{Users: users}, nil
}

func (u *UserRepo) CheckUniques(ctx context.Context, req *pb.CheckUniqReq) (*pb.CheckUniqRes, error) {
	filter := bson.M{req.Field: req.Value}
	count, err := u.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	if count != 0 {
		return &pb.CheckUniqRes{Code: 0}, nil
	}

	num := rand.Int31() % 1000000
	return &pb.CheckUniqRes{Code: num}, nil
}

func (u *UserRepo) Login(ctx context.Context, req *pb.LoginReq) (*pb.LoginRes, error) {
	var user pb.User
	filter := bson.M{"email": req.Email}
	err := u.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &pb.LoginRes{
		Password: user.Password,
		Role:     user.Role,
		Id:       user.Id,
		Email:    user.Email,
	}, nil
}

func (u *UserRepo) UpdateRole(ctx context.Context, req *pb.UpdateRoleReq) (*pb.UpdateRoleRes, error) {
	filter := bson.M{"id": req.Id}
	update := bson.M{"$set": bson.M{"role": req.NewRole}}

	_, err := u.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateRoleRes{Stats: true}, nil
}

func (u *UserRepo) UpdateEmail(ctx context.Context, req *pb.UpdateEmailReq) (*pb.UpdateEmailRes, error) {
	filter := bson.M{"id": req.Id}
	update := bson.M{"$set": bson.M{"email": req.Email}}

	_, err := u.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateEmailRes{Status: true}, nil
}

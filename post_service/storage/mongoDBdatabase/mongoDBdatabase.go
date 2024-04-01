package mongoDBdatabase

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	pb "post_service/genproto/post"
	"time"
)

type PostRepo struct {
	collection *mongo.Collection
}

// NewPostRepo ...
func NewPostRepo(client *mongo.Client, dbName, collectionName string) *PostRepo {
	return &PostRepo{
		collection: client.Database(dbName).Collection(collectionName),
	}
}

func (p *PostRepo) Create(ctx context.Context, req *pb.CreateReq) (*pb.Post, error) {
	post := &pb.Post{
		Id:        req.Id,
		Title:     req.Title,
		Content:   req.Content,
		UserId:    req.UserId,
		Category:  req.Category,
		Likes:     0,
		Dislikes:  0,
		Views:     0,
		CreatedAt: time.Now().String(),
		UpdatedAt: time.Now().String(),
	}

	_, err := p.collection.InsertOne(ctx, post)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (p *PostRepo) GetPost(ctx context.Context, req *pb.GetReq) (*pb.Post, error) {
	var post pb.Post
	filter := bson.M{"id": req.PostId, "deleted_at": nil}
	err := p.collection.FindOne(ctx, filter).Decode(&post)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (p *PostRepo) SearchPost(ctx context.Context, req *pb.SearchReq) (*pb.PostsRes, error) {
	offset := req.Limit * (req.Page - 1)
	var posts []*pb.Post
	filter := bson.M{req.Field: req.Value, "deleted_at": nil}
	opts := options.Find().SetLimit(int64(req.Limit)).SetSkip(int64(offset))
	cursor, err := p.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var post pb.Post
		if err := cursor.Decode(&post); err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}

	return &pb.PostsRes{Posts: posts}, nil
}

func (p *PostRepo) UpdatePost(ctx context.Context, req *pb.UpdatePostReq) (*pb.Post, error) {
	filter := bson.M{"id": req.Id, "deleted_at": nil}
	update := bson.M{
		"$set": bson.M{
			"title":     req.Title,
			"content":   req.Content,
			"category":  req.Category,
			"updatedAt": time.Now().Unix(),
		},
	}

	var updatedPost pb.Post
	err := p.collection.FindOneAndUpdate(ctx, filter, update).Decode(&updatedPost)
	if err != nil {
		return nil, err
	}

	return &updatedPost, nil
}

func (p *PostRepo) DeletePost(ctx context.Context, req *pb.DeletePostReq) (*pb.DeletePostRes, error) {
	filter := bson.M{"id": req.Id}
	update := bson.M{"$set": bson.M{"deleted_at": time.Now()}}

	_, err := p.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return &pb.DeletePostRes{Status: true}, nil
}

func (p *PostRepo) PostClickLike(ctx context.Context, req *pb.ClickReq) (*pb.PostLike, error) {
	var status bool
	filter := bson.M{"post_id": req.PostId, "user_id": req.UserId}
	err := p.collection.FindOne(ctx, filter).Decode(&status)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			_, err := p.collection.InsertOne(ctx, bson.M{"post_id": req.PostId, "user_id": req.UserId, "status": true})

			if err != nil {
				return nil, err
			}

			_, err = p.collection.UpdateOne(ctx, bson.M{"id": req.PostId}, bson.M{"$inc": bson.M{"likes": 1}})

			if err != nil {
				return nil, err
			}
			return &pb.PostLike{Like: true}, nil
		} else {
			return nil, err
		}
	}
	if status == true {
		status = false
		_, err = p.collection.UpdateOne(ctx, bson.M{"id": req.PostId}, bson.M{"$inc": bson.M{"likes": -1}})
	} else {
		status = true
		_, err = p.collection.UpdateOne(ctx, bson.M{"id": req.PostId}, bson.M{"$inc": bson.M{"likes": 1}})
	}
	_, err = p.collection.UpdateOne(ctx, filter, bson.M{"$set": bson.M{"status": status}})

	if err != nil {
		return nil, err
	}

	return &pb.PostLike{Like: status}, nil
}

func (p *PostRepo) PostClickDisLike(ctx context.Context, req *pb.ClickReq) (*pb.PostLike, error) {
	var status bool
	filter := bson.M{"post_id": req.PostId, "user_id": req.UserId}
	err := p.collection.FindOne(ctx, filter).Decode(&status)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			_, err := p.collection.InsertOne(ctx, bson.M{"post_id": req.PostId, "user_id": req.UserId, "status": true})

			if err != nil {
				return nil, err
			}

			_, err = p.collection.UpdateOne(ctx, bson.M{"id": req.PostId}, bson.M{"$inc": bson.M{"dislikes": 1}})

			if err != nil {
				return nil, err
			}
			return &pb.PostLike{Like: true}, nil
		} else {
			return nil, err
		}
	}
	if status == true {
		status = false
		_, err = p.collection.UpdateOne(ctx, bson.M{"id": req.PostId}, bson.M{"$inc": bson.M{"dislikes": -1}})
	} else {
		status = true
		_, err = p.collection.UpdateOne(ctx, bson.M{"id": req.PostId}, bson.M{"$inc": bson.M{"dislikes": 1}})
	}

	_, err = p.collection.UpdateOne(ctx, filter, bson.M{"$set": bson.M{"status": status}})

	if err != nil {
		return nil, err
	}

	return &pb.PostLike{Like: status}, nil
}

func (p *PostRepo) Views(ctx context.Context, req *pb.ViewReq) (*pb.ViewRes, error) {
	var viewExists bool
	filter := bson.M{"user_id": req.UserId, "post_id": req.PostId}
	err := p.collection.FindOne(ctx, filter).Decode(&viewExists)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			_, err := p.collection.InsertOne(ctx, bson.M{"user_id": req.UserId, "post_id": req.PostId})
			if err != nil {
				return nil, err
			}

			_, err = p.collection.UpdateOne(ctx, bson.M{"id": req.PostId}, bson.M{"$inc": bson.M{"views": 1}})
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return &pb.ViewRes{}, nil
}

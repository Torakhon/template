package mongoDBdatabase

import (
	pb "comment_service/genproto/comment"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type CommentRepo struct {
	collection *mongo.Collection
}

// NewCommentRepo ...
func NewCommentRepo(client *mongo.Client, dbName, collectionName string) *CommentRepo {
	return &CommentRepo{
		collection: client.Database(dbName).Collection(collectionName),
	}
}

func (r *CommentRepo) Create(ctx context.Context, req *pb.CreateReq) (*pb.Comment, error) {
	comment := &pb.Comment{
		CommentId: req.CommentId,
		PostId:    req.PostId,
		UserId:    req.UserId,
		Content:   req.Content,
	}

	_, err := r.collection.InsertOne(ctx, comment)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (r *CommentRepo) GetCommentsByPostId(ctx context.Context, req *pb.GetByPostIdReq) (*pb.GetByIdCommentsRes, error) {
	var comments []*pb.Comment
	filter := bson.M{"post_id": req.PostId}
	opts := options.Find().SetLimit(int64(req.Limit)).SetSkip(int64(req.Limit * (req.Page - 1)))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var comment pb.Comment
		if err := cursor.Decode(&comment); err != nil {
			return nil, err
		}
		comments = append(comments, &comment)
	}

	return &pb.GetByIdCommentsRes{Comments: comments}, nil
}

func (r *CommentRepo) GetCommentsByOwnerId(ctx context.Context, req *pb.GetByOwnerIdReq) (*pb.GetByIdCommentsRes, error) {
	var comments []*pb.Comment
	filter := bson.M{"user_id": req.OwnerId}
	opts := options.Find().SetLimit(int64(req.Limit)).SetSkip(int64(req.Limit * (req.Page - 1)))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var comment pb.Comment
		if err := cursor.Decode(&comment); err != nil {
			return nil, err
		}
		comments = append(comments, &comment)
	}

	return &pb.GetByIdCommentsRes{Comments: comments}, nil
}

func (r *CommentRepo) UpdateComment(ctx context.Context, req *pb.UpdateCommentReq) (*pb.Comment, error) {
	filter := bson.M{"comment_id": req.CommentId, "user_id": req.UserId}
	update := bson.M{"$set": bson.M{"content": req.NewContent}}

	var updatedComment pb.Comment
	err := r.collection.FindOneAndUpdate(ctx, filter, update).Decode(&updatedComment)
	if err != nil {
		return nil, err
	}

	return &updatedComment, nil
}

func (r *CommentRepo) DeleteComment(ctx context.Context, req *pb.DeleteCommentReq) (*pb.DeleteRes, error) {
	filter := bson.M{"comment_id": req.CommentId, "user_id": req.UserId}
	update := bson.M{"$set": bson.M{"deleted_at": time.Now()}}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return &pb.DeleteRes{Status: false}, err
	}

	return &pb.DeleteRes{Status: true}, nil
}

func (r *CommentRepo) CommentClickLike(ctx context.Context, req *pb.ClickReq) (*pb.CommentLike, error) {
	var status bool
	filter := bson.M{"comment_id": req.CommentId, "user_id": req.UserId}
	err := r.collection.FindOne(ctx, filter).Decode(&status)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			_, err := r.collection.InsertOne(ctx, bson.M{"comment_id": req.CommentId, "user_id": req.UserId, "status": true})

			if err != nil {
				return nil, err
			}

			_, err = r.collection.UpdateOne(ctx, bson.M{"comment_id": req.CommentId}, bson.M{"$inc": bson.M{"likes": 1}})

			if err != nil {
				return nil, err
			}
			return &pb.CommentLike{Like: true}, nil
		} else {
			return nil, err
		}
	}
	if status == true {
		status = false
		_, err = r.collection.UpdateOne(ctx, bson.M{"comment_id": req.CommentId}, bson.M{"$inc": bson.M{"likes": -1}})
	} else {
		status = true
		_, err = r.collection.UpdateOne(ctx, bson.M{"comment_id": req.CommentId}, bson.M{"$inc": bson.M{"likes": 1}})
	}
	_, err = r.collection.UpdateOne(ctx, filter, bson.M{"$set": bson.M{"status": status}})

	if err != nil {
		return nil, err
	}

	return &pb.CommentLike{Like: status}, nil
}

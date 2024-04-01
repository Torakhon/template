package postgres

import (
	pb "comment_service/genproto/comment"
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
)

type CommentRepo struct {
	db *sqlx.DB
}

// NewCommentRepo ...
func NewCommentRepo(db *sqlx.DB) *CommentRepo {
	return &CommentRepo{db: db}
}

func (r *CommentRepo) Create(ctx context.Context, req *pb.CreateReq) (*pb.Comment, error) {
	var comment pb.Comment
	query := `INSERT INTO comments (comment_id,post_id,user_id,content) VALUES 
				($1,$2,$3,$4) RETURNING comment_id,post_id,user_id,content`

	err := r.db.QueryRow(query, req.CommentId, req.PostId, req.UserId, req.Content).Scan(
		&comment.CommentId,
		&comment.PostId,
		&comment.UserId,
		&comment.Content,
	)
	if err != nil {
		return nil, err
	}

	return &comment, nil
}

func (r *CommentRepo) GetCommentsByPostId(ctx context.Context, req *pb.GetByPostIdReq) (*pb.GetByIdCommentsRes, error) {
	var comments pb.GetByIdCommentsRes
	offset := req.Limit * (req.Page - 1)
	query := `SELECT 
					comment_id,
					post_id,
					user_id,
					content,
					likes,
					created_at,
					updated_at FROM comments
					WHERE post_id = $1 LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(query, req.PostId, req.Limit, offset)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var comment pb.Comment
		err := rows.Scan(
			&comment.CommentId,
			&comment.PostId,
			&comment.UserId,
			&comment.Content,
			&comment.Likes,
			&comment.CreatedAt,
			&comment.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		comments.Comments = append(comments.Comments, &comment)
	}

	return &comments, nil
}

func (r *CommentRepo) GetCommentsByOwnerId(ctx context.Context, req *pb.GetByOwnerIdReq) (*pb.GetByIdCommentsRes, error) {
	var comments pb.GetByIdCommentsRes
	offset := req.Limit * (req.Page - 1)
	query := `SELECT 
					comment_id,
					post_id,
					user_id,
					content,
					likes,
					created_at,
					updated_at FROM comments
					WHERE user_id = $1 LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(query, req.OwnerId, req.Limit, offset)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var comment pb.Comment
		err := rows.Scan(
			&comment.CommentId,
			&comment.PostId,
			&comment.UserId,
			&comment.Content,
			&comment.Likes,
			&comment.CreatedAt,
			&comment.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		comments.Comments = append(comments.Comments, &comment)
	}

	return &comments, nil
}

func (r *CommentRepo) UpdateComment(ctx context.Context, req *pb.UpdateCommentReq) (*pb.Comment, error) {
	var comment pb.Comment
	query := `UPDATE comments SET
				content = $1
				WHERE comment_id = $2 and user_id = $3 RETURNING
				comment_id,post_id,user_id,content,likes`

	err := r.db.QueryRow(query, req.NewContent, req.CommentId, req.UserId).Scan(
		&comment.CommentId,
		&comment.PostId,
		&comment.UserId,
		&comment.Content,
		&comment.Likes,
	)

	if err != nil {
		return nil, err
	}

	return &comment, nil
}

func (r *CommentRepo) DeleteComment(ctx context.Context, req *pb.DeleteCommentReq) (*pb.DeleteRes, error) {
	query := `UPDATE comments SET
				deleted_at = now()
				WHERE comment_id = $1 and user_id = $2`

	_, err := r.db.Exec(query, req.CommentId, req.UserId)
	if err != nil {
		return &pb.DeleteRes{Status: false}, err
	}

	return &pb.DeleteRes{Status: true}, nil
}

func (r *CommentRepo) CommentClickLike(ctx context.Context, req *pb.ClickReq) (*pb.CommentLike, error) {
	var status bool
	err := r.db.QueryRow(`SELECT status FROM comment_like WHERE comment_id = $1 and user_id = $2`, req.CommentId, req.UserId).Scan(
		&status,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_, err := r.db.Exec(`INSERT INTO comment_like (comment_id,user_id,status) VALUES 
										($1,$2,$3)`, req.CommentId, req.UserId, true)

			if err != nil {
				return nil, err
			}

			_, err = r.db.Exec(`UPDATE comments SET likes = likes +1 WHERE comment_id = $1`, req.CommentId)

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
		_, err = r.db.Exec(`UPDATE comments SET likes = likes -1 WHERE comment_id = $1`, req.CommentId)
	} else {
		status = true
		_, err = r.db.Exec(`UPDATE comments SET likes = likes +1 WHERE comment_id = $1`, req.CommentId)
	}
	_, err = r.db.Exec(`UPDATE comment_like SET
										status = $1`, status)

	if err != nil {
		return nil, err
	}

	return &pb.CommentLike{Like: status}, nil
}

package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	pb "post_service/genproto/post"
)

type PostRepo struct {
	db *sqlx.DB
}

// NewPostRepo ...
func NewPostRepo(db *sqlx.DB) *PostRepo {
	return &PostRepo{db: db}
}

func (p *PostRepo) Create(ctx context.Context, req *pb.CreateReq) (*pb.Post, error) {
	var post pb.Post
	query := `INSERT INTO posts(id,title,content,user_id,category) VALUES ($1,$2,$3,$4,$5) RETURNING 
				id,title,content,user_id,category`
	err := p.db.QueryRow(query, req.Id, req.Title, req.Content, req.UserId, req.Category).Scan(
		&post.Id,
		&post.Title,
		&post.Content,
		&post.UserId,
		&post.Category,
	)

	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (p *PostRepo) GetPost(ctx context.Context, req *pb.GetReq) (*pb.Post, error) {
	var post pb.Post
	query := `SELECT id,title,content,user_id,category,likes,dislikes,views,created_at,updated_at FROM posts
				WHERE id = $1 and deleted_at IS NULL `
	err := p.db.QueryRow(query, req.PostId).Scan(
		&post.Id,
		&post.Title,
		&post.Content,
		&post.UserId,
		&post.Category,
		&post.Likes,
		&post.Dislikes,
		&post.Views,
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (p *PostRepo) SearchPost(ctx context.Context, req *pb.SearchReq) (*pb.PostsRes, error) {
	offset := req.Limit * (req.Page - 1)
	var posts pb.PostsRes
	query := fmt.Sprintf(`SELECT id,title,content,user_id,category,likes,dislikes,views,created_at,updated_at FROM posts
				WHERE %s = $1  and deleted_at IS NULL LIMIT $2 OFFSET $3`, req.Field)

	rows, err := p.db.Query(query, req.Value, req.Limit, offset)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var post pb.Post
		err := rows.Scan(
			&post.Id,
			&post.Title,
			&post.Content,
			&post.UserId,
			&post.Category,
			&post.Likes,
			&post.Dislikes,
			&post.Views,
			&post.CreatedAt,
			&post.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		posts.Posts = append(posts.Posts, &post)
	}

	return &posts, nil
}

func (p *PostRepo) UpdatePost(ctx context.Context, req *pb.UpdatePostReq) (*pb.Post, error) {
	var post pb.Post
	query := `UPDATE posts SET
				title = $1,
				content = $2,
				category = $3 WHERE id = $4 and deleted_at IS NULL
				RETURNING 
				id,title,content,user_id,category,likes,dislikes`

	err := p.db.QueryRow(query, req.Title, req.Content, req.Category, req.Id).Scan(
		&post.Id,
		&post.Title,
		&post.Content,
		&post.UserId,
		&post.Category,
		&post.Likes,
		&post.Dislikes,
	)

	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (p *PostRepo) DeletePost(ctx context.Context, req *pb.DeletePostReq) (*pb.DeletePostRes, error) {
	query := `UPDATE posts SET deleted_at = now() WHERE id = $1`
	_, err := p.db.Exec(query, req.Id)
	if err != nil {
		return &pb.DeletePostRes{
			Status: false,
		}, nil
	}

	return &pb.DeletePostRes{
		Status: true,
	}, nil
}

func (p *PostRepo) PostClickLike(ctx context.Context, req *pb.ClickReq) (*pb.PostLike, error) {
	var status bool
	err := p.db.QueryRow(`SELECT status FROM post_like WHERE post_id = $1 and user_id = $2`, req.PostId, req.UserId).Scan(
		&status,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_, err := p.db.Exec(`INSERT INTO post_like (post_id,user_id,status) VALUES 
										($1,$2,$3)`, req.PostId, req.UserId, true)

			if err != nil {
				return nil, err
			}

			_, err = p.db.Exec(`UPDATE posts SET likes = likes +1 WHERE id = $1`, req.PostId)

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
		_, err = p.db.Exec(`UPDATE posts SET likes = likes -1 WHERE id = $1`, req.PostId)
	} else {
		status = true
		_, err = p.db.Exec(`UPDATE posts SET likes = likes +1 WHERE id = $1`, req.PostId)
	}
	_, err = p.db.Exec(`UPDATE post_like SET
										status = $1`, status)

	if err != nil {
		return nil, err
	}

	return &pb.PostLike{Like: status}, nil

}

func (p *PostRepo) PostClickDisLike(ctx context.Context, req *pb.ClickReq) (*pb.PostLike, error) {
	var status bool
	err := p.db.QueryRow(`SELECT status FROM post_like WHERE dislike = 'dislike' and post_id = $1 and user_id = $2`, req.PostId, req.UserId).Scan(
		&status,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_, err := p.db.Exec(`INSERT INTO post_like (dislike,post_id,user_id,status) VALUES 
										($1,$2,$3,$4)`, "dislike", req.PostId, req.UserId, true)

			if err != nil {
				return nil, err
			}

			_, err = p.db.Exec(`UPDATE posts SET dislikes = dislikes +1 WHERE id = $1`, req.PostId)

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
		_, err = p.db.Exec(`UPDATE posts SET dislikes = dislikes -1 WHERE id = $1`, req.PostId)
	} else {
		status = true
		_, err = p.db.Exec(`UPDATE posts SET dislikes = dislikes +1 WHERE id = $1`, req.PostId)
	}

	_, err = p.db.Exec(`UPDATE post_like SET
										status = $1 WHERE dislike = 'dislike' and post_id = $2 and user_id = $3 `, status, req.PostId, req.UserId)

	if err != nil {
		return nil, err
	}

	return &pb.PostLike{Like: status}, nil

}

func (p *PostRepo) Views(ctx context.Context, req *pb.ViewReq) (*pb.ViewRes, error) {
	var (
		userId string
		postId string
	)

	err := p.db.QueryRow(`SELECT user_id, post_id FROM views WHERE user_id = $1 AND post_id = $2`, req.UserId, req.PostId).Scan(
		&userId, &postId,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			query := `INSERT INTO views (user_id, post_id) VALUES ($1, $2)`
			_, err = p.db.Exec(query, req.UserId, req.PostId)
			if err != nil {
				return nil, err
			}
			query = `UPDATE posts SET views = views + 1 WHERE id = $1`
			_, err = p.db.Exec(query, req.PostId)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return &pb.ViewRes{}, nil
}

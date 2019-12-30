package edc

import (
	"context"
	"time"
)

// Post - struct for post
type Post struct {
	ID        int64  `sql:"id"         json:"id"   form:"id"   query:"id"`
	Name      string `sql:"name"       json:"name" form:"name" query:"name"`
	GO        bool   `sql:"go"         json:"go"   form:"go"   query:"go"`
	Note      string `sql:"note"       json:"note" form:"note" query:"note"`
	CreatedAt string `sql:"created_at" json:"-"`
	UpdatedAt string `sql:"updated_at" json:"-"`
}

// PostList - struct for post list
type PostList struct {
	ID   int64  `sql:"id"   json:"id"   form:"id"   query:"id"`
	Name string `sql:"name" json:"name" form:"name" query:"name"`
	GO   bool   `sql:"go"   json:"go"   form:"go"   query:"go"`
	Note string `sql:"note" json:"note" form:"note" query:"note"`
}

// PostGet - get one post by id
func PostGet(id int64) (Post, error) {
	var post Post
	if id == 0 {
		return post, nil
	}
	post.ID = id
	err := pool.QueryRow(context.Background(), `
		SELECT
			name,
			go,
			note,
			created_at,
			updated_at
		FROM
			kinds
		WHERE
			id = $1
	`, id).Scan(&post.Name, &post.GO, &post.Note, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		errmsg("PostGet QueryRow", err)
	}
	return post, nil
}

// PostListGet - get all post for list
func PostListGet() ([]PostList, error) {
	var posts []PostList
	rows, err := pool.Query(context.Background(), `
		SELECT
			id,
			name,
			go,
			note
		FROM
			kinds
		ORDER BY
			name ASC
	`)
	if err != nil {
		errmsg("PostListGet Query", err)
	}
	for rows.Next() {
		var post PostList
		err := rows.Scan(&post.ID, &post.Name, &post.GO, &post.Note)
		if err != nil {
			errmsg("PostListGet Scan", err)
			return posts, err
		}
		posts = append(posts, post)
	}
	return posts, rows.Err()
}

// PostSelectGet - get all post for select
func PostSelectGet(g bool) ([]SelectItem, error) {
	var posts []SelectItem
	rows, err := pool.Query(context.Background(), `
		SELECT
			id,
			name
		FROM
			posts
		WHERE
			go = $1
		ORDER BY
			name ASC
	`, g)
	if err != nil {
		errmsg("PostSelectGet Query", err)
	}
	for rows.Next() {
		var post SelectItem
		err := rows.Scan(&post.ID, &post.Name)
		if err != nil {
			errmsg("PostSelectGet Scan", err)
			return posts, err
		}
		posts = append(posts, post)
	}
	return posts, rows.Err()
}

// PostInsert - create new post
func PostInsert(post Post) (int64, error) {
	err := pool.QueryRow(context.Background(), `
		INSERT INTO posts
		(
			name,
			go,
			note,
			created_at,
			updated_at
		)
		VALUES
		(
			$1,
			$2,
			$3,
			$4,
			$5
		)
		RETURNING
			id
	`, post.Name, post.GO, post.Note, time.Now(), time.Now()).Scan(&post.ID)
	if err != nil {
		errmsg("PostInsert QueryRow", err)
	}
	return post.ID, nil
}

// PostUpdate - save post changes
func PostUpdate(post Post) error {
	_, err := pool.Exec(context.Background(), `
		UPDATE posts SET
			name = $2,
			go = $3,
			note = $4,
			updated_at = $5
		WHERE
			id = $1
	`, post.ID, post.Name, post.GO, post.Note, time.Now())
	if err != nil {
		errmsg("UpdatePost update", err)
	}
	return err
}

// PostDelete - delete post by id
func PostDelete(id int64) error {
	if id == 0 {
		return nil
	}
	_, err := pool.Exec(context.Background(), `
		DELETE FROM
			posts
		WHERE
			id = $1
	`, id)
	if err != nil {
		errmsg("DeletePost Exec", err)
	}
	return err
}

func postCreateTable() error {
	str := `
		CREATE TABLE IF NOT EXISTS
			posts (
				id BIGSERIAL PRIMARY KEY,
				name TEXT,
				go BOOL NOT NULL DEFAULT FALSE,
				note TEXT,
				created_at TIMESTAMP without time zone,
				updated_at TIMESTAMP without time zone default now(),
				UNIQUE (name, go)
			)
	`
	_, err := pool.Exec(context.Background(), str)
	if err != nil {
		errmsg("postCreateTable exec", err)
	}
	return err
}

package edc

import "context"

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
	err := pool.QueryRow(context.Background(), &post).
		Where("id = ?", id).
		Select()
	if err != nil {
		errmsg("GetPost select", err)
	}
	return post, nil
}

// PostListGet - get all post for list
func PostListGet() ([]PostList, error) {
	var posts []PostList
	err := pool.QueryRow(context.Background(), &Post{}).
		Column("id", "name", "go", "note").
		Order("name ASC").
		Select(&posts)
	if err != nil {
		errmsg("GetPostListAll select", err)
	}
	return posts, nil
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
	err := pool.Insert(&post)
	if err != nil {
		errmsg("CreatePost insert", err)
	}
	return post.ID, nil
}

// PostUpdate - save post changes
func PostUpdate(post Post) error {
	err := pool.Update(&post)
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
	_, err := pool.QueryRow(context.Background(), &Post{}).
		Where("id = ?", id).
		Delete()
	if err != nil {
		errmsg("DeletePost delete", err)
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
				updated_at
 TIMESTAMP without time zone default now(),
				UNIQUE (name, go)
			)
	`
	_, err := pool.Exec(str)
	if err != nil {
		errmsg("postCreateTable exec", err)
	}
	return err
}

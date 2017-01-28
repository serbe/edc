package edc

// Post - struct for post
type Post struct {
	ID        int64  `sql:"id" json:"id"`
	Name      string `sql:"name" json:"name"`
	GO        bool   `sql:"go" json:"go"`
	Note      string `sql:"note, null" json:"note"`
	CreatedAt string `sql:"created_at" json:"created_at"`
	UpdatedAt string `sql:"updated_at" json:"updated_at"`
}

// PostList - struct for post list
type PostList struct {
	ID   int64  `sql:"id" json:"id"`
	Name string `sql:"name" json:"name"`
	GO   bool   `sql:"go" json:"go"`
	Note string `sql:"note, null" json:"note"`
}

// GetPost - get one post by id
func (e *Edb) GetPost(id int64) (Post, error) {
	var post Post
	if id == 0 {
		return post, nil
	}
	err := e.db.Model(&post).Where("id = ?", id).Select()
	if err != nil {
		errmsg("GetPost select", err)
	}
	return post, nil
}

// GetPostList - get all post for list
func (e *Edb) GetPostList() ([]PostList, error) {
	var posts []PostList
	_, err := e.db.Query(&posts, `
		SELECT
			id,
			name,
			go,
			note
		FROM
			posts
		ORDER BY
			name ASC`)
	if err != nil {
		errmsg("GetPostList query", err)
	}
	return posts, nil
}

// GetPostSelect - get all post for select
func (e *Edb) GetPostSelect(g bool) ([]SelectItem, error) {
	var posts []SelectItem
	_, err := e.db.Query(&posts, `
		SELECT
			id,
			name
		FROM
			posts
		WHERE
			go = ?
		ORDER BY
			name ASC
	`, g)
	if err != nil {
		errmsg("GetPostSelect query", err)
	}
	return posts, nil
}

// CreatePost - create new post
func (e *Edb) CreatePost(post Post) (int64, error) {
	err := e.db.Insert(&post)
	if err != nil {
		errmsg("CreatePost insert", err)
	}
	return post.ID, nil
}

// UpdatePost - save post changes
func (e *Edb) UpdatePost(post Post) error {
	err := e.db.Update(&post)
	if err != nil {
		errmsg("UpdatePost update", err)
	}
	return err
}

// DeletePost - delete post by id
func (e *Edb) DeletePost(id int64) error {
	if id == 0 {
		return nil
	}
	_, err := e.db.Model(&Post{}).Where("id = ?", id).Delete()
	if err != nil {
		errmsg("DeletePost delete", err)
	}
	return err
}

func (e *Edb) postCreateTable() error {
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
	_, err := e.db.Exec(str)
	if err != nil {
		errmsg("postCreateTable exec", err)
	}
	return err
}

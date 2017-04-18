package edc

// Post - struct for post
type Post struct {
	ID        int64  `sql:"id"         json:"id"`
	Name      string `sql:"name"       json:"name"`
	GO        bool   `sql:"go"         json:"go"`
	Note      string `sql:"note, null" json:"note"`
	CreatedAt string `sql:"created_at" json:"-"`
	UpdatedAt string `sql:"updated_at" json:"-"`
}

// PostList - struct for post list
type PostList struct {
	ID   int64  `sql:"id"         json:"id"`
	Name string `sql:"name"       json:"name"`
	GO   bool   `sql:"go"         json:"go"`
	Note string `sql:"note, null" json:"note"`
}

// GetPost - get one post by id
func (e *Edb) GetPost(id int64) (Post, error) {
	var post Post
	if id == 0 {
		return post, nil
	}
	err := e.db.Model(&post).
		Where("id = ?", id).
		Select()
	if err != nil {
		errmsg("GetPost select", err)
	}
	return post, nil
}

// GetPostList - get all post for list
func (e *Edb) GetPostList() ([]PostList, error) {
	var posts []PostList
	err := e.db.Model(&Post{}).
		Column("id", "name", "go", "note").
		Order("name ASC").
		Select(&posts)
	if err != nil {
		errmsg("GetPostList query", err)
	}
	return posts, nil
}

// GetPostSelect - get post for select
func (e *Edb) GetPostSelect(id int64) (SelectItem, error) {
	var post SelectItem
	if id == 0 {
		return post, nil
	}
	err := e.db.Model(&Post{}).
		Column("id", "name").
		Where("go = false AND id = ?", id).
		Order("name ASC").
		Select(&post)
	if err != nil {
		errmsg("GetPostSelect query", err)
	}
	return post, nil
}

// GetPostGOSelect - get post go for select
func (e *Edb) GetPostGOSelect(id int64) (SelectItem, error) {
	var post SelectItem
	if id == 0 {
		return post, nil
	}
	err := e.db.Model(&Post{}).
		Column("id", "name").
		Where("go = true AND id = ?", id).
		Order("name ASC").
		Select(&post)
	if err != nil {
		errmsg("GetPostGOSelect query", err)
	}
	return post, nil
}

// GetPostSelectAll - get all post for select
func (e *Edb) GetPostSelectAll(g bool) ([]SelectItem, error) {
	var posts []SelectItem
	err := e.db.Model(&Post{}).
		Column("id", "name").
		Where("go = ?", g).
		Order("name ASC").
		Select(&posts)
	if err != nil {
		errmsg("GetPostSelectAll query", err)
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
	_, err := e.db.Model(&Post{}).
		Where("id = ?", id).
		Delete()
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

package edc

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

// PostListGet - get post for list by id
func PostListGet(id int64) (PostList, error) {
	var post PostList
	err := pool.QueryRow(context.Background(), &Post{}).
		Column("id", "name", "go", "note").
		Where("id = ?", id).
		Select(&post)
	if err != nil {
		errmsg("GetPostList select", err)
	}
	return post, nil
}

// PostListAllGet - get all post for list
func PostListAllGet() ([]PostList, error) {
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

// PostSelectGet - get post for select
func PostSelectGet(id int64) (SelectItem, error) {
	var post SelectItem
	if id == 0 {
		return post, nil
	}
	err := pool.QueryRow(context.Background(), &Post{}).
		Column("id", "name").
		Where("go = false AND id = ?", id).
		Order("name ASC").
		Select(&post)
	if err != nil {
		errmsg("GetPostSelect query", err)
	}
	return post, nil
}

// PostGOSelectGet - get post go for select
func PostGOSelectGet(id int64) (SelectItem, error) {
	var post SelectItem
	if id == 0 {
		return post, nil
	}
	err := pool.QueryRow(context.Background(), &Post{}).
		Column("id", "name").
		Where("go = true AND id = ?", id).
		Order("name ASC").
		Select(&post)
	if err != nil {
		errmsg("GetPostGOSelect query", err)
	}
	return post, nil
}

// PostSelectGet - get all post for select
func PostSelectGet(g bool) ([]SelectItem, error) {
	var posts []SelectItem
	err := pool.QueryRow(context.Background(), &Post{}).
		Column("id", "name").
		Where("go = ?", g).
		Order("name ASC").
		Select(&posts)
	if err != nil {
		errmsg("GetPostSelectAll query", err)
	}
	return posts, nil
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

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

// GetPost - get one post by id
func GetPost(id int64) (Post, error) {
	var post Post
	if id == 0 {
		return post, nil
	}
	err := pool.Model(&post).
		Where("id = ?", id).
		Select()
	if err != nil {
		errmsg("GetPost select", err)
	}
	return post, nil
}

// GetPostList - get post for list by id
func GetPostList(id int64) (PostList, error) {
	var post PostList
	err := pool.Model(&Post{}).
		Column("id", "name", "go", "note").
		Where("id = ?", id).
		Select(&post)
	if err != nil {
		errmsg("GetPostList select", err)
	}
	return post, nil
}

// GetPostListAll - get all post for list
func GetPostListAll() ([]PostList, error) {
	var posts []PostList
	err := pool.Model(&Post{}).
		Column("id", "name", "go", "note").
		Order("name ASC").
		Select(&posts)
	if err != nil {
		errmsg("GetPostListAll select", err)
	}
	return posts, nil
}

// GetPostSelect - get post for select
func GetPostSelect(id int64) (SelectItem, error) {
	var post SelectItem
	if id == 0 {
		return post, nil
	}
	err := pool.Model(&Post{}).
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
func GetPostGOSelect(id int64) (SelectItem, error) {
	var post SelectItem
	if id == 0 {
		return post, nil
	}
	err := pool.Model(&Post{}).
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
func GetPostSelectAll(g bool) ([]SelectItem, error) {
	var posts []SelectItem
	err := pool.Model(&Post{}).
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
func CreatePost(post Post) (int64, error) {
	err := pool.Insert(&post)
	if err != nil {
		errmsg("CreatePost insert", err)
	}
	return post.ID, nil
}

// UpdatePost - save post changes
func UpdatePost(post Post) error {
	err := pool.Update(&post)
	if err != nil {
		errmsg("UpdatePost update", err)
	}
	return err
}

// DeletePost - delete post by id
func DeletePost(id int64) error {
	if id == 0 {
		return nil
	}
	_, err := pool.Model(&Post{}).
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
				updated_at TIMESTAMP without time zone default now(),
				UNIQUE (name, go)
			)
	`
	_, err := pool.Exec(str)
	if err != nil {
		errmsg("postCreateTable exec", err)
	}
	return err
}

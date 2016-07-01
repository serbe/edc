package edc

import "fmt"

// Post - struct for post
type Post struct {
	ID    int64  `sql:"id" json:"id"`
	Name  string `sql:"name" json:"name"`
	GO    bool   `sql:"go" json:"go"`
	Notes string `sql:"notes, null" json:"notes"`
}

// GetPost - get one post dy id
func (e *EDc) GetPost(id int64) (post Post, err error) {
	if id == 0 {
		return post, nil
	}
	_, err = e.db.QueryOne(&post, "SELECT * FROM posts WHERE id = ? LIMIT 1", id)
	if err != nil {
		return post, fmt.Errorf("GetPost e.db.QueryRow Scan: %s", err)
	}
	return
}

// GetPostAll - get all post
func (e *EDc) GetPostAll() (posts []Post, err error) {
	_, err = e.db.Query(&posts, "SELECT * FROM posts")
	if err != nil {
		return posts, fmt.Errorf("GetPostAll e.db.Query: %s", err)
	}
	return
}

// GetPostNoGOAll - get all post with no go
func (e *EDc) GetPostNoGOAll() (posts []Post, err error) {
	_, err = e.db.Query(&posts, "SELECT * FROM posts WHERE go = ?", false)
	if err != nil {
		return posts, fmt.Errorf("GetPostNoGOAll e.db.Query: %s", err)
	}
	return
}

// GetPostGOAll - get all post with go
func (e *EDc) GetPostGOAll() (posts []Post, err error) {
	_, err = e.db.Query(&posts, "SELECT * FROM posts WHERE go = ?", true)
	if err != nil {
		return posts, fmt.Errorf("GetPostGOAll e.db.Query: %s", err)
	}
	return
}

// CreatePost - create new post
func (e *EDc) CreatePost(post Post) (err error) {
	err = e.db.Create(&post)
	if err != nil {
		return fmt.Errorf("CreatePost e.db.Exec: %s", err)
	}
	return
}

// UpdatePost - save post changes
func (e *EDc) UpdatePost(post Post) (err error) {
	err = e.db.Update(&post)
	if err != nil {
		return fmt.Errorf("UpdatePost e.db.Exec: %s", err)
	}
	return
}

// DeletePost - delete post by id
func (e *EDc) DeletePost(id int64) error {
	if id == 0 {
		return nil
	}
	_, err := e.db.Exec("DELETE FROM posts WHERE id=?", id)
	if err != nil {
		return fmt.Errorf("DeletePost e.db.Exec: %s", err)
	}
	return nil
}

func (e *EDc) postCreateTable() (err error) {
	str := `CREATE TABLE IF NOT EXISTS posts (id BIGSERIAL PRIMARY KEY, name TEXT, go BOOL NOT NULL DEFAULT FALSE, notes TEXT, UNIQUE (name, go))`
	_, err = e.db.Exec(str)
	if err != nil {
		return fmt.Errorf("postCreateTable e.db.Exec: %s", err)
	}
	return
}

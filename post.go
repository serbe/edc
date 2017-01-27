package edc

import (
	"database/sql"
	"log"
)

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

func scanPost(row *sql.Row) (Post, error) {
	var (
		sID   sql.NullInt64
		sName sql.NullString
		sGo   sql.NullBool
		sNote sql.NullString
		post  Post
	)
	err := row.Scan(&sID, &sName, &sGo, &sNote)
	if err != nil {
		log.Println("scanPost row.Scan ", err)
		return post, err
	}
	post.ID = n2i(sID)
	post.Name = n2s(sName)
	post.GO = n2b(sGo)
	post.Note = n2s(sNote)
	return post, nil
}

func scanPosts(rows *sql.Rows, opt string) ([]Post, error) {
	var posts []Post
	for rows.Next() {
		var (
			sID   sql.NullInt64
			sName sql.NullString
			sGo   sql.NullBool
			sNote sql.NullString
			post  Post
		)
		err := rows.Scan(&sID, &sName, &sGo, &sNote)
		if err != nil {
			log.Println("scanPosts rows.Scan ", err)
			return posts, err
		}
		post.Name = n2s(sName)
		post.GO = n2b(sGo)
		post.Note = n2s(sNote)
		post.ID = n2i(sID)
		posts = append(posts, post)
	}
	err := rows.Err()
	if err != nil {
		log.Println("scanPosts rows.Err ", err)
	}
	return posts, err
}

func scanPostsList(rows *sql.Rows) ([]PostList, error) {
	var posts []PostList
	for rows.Next() {
		var (
			sID   sql.NullInt64
			sName sql.NullString
			sGo   sql.NullBool
			sNote sql.NullString
			post  PostList
		)
		err := rows.Scan(&sID, &sName, &sGo, &sNote)
		if err != nil {
			log.Println("scanPostsList rows.Scan ", err)
			return posts, err
		}
		post.ID = n2i(sID)
		post.Name = n2s(sName)
		post.GO = n2b(sGo)
		post.Note = n2s(sNote)
		posts = append(posts, post)
	}
	err := rows.Err()
	if err != nil {
		log.Println("scanPostsList rows.Err ", err)
	}
	return posts, err
}

func scanPostsSelect(rows *sql.Rows) ([]SelectItem, error) {
	var posts []SelectItem
	for rows.Next() {
		var (
			sID   sql.NullInt64
			sName sql.NullString
			post  SelectItem
		)
		err := rows.Scan(&sID, &sName)
		if err != nil {
			log.Println("scanPostsSelect rows.Scan ", err)
			return posts, err
		}
		post.ID = n2i(sID)
		post.Name = n2s(sName)
		posts = append(posts, post)
	}
	err := rows.Err()
	if err != nil {
		log.Println("scanPostsSelect rows.Err ", err)
	}
	return posts, err
}

// GetPost - get one post by id
func (e *Edb) GetPost(id int64) (Post, error) {
	var post Post
	if id == 0 {
		return post, nil
	}
	err := e.db.Model(&post).Where(`id = ?`, id).Select()
	if err != nil {
		log.Println("GetPost e.db.Select ", err)
		return post, err
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
		log.Println("GetPostList e.db.Query ", err)
		return posts, err
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
		log.Println("GetPostSelect e.db.Query ", err)
		return posts, err
	}
	return posts, nil
}

// CreatePost - create new post
func (e *Edb) CreatePost(post Post) (int64, error) {
	err := e.db.Insert(&post)
	if err != nil {
		log.Println("CreatePost e.db.Insert ", err)
		return 0, err
	}
	return post.ID, nil
}

// UpdatePost - save post changes
func (e *Edb) UpdatePost(post Post) error {
	err := e.db.Update(&post)
	if err != nil {
		log.Println("UpdatePost e.db.Update ", err)
		return err
	}
	return err
}

// DeletePost - delete post by id
func (e *Edb) DeletePost(id int64) error {
	if id == 0 {
		return nil
	}
	_, err := e.db.Exec(`
		DELETE FROM
			posts
		WHERE
			id = $1
	`, id)
	if err != nil {
		log.Println("DeletePost e.db.Exec ", id, err)
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
				updated_at TIMESTAMP without time zone,
				UNIQUE (name, go)
			)
	`
	_, err := e.db.Exec(str)
	if err != nil {
		log.Println("postCreateTable e.db.Exec ", err)
	}
	return err
}

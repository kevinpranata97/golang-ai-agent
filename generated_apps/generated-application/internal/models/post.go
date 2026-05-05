package models

import (
	"time"
	"database/sql"
)

// Post represents the Post entity
type Post struct {
	Id int `json:"id" validate:"required"`
	Title string `json:"title" validate:"required"`
	Content string `json:"content" validate:"required"`
	Author_id int `json:"author_id" validate:"required"`
	Published bool `json:"published" validate:"required"`
	Created_at time.Time `json:"created_at" validate:"required"`
}

// CreatePost creates a new Post in the database
func CreatePost(db *sql.DB, post *Post) error {
	query := `INSERT INTO posts (title, content, author_id, published) VALUES (?, ?, ?, ?)`
	
	result, err := db.Exec(query, post.Title, post.Content, post.Author_id, post.Published)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	post.ID = int(id)
	return nil
}

// GetPostByID retrieves a Post by ID
func GetPostByID(db *sql.DB, id int) (*Post, error) {
	post := &Post{}
	query := `SELECT id, title, content, author_id, published, created_at FROM posts WHERE id = ?`
	
	err := db.QueryRow(query, id).Scan(&post.Id&post.Title&post.Content&post.Author_id&post.Published&post.Created_at)
	if err != nil {
		return nil, err
	}

	return post, nil
}

// GetAllPosts retrieves all Posts
func GetAllPosts(db *sql.DB) ([]Post, error) {
	query := `SELECT id, title, content, author_id, published, created_at FROM posts`
	
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		post := Post{}
		err := rows.Scan(&post.Id&post.Title&post.Content&post.Author_id&post.Published&post.Created_at)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

// UpdatePost updates a Post in the database
func UpdatePost(db *sql.DB, post *Post) error {
	query := `UPDATE posts SET title = ?, content = ?, author_id = ?, published = ? WHERE id = ?`
	
	_, err := db.Exec(query, post.Title, post.Content, post.Author_id, post.Published, post.ID)
	return err
}

// DeletePost deletes a Post from the database
func DeletePost(db *sql.DB, id int) error {
	query := `DELETE FROM posts WHERE id = ?`
	
	_, err := db.Exec(query, id)
	return err
}

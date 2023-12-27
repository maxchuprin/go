package storage

import "time"

// Post - публикация.
type Post struct {
	ID          int       `json:"id" bson:"id"`
	Title       string    `json:"title" bson:"title"`
	Content     string    `json:"content" bson:"content"`
	AuthorID    int       `json:"author_id" bson:"author_id"`
	AuthorName  string    `json:"author_name" bson:"author_name"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
	PublishedAt time.Time `json:"published_at" bson:"published_at"`
}

type Author struct {
	ID   int    `bson:"_id"`
	Name string `bson:"name"`
}

// Interface задаёт контракт на работу с БД.
type Interface interface {
	Posts() ([]Post, error) // получение всех публикаций
	AddPost(Post) error     // создание новой публикации
	UpdatePost(Post) error  // обновление публикации
	DeletePost(Post) error  // удаление публикации по ID
}

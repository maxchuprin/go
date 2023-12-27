package mongodb

import (
	"GoNews/pkg/storage"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const (
	dbName            = "posts"
	collectName       = "posts"
	authorsCollection = "authors"
)

type Storage struct {
	db *mongo.Database
}

type Author struct {
	ID   int    `json:"id" bson:"_id"`
	Name string `json:"name" bson:"name"`
}

func Init(constr string) (*Storage, error) {
	mongoOpts := options.Client().ApplyURI(constr)
	client, err := mongo.Connect(context.Background(), mongoOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect tot mongoDB due error: %v", err)
	}

	if err = client.Ping(context.Background(), nil); err != nil {
		return nil, fmt.Errorf("failed to connect tot mongoDB due error: %v", err)
	}

	s := Storage{
		db: client.Database(dbName),
	}

	return &s, nil
}

func (s Storage) Posts() (posts []storage.Post, err error) {
	cursor, err := s.db.Collection(collectName).Find(context.Background(), bson.M{})
	if cursor.Err() != nil {
		return posts, fmt.Errorf("failed to find all posts due to error: %v", err)
	}

	if err = cursor.All(context.Background(), &posts); err != nil {
		return posts, fmt.Errorf("failed to read all documents from cursor. error: %v", err)
	}

	return posts, nil
}

func (s Storage) AddPost(post storage.Post) error {

	author, err := s.getAuthorByID(post.AuthorID)
	if err != nil {
		return fmt.Errorf("failed to find author to error: %v", err)
	}

	post.CreatedAt = time.Now()
	post.AuthorName = author.Name
	_, err = s.db.Collection(collectName).InsertOne(context.Background(), post)
	if err != nil {
		return fmt.Errorf("failed to create post due to error: %v", err)
	}

	return nil
}

func (s Storage) UpdatePost(post storage.Post) error {
	author, err := s.getAuthorByID(post.AuthorID)
	if err != nil {
		return fmt.Errorf("failed to find author to error: %v", err)
	}

	post.AuthorName = author.Name
	_, err = s.db.Collection(collectName).UpdateOne(context.Background(), bson.M{"id": post.ID}, bson.M{"$set": post})
	if err != nil {
		return fmt.Errorf("failed to update due to error: %v", err)
	}

	return nil
}

func (s Storage) DeletePost(post storage.Post) error {
	_, err := s.db.Collection(collectName).DeleteOne(context.Background(), bson.M{"id": post.ID})
	if err != nil {
		return fmt.Errorf("failed to delete post due to error: %v", err)
	}

	return nil
}

func (s *Storage) getAuthorByID(authorID int) (*Author, error) {
	var author Author
	err := s.db.Collection(authorsCollection).FindOne(context.Background(), bson.M{"_id": authorID}).Decode(&author)
	if err != nil {
		return nil, err
	}
	return &author, nil
}

// AddAuthors своего рода заглушка, для проверки работы создания и апдейта постов
//func (s Storage) AddAuthors() error {
//	authors := []interface{}{
//		Author{ID: 1, Name: "Max"},
//		Author{ID: 2, Name: "John"},
//		Author{ID: 3, Name: "Alice"},
//	}
//	_, err := s.db.Collection("authors").InsertMany(context.Background(), authors)
//	if err != nil {
//		return fmt.Errorf("failed to create authors due to error %v", err)
//	}
//
//	return nil
//}

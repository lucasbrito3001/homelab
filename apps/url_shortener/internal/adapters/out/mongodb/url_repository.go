package mongodb

import (
	"context"
	"fmt"
	"strings"

	"github.com/lucasbrito3001/url_shortner/internal/app/ports"
	"github.com/lucasbrito3001/url_shortner/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDbUrlRepository struct {
	collection *mongo.Collection
}

func NewMongoDbUrlRepository(database *mongo.Database) ports.UrlRepository {
	return &MongoDbUrlRepository{
		collection: database.Collection("urls"),
	}
}

func (r *MongoDbUrlRepository) Save(ctx context.Context, url *domain.ShortenedUrl) error {
	shortenedUrlDocument := fromDomain(url)
	_, err := r.collection.InsertOne(ctx, shortenedUrlDocument)
	return err
}

func (r *MongoDbUrlRepository) FindByCode(ctx context.Context, code domain.Code) (*domain.ShortenedUrl, error) {
	var shortenedUrlDocument shortenedUrlDocument

	filter := bson.M{"code": strings.TrimSpace(string(code))}

	err := r.collection.FindOne(ctx, filter).Decode(&shortenedUrlDocument)
	if err != nil {
		fmt.Println("error getting document on mongo: ", err.Error())
		return nil, err
	}

	return shortenedUrlDocument.toDomain(), nil
}

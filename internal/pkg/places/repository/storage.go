package repository

import (
	"context"
	"errors"
	"log"

	storage2 "cloud.google.com/go/storage"
	"firebase.google.com/go/storage"
	"gitlab.com/v.rianov/favs-backend/internal/models"
	"google.golang.org/api/iterator"
)

type StorageRepository struct {
	cl     *storage.Client
	bucket string
}

func NewStorageRepository(cl *storage.Client, bucket string) StorageRepository {
	return StorageRepository{
		cl:     cl,
		bucket: bucket,
	}
}

func (r StorageRepository) GetPlacePhotoURLs(ctx context.Context, object string) ([]string, models.Status) {
	b, err := r.cl.Bucket(r.bucket)
	if err != nil {
		log.Println("Error while getting bucket: ", err)
		return nil, models.Status{models.InternalError, err.Error()}
	}
	if b == nil {
		log.Println("Bucket is nil")
		return nil, models.Status{models.InternalError, "Bucket is nil"}
	}

	objects := b.Objects(ctx, &storage2.Query{Prefix: object})
	urls := make([]string, 0)
	for {
		attrs, err := objects.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			log.Println("Error while getting object: ", err)
			return nil, models.Status{models.InternalError, err.Error()}
		}
		log.Println("Got object: ", attrs.ContentDisposition, attrs.Name)
		urls = append(urls, attrs.MediaLink)
	}
	return urls, models.Status{models.OK, "OK"}
}

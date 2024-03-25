package repository

import (
	"context"
	"errors"
	"log"

	"cloud.google.com/go/storage"
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
	b := r.cl.Bucket(r.bucket)
	if b == nil {
		log.Println("Bucket is nil")
		return nil, models.Status{models.InternalError, "Bucket is nil"}
	}

	objects := b.Objects(ctx, &storage.Query{Prefix: object})
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

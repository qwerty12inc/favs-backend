package repository

import (
	"context"
	"errors"
	"time"

	storage2 "cloud.google.com/go/storage"
	"firebase.google.com/go/storage"
	"github.com/labstack/gommon/log"
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

func (r StorageRepository) GenerateSignedURL(ctx context.Context, object string) (string, models.Status) {
	b, err := r.cl.Bucket(r.bucket)
	if err != nil {
		log.Error("Error while getting bucket: ", err)
		return "", models.Status{models.InternalError, err.Error()}
	}
	if b == nil {
		log.Error("Bucket is nil")
		return "", models.Status{models.InternalError, "Bucket is nil"}
	}

	o := b.Object(object)
	_, err = o.Attrs(ctx)
	if err != nil {
		log.Error("Error while getting object: ", err)
		return "", models.Status{models.InternalError, err.Error()}
	}

	url, err := b.SignedURL(object, &storage2.SignedURLOptions{
		Method:  "GET",
		Expires: time.Now().Add(time.Hour),
	})
	if err != nil {
		log.Error("Error while generating signed url: ", err)
		return "", models.Status{models.InternalError, err.Error()}
	}

	return url, models.Status{models.OK, "OK"}
}

func (r StorageRepository) GetPlacePhotoURLs(ctx context.Context, object string) ([]string, models.Status) {
	b, err := r.cl.Bucket(r.bucket)
	if err != nil {
		log.Error("Error while getting bucket: ", err)
		return nil, models.Status{models.InternalError, err.Error()}
	}
	if b == nil {
		log.Error("Bucket is nil")
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
			log.Error("Error while getting object: ", err)
			return nil, models.Status{models.InternalError, err.Error()}
		}
		log.Debug("Got object: ", attrs.ContentDisposition, attrs.Name)
		urls = append(urls, attrs.Name)
	}
	return urls, models.Status{models.OK, "OK"}
}

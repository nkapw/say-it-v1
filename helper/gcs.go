package helper

import (
	"cloud.google.com/go/storage"
	"context"
	"google.golang.org/api/option"
)

func CreateGCSClient() (*storage.Client, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile("./credentials.json"))
	if err != nil {
		return nil, err
	}
	return client, nil
}

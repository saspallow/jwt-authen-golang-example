package api

import (
	"context"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

// Config API type
type Config struct {
	ServiceAccountJSON []byte
	ProjectID          string
}

var (
	client *firestore.Client
)

// Init api
func Init(cfg Config) error {

	ctx := context.Background()

	sa := option.WithCredentialsFile("service-account.json")
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	return err
}

func getContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 30*time.Second)
}

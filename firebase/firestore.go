package firebase

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"

	"locator-backend/model"
)

var FirestoreClient *firestore.Client

func InitFirestore() {
	ctx := context.Background()
	sa := option.WithCredentialsFile("serviceAccountKey.json")

	client, err := firestore.NewClient(ctx, "locator-dccf6", sa)
	if err != nil {
		log.Fatalf("Failed to connect to Firestore: %v", err)
	}

	FirestoreClient = client
}

func GetUsersFromFirestore() ([]model.User, error) {
	ctx := context.Background()
	sa := option.WithCredentialsFile("path/to/serviceAccountKey.json") // File JSON dari Firebase Project

	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		return nil, err
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	var users []model.User
	docs, err := client.Collection("users").Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	for _, doc := range docs {
		var user model.User
		if err := doc.DataTo(&user); err != nil {
			log.Println("Error mapping user:", err)
			continue
		}
		users = append(users, user)
	}

	return users, nil
}
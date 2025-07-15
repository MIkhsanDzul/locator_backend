package firebase

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
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

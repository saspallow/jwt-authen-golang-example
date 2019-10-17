package api

import (
	"jwt-authen-golang-example/model"
	"log"

	"google.golang.org/api/iterator"
)

const kindUser = "User"

// FindUser from firestore
func FindUser(username, password string) (*model.User, error) {
	ctx, cancel := getContext()
	defer cancel()

	var user model.User
	iter := client.Collection(kindUser).Where("Username", "==", username).Limit(1).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Println("test 2")
			return nil, err
		}

		doc.DataTo(&user)
		user.SetKey(doc.Ref)
	}

	log.Println(user)
	if !user.ComparePassword(password) {
		// wrong password return like user not found
		return nil, nil
	}
	return &user, nil
}

// SaveUser to firestore
func SaveUser(user *model.User) error {
	ctx, cancel := getContext()
	defer cancel()

	var err error
	user.Stamp()
	key := user.Key()
	if key == nil {
		key = client.Collection(kindUser).NewDoc()
	}

	_, err = client.Collection(kindUser).Doc(key.ID).Set(ctx, user)
	if err != nil {
		return err
	}
	user.SetKey(key)
	return nil
}

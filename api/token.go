package api

import (
	"jwt-authen-golang-example/model"
	"time"
	"errors"
	"google.golang.org/api/iterator"
)

const kindToken = "Token"

// CreateToken save new token to database
func CreateToken(token string, userID string) error {
	ctx, cancel := getContext()
	defer cancel()

	var err error
	tk := &model.Token{
		Token:  token,
		UserID: userID,
	}
	tk.Stamp()
	key := client.Collection(kindUser).NewDoc()
	_, err = client.Collection(kindToken).Doc(key.ID).Set(ctx, tk)
	if err != nil {
		return err
	}
	tk.SetKey(key)
	return nil
}

func getToken(token string) (*model.Token, error) {
	ctx, cancel := getContext()
	defer cancel()

	var tk model.Token
	iter := client.Collection(kindToken).Where("Token", "==", token).Limit(1).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		doc.DataTo(&tk)
		tk.SetKey(doc.Ref)
	}

	if &tk == nil {
		return nil, errors.New("Not found")
	}

	return &tk, nil
}

// DeleteToken delete a token from firestore
func DeleteToken(token string) error {
	tk, err := getToken(token)
	if err != nil {
		return err
	}
	ctx, cancel := getContext()
	defer cancel()

	_, err = client.Collection(kindToken).Doc(tk.Key().ID).Delete(ctx)
	if err != nil {
		return err
	}
	return nil
}

// ValidateToken validate and update token last access timestamp
func ValidateToken(token string, userID string, expiresInFromLastAccess time.Duration) (bool, error) {
	tk, err := getToken(token)
	if err != nil {
		return false, err
	}
	if tk == nil || tk.UserID != userID {
		return false, nil
	}
	if time.Now().After(tk.LastAccessAt.Add(expiresInFromLastAccess)) {
		// token expired
		// remove expired token from database
		go DeleteToken(token)
		return false, nil
	}
	tk.Stamp()
	go func(tk model.Token) {
		ctx, cancel := getContext()
		defer cancel()
		client.Collection(kindToken).Doc(tk.Key().ID).Set(ctx, tk)
	}(*tk)
	return true, nil
}

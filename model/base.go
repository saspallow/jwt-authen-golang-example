package model

import (
	"cloud.google.com/go/firestore"
)

// Base type provides datastofirestorere-based model
type Base struct {
	key *firestore.DocumentRef
	ID  string `firestore:"-" json:"id"`
}

// Key return firestore key or nil
func (x *Base) Key() *firestore.DocumentRef {
	return x.key
}

// SetKey sets key and id to new given key
func (x *Base) SetKey(key *firestore.DocumentRef) {
	x.key = key
	x.ID = key.ID
}

package store

import (
	"context"
	"log"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/protobuf/proto"
)

type Store[T proto.Message] struct {
	locaColl   *mongo.Collection
	protoField string
}

// add your mongo uri, and collection name
// connect to your proto.Message type
// e.g. store.Connect[*proto.Message]("mongodb://localhost:27017", "info")
func Connect[T proto.Message](uri string, opts ...ClientOption) Store[T] {
	var err error

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	msg := *new(T)
	pbName := string(msg.ProtoReflect().Descriptor().Name())
	pbName = strings.ToLower(pbName)

	db := client.Database("info")
	db.Collection(pbName)

	return Store[T]{
		locaColl:   db.Collection(pbName),
		protoField: pbName,
	}
}

// clientOpts is used to pass data the ClientOption function
type clientOpts struct {
	app string
	db  *mongo.Database
}

// ClientOption is used to pass optional arguments when creating a Client.
type ClientOption func(*clientOpts) error

// App is the name of the Mongo App to connect to .
func WithApp(a string) ClientOption {
	return func(c *clientOpts) error {
		c.app = a
		return nil
	}
}

// set reference to Mongo database
func (s *Store[T]) setDB(db *mongo.Database) error {
	// apply configuration (indexes, etc.)
	s.locaColl = db.Collection(s.protoField)

	// track data changes for auditing purposes
	return nil
}

// WithStore adds a Store reference to the client.
func WithStore[T proto.Message](s *Store[T]) ClientOption {
	return func(c *clientOpts) error {
		return s.setDB(c.db)
	}
}

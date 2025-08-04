package mongo

import (
	"context"
	"fmt"
	"time"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// --- Interfaces ---

type Database interface {
	Collection(string) Collection
	Client() Client
}

type Collection interface {
	FindOne(ctx context.Context, filter any) SingleResult
	InsertOne(ctx context.Context, document any) (*mongo.InsertOneResult, error)
	InsertMany(ctx context.Context, documents []any) (*mongo.InsertManyResult, error)
	DeleteOne(ctx context.Context, filter any) (int64, error)
	Find(ctx context.Context, filter any, opts ...*options.FindOptions) (Cursor, error)
	CountDocuments(ctx context.Context, filter any, opts ...*options.CountOptions) (int64, error)
	Aggregate(ctx context.Context, pipeline any) (Cursor, error)
	UpdateOne(ctx context.Context, filter, update any, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	UpdateMany(ctx context.Context, filter, update any, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
}

type Client interface {
	Database(name string) Database
	Connect(ctx context.Context) error
	Disconnect(ctx context.Context) error
	StartSession() (mongo.Session, error)
	UseSession(ctx context.Context, fn func(mongo.SessionContext) error) error
	Ping(ctx context.Context) error
}

type SingleResult interface {
	Decode(v any) error
}

type Cursor interface {
	Close(ctx context.Context) error
	Next(ctx context.Context) bool
	Decode(v any) error
	All(ctx context.Context, results any) error
	Err() error
}

// --- Struct Wrappers ---

type mongoClient struct {
	cl *mongo.Client
}

type mongoDatabase struct {
	db *mongo.Database
}

type mongoCollection struct {
	coll *mongo.Collection
}

type mongoSingleResult struct {
	sr *mongo.SingleResult
}

type mongoCursor struct {
	mc *mongo.Cursor
}

// --- Utility Function ---

func ErrNoDocuments() error {
	return mongo.ErrNoDocuments
}

// --- Factory ---

func NewClient(uri string) (Client, error) {
	time.Local = time.UTC
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}
	log.Println("MongoDB connection established successfully")
	return &mongoClient{cl: client}, nil
}

// --- Client Methods ---

func (mc *mongoClient) Connect(ctx context.Context) error {
	// Connection is already done in NewClient
	return nil
}

func (mc *mongoClient) Disconnect(ctx context.Context) error {
	return mc.cl.Disconnect(ctx)
}

func (mc *mongoClient) Ping(ctx context.Context) error {
	return mc.cl.Ping(ctx, readpref.Primary())
}

func (mc *mongoClient) Database(name string) Database {
	return &mongoDatabase{db: mc.cl.Database(name)}
}

func (mc *mongoClient) StartSession() (mongo.Session, error) {
	return mc.cl.StartSession()
}

func (mc *mongoClient) UseSession(ctx context.Context, fn func(mongo.SessionContext) error) error {
	return mc.cl.UseSession(ctx, fn)
}

// --- Database Methods ---

func (md *mongoDatabase) Collection(name string) Collection {
	return &mongoCollection{coll: md.db.Collection(name)}
}

func (md *mongoDatabase) Client() Client {
	return &mongoClient{cl: md.db.Client()}
}

// --- Collection Methods ---

func (mc *mongoCollection) FindOne(ctx context.Context, filter any) SingleResult {
	return &mongoSingleResult{sr: mc.coll.FindOne(ctx, filter)}
}

func (mc *mongoCollection) InsertOne(ctx context.Context, document any) (*mongo.InsertOneResult, error) {
	res, err := mc.coll.InsertOne(ctx, document)
	return res, err
}

func (mc *mongoCollection) InsertMany(ctx context.Context, documents []any) (*mongo.InsertManyResult, error) {
	res, err := mc.coll.InsertMany(ctx, documents)
	return res, err
}

func (mc *mongoCollection) DeleteOne(ctx context.Context, filter any) (int64, error) {
	res, err := mc.coll.DeleteOne(ctx, filter)
	return res.DeletedCount, err
}

func (mc *mongoCollection) Find(ctx context.Context, filter any, opts ...*options.FindOptions) (Cursor, error) {
	cursor, err := mc.coll.Find(ctx, filter, opts...)
	return &mongoCursor{mc: cursor}, err
}

func (mc *mongoCollection) CountDocuments(ctx context.Context, filter any, opts ...*options.CountOptions) (int64, error) {
	return mc.coll.CountDocuments(ctx, filter, opts...)
}

func (mc *mongoCollection) Aggregate(ctx context.Context, pipeline any) (Cursor, error) {
	cursor, err := mc.coll.Aggregate(ctx, pipeline)
	return &mongoCursor{mc: cursor}, err
}

func (mc *mongoCollection) UpdateOne(ctx context.Context, filter, update any, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return mc.coll.UpdateOne(ctx, filter, update, opts...)
}

func (mc *mongoCollection) UpdateMany(ctx context.Context, filter, update any, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return mc.coll.UpdateMany(ctx, filter, update, opts...)
}

// --- SingleResult Methods ---

func (sr *mongoSingleResult) Decode(v any) error {
	return sr.sr.Decode(v)
}

// --- Cursor Methods ---

func (c *mongoCursor) Close(ctx context.Context) error {
	return c.mc.Close(ctx)
}

func (c *mongoCursor) Next(ctx context.Context) bool {
	return c.mc.Next(ctx)
}

func (c *mongoCursor) Decode(v any) error {
	return c.mc.Decode(v)
}

func (c *mongoCursor) All(ctx context.Context, results any) error {
	return c.mc.All(ctx, results)
}

func (c *mongoCursor) Err() error {
	return c.mc.Err()
}

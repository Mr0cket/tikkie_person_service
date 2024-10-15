package mongo

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongo struct {
	client *mongo.Client
	db     string
	ctx    context.Context
}

type MongoUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Database client
func NewClient(ctx context.Context, mongoUri string, database string, mongoUser MongoUser) *Mongo {
	connectionString := fmt.Sprintf("mongodb://%s:%s@%s/?retryWrites=false&directConnection=true", mongoUser.Username, mongoUser.Password, mongoUri)
	fmt.Printf("Mongo connectionString: mongodb://<credentials>@%s/?retryWrites=false&directConnection=truen\n", mongoUri)
	opts := options.Client().ApplyURI(connectionString).SetMaxConnIdleTime(time.Second * 30)
	mongoClient, err := mongo.Connect(ctx, opts)
	if err != nil {
		log.Fatal(err)
	}
	if err := mongoClient.Ping(context.TODO(), nil); err != nil {
		panic(err)
	}
	log.Println("Connected to MongoDB")

	return &Mongo{
		client: mongoClient,
		db:     database,
		ctx:    ctx,
	}
}

func (m *Mongo) FindDocument(collection string, query bson.M) mongo.SingleResult {
	return *m.client.Database(m.db).Collection(collection).FindOne(m.ctx, query)
}

func (m *Mongo) InsertDocument(collection string, document interface{}) (mongo.InsertOneResult, error) {
	result, err := m.client.Database(m.db).Collection(collection).InsertOne(m.ctx, document)
	return *result, err
}

// As there is no requirement to filter the list of documents,
// We can simplify the function by listing all documents in the collection
func (m *Mongo) ListDocuments(collection string) mongo.Cursor {

	cursor, err := m.client.Database(m.db).Collection(collection).Find(m.ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	return *cursor
}

func (m *Mongo) Close() {
	defer func() {
		if err := m.client.Disconnect(m.ctx); err != nil {
			panic(err)
		}
	}()
}

// Returns true if the connection to MongoDB is successful
func (m *Mongo) Test() bool {
	if err := m.client.Ping(context.TODO(), nil); err != nil {
		log.Printf("Error pinging MongoDB: %v\n", err)
		return false
	}
	log.Println("Successfully connected to MongoDB")
	return true
}

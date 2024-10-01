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
	logger *log.Logger
	client *mongo.Client
	db     string
}

// Database client
func NewClient(mongoUri string, database string) *Mongo {
	logger := log.New(log.Writer(), "mongo", log.Flags())
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(mongoUri).SetServerAPIOptions(serverAPI).SetMaxConnIdleTime(time.Duration(10) * time.Second)
	fmt.Println(opts.MaxConnIdleTime)
	mongoClient, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		logger.Fatal(err)
	}
	db := mongoClient.Database(database)

	// Test MongoDB connection
	var result bson.M
	if err := db.RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&result); err != nil {
		panic(err)
	}
	logger.Printf("Connected to MongoDB: %+v\n", result["ok"])

	return &Mongo{
		logger: logger,
		client: mongoClient,
		db:     database,
	}
}

func (m *Mongo) FindDocument(collection string, query bson.M) mongo.SingleResult {
	return *m.client.Database(m.db).Collection(collection).FindOne(context.TODO(), query)
}

func (m *Mongo) InsertDocument(collection string, document interface{}) (mongo.InsertOneResult, error) {
	result, err := m.client.Database(m.db).Collection(collection).InsertOne(context.TODO(), document)
	return *result, err
}

// As there is no requirement to filter the list of documents,
// We can simplify the function by listing all documents in the collection
func (m *Mongo) ListDocuments(collection string) mongo.Cursor {

	cursor, err := m.client.Database(m.db).Collection(collection).Find(context.TODO(), bson.D{})
	if err != nil {
		m.logger.Fatal(err)
	}
	return *cursor
}

func (m *Mongo) Close() {
	defer func() {
		if err := m.client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
}

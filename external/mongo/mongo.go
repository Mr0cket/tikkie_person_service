package mongo

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongo struct {
	logger *log.Logger
	db     *mongo.Database
}

// Database client
func NewClient(mongoUri string, database string) *Mongo {
	logger := log.New(log.Writer(), "mongo", log.Flags())
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(mongoUri).SetServerAPIOptions(serverAPI)
	mongoClient, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		logger.Fatal(err)
	}
	defer func() {
		if err = mongoClient.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	db := mongoClient.Database(database)

	// Test MongoDB connection
	var result bson.M
	if err := db.RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&result); err != nil {
		panic(err)
	}
	logger.Printf("Connected to MongoDB: %+v\n", result["ok"])

	return &Mongo{
		logger: logger,
		db:     db,
	}
}

func (m *Mongo) FindDocument(collection string, query bson.M) mongo.SingleResult {
	return *m.db.Collection(collection).FindOne(context.TODO(), query)
}

func (m *Mongo) InsertDocument(collection string, document interface{}) (mongo.InsertOneResult, error) {
	result, err := m.db.Collection(collection).InsertOne(context.TODO(), document)
	return *result, err
}

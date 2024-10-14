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
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	connectionString := fmt.Sprintf("mongodb://%s:%s@%s/?tls=true&tlsCAFile=global-bundle.pem&retryWrites=false", mongoUser.Username, mongoUser.Password, mongoUri)
	opts := options.Client().ApplyURI(connectionString).SetServerAPIOptions(serverAPI).SetMaxConnIdleTime(time.Duration(10) * time.Second)
	fmt.Println(opts.MaxConnIdleTime)
	mongoClient, err := mongo.Connect(ctx, opts)
	if err != nil {
		log.Fatal(err)
	}
	db := mongoClient.Database(database)

	// Test MongoDB connection
	var result bson.M
	if err := db.RunCommand(ctx, bson.D{{"ping", 1}}).Decode(&result); err != nil {
		panic(err)
	}
	log.Printf("Connected to MongoDB: %+v\n", result["ok"])

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

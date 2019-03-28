package sandboxmongodb

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	// "github.com/mongodb/mongo-go-driver/bson"
	//documentation_examples "github.com/mongodb/mongo-go-driver/examples/documentation_examples"
	documentation_examples "go.mongodb.org/mongo-driver/examples/documentation_examples"
	//"github.com/mongodb/mongo-go-driver/mongo"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestCursor(t *testing.T) {

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://mongo1:27017,mongo2:27018,mongo3:27019/?replicaSet=my-mongo-set"))

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)

	// db := client.Database("documentation_examples")
	db := DocumentationDatabase(client)
	documentation_examples.InsertExamples(t, db)

	// we are going to monitor the inventory collection
	// for changes
	coll := db.Collection("inventory")

	ctx, _ = context.WithTimeout(context.Background(), 30*time.Second)
	cur, err := coll.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var result bson.M
		err := cur.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		log.Print("A change?")
		// do something with result....
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	t.Log("end forever")
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
}

func TestWatchCursor(t *testing.T) {

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://mongo1:27017,mongo2:27018,mongo3:27019/?replicaSet=my-mongo-set"))

	ctx, _ := context.WithTimeout(context.Background(), 1000*time.Second)
	err = client.Connect(ctx)

	db := DocumentationDatabase(client)
	documentation_examples.InsertExamples(t, db)

	// we are going to monitor the inventory collection
	// for changes
	coll := db.Collection("inventory")

	cur, err := coll.Watch(ctx, pipeline) // Watch can only be done against a MongoDB REPLICA set
	if err != nil {
		fmt.Println(err)
		t.Log(err)
		return
	}
	defer cur.Close(ctx)

	for ever := false; !ever; { // forever...
		for cur.Next(ctx) {
			t.Log("something in collection changed!")

			var result bson.M
			if err := cur.Decode(&result); err != nil {
				log.Fatal(err)
			}
			// do something with result
			t.Log(result)
		}
		time.Sleep(1 * time.Second) // should be put in async process
		t.Log("A change..")
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
}

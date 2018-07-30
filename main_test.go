package sandboxmongodb

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	documentation_examples "github.com/mongodb/mongo-go-driver/examples/documentation_examples"
	"github.com/mongodb/mongo-go-driver/mongo"
)

func TestCursor(t *testing.T) {

	client, err := mongo.Connect(context.Background(), "mongodb://mongo1:27017,mongo2:27018,mongo3:27019/?replicaSet=my-mongo-set", nil)

	// db := client.Database("documentation_examples")
	db := DocumentationDatabase(client)
	documentation_examples.InsertExamples(t, db)

	// we are going to monitor the inventory collection
	// for changes
	coll := db.Collection("inventory")

	ctx := context.Background()

	var pipeline interface{} // set up pipeline

	cur, err := coll.Watch(ctx, pipeline) // can only be done against a MongoDB REPLICA set
	if err != nil {
		// Handle err
		fmt.Println(err)
		t.Log(err)
		return
	}
	defer cur.Close(ctx)

	for ever := false; !ever; { // forever... hmmm...
		for cur.Next(ctx) {
			t.Log("something in collection changed!")
			elem := bson.NewDocument()
			if err := cur.Decode(elem); err != nil {
				log.Fatal(err)
			}
			// do something with elem or for example update cache :)

		}
		time.Sleep(1 * time.Second) // should be put in async process
		t.Log("olee..")
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
}

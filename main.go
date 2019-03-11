package sandboxmongodb

import (
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {

}

// DocumentationDatabase ...
func DocumentationDatabase(mgoClient *mongo.Client) *mongo.Database {

	//client, err := mongo.Connect(context.Background(), "mongodb://mongo1:30001/", nil)

	db := mgoClient.Database("documentation_examples")
	return db
}

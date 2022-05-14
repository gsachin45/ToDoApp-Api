package controller

import (
	"ToDoApp/model"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func goDotEnvInit() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("err")
	}
}

//collection object/instance
var collection *mongo.Collection

//connect with mongoDB

func init() {
	goDotEnvInit()

	connectionString := os.Getenv("MongoDB_CString")
	dbName := os.Getenv("DB_NAME")
	colName := os.Getenv("COL_NAME")

	//client option
	clientOption := options.Client().ApplyURI(connectionString)

	//connect to mongodb
	client, err := mongo.Connect(context.TODO(), clientOption)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("MongoDB connection success")

	collection = client.Database(dbName).Collection(colName)

	fmt.Println("Collection instance is ready")

}

func insertOneTask(task model.ToDoList) {
	inserted, err := collection.InsertOne(context.Background(), task)

	if err != nil {
		log.Fatal(err)

	}
	fmt.Println("One Task added", inserted.InsertedID)
}

func undoTask(task_id string) {
	fmt.Println(task_id)
	id, _ := primitive.ObjectIDFromHex(task_id)

	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"status": false}}

	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("modified count: ", result.ModifiedCount)
}

func taskCompleted(task_id string) {
	fmt.Println(task_id)
	id, _ := primitive.ObjectIDFromHex(task_id)

	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"status": true}}

	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("modified count: ", result.ModifiedCount)
}

func deleteOneTask(task_id string) {
	fmt.Println(task_id)
	id, _ := primitive.ObjectIDFromHex(task_id)
	filter := bson.M{"_id": id}

	result, err := collection.DeleteOne(context.Background(), filter)

	if err != nil {
		log.Fatal(err)

	}
	fmt.Println("Task deleted", result.DeletedCount)

}

func deleteAllTask() int64 {
	result, err := collection.DeleteMany(context.Background(), bson.D{{}}, nil)
	if err != nil {
		log.Fatal(err)

	}
	fmt.Println("Number of task deleted", result.DeletedCount)
	return result.DeletedCount
}

func getAllTask() []primitive.M {

	cur, err := collection.Find(context.Background(), bson.D{{}})

	if err != nil {
		log.Fatal(err)

	}
	var tasks []primitive.M

	for cur.Next(context.Background()) {
		var task bson.M
		err := cur.Decode(&task)
		if err != nil {
			log.Fatal(err)
		}
		tasks = append(tasks, task)

	}
	defer cur.Close(context.Background())
	return tasks
}

func GetAllTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	allTask := getAllTask()

	json.NewEncoder(w).Encode(allTask)

}
func CreateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	var task model.ToDoList
	_ = json.NewDecoder(r.Body).Decode(&task)
	insertOneTask(task)
	json.NewEncoder(w).Encode(task)
}
func UndoTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	w.Header().Set("Allow-Control-Allow-Methods", "PUT")

	params := mux.Vars(r)
	undoTask(params["id"])
	json.NewEncoder(w).Encode(params["id"])

}
func TaskCompleted(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	w.Header().Set("Allow-Control-Allow-Methods", "PUT")

	params := mux.Vars(r)
	taskCompleted(params["id"])
	json.NewEncoder(w).Encode(params["id"])

}
func DeleteOneTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	w.Header().Set("Allow-Control-Allow-Methods", "DELETE")

	params := mux.Vars(r)
	deleteOneTask(params["id"])
	json.NewEncoder(w).Encode(params["id"])

}
func DeleteAllTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	w.Header().Set("Allow-Control-Allow-Methods", "DELETE")

	count := deleteAllTask()

	json.NewEncoder(w).Encode(count)

}

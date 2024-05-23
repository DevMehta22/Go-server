package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/DevMehta22/mongoapi/model"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func goDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")
  
	if err != nil {
	  log.Fatalf("Error loading .env file")
	}
  
	return os.Getenv(key)
  }

//IMP
var collection *mongo.Collection

//connection with mongoDB

func init()  {
 connectionString := goDotEnvVariable("MONGO_URI")
 dbName := goDotEnvVariable("DB_NAME")
 colName := goDotEnvVariable("COLLECTION_NAME")
	//client option
	clientOption := options.Client().ApplyURI(connectionString)

	//connect to mongodb
	client,err := mongo.Connect(context.TODO(),clientOption)
	
	if err!=nil{
		log.Fatal(err)
	}
	fmt.Println("Connected to DB!")
	collection = client.Database(dbName).Collection(colName)

	//collection instance
	fmt.Println("Collection instance is ready")
}

//MONGODB helpers - file

//insert 1 record
func insertMovie(movie model.Netflix)  {
	inserted,err := collection.InsertOne(context.Background(),movie)

	if err !=nil{
		log.Fatal(err)
	}
	fmt.Println("Inserted Movie: ",inserted.InsertedID)
}

//update record
func updateMovie(movieID string)  {
	id,_ := primitive.ObjectIDFromHex(movieID)
	filter := bson.M{"_id":id}
	update := bson.M{"$set":bson.M{"watched":true}}

	result,err := collection.UpdateOne(context.Background(),filter,update)
	if err!=nil{
		log.Fatal(err)
	}
	fmt.Println("Updated Movie: ",result.ModifiedCount)
}

//delete record
func deleteMovie(movieID string)  {
	id,_ := primitive.ObjectIDFromHex(movieID)
	filter := bson.M{"_id":id}
	result,err := collection.DeleteOne(context.Background(),filter)
	if err!=nil{
		log.Fatal(err)
		}
	fmt.Println("Deleted Movie: ",result.DeletedCount)
}

//delete all
func deleteAllMovies() int64 {
	result,err := collection.DeleteMany(context.Background(),bson.D{{}})
	if err!=nil{
		log.Fatal(err)
		}
	fmt.Println("Deleted All Movies: ",result.DeletedCount)
	return result.DeletedCount
}

//find all records
func findMovies() []primitive.M{
	cursor,err := collection.Find(context.Background(),bson.D{{}})
	if err!=nil{
		log.Fatal(err)
		}
	
	var movies []primitive.M
	
	defer cursor.Close(context.Background())
	
	for cursor.Next(context.Background()){
		var movie bson.M
		err := cursor.Decode(&movie)
		if err!=nil{
			log.Fatal(err)
			}
		movies = append(movies, movie)
		}
	return movies
}

//Actual controller - file

func GetAllMovies(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-Type","application/json")
	movies := findMovies()
	json.NewEncoder(w).Encode(movies)
}

func CreateMovie(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-Type","application/json")
	w.Header().Set("Allow-Control-Allow-Methods","POST")

	var movie model.Netflix
	_ = json.NewDecoder(r.Body).Decode(&movie)
	insertMovie(movie)
	json.NewEncoder(w).Encode(movie)
}

func MarkAsWatched(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-Type","application/json")
	w.Header().Set("Allow-Control-Allow-Methods","PUT")

	params := mux.Vars(r)
	updateMovie(params["id"])
	json.NewEncoder(w).Encode("Movie Updated:"+params["id"])
}

func DeleteMove(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-Type","application/json")
	w.Header().Set("Allow-Control-Allow-Methods","DELETE")

	params := mux.Vars(r)

	deleteMovie(params["id"])
	json.NewEncoder(w).Encode("Data deleted:"+params["id"])
}

func DeleteMovies(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-Type","application/json")
	w.Header().Set("Allow-Control-Allow-Methods","DELETE")

	count := deleteAllMovies()
	json.NewEncoder(w).Encode(count)
}

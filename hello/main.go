package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Trainer struct {
	Name string `bson: name`
	Age  int    `bson: age`
	City string `bson: city`
}

var trainers []interface{}

func databaseInit() (*mongo.Client, error) {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	return client, err
}

var client, err = databaseInit()
var collection = client.Database("test").Collection("trainers")

func trainersInit(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Header().Set("Content-type", "application/json")
		trainers = append(trainers, Trainer{"Ash", 10, "Pallet Town test"})
		trainers = append(trainers, Trainer{"Misty", 10, "Cerulean City"})
		trainers = append(trainers, Trainer{"Brock", 15, "Pewter City"})

		insertManyResult, err := collection.InsertMany(context.TODO(), trainers)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Inserted multiple documents: ", insertManyResult.InsertedIDs)
		json.NewEncoder(w).Encode(trainers)
		return
	}
}

func getTrainers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "json/application")
	var t Trainer
	cursor, err := collection.Find(context.TODO(), bson.D{})

	if err != nil {
		fmt.Println("Error")
	} else {
		for cursor.Next(context.TODO()) {
			err := cursor.Decode(&t)

			if err != nil {

			} else {
				json.NewEncoder(w).Encode(t)
			}
		}

	}
	fmt.Println(t)

	return
}

func getTrainer(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Header().Set("Content-type", "json/application")
		var t Trainer
		name := r.URL.Query().Get("name")
		cursor, err := collection.Find(context.TODO(), bson.M{"name": name})

		if err != nil {
			log.Println(err)
		} else {
			for cursor.Next(context.TODO()) {
				err := cursor.Decode(&t)
				if err != nil {
					fmt.Println("Error")
				}
				json.NewEncoder(w).Encode(trainers)
			}
		}
	}
	http.NotFound(w, r)
	return
}
func createTrainer(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		w.Header().Set("Content-type", "application/json")
		var t Trainer
		json.NewDecoder(r.Body).Decode(&t)
		trainers = append(trainers, t)
		_, err = collection.InsertOne(context.TODO(), trainers[len(trainers)-1])
		fmt.Println(trainers[len(trainers)-1])
		json.NewEncoder(w).Encode(&t)
		return
	}

	http.NotFound(w, r)
	return
}

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("/", getTrainers)
	mux.HandleFunc("/init", trainersInit)
	mux.HandleFunc("/trainer", getTrainer)
	mux.HandleFunc("/AddTrainer", createTrainer)
	log.Fatal(http.ListenAndServe(":4000", mux))

}

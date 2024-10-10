package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Person struct {
	Name string `json:"name"`
	City string `json:"city"`
	Age  string `json:"age"`
}

var userCollection = db().Database("MuxDB").Collection("users")

// CreateProfile Function
func createProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user Person

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		log.Fatal(" error during the decoding of the creation function ", err)
	}

	insertResult, err := userCollection.InsertOne(context.TODO(), user)

	if err != nil {
		log.Fatal(" error during the insertion of the creatin function ", err)
	}

	json.NewEncoder(w).Encode(insertResult.InsertedID)
}

// getAllUsers Function
func getAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Find all users
	cursor, err := userCollection.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Fatal("error while getting all users:", err)
		return
	}
	defer cursor.Close(context.TODO())

	// Create a slice to hold the results
	var users []Person
	if err = cursor.All(context.TODO(), &users); err != nil {
		log.Fatal("error decoding users:", err)
		return
	}

	// Encode the result to JSON
	json.NewEncoder(w).Encode(users)
}

// getUserProfile Function
func getUserProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // Get the URL parameters
	id := params["id"]    // Assuming you passed the ID as a parameter

	// Convert the string ID to a MongoDB ObjectID
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Fatal("error at id conversion in get req in userprofile", err)
	}

	var user Person
	err = userCollection.FindOne(context.TODO(), bson.D{{Key: "_id", Value: _id}}).Decode(&user)
	if err != nil {
		log.Fatal("error at finding profile in userprofile", err)
	}

	json.NewEncoder(w).Encode(user) // Encode the user profile as JSON
}

// deleteProfile Functions
func deleteProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)["id"]
	_id, err := primitive.ObjectIDFromHex(params)
	if err != nil {
		log.Fatal("error at id conversion in delete req", err)
	}

	res, err := userCollection.DeleteOne(context.TODO(), bson.D{{Key: "_id", Value: _id}})
	if err != nil {
		log.Fatal("error at deleting profile in delete req", err)
	}

	if res.DeletedCount == 0 {
		http.Error(w, "No profile found with that ID", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Profile deleted successfully"})
}

// updateProfile Function
func updateProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	type updateBody struct {
		Name string `json:"name"` //value that has to be matched
		City string `json:"city"` // value that has to be modified
	}
	var body updateBody
	e := json.NewDecoder(r.Body).Decode(&body)
	if e != nil {
		fmt.Print(e)
	}
	filter := bson.D{{Key: "name", Value: body.Name}} // converting value to BSON
	after := options.After                            // for returning updated document
	returnOpt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "city", Value: body.City}}}}
	updateResult := userCollection.FindOneAndUpdate(context.TODO(), filter, update, &returnOpt)
	var result primitive.M
	_ = updateResult.Decode(&result)
	json.NewEncoder(w).Encode(result)
}

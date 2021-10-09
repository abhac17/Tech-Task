package main

import (
	"context"
	"time"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"encoding/json"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
	"github.com/dgrijalva/jwt-go"
)


var SECRET_KEY = []byte("gosecretkey")

type User struct{
	Id string `json:"id" bson:"id"`
	UserId string `json:"userId" bson:"userId"`
	Name string `json:"name" bson:"name"`
	Email string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}



type Post struct{
	Id string `json:"id" bson:"id"`
	Caption string `json:"caption" bson:"caption"`
	ImageUrl string `json:"imageUrl" bson:"imageUrl"`
	TimeStamp string `json:"timeStamp" bson:"timeStamp"`
}

var client *mongo.Client

func getHash(pwd []byte) string {
    hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
    if err != nil {
        log.Println(err)
    }
    return string(hash)
}

func GenerateJWT()(string,error){
	token:= jwt.New(jwt.SigningMethodHS256)
	tokenString, err :=  token.SignedString(SECRET_KEY)
	if err !=nil{
		log.Println("Error in JWT token generation")
		return "",err
	}
	return tokenString, nil
}

func userSignup(response http.ResponseWriter, request *http.Request){
	response.Header().Set("Content-Type","application/json")
	var user User
	json.NewDecoder(request.Body).Decode(&user)
	user.Password = getHash([]byte(user.Password))
	collection := client.Database("GODB").Collection("user")
	ctx,_ := context.WithTimeout(context.Background(), 10*time.Second)
	result,_ := collection.InsertOne(ctx,user)
	json.NewEncoder(response).Encode(result)
}


func postCreation(response http.ResponseWriter, request *http.Request){
	response.Header().Set("Content-Type","application/json")
	var post Post
	json.NewDecoder(request.Body).Decode(&post)
	currentTime:= time.Now()
	post.TimeStamp = currentTime.String()
	collection := client.Database("GODB").Collection("post")
	ctx,_ := context.WithTimeout(context.Background(), 10*time.Second)
	result,_ := collection.InsertOne(ctx,post)
	json.NewEncoder(response).Encode(result)
}


func getUser(response http.ResponseWriter, request *http.Request){
	response.Header().Set("Content-Type","application/json")
	id := mux.Vars(r)["id"]
	var user User
	
	result:= client.Database("GODB").Collection("user").Find(context.Background(), bson.M{"id": id})

	json.NewEncoder(response).Encode(result)
}


func getPost(response http.ResponseWriter, request *http.Request){
	response.Header().Set("Content-Type","application/json")
	id := mux.Vars(r)["id"]
	var user User
	
	result:= client.Database("GODB").Collection("POST").Find(context.Background(), bson.M{"id": id})

	json.NewEncoder(response).Encode(result)
}



func getPostUser(response http.ResponseWriter, request *http.Request){
	response.Header().Set("Content-Type","application/json")
	id := mux.Vars(r)["id"]
	var user User
	
	result:= client.Database("GODB").Collection("user").Find(context.Background(), bson.M{"userId": id})

	json.NewEncoder(response).Encode(result)
}


func main(){
	log.Println("Starting the application")

	router:= mux.NewRouter()
	ctx,_ := context.WithTimeout(context.Background(), 10*time.Second)
	client,_= mongo.Connect(ctx,options.Client().ApplyURI("mongodb://localhost:27017"))

	router.HandleFunc("/users",userSignup).Methods("POST")
	router.HandleFunc("/post",postCreation).Methods("POST")


	router.HandleFunc("/users/{id}",getUser).Methods("GET")
	router.HandleFunc("/posts/{id}",postCreation).Methods("GET")
	
	router.HandleFunc("/posts/users/{id}",getPostUser).Methods("GET")
	
	
	
	log.Fatal(http.ListenAndServe(":8080", router))

}
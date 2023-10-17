package beurse

import (
	"fmt"
	"testing"

	"github.com/aiteung/atdb"
	"github.com/whatsauth/watoken"
	"go.mongodb.org/mongo-driver/bson"
)

func TestGeneratePasswordHash(t *testing.T) {
	password := "secret"
	hash, _ := HashPassword(password) // ignore error for the sake of simplicity

	fmt.Println("Password:", password)
	fmt.Println("Hash:    ", hash)

	match := CheckPasswordHash(password, hash)
	fmt.Println("Match:   ", match)
}
func TestGeneratePrivateKeyPaseto(t *testing.T) {
	privateKey, publicKey := watoken.GenerateKey()
	fmt.Println("Ini Private", privateKey)
	fmt.Println("Ini Public", publicKey)
	hasil, err := watoken.Encode("urse", privateKey)
	fmt.Println("Ini Hasil", hasil, err)
}

func TestHashFunction(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "urse")
	var userdata User
	userdata.Username = "dito"
	userdata.Password = "secret"

	filter := bson.M{"username": userdata.Username}
	res := atdb.GetOneDoc[User](mconn, "user", filter)
	fmt.Println("Mongo User Result: ", res)
	hash, _ := HashPassword(userdata.Password)
	fmt.Println("Hash Password : ", hash)
	match := CheckPasswordHash(userdata.Password, res.Password)
	fmt.Println("Match:   ", match)

}

func TestIsPasswordValid(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "urse")
	var userdata User
	userdata.Username = "dani"
	userdata.Password = "secretoo"

	anu := IsPasswordValid(mconn, "user", userdata)
	fmt.Println(anu)
}

// func TestInsertUser(t *testing.T){
// 	mconn := SetConnection("MONGOSTRING", "urse")
// 	var userdata User
// 	userdata.Username = "fatwa"
// 	userdata.Password = "secretcuy"
// 	userdata.Role = "admin"

// 	nama := InsertUser(mconn, "user", userdata)
// 	fmt.Println(nama)
// }

// func TestInsertUser(t *testing.T){
// 	var userdata User
// 	userdata.Username = "ade"
// 	userdata.Password = "secretoo"

// 	nama := InsertUser("MONGOSTRING", "urse", "user", userdata)
// 	fmt.Println(nama)
// }
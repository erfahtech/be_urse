package beurse

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/aiteung/atdb"
	"github.com/whatsauth/watoken"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GCFReturnStruct(DataStuct any) string {
	jsondata, _ := json.Marshal(DataStuct)
	return string(jsondata)
}

func GCFPostHandler(PASETOPRIVATEKEYENV, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	var Response Credential
	Response.Status = false
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	var datauser User
	err := json.NewDecoder(r.Body).Decode(&datauser)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
	} else {
		if IsPasswordValid(mconn, collectionname, datauser) {
			Response.Status = true
			tokenstring, err := watoken.Encode(datauser.Email, os.Getenv(PASETOPRIVATEKEYENV))
			if err != nil {
				Response.Message = "Gagal Encode Token : " + err.Error()
			} else {
				Response.Message = "Selamat Datang " + datauser.Username
				Response.Token = tokenstring
			}
		} else {
			Response.Message = "Email atau Password Salah"
		}
	}

	return GCFReturnStruct(Response)
}


func InsertUser(r *http.Request) string {
	var Response Credential
	var userdata User
	err := json.NewDecoder(r.Body).Decode(&userdata)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
		return GCFReturnStruct(Response)
	}
	hash, _ := HashPassword(userdata.Password)
	userdata.Password = hash
	atdb.InsertOneDoc(SetConnection("MONGOSTRING", "db_urse"), "user", userdata)
	Response.Status = true
	Response.Message = "Akun berhasil dibuat untuk username: " + userdata.Username
	return GCFReturnStruct(Response)
}

func InsertDevice(r *http.Request) string {
	var Response Credential
	var devicedata Device
	var datauser User
	err := json.NewDecoder(r.Body).Decode(&devicedata)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
		return GCFReturnStruct(Response)
	}

	tokenString, err := watoken.Decode(datauser.Email, os.Getenv("PASETOPRIVATEKEYENV"))
	    if err != nil {
        Response.Message = "Error decoding token: " + err.Error()
        return GCFReturnStruct(Response)
    }

	devicedata.Email = tokenString.Id
	devicedata.ID = primitive.NewObjectID() // generate new ObjectID for device
	mconn := SetConnection("MONGOSTRING", "db_urse")
	atdb.InsertOneDoc(mconn, "devices", devicedata)
	Response.Status = true
	Response.Message = "Device berhasil ditambahkan dengan nama: " + devicedata.Name
	return GCFReturnStruct(Response)
}



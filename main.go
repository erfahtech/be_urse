package beurse

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/aiteung/atdb"
	"github.com/whatsauth/watoken"
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
	atdb.InsertOneDoc(SetConnection("MONGOSTRING", "urse"), "db_user", userdata)
	Response.Status = true
	Response.Message = "Akun berhasil dibuat untuk username: " + userdata.Username
	return GCFReturnStruct(Response)
}

func signUpUser(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) (string, error) {
    var Response Credential
    var userdata User
	err := json.NewDecoder(r.Body).Decode(&userdata)
	if err != nil {
		Response.Message = "error parsing application/json: "
		return GCFReturnStruct(Response), err
	}
	hash, err := HashPassword(userdata.Password)
	if err != nil {
		Response.Message = "error hashing password: " 
		return GCFReturnStruct(Response), err
	}
	userdata.Password = hash
	if err := atdb.InsertOneDoc(SetConnection(MONGOCONNSTRINGENV, dbname), collectionname, userdata); err != nil {
		Response.Message = "error inserting user data: "
		return GCFReturnStruct(Response), err.(error)
	}
		Response.Status = true
		Response.Message = "Akun berhasil dibuat untuk username: " + userdata.Username
		return GCFReturnStruct(Response), nil
}
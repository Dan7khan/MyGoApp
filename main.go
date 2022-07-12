package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
	"os"

	"github.com/boltdb/bolt"
	"github.com/gorilla/mux"
)

type Promotions struct {
	Id              string `json:"Id"`
	Price           string `json:"Price"`
	Expiration_Date string `json:"Expiration_Date"`
}
const b = "MyBucket1"

func main() {
	//port := os.Getenv("PORT")
	port := os.Getenv("PORT")
	fileServer := http.FileServer(http.Dir("./static"))
	myRouter := mux.NewRouter().StrictSlash(true)
	// replace http.HandleFunc with myRouter.HandleFunc
	myRouter.Handle("/",fileServer)
	myRouter.HandleFunc("/upload", uploadfunc)
	myRouter.HandleFunc("/promotion/{id}", retrieveValue)
	myRouter.HandleFunc("/uploadFile", uploadFile)
	// finally, instead of passing in nil, we want
	// to pass in our newly created router as the second
	// argument
	fmt.Printf("starting the server at port:"+port)
	log.Fatal(http.ListenAndServe(":"+port, myRouter))
}

func uploadfunc(w http.ResponseWriter, r *http.Request){
	fmt.Printf("entering upload func")
	if r.URL.Path !="/upload"{
		http.Error(w,"404 not found", http.StatusNotFound)
		return
	}
	p := "./static/upload.html"
	http.ServeFile(w, r, p)
	
}

func uploadFile(w http.ResponseWriter, r *http.Request){
	fmt.Println("File Upload Endpoint Hit")

    // Parse our multipart form, 10 << 20 specifies a maximum
    // upload of 10 MB files.
    r.ParseMultipartForm(10 << 20)
    // FormFile returns the first file for the given key `myFile`
    // it also returns the FileHeader so we can get the Filename,
    // the Header and the size of the file
    file, handler, err := r.FormFile("myFile")
    if err != nil {
        fmt.Println("Error Retrieving the File")
        fmt.Println(err)
        return
    }
    defer file.Close()
	if !strings.Contains(handler.Filename,".csv"){
		fmt.Fprintf(w,"Please upload a csv file")
		return
	}
    fmt.Printf("Uploaded File: %+v\n", handler.Filename)
    fmt.Printf("File Size: %+v\n", handler.Size)
    fmt.Printf("MIME Header: %+v\n", handler.Header)

    // // Create a temporary file within our temp-images directory that follows
    // // a particular naming pattern
    // tempFile, err := ioutil.TempFile("temp-images", "upload-*.png")
    // if err != nil {
    //     fmt.Println(err)
    // }
    // defer tempFile.Close()

    // read all of the contents of our uploaded file into a
    // byte array
    db := open("my.db")
	defer db.Close()
	// file, error := os.Open(file)
	// if error != nil {
	// 	log.Fatal((error))
	// 	fmt.Println("Please make sure to keep the file in the proper root folder")
	// }
	// defer file.Close()
	var count = 1
	reader := csv.NewReader(file)
	for {
		rec, error := reader.Read()
		if error == io.EOF {
			break
		}
		if error != nil {
			log.Fatal(error)
		}

		count = count + 1
		set(db, b, rec[0], rec[1]+","+rec[2])
	}
    // return that we have successfully uploaded our file!
    fmt.Fprintf(w, "Successfully Uploaded File\n")
}

func retrieveValue(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]
	if len(key) == 0 {
		http.NotFound(w, r)
	}

	db := open("my.db")
	defer db.Close()
	v := get(db, b, key)
	strArry := strings.Split(v, ",")
	var rec Promotions
	rec.Id = key
	rec.Price = strArry[0]
	rec.Expiration_Date = strArry[1]
	json.NewEncoder(w).Encode(rec)

}

// func InitializeFileInfoStorage(file FormFile) {
	
// }



func open(file string) *bolt.DB {
	db, err := bolt.Open(file, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		//handle error
		log.Fatal(err)
	}
	return db
}

func set(db *bolt.DB, bucket, key, value string) {
	db.Update(func(tx *bolt.Tx) error {
		b, _ := tx.CreateBucketIfNotExists([]byte(bucket))
		err := b.Put([]byte(key), []byte(value))
		return err
	})
}

func get(db *bolt.DB, bucket, key string) string {
	s := ""
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		s = string(b.Get([]byte(key)))
		return nil
	})
	return s
}

// func homePage(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "Welcome to the HomePage!")
// 	fmt.Println("Endpoint Hit: homePage")
// }


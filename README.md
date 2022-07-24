To run the application locally, comment the code on line 27 which is for obtaining the port in deployed environment. 
Replace the code at line 36 and 37 in main.go file with the following. 

'fmt.Printf("starting the server at port:10000")
log.Fatal(http.ListenAndServe(":10000", myRouter))'

build the application in VS code using the command 'go build'
Run the application using the command 'go run main.go'
To browse the application, open webbrowser and type 'localhost:10000'

'localhost:10000/' -> navigates to default index.html page
'localhost:10000/upload' -> navigates to the upload page, where you have to upload the promotions.csv file and click upload. 
'localhost:10000/promotions/{id}' -> replace {id} with the id of the entry you would like to retireve. 
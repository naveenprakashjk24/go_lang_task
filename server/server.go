package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	taskdata "github.com/naveenprakashjk24/go_lang_task"

	"github.com/gorilla/mux"
)

type UserInput struct{}

func (handler *UserInput) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 2048666))
	if err != nil {
		log.Println(err)
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		log.Println(err)
		panic(err)
	}
	inputCollection := taskdata.ProcessJsonInput(body)
	if len(inputCollection) <= 0 {
		http.Error(w, "{\"error\": \"empty input list submitted.\"}", 500)
		return
	}

	inputMap := make(map[int]taskdata.InputData)
	for k, v := range inputCollection {
		inputMap[k] = v
	}
	//The number of jobs we have to execute is the count of elements in our input map
	numberOfJobs := len(inputMap)
	numberOfWorkers := 100

	//Buffered chanels
	jobs := make(chan taskdata.InputData, 1000)     //Chanel to send in the jobs
	results := make(chan taskdata.OutputData, 1000) //Chanel to receive job results

	for w := 1; w <= numberOfWorkers; w++ {
		go taskdata.NewWorker(jobs, results)
	}

	for j := 0; j <= numberOfJobs-1; j++ {
		job := inputMap[j]
		jobs <- job
	}
	//close the jobs channel
	close(jobs)

	//The res map will be used to collect results of our workers...
	res := make(map[int]taskdata.OutputData)

	//Collect worker results
	for a := 0; a <= numberOfJobs-1; a++ {
		r := <-results
		res[r.Index] = r
	}
	close(results)

	//Set response header information
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)

	//Create json with all results...
	response, err := taskdata.GenerateJsonOutput(res)
	if err != nil {
		log.Println(err)
		panic(err)
	}

	fmt.Fprintf(w, string(response))

}

func main() {
	//The port our server will be litening on
	port := 8080

	//Our fancy Gorilla mux router
	router := mux.NewRouter()

	//Our Routes :
	//We're only interested in POST requests on /calcs URL
	router.Handle("/userdetails", &UserInput{}).Methods("POST")

	log.Printf("Starting server. Listening on port %d", port)
	err := http.ListenAndServe(":"+strconv.Itoa(port), router)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

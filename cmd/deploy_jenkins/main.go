package main

import (
	"GolangProj/pkg/jenkins"
	"fmt"
)

func main() {
	client, err := jenkins.LoadConfig("configs/jenkins_credentials.yaml")
	if err != nil {
		fmt.Printf("Failed to load config: %v\\n", err)
		return
	}

	fetchResponse, err := client.FetchAllJobs()
	if err != nil {
		fmt.Println("Err fetching jobs: ", err)
		return
	}
	// Fetch data from the pointer (Go can do itself) -DEREFERENCE
	for _, job := range fetchResponse.Jobs {
		fmt.Println("Job:", job.Name)
	}

	//test, err := client.CheckExistJobs("test.txt")
}

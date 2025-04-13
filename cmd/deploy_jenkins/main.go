package main

import (
	"GolangProj/pkg/jenkins"
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	// Init Client struct
	client, err := jenkins.LoadConfig("../../configs/jenkins_credentials.yaml")
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		return
	}
	// END Init Client struct

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\n===== Jenkins Management Menu =====")
		fmt.Println("1. Fetch all jobs")
		fmt.Println("2. Check if a job(s) exists")
		fmt.Println("3. Deploy job(s)")
		fmt.Println("4. Exit")
		fmt.Print("Select an option: ")

		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			// Fetch all jobs
			fetchResponse, err := client.FetchAllJobs()
			if err != nil {
				fmt.Println("Err fetching jobs: ", err)
				continue
			}
			for _, job := range fetchResponse.Jobs {
				fmt.Println("Job:", job.Name)
			}
		case "2":
			// Check if a job exists
			fmt.Print("Enter job name or path to file: ")
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)

			checkedJobs, err := client.CheckExistJobs(input)
			if err != nil {
				fmt.Println("Error checking jobs: ", err)
				continue
			}
			for _, job := range checkedJobs {
				fmt.Printf("Job: %s, Exist: %t\n", job.Name, job.Exist)
			}
		case "3":
			fmt.Print("Enter job name or path to file: ")
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)

			checkedJobs, err := client.CheckExistJobs(input)
			if err != nil {
				fmt.Println("Error checking jobs: ", err)
				continue
			}

			for _, job := range checkedJobs {
				if job.Exist {
					err := client.DeployItem(job.Name)
					if err != nil {
						fmt.Printf("Error deploying job %s: %v\n", job.Name, err)
					} else {
						fmt.Printf("Job %s deployed successfully\n", job.Name)
					}
				}
			}
		case "4":
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid option. Please choose again.")
		}
	}
}

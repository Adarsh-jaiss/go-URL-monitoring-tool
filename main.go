package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"time"
	// "net/smtp"
)

// const URL = "https://adarshjaiswal.me"

// Checking the status of URL
func checkURL(URL string){
	response,err := http.Get(URL)
	CheckErr(err)
	

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		fmt.Printf("%v Status code : %d\n",URL,response.StatusCode)

	}else {
		fmt.Printf(" %v Status : OK\n", URL)
	}


}


func main (){
	fmt.Print("Welcome to Url-monitoring application\n")

	// opening the file containing the URLs
	
	file,err := os.Open("urls.txt")
	CheckErr(err)
	defer file.Close()

	// Reading the URls 
	urls := make([] string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan(){
		urls = append(urls,scanner.Text())
	}

	// checking for errs in reading the file
	if err := scanner.Err(); err!= nil{
		fmt.Println("Error reading file:",err)
		return
	}
	
	// Assinging the intervals after each cycle of URL checking

	interval := 10 *time.Second
	
	
	for {
		for _, URL := range urls {
			checkURL(URL)
            
        }
        time.Sleep(interval)
    }


}

// A function that checks error
func CheckErr(err error){

	if err!= nil {
		panic(err)
	}
}

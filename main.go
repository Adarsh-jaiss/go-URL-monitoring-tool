package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"time"
	"net/smtp"
)

const (
	// URL        = "https://adarshjaiswal.me"
	SMTPServer = "smtp.gmail.com" // Update with your SMTP server details
	SMTPPort   = 587               // Update with your SMTP server port
	FromEmail  = "test.my.softwar@gmail.com" // Update with your email address
	Password   = ""    // Update with your email password
	ToEmail    = "" // Update with the recipient's email address
)

func checkURL(URL string, ch chan<- bool) {
	response, err := http.Get(URL)
	if err != nil {
		fmt.Printf("Error checking %v: %v\n", URL, err)
		ch <- false
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		fmt.Printf("%v Status code : %d\n", URL, response.StatusCode)
		ch <- false
		return
	}

	fmt.Printf(" %v Status : OK\n", URL)
	ch <- true
}

func sendEmail(subject, body string) error {
	auth := smtp.PlainAuth("", FromEmail, Password, SMTPServer)
	msg := fmt.Sprintf("Subject: %s\r\n\r\n%s", subject, body)

	err := smtp.SendMail(fmt.Sprintf("%s:%d", SMTPServer, SMTPPort), auth, FromEmail, []string{ToEmail}, []byte(msg))
	if err != nil {
		return err
	}
	return nil
}

func main() {
	fmt.Print("Welcome to URL-monitoring application\n")

	// Opening the file containing the URLs
	file, err := os.Open("urls.txt")
	CheckErr(err)
	defer file.Close()

	// Reading the URLs
	urls := make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		urls = append(urls, scanner.Text())
	}

	// Checking for errors in reading the file
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Assigning the intervals after each cycle of URL checking
	interval := 10 * time.Second

	for {
		ch := make(chan bool, len(urls))
		for _, URL := range urls {
			go checkURL(URL, ch)
		}

		allOK := true
		for i := 0; i < len(urls); i++ {
			if !<-ch {
				allOK = false
			}
		}

		if !allOK {
			// Send an email notification when at least one URL is down
			subject := "URL Monitor Notification"
			body := "One or more URLs are down."
			err := sendEmail(subject, body)
			if err != nil {
				fmt.Printf("Error sending email notification: %v\n", err)
			}
		}

		time.Sleep(interval)
	}
}

// A function that checks error
func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

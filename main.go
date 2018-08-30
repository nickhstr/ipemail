package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/smtp"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/joho/godotenv"
	"github.com/jordan-wright/email"
)

var lastIPDir string
var lastIPFile string
var pathToIPFile string

func main() {
	envFile := os.Getenv("ENV_FILE")
	if envFile == "" {
		envFile = ".env"
	}

	if err := godotenv.Load(envFile); err != nil {
		fmt.Println("Failed to load environment variables")
		return
	}

	lastIPDir = os.Getenv("LAST_IP_DIR")
	lastIPFile = os.Getenv("LAST_IP_FILE")
	pathToIPFile = filepath.Join(lastIPDir, lastIPFile)

	ipAddress, err := getIPAddress()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Current IP Address: %s\n", ipAddress)

	isNew := isNewIPAddress(ipAddress)
	if !isNew {
		fmt.Println("IP Address has not changed")
		return
	}

	fmt.Println("IP Address has been updated")

	if err = setNewIPAddress(ipAddress); err != nil {
		fmt.Println(err)
		return
	}

	emailMessage := fmt.Sprintf("Updated address: %s", ipAddress)

	if err = sendEmail([]byte(emailMessage)); err != nil {
		fmt.Println(err)
		return
	}
}

func getIPAddress() ([]byte, error) {
	var ipAddress []byte

	requestURL := "http://checkip.dyndns.org"
	ipRegex, err := regexp.Compile(`\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}`)
	if err != nil {
		fmt.Println("Cannot compile regular expression")
		return ipAddress, err
	}

	resp, err := http.Get(requestURL)
	if err != nil {
		fmt.Println("Could not fetch request")
		return ipAddress, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Could not read body")
		return ipAddress, err
	}

	ipAddress = ipRegex.Find(body)

	return ipAddress, nil
}

func isNewIPAddress(ipAddress []byte) bool {
	rawLastAddress, err := ioutil.ReadFile(pathToIPFile)
	if err != nil {
		fmt.Printf("Failed to read '%s', creating new file\n", pathToIPFile)

		err = os.MkdirAll(lastIPDir, os.ModePerm)
		if err != nil {
			fmt.Printf("Failed to create dir: %s\n", lastIPDir)
		}

		err = ioutil.WriteFile(pathToIPFile, ipAddress, 0644)
		if err != nil {
			fmt.Printf("Failed to create %s\n", lastIPFile)
		}

		return true
	}

	lastAddress := strings.TrimSpace(string(rawLastAddress))

	return lastAddress != string(ipAddress)
}

func setNewIPAddress(ipAddress []byte) error {
	err := ioutil.WriteFile(pathToIPFile, ipAddress, 0644)
	if err != nil {
		fmt.Println("Unable to write IP address")
		return err
	}

	return nil
}

func sendEmail(message []byte) error {
	fromAddress := os.Getenv("EMAIL_FROM_ADDRESS")
	fromUser := os.Getenv("EMAIL_FROM_USER")
	fromPassword := os.Getenv("EMAIL_FROM_PASSWORD")
	toAddress := os.Getenv("EMAIL_TO_ADDRESS")

	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", fromUser, fromAddress)
	e.To = []string{toAddress}
	e.Subject = "New IP Address for Raspberry Pi"
	e.Text = message

	err := e.Send("smtp.gmail.com:587", smtp.PlainAuth("", fromAddress, fromPassword, "smtp.gmail.com"))
	if err != nil {
		fmt.Println("Failed to send email")
		return err
	}

	return err
}

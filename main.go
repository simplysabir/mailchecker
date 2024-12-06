package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/smtp"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("domain, hasMX, hasSPF, spfRecords, hasDMARC, dmarcRecords, emailValid\n")

	for scanner.Scan() {
		email := scanner.Text()
		checkEmail(email)
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error: could not read from input: %v\n", err)
	}
}

func checkEmail(email string) {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		log.Printf("Invalid email format: %s\n", email)
		return
	}
	domain := parts[1]

	var hasMX, hasSPF, hasDMARC, emailValid bool
	var spfRecords, dmarcRecords string

	// Check MX records
	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		log.Printf("Error: could not find MX records for %s\n", domain)
	}
	if len(mxRecords) > 0 {
		hasMX = true
	}

	// Check SPF records
	txtRecords, err := net.LookupTXT(domain)
	if err != nil {
		log.Printf("Error: could not find TXT records for %s\n", domain)
	}
	for _, record := range txtRecords {
		if strings.HasPrefix(record, "v=spf1") {
			hasSPF = true
			spfRecords += record
			break
		}
		if strings.HasPrefix(record, "v=DMARC1") {
			hasDMARC = true
			dmarcRecords += record
		}
	}

	// Validate email address via SMTP
	if hasMX {
		emailValid = validateEmailSMTP(email, mxRecords)
	}

	fmt.Printf("%s, %t, %t, %s, %t, %s, %t\n",
		email, hasMX, hasSPF, spfRecords, hasDMARC, dmarcRecords, emailValid)
}

func validateEmailSMTP(email string, mxRecords []*net.MX) bool {
	// Try connecting to the first available MX server
	mxHost := mxRecords[0].Host
	client, err := smtp.Dial(mxHost + ":25")
	if err != nil {
		log.Printf("Error: could not connect to SMTP server %s: %v\n", mxHost, err)
		return false
	}
	defer client.Close()

	// Introduce ourselves to the SMTP server
	if err := client.Hello("localhost"); err != nil {
		log.Printf("Error: could not say HELO to SMTP server %s: %v\n", mxHost, err)
		return false
	}

	// Set the sender and recipient
	if err := client.Mail("test@example.com"); err != nil {
		log.Printf("Error: could not set sender: %v\n", err)
		return false
	}
	if err := client.Rcpt(email); err != nil {
		log.Printf("Error: email address %s is invalid: %v\n", email, err)
		return false
	}

	// Close the connection
	return true
}

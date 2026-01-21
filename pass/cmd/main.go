package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/polishedfeedback/go-utils/pass/internal/storage"
)

func printUsage() {
	fmt.Println("Password Manager CLI")
	fmt.Println("\nUsage:")
	fmt.Println("  pass add <name> --url <url> --email <email> --password <password>")
	fmt.Println("  pass get <name>")
	fmt.Println("  pass list")
	fmt.Println("  pass delete <name>")
	fmt.Println("\nExamples:")
	fmt.Println("  pass add google --url https://google.com --email john@gmail.com --password xyz123")
	fmt.Println("  pass get google")
	fmt.Println("  pass list")
	fmt.Println("  pass delete google")
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	store, err := storage.NewSQLiteStorage("passwords.db")
	if err != nil {
		log.Fatalf("error connecting to the database: %v", err)
	}

	command := os.Args[1]

	switch command {
	case "add":
		addCmd := flag.NewFlagSet("add", flag.ExitOnError)
		url := addCmd.String("url", "", "url to add")
		email := addCmd.String("email", "", "email to add")
		password := addCmd.String("password", "", "password to add")

		addCmd.Parse(os.Args[2:])

		if addCmd.NArg() < 1 {
			fmt.Println("error: url, email, password are required")
			printUsage()
			return
		}
		name := addCmd.Arg(0)
		err = store.AddPassword(name, *url, *email, *password)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		fmt.Printf("✓ Added password for %s\n", name)

	case "get":
		getCmd := flag.NewFlagSet("get", flag.ExitOnError)
		getCmd.Parse(os.Args[2:])

		if getCmd.NArg() < 1 {
			fmt.Println("name is required")
			printUsage()
			return
		}
		name := getCmd.Arg(0)
		cred, err := store.GetPassword(name)
		if err != nil {
			log.Fatalf("error getting the password: %v", err)
		}
		fmt.Printf("\nName: %s\n", cred.Name)
		fmt.Printf("URL: %s\n", cred.URL)
		fmt.Printf("Email: %s\n", cred.Email)
		fmt.Printf("Password: %s\n", cred.Password)

	case "list":
		creds, err := store.ListPasswords()
		if err != nil {
			log.Fatalf("error getting the password list: %v", err)
		}
		if len(creds) == 0 {
			fmt.Println("No passwords saved")
			return
		}
		fmt.Println("\nSaved Passwords:")
		for _, cred := range creds {
			fmt.Printf("  %s - %s (%s)\n", cred.Name, cred.Email, cred.URL)
		}

	case "delete":
		deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
		deleteCmd.Parse(os.Args[2:])

		if deleteCmd.NArg() < 1 {
			fmt.Println("name is required")
			printUsage()
			return
		}

		name := deleteCmd.Arg(0)
		err := store.DeletePassword(name)
		if err != nil {
			log.Fatalf("error deleting the password: %v", err)
		}
		fmt.Printf("✓ Deleted password for %s\n", name)
	}
}

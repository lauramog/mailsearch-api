package main

import (
	"bufio"
	"context"
	"github.com/joho/godotenv"
	client "github.com/zinclabs/sdk-go-zincsearch"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	mailDirPath := os.Getenv("MAIL_DIR_PATH")
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(err)
	}

	dirEntries, err := os.ReadDir(mailDirPath)
	if err != nil {
		log.Fatal("cannot open maildirectory", err)
	}
	log.Print("start reading inbox ")
	var allEmails [][]map[string]interface{}
	for _, userInbox := range dirEntries {
		inboxEntries, err := os.ReadDir(filepath.Join(mailDirPath, userInbox.Name(), "inbox"))
		if os.IsNotExist(err) {
			log.Printf("no inbox for user %s", userInbox.Name())
			continue
		}

		var emails []map[string]interface{}
		for _, inboxFile := range inboxEntries {
			email, err := extractEmail(mailDirPath, userInbox, inboxFile)
			if err != nil {
				log.Print(err)
			}
			emails = append(emails, email)
		}

		allEmails = append(allEmails, emails)
	}

	ctx := context.WithValue(context.Background(), client.ContextBasicAuth, client.BasicAuth{
		UserName: os.Getenv("UserName"),
		Password: os.Getenv("Password"),
	})

	configuration := client.NewConfiguration()
	configuration.Servers = client.ServerConfigurations{
		client.ServerConfiguration{
			URL: "http://localhost:4080",
		},
	}
	apiClient := client.NewAPIClient(configuration)
	log.Print("starting to index")
	for _, emails := range allEmails {
		query := client.NewMetaJSONIngest()
		query.SetIndex("inbox")
		query.SetRecords(emails)
		_, _, err := apiClient.Document.Bulkv2(ctx).Query(*query).Execute()
		if err != nil {
			log.Fatal(err)
		}
	}
	log.Printf("index done for emails of %d inboxes", len(allEmails))
}

func extractValue(line, word string) (string, bool) {
	_, after, found := strings.Cut(line, word)
	return after, found

}
func extractEmail(mailDirPath string, userInbox os.DirEntry, inboxFile os.DirEntry) (map[string]interface{}, error) {
	file, err := os.Open(filepath.Join(mailDirPath, userInbox.Name(), "inbox", inboxFile.Name()))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	email := make(map[string]interface{})
	email["Username"] = userInbox.Name()

	scan := bufio.NewScanner(file)
	for scan.Scan() {
		line := scan.Text()

		from, foundFrom := extractValue(line, "From:")
		if foundFrom {
			email["From"] = from
			continue
		}
		to, foundTo := extractValue(line, "To:")
		if foundTo {
			email["to"] = to
			continue
		}
		subject, foundSub := extractValue(line, "Subject:")
		if foundSub {
			email["subject"] = subject
			break
		}
	}
	return email, nil
}

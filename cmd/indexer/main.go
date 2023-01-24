package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/joho/godotenv"
	client "github.com/zinclabs/sdk-go-zincsearch"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	mailDirPath := os.Getenv("MAIL_DIR_PATH")
	envErr := godotenv.Load(".env")
	if envErr != nil {
		fmt.Printf("could not load .env file ")
		os.Exit(1)
	}

	dirEntries, err := os.ReadDir(mailDirPath)

	if err != nil {
		log.Print("error:", err)
	}
	var allInboxFilesAllUsers [][]map[string]interface{}

	for _, entryUser := range dirEntries {
		inboxUser, err := os.ReadDir(filepath.Join(mailDirPath, entryUser.Name(), "inbox"))
		if os.IsNotExist(err) {
			log.Print("user %s  , err:%s", entryUser.Name(), err)
			continue
		}

		var allInboxFilesUser []map[string]interface{}
		for _, inboxFiles := range inboxUser {
			files, err := os.Open(filepath.Join(mailDirPath, entryUser.Name(), "inbox", inboxFiles.Name()))
			if err != nil {
				log.Print(err)
			}
			inboxFileUser := make(map[string]interface{})
			scan := bufio.NewScanner(files)

			for scan.Scan() {
				line := scan.Text()

				from, foundFrom := extractValue(line, "From:")
				if foundFrom {
					inboxFileUser["From"] = from
					continue
				}
				to, foundTo := extractValue(line, "To:")
				if foundTo {
					inboxFileUser["to"] = to
					continue
				}
				subject, foundSub := extractValue(line, "Subject:")
				if foundSub {
					inboxFileUser["subject"] = subject
					break
				}
			}
			allInboxFilesUser = append(allInboxFilesUser, inboxFileUser)
		}
		allInboxFilesAllUsers = append(allInboxFilesAllUsers, allInboxFilesUser)
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
	for _, allEmailsUser := range allInboxFilesAllUsers {
		query := client.NewMetaJSONIngest()
		query.SetIndex("inbox")
		query.SetRecords(allEmailsUser)
		_, _, err := apiClient.Document.Bulkv2(ctx).Query(*query).Execute()
		if err != nil {
			log.Fatal(err)
		}

	}

}

func extractValue(line, word string) (string, bool) {
	_, after, found := strings.Cut(line, word)
	return after, found

}

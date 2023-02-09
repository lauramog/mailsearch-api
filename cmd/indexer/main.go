package main

import (
	"context"
	"github.com/joho/godotenv"
	client "github.com/zinclabs/sdk-go-zincsearch"
	"io"
	"log"
	"net/http"
	_ "net/http/pprof"
	"net/mail"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	mailDirPath := os.Getenv("MAIL_DIR_PATH")
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(err)
	}

	log.Print("start reading inbox ")

	allEmails, err := parseEmails(mailDirPath)
	if err != nil {
		log.Fatal(err)
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

func parseEmails(mailDirPath string) ([][]map[string]interface{}, error) {
	var allEmails [][]map[string]interface{}

	dirEntries, err := os.ReadDir(mailDirPath)
	if err != nil {
		return allEmails, err
	}

	for _, userInbox := range dirEntries {
		inboxEntries, err := os.ReadDir(filepath.Join(mailDirPath, userInbox.Name(), "inbox"))
		if os.IsNotExist(err) {
			log.Printf("no inbox for user %s", userInbox.Name())
			continue
		}

		var emails []map[string]interface{}
		for _, inboxFile := range inboxEntries {
			if inboxFile.IsDir() {
				continue
			}
			email, err := extractEmail(mailDirPath, userInbox, inboxFile)
			if err != nil {
				log.Printf("cannot parse email from inbox %s: %s", userInbox.Name(), err)
			}
			emails = append(emails, email)
		}

		allEmails = append(allEmails, emails)
	}

	return allEmails, nil
}

func extractEmail(mailDirPath string, userInbox os.DirEntry, inboxFile os.DirEntry) (map[string]interface{}, error) {
	file, err := os.ReadFile(filepath.Join(mailDirPath, userInbox.Name(), "inbox", inboxFile.Name()))
	if err != nil {
		return nil, err
	}
	email := make(map[string]interface{})
	email["Username"] = userInbox.Name()

	r := strings.NewReader(string(file))
	m, err := mail.ReadMessage(r)
	if err != nil {
		return nil, err
	}
	header := m.Header

	from := header.Get("From")
	email["From"] = from
	to := header.Get("To")
	email["To"] = to
	subject := header.Get("Subject")
	email["Subject"] = subject

	body, err := io.ReadAll(m.Body)
	if err != nil {
		return nil, err
	}
	email["Message"] = string(body)

	return email, nil
}

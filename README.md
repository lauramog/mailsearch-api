# Backend: email searching engine 

tool(indexer)  to process  the Enron Email Dataset (download it [here](http://www.cs.cmu.edu/~enron/enron_mail_20110402.tgz)),
indexing it in [zincsearch](https://zincsearch.com/), this tools parses the emails in the inbox directory, extracting the fields: to, from, 
subject and the body of the email. In the directory api you find a web server to expose the API.   

## :wrench: Installation 
Ensure the following is installed:

* [Go is installed locally](https://go.dev/doc/install) 
* [docker](https://www.docker.com/get-started/)

Then clone the repo locally: `https://github.com/lauramog/mailsearch-api.git`, you can find a User interface for this 
API [here](https://github.com/lauramog/mailsearch-ui)



## Run 

Connect with ZincSearch

```shell
docker run -p 4080:4080 -e ZINC_FIRST_ADMIN_USER=admin -e ZINC_FIRST_ADMIN_PASSWORD=Complexpass#123 --name zinc public.ecr.aws/zinclabs/zinc:latest
```

Index data in ZincSearch

```shell
├── cmd
  ├──indexer
MAIL_DIR_PATH="/path/to/root/folder/of/data" go run main.go
```

*search in ZincSearch UI http://localhost:4080/* 


Start the server 

```shell
├── cmd
  ├──api
  go run main.go
```

request to the server 

```shell
curl "http://localhost:port/search?term=enter&a&word/"
```

use the user interface provided [here](https://github.com/lauramog/mailsearch-ui)
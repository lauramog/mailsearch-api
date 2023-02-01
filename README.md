# Backend: email searching engine 

tool to process the Enron Email Dataset (download it [here](http://www.cs.cmu.edu/~enron/enron_mail_20110402.tgz)) 
and index it in [zincsearch](https://zincsearch.com/).

## Install 
Ensure the following is installed:

* [Go is installed locally](https://go.dev/doc/install) 
* [docker](https://www.docker.com/get-started/)

Then clone the repo locally: `https://github.com/lauramog/mailsearch-api.git`

## Test

```shell
go test ./...
```

## Run 

Connect with ZincSearch

```shell
docker run -p 4080:4080 -e ZINC_FIRST_ADMIN_USER=admin -e ZINC_FIRST_ADMIN_PASSWORD=Complexpass#123 --name zinc public.ecr.aws/zinclabs/zinc:latest
```

index data in ZincSearch

```shell
├── cmd
  ├──indexer
MAIL_DIR_PATH="/path/to/root/folder/of/data" go run main.go
```
start the server 

```shell
├── cmd
  ├──api
  go run main.go
```

*search in ZincSearch UI*
http://localhost:4080/

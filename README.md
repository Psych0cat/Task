## Quickstart:

``` 
git clone https://github.com/Psych0cat/Task.git
cd Task
docker-compose up --build
 ```
## Request example: 
``` 
curl --location --request POST '0.0.0.0:8080/transaction' \
--header 'Source-Type: server' \
--header 'Host: 127.0.0.1' \
--header 'Content-Type: application/json' \
--data-raw '{"state": "win", "amount": "10.15", "transactionId": "1k1"}'
``` 
## EVN location:
All defaults are set to run the app and DB from the box.

`.env.example`

## Logs:
 Log files are inside the app container, see naming in .env.example
## Tests:
 Current tests are end-to-end, without DB mocking, so app need to be up and running.

 `go test Router/main_test.go`



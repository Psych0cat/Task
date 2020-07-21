## Quickstart:

``` 
git clone https://github.com/Psych0cat/Task.git
docker-compose up --build
 ```
## EVN location:
All defaults are set to run the app and DB from the box.

`.env.example`

## Logs:
 Log files are inside the app container, see naming in .env.example
## Tests:
 Current tests are end-to-end, without DB mocking, so app need to be up and running.

 `go test Router/main_test.go`



# GO API

Following a udemy RESTful API course

## Installed packages

go get gopkg.in/mgo.v2/bson

go get github.com/asdine/storm

## Postman

Import the json file

Create a environment with:
host localhost:11111

## Benchmark

go test -bench .

## TODO

- Use gomod instead of go get
- Implement basic metrics using prometheus
- Show cache usage
  - Add a debug log so we see what's in the cache
  - Add prom metrics so we know how often the cache get's hit

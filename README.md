# GO API

Following a udemy RESTful API course

## Installed packages

To add extra packages

go get gopkg.in/mgo.v2/bson

go get github.com/asdine/storm

### mod

Initialize go mod, it will parese all req from your code.

go mod init

### vendor

Instead of grabbing the packages manually let's use go mod vendor.

go mod vendor

## Postman

Import the json file

Create a environment with:
host localhost:11111

## Benchmark

go test -bench .

## GO Documentation

Apparently there is something called godoc
https://blog.golang.org/godoc

Have some issues to install go-toolset on rhel8 with a self built go 1.13
I probably just have to build the bin my self but currently don't want to give that prio.

As a reminder to myself i created doc.go which apparently a way to write documentation about your package :).

## TODO

- [x] Use gomod instead of go get
- [x] Implement basic metrics using prometheus
- [ ] Show cache usage
  - [ ] Add a debug log so we see what's in the cache
  - [ ] Add prom metrics so we know how often the cache get's hit
- [ ] Only use echo
- [ ] Try to refactor code to look more like a echo project

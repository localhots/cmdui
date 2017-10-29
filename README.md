# Command UI

A web UI for go-api commands.

# Dependencies

* npm
* go
* dep

If you are using a Mac and have homebrew installed, you can run this:

```
brew install go # If you don't have Go installed
go get -u github.com/golang/dep/cmd/dep
brew install npm
```

# Installation

```
make install
```

Import `schema.sql` into a MySQL database.

# Starting

First session:

```
cd backend && go run main.go
```

Second session:

```
cd frontend && npm start
```

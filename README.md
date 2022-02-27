# Simple Golang Practice Exercises

Each directory is an independent go exercise. Each directory therefore represents a package that is a module, and hence its own binary executable.

`cd` into a top level directory, and run `go build` to create the executable.

### GO Minion CLI Tool

- `gominion` is a golang cli todo list tool.
- from inside the project's root dir, you can install into $GOPATH by running `go install .`.  Then the tool can be run with `$> gominion.....`
- alternatively it can be built in the project dir with `go build .`. Then from inside the directory, the tool can be invoked with `$> ./gominion...`
- it can be run without pre-building by simply doing `go run main.go <insert command>`

# Useful references

- [Awesome Go](https://github.com/avelino/awesome-go) - a github repo filled with curated go tools, frameworks, packages.
- [Rest API + Geth project](https://hackernoon.com/create-an-api-to-interact-with-ethereum-blockchain-using-golang-part-1-sqf3z7z)
- [Tutorial on Web3 + Geth](https://medium.com/@m.vanderwijden1/web3-go-part-1-31c68c68e20e)

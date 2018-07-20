# most-common-trigrams
Computes the 100 most common trigrams (sequences of three words)

## How to build and test:
After installing go in a standard way, these commans should work out of the box:
```
go get -d github.com/rchernobelskiy/most-common-trigrams
go build -o main github.com/rchernobelskiy/most-common-trigrams
go test github.com/rchernobelskiy/most-common-trigrams
```

## How to use locally:
Build (or use the bundled static linux binary):
```
go build -o main github.com/rchernobelskiy/most-common-trigrams
```
Send input on stdin:
```
echo hello foo bar baz one foo bar baz two | ./main
```
Send input from file(s):
```
./main LICENSE README.md
```

# boltdb-bench
BoltDB random read bench test 

Before run the benchmark test, please uncomment the code to generate data in `main.go`:

```Go
    // generate test db data
    genData(db)
```
That method will generate `600000 * 1000` records into the boltdb.
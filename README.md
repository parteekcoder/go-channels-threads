## Ticker

This project demonstrates:
-  creation of goroutines(threads) in golang 
-  the communication accross goroutines using channels
-  Listening and sending data to multiple channels from a goroutine
-  Writing data to a Redis container

This project is developed to learn goroutines and channels in Go. 

Take aways:
- In golang , creation of threads(goroutines) is very easy
- Go can create large number of threads efficiently
- Go manage threads in goroutines
- Communication between goroutines using channels is very easy

## Run this project

- Start a redis container using docker

  ```
    docker run --name my-redis-container -p 6379:6379 -d redis:latest
  ```

- Install the Go dependencies in the project
  
- Run the main file

```
go run main.go
```
package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {
	redisName := os.Getenv("REDIS_NAME")
	redisAddress := os.Getenv("REDIS_" + redisName + "_ADDRESS")
	redisPassword := os.Getenv("REDIS_" + redisName + "_PASSWORD")
	redisPort := os.Getenv("REDIS_" + redisName + "_PORT")
	redisDb := os.Getenv("REDIS_" + redisName + "_DB")

	if redisAddress == "" {
		redisAddress = "localhost"
	}

	if redisPort == "" {
		redisPort = "6379"
	}

	if redisDb == "" {
		redisDb = "0"
	}
	redisDbToUse, err := strconv.Atoi(redisDb)
	if err != nil {
		fmt.Println("Can convert port in integer '" + redisDb + "' :")
		fmt.Println(nil)
	}

	fmt.Println("Redis name: " + redisName)
	fmt.Println("Redis address: " + redisAddress)
	fmt.Println("Redis port: " + redisPort)
	if redisPassword == "" {
		fmt.Println("Redis without password")
	} else {
		fmt.Println("Redis with password")
	}
	fmt.Println("Redis DB: " + redisDb)

	client := redis.NewClient(&redis.Options{
		Addr:     redisAddress + ":" + redisPort,
		Password: redisPassword, // no password set
		DB:       redisDbToUse,  // use default DB
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)

	http.HandleFunc("/ping-redis", func(w http.ResponseWriter, r *http.Request) {
		client := redis.NewClient(&redis.Options{
			Addr:     redisAddress + ":" + redisPort,
			Password: redisPassword, // no password set
			DB:       redisDbToUse,  // use default DB
		})

		pong, err := client.Ping().Result()
		// fmt.Println(pong, err)
		fmt.Fprintf(w, "%s %s\n", pong, err)
		// Output: PONG <nil>

		fmt.Fprintf(w, "Redis name: %s\n", redisName)
		fmt.Fprintf(w, "Redis address: %s\n", redisAddress)
		fmt.Fprintf(w, "Redis port: %s\n", redisPort)
		if redisPassword == "" {
			fmt.Fprintf(w, "Redis without password\n")
		} else {
			fmt.Fprintf(w, "Redis with password\n")
		}
		fmt.Fprintf(w, "Redis DB: %s\n", redisDb)
	})

	log.Fatal(http.ListenAndServe(":8081", nil))

}

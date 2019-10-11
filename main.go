package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis"
)

const (
	KEY_COUNT = int(1e5)
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		msg := fmt.Sprintf(
			`Possible routes:
- /set   (will execute %d set operations)
- /get   (will execute %d get operations, you MUST call /set first)
- /stats (will show redis stats)
`, KEY_COUNT, KEY_COUNT)

		fmt.Fprintf(w, msg)

		w.Header().Set("Content-Type", "text/plain")
	})

	http.HandleFunc("/set", func(w http.ResponseWriter, r *http.Request) {
		client, err := newRedisClient()
		if err != nil {
			httpError(w, err)
			return
		}

		start := time.Now()

		for i := 0; i < KEY_COUNT; i++ {
			err := client.Set(redisKey(i), "1", 0).Err()

			if err != nil {
				httpError(w, err)
				return
			}
		}

		msg := fmt.Sprintf(
			"SET %d keys in %s",
			int(KEY_COUNT),
			time.Now().Sub(start),
		)

		w.Write([]byte(msg))
	})

	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		client, err := newRedisClient()
		if err != nil {
			httpError(w, err)
			return
		}

		start := time.Now()

		for i := 0; i < KEY_COUNT; i++ {
			err := client.Get(redisKey(i)).Err()

			if err != nil {
				httpError(w, err)
				return
			}
		}

		msg := fmt.Sprintf(
			"GET %d keys in %s",
			int(KEY_COUNT),
			time.Now().Sub(start),
		)

		w.Write([]byte(msg))
	})

	http.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		client, err := newRedisClient()
		if err != nil {
			httpError(w, err)
			return
		}

		buff := new(bytes.Buffer)
		for _, s := range []string{"stats", "keyspace"} {
			values, err := client.Info(s).Result()
			if err != nil {
				httpError(w, err)
				return
			}

			buff.WriteString(values)
			buff.WriteString("\n")
		}

		w.Write(buff.Bytes())
	})

	log.Fatal(
		http.ListenAndServe(":8081", nil),
	)
}

// ---

type connInfo struct {
	addr   string
	passwd string
	port   string
	db     int
}

func newConnInfo() (*connInfo, error) {
	instanceName := os.Getenv("REDIS_INSTANCE_NAME")
	if instanceName == "" {
		return nil, errors.New("No redis db instance name provided")
	}

	addr := getEnv(
		instanceName, "ADDRESS", "",
	)

	if addr == "" {
		return nil, errors.New("No redis server addr provided")
	}

	passwd := getEnv(
		instanceName, "PASSWORD", "",
	)

	if passwd == "" {
		return nil, errors.New("No redis password provided")
	}

	db, err := strconv.Atoi(
		getEnv(instanceName, "DB", "0"),
	)

	if err != nil {
		return nil, err
	}

	return &connInfo{
		addr:   addr,
		passwd: passwd,
		port:   getEnv(instanceName, "PORT", "6379"),
		db:     db,
	}, nil
}

func (ci *connInfo) dump() {
	fmt.Println("\n\n")
	fmt.Println("Conn info")
	fmt.Println("addr", ci.addr)
	fmt.Println("passwd", ci.passwd)
	fmt.Println("port", ci.passwd)
	fmt.Println("passwd", ci.passwd)
	fmt.Println("db", ci.db)
	fmt.Println("\n\n")
}

func (ci *connInfo) fullAddr() string {
	return fmt.Sprintf(
		"%s:%s", ci.addr, ci.port,
	)
}

// ---

func getEnv(instanceName, what, dfault string) string {
	v := os.Getenv(
		fmt.Sprintf(
			"REDIS_%s_%s",
			strings.ToUpper(instanceName),
			what,
		),
	)

	if len(v) == 0 {
		return dfault
	}

	return v
}

func fatal(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func newRedisClient() (*redis.Client, error) {
	ci, err := newConnInfo()
	if err != nil {
		return nil, err
	}

	//ci.dump()

	return redis.NewClient(
		&redis.Options{
			Addr:     ci.fullAddr(),
			Password: ci.passwd,
			DB:       ci.db,
			//TLSConfig: &tls.Config{
			//	InsecureSkipVerify: true,
			//},
		},
	), nil
}

func httpError(w http.ResponseWriter, e error) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(e.Error()))
}

func redisKey(i int) string {
	return fmt.Sprintf("k-%d", i)
}

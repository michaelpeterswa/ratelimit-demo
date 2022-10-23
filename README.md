# ratelimit-demo
performant rate limiter demo using the sliding window algorithim with a redis backing

## Running the demo
```bash
$ docker-compose up --build
$ curl --request POST \
  --url http://localhost:8080/ \
  --header 'Content-Type: application/json' \
  --data '{
	"id": "asdf123"
    }'
```


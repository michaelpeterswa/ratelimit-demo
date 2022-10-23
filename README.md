# ratelimit-demo
performant rate limiter demo using the sliding window algorithim with a redis backing

## running the demo
```bash
$ docker-compose up --build
$ curl --request POST \
  --url http://localhost:8080/ \
  --header 'Content-Type: application/json' \
  --data '{
	"id": "asdf123"
    }'
```

## notes
for maximum performance, consider using KeyDB instead of redis. KeyDB is a fork of redis that is optimized for speed. it is a drop in replacement for redis, so you can use the same commands and libraries.
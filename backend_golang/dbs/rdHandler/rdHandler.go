package rdHandler

import (
  "github.com/go-redis/redis"
)

type RedisdbWrite struct {}
type RedisdbRead struct {}

// Write:
//  redis-master is the the DNS name given to the backend Service
//  The DNS name is 'redis-master', which is the value of the name field in the preceding Service configuration file.
var redisdbWrite = redis.NewClient(&redis.Options{
  Addr: "redis-master:6379", // use default Addr
  Password: "",               // no password set
  DB: 0,                // use default DB
})
// Read
var redisdbRead = redis.NewClient(&redis.Options{
  Addr: "redis-slave:6379",
  Password: "",
  DB: 0, // use default DB
})

func (handler RedisdbWrite) Set(key *string, val *string) error {
  // Testing of Redis write/read
  err := redisdbWrite.Set(*key, *val, 0).Err()
  if err != nil {
    panic(err)
  }
  return err
}

func (handler RedisdbRead) Get(key *string) (string, error) {
  // Testing of Redis write/read
  val, err := redisdbRead.Get(*key).Result()
  if err != nil {
    panic(err)
  }
  return val, err
}
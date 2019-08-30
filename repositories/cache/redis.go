package cache

import (
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"github.com/crazyfacka/yanbapp-api/domain/users"
)

// Redis repository struct
type Redis struct {
	client *redis.Client
}

func (r *Redis) updateTTL(u *users.User) error {
	return r.client.Expire(u.SessionKey, users.SessionDuration).Err()
}

// LoadSession loads a Redis session into user domain
func (r *Redis) LoadSession(u *users.User) error {
	var (
		id   int
		data map[string]string
		err  error
	)

	if data, err = r.client.HGetAll(u.SessionKey).Result(); err != nil {
		return err
	}

	if id, err = strconv.Atoi(data["id"]); err != nil {
		return err
	}

	u.ID = uint(id)

	return r.updateTTL(u)
}

// CreateSession creates a new user session in the cache system
func (r *Redis) CreateSession(u *users.User) error {
	data := map[string]interface{}{
		"id":      u.ID,
		"created": time.Now().Unix(),
	}

	u.GenRandSessionString()

	if err := r.client.HMSet(u.SessionKey, data).Err(); err != nil {
		return err
	}

	log.Printf("saved session %s for user ID %d", u.SessionKey, u.ID)

	return r.updateTTL(u)
}

// Close closes the Redis connection
func (r *Redis) Close() error {
	return r.client.Close()
}

// NewRedis initializes a Redis repository
func NewRedis(host string, port int, database int) *Redis {
	client := redis.NewClient(&redis.Options{
		Addr: host + ":" + strconv.Itoa(port),
		DB:   database,
	})

	err := client.Ping().Err()
	if err != nil {
		log.Fatalln("error pinging redis")
	}

	log.Println("connected to Redis")

	return &Redis{
		client: client,
	}
}

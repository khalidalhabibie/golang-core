package redis

import (
	"fmt"
	"log"
	"time"

	redisClient "github.com/go-redis/redis"
	"github.com/go-redsync/redsync/v3"
	"github.com/go-redsync/redsync/v3/redis"
	"github.com/go-redsync/redsync/v3/redis/goredis"

	"gitlab.com/depatu/core/utils/errors"
)

const (
	RedisTopupPrefix              = "top_up:"
	RedisInvoicePrefix            = "invoice:"
	RedisBnibProductPrefix        = "bnib_product:"
	RedisBnibDirectBuyOrderPrefix = "direct_bnib_buy_order:"
	RedisBnibDirectProductPrefix  = "direct_bnib_product:"
	RedisBnibBuyOrderPrefix       = "bnib_buy_order:"
	RedisBnibProductPricePrefix   = "bnib_product_price:"
	RedisBnibBuyOrderPricePrefix  = "bnib_buy_order_price:"
	RedisBnibTransactionPrefix 	  = "bnib_transaction:"
	RedisLock                     = "lock:"
	RedisSellingHistoryPrefix	  = "selling_history:"
	// mercury
	RedisRetailPrefix        = "retail:"
	RedisStoreClosedCooldown = "store_closed_cooldown:"
	// covid
	RedisFolderNamePrefix        = "folder_name:"
	RedisAuthenticationPrefix    = "authentication:"
	RedisPackageSetPrefix        = "package_set:"
	RedisLegitCheckInvoicePrefix = "legit_check_invoice:"
	RedisPackageInvoicePrefix    = "package_invoice:"
	RedisPaymentGatewayPrefix    = "payment_gateway:"
)

type Credentials struct {
	Host     string
	Port     string
	Password string
}

type Client interface {
	Get(prefix string, key string) string
	Set(prefix string, key string, value string, expirationTime time.Duration) error
	Delete(prefix string, key string) error
	Ping() error
	Close() error
	NewMutex(key string) *redsync.Mutex
	SAdd(prefix string, key string, members interface{}) error
	SRem(prefix string, key string, members interface{}) error
	SMembers(prefix string, key string) ([]string, error)
	//CreateLock(key string) error
	//CheckLock(key string) (bool, error)
	//ReleaseLock(key string) error
}

type Redis struct {
	client  *redisClient.Client
	redsync *redsync.Redsync
}

func NewClient(credentials Credentials, appEnv string) Client {
	client := redisClient.NewClient(&redisClient.Options{
		Addr:     fmt.Sprintf("%s:%s", credentials.Host, credentials.Port),
		Password: credentials.Password,
		DB:       0,
	})
	status := client.Ping()
	if status.Err() != nil {
		if appEnv != "development" {
			log.Panic(status.Err())
		} else {
			log.Println("warning: redis not connected")
		}
	}

	pool := goredis.NewGoredisPool(client)
	rs := redsync.New([]redis.Pool{pool})

	return &Redis{
		client:  client,
		redsync: rs,
	}
}

func (r *Redis) Get(prefix string, key string) string {
	if r.client == nil {
		log.Println("warning: redis not connected")
		return ""
	}

	val, err := r.client.Get(fmt.Sprint(prefix, key)).Result()
	if err != nil {
		fmt.Println(err)
	}
	return val
}

func (r *Redis) Set(prefix string, key string, value string, expirationTime time.Duration) error {
	if r.client == nil {
		log.Println("warning: redis not connected")
		return errors.ErrUnprocessableEntity
	}

	err := r.client.Set(fmt.Sprint(prefix, key), value, expirationTime).Err()
	if err != nil {
		fmt.Println(err)
		return errors.ErrUnprocessableEntity
	}

	return nil
}

func (r *Redis) Delete(prefix string, key string) error {
	if r.client == nil {
		log.Println("warning: redis not connected")
		return errors.ErrUnprocessableEntity
	}

	err := r.client.Del(fmt.Sprint(prefix, key)).Err()
	if err != nil {
		fmt.Println(err)
	}
	return err
}

func (r *Redis) Ping() error {
	pong, err := r.client.Ping().Result()
	if err != nil {
		log.Println("error pinging redis:", err)
		return errors.ErrUnprocessableEntity

	}

	fmt.Println("connected to redis:", pong)
	return nil
}

func (r *Redis) Close() error {
	err := r.client.Close()
	if err != nil {
		log.Println("error closing redis:", err)
		return err
	}
	return nil
}

func (r *Redis) NewMutex(key string) *redsync.Mutex {
	return r.redsync.NewMutex(key)
}

// Add one or more members to a set
func (r *Redis) SAdd(prefix string, key string, members interface{}) error {
	if r.client == nil {
		log.Println("warning: redis not connected")
		return errors.ErrUnprocessableEntity
	}

	err := r.client.SAdd(fmt.Sprint(prefix, key), members).Err()
	if err != nil {
		fmt.Println(err)
	}
	return err
}

// Remove one or more members from a set
func (r *Redis) SRem(prefix string, key string, members interface{}) error {
	if r.client == nil {
		log.Println("warning: redis not connected")
		return errors.ErrUnprocessableEntity
	}

	err := r.client.SRem(fmt.Sprint(prefix, key), members).Err()
	if err != nil {
		fmt.Println(err)
		return errors.ErrUnprocessableEntity
	}
	return err
}

// Get all the members in a set
func (r *Redis) SMembers(prefix string, key string) ([]string, error) {
	if r.client == nil {
		log.Println("warning: redis not connected")
		return nil, errors.ErrUnprocessableEntity
	}

	members, err := r.client.SMembers(fmt.Sprint(prefix, key)).Result()
	if err != nil {
		fmt.Println(err)
		return nil, errors.ErrUnprocessableEntity
	}
	return members, nil
}

// Set a key's time to live in seconds
func (r *Redis) Expire(key string, expiration time.Duration) error {
	if r.client == nil {
		log.Println("warning: redis not connected")
		return errors.ErrUnprocessableEntity
	}

	err := r.client.Expire(key, expiration).Err()
	if err != nil {
		fmt.Println(err)
		return errors.ErrUnprocessableEntity
	}
	return nil
}

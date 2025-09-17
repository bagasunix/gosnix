package middlewares

import (
	"context"
	"fmt"
	"log"
	"math"
	"strconv"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"

	"github.com/bagasunix/gosnix/pkg/configs"
)

// Local storage untuk fallback ketika Redis down
type localRateLimitData struct {
	sync.RWMutex
	tokens map[string]float64
	ts     map[string]int64
	counts map[string]int64
}

var localData = &localRateLimitData{
	tokens: make(map[string]float64),
	ts:     make(map[string]int64),
	counts: make(map[string]int64),
}

func Limitter(redisClient *redis.Client, cfg *configs.Cfg) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Jika Redis client nil, gunakan hanya local storage
		if redisClient == nil {
			log.Printf("Redis client is nil, using local storage only")
			return handleLocalRateLimit(c, cfg)
		}

		ctx := context.Background()
		ip := c.IP()
		now := time.Now().Unix()

		// Cek ketersediaan Redis
		redisAvailable := true
		if _, err := redisClient.Ping(ctx).Result(); err != nil {
			redisAvailable = false
			log.Printf("Redis tidak tersedia, menggunakan fallback local storage: %v", err)
		}

		// Jika Redis tidak tersedia, gunakan local storage
		if !redisAvailable {
			return handleLocalRateLimit(c, cfg)
		}

		// ----- TOKEN BUCKET dengan Redis -----
		bucketKey := fmt.Sprintf("bucket:%s", ip)
		refillRate := float64(cfg.Server.RateLimiter.Limit) / cfg.Server.RateLimiter.Duration.Seconds()
		capacity := cfg.Server.RateLimiter.Limit

		pipe := redisClient.TxPipeline()
		vals, err := redisClient.HMGet(ctx, bucketKey, "tokens", "ts").Result()
		if err != nil {
			log.Printf("Error mengambil data dari Redis: %v", err)
			return handleLocalRateLimit(c, cfg)
		}

		tokens := float64(capacity)
		lastTs := now

		if vals[0] != nil {
			tokens, _ = strconv.ParseFloat(fmt.Sprint(vals[0]), 64)
		}

		if vals[1] != nil {
			lastTs, _ = strconv.ParseInt(fmt.Sprint(vals[1]), 10, 64)
		}

		// Hitung refill
		elapsed := now - lastTs
		tokens = math.Min(float64(capacity), tokens+float64(elapsed)*refillRate)
		lastTs = now

		// Ambil 1 token
		if tokens < 1 {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{"error": "token bucket empty"})
		}
		tokens--

		// Simpan kembali
		pipe.HSet(ctx, bucketKey, "tokens", tokens, "ts", lastTs)
		pipe.Expire(ctx, bucketKey, cfg.Server.RateLimiter.Duration*2)
		if _, err := pipe.Exec(ctx); err != nil {
			log.Printf("Error menyimpan data ke Redis: %v", err)
		}

		// ----- SLIDING WINDOW dengan Redis -----
		winSize := int64(cfg.Server.RateLimiter.Duration.Seconds())
		currWin := now / winSize
		prevWin := currWin - 1
		keyCurr := fmt.Sprintf("sw:%s:%d", ip, currWin)
		keyPrev := fmt.Sprintf("sw:%s:%d", ip, prevWin)

		currCount, err := redisClient.Incr(ctx, keyCurr).Result()
		if err != nil {
			log.Printf("Error incrementing counter: %v", err)
			return handleLocalRateLimit(c, cfg)
		}

		if currCount == 1 {
			redisClient.Expire(ctx, keyCurr, cfg.Server.RateLimiter.Duration*2)
		}

		prevCount, err := redisClient.Get(ctx, keyPrev).Int64()
		if err != nil && err != redis.Nil {
			log.Printf("Error mengambil previous count: %v", err)
			prevCount = 0
		}

		elapsedWin := now % winSize
		est := float64(currCount) + float64(prevCount)*(1-float64(elapsedWin)/float64(winSize))

		if est > float64(cfg.Server.RateLimiter.Limit) {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{"error": "sliding window exceeded"})
		}

		// ----- ADAPTIVE dengan Redis -----
		if est > float64(cfg.Server.RateLimiter.Limit*3) {
			err := redisClient.Set(ctx, fmt.Sprintf("penalty:%s", ip), "1", 5*time.Minute).Err()
			if err != nil {
				log.Printf("Error menyimpan penalty: %v", err)
			}
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{"error": "adaptive block triggered"})
		}

		return c.Next()
	}
}

// handleLocalRateLimit menangani rate limiting menggunakan local storage
func handleLocalRateLimit(c *fiber.Ctx, cfg *configs.Cfg) error {
	ip := c.IP()
	now := time.Now().Unix()
	bucketKey := fmt.Sprintf("bucket:%s", ip)
	refillRate := float64(cfg.Server.RateLimiter.Limit) / cfg.Server.RateLimiter.Duration.Seconds()
	capacity := cfg.Server.RateLimiter.Limit

	// Gunakan local storage untuk token bucket
	localData.Lock()
	defer localData.Unlock()

	// Ambil atau inisialisasi data dari local storage
	tokens, exists := localData.tokens[bucketKey]
	if !exists {
		tokens = float64(capacity)
	}

	lastTs, exists := localData.ts[bucketKey]
	if !exists {
		lastTs = now
	}

	// Hitung refill
	elapsed := now - lastTs
	tokens = math.Min(float64(capacity), tokens+float64(elapsed)*refillRate)
	lastTs = now

	// Ambil 1 token
	if tokens < 1 {
		return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{"error": "token bucket empty (local)"})
	}
	tokens--

	// Simpan kembali ke local storage
	localData.tokens[bucketKey] = tokens
	localData.ts[bucketKey] = lastTs

	// Bersihkan data lama dari local storage
	go func() {
		time.Sleep(cfg.Server.RateLimiter.Duration * 2)
		localData.Lock()
		delete(localData.tokens, bucketKey)
		delete(localData.ts, bucketKey)
		localData.Unlock()
	}()

	// Sliding window dengan local storage
	winSize := int64(cfg.Server.RateLimiter.Duration.Seconds())
	currWin := now / winSize
	prevWin := currWin - 1
	keyCurr := fmt.Sprintf("sw:%s:%d", ip, currWin)
	keyPrev := fmt.Sprintf("sw:%s:%d", ip, prevWin)

	// Current window
	if _, exists := localData.counts[keyCurr]; !exists {
		localData.counts[keyCurr] = 1
	} else {
		localData.counts[keyCurr]++
	}
	currCount := localData.counts[keyCurr]

	// Previous window
	prevCount, exists := localData.counts[keyPrev]
	if !exists {
		prevCount = 0
	}

	elapsedWin := now % winSize
	est := float64(currCount) + float64(prevCount)*(1-float64(elapsedWin)/float64(winSize))

	if est > float64(cfg.Server.RateLimiter.Limit) {
		return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{"error": "sliding window exceeded (local)"})
	}

	// Adaptive dengan local storage
	if est > float64(cfg.Server.RateLimiter.Limit*3) {
		penaltyKey := fmt.Sprintf("penalty:%s", ip)
		localData.counts[penaltyKey] = 1

		// Bersihkan penalty setelah 5 menit
		go func() {
			time.Sleep(5 * time.Minute)
			localData.Lock()
			delete(localData.counts, penaltyKey)
			localData.Unlock()
		}()

		return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{"error": "adaptive block triggered (local)"})
	}

	return c.Next()
}

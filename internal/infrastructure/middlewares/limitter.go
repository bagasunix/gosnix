package middlewares

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"

	"github.com/bagasunix/gosnix/pkg/configs"
)

func Limitter(redisClient *redis.Client, cfg *configs.Cfg) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := context.Background()
		ip := c.IP()
		now := time.Now().Unix()

		// ----- TOKEN BUCKET -----
		bucketKey := fmt.Sprintf("bucket:%s", ip)
		// isi ulang token (refill_rate token per detik)
		refillRate := float64(cfg.Server.RateLimiter.Limit) / cfg.Server.RateLimiter.Duration.Seconds()
		capacity := cfg.Server.RateLimiter.Limit

		pipe := redisClient.TxPipeline()
		// Ambil token dan timestamp terakhir
		vals, _ := redisClient.HMGet(ctx, bucketKey, "tokens", "ts").Result()
		tokens, _ := strconv.ParseFloat(fmt.Sprint(vals[0]), 64)
		lastTs, _ := strconv.ParseInt(fmt.Sprint(vals[1]), 10, 64)

		if lastTs == 0 { // belum ada data, inisialisasi
			tokens = float64(capacity)
			lastTs = now
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
		_, _ = pipe.Exec(ctx)

		// ----- SLIDING WINDOW -----
		winSize := int64(cfg.Server.RateLimiter.Duration.Seconds())
		currWin := now / winSize
		prevWin := currWin - 1
		keyCurr := fmt.Sprintf("sw:%s:%d", ip, currWin)
		keyPrev := fmt.Sprintf("sw:%s:%d", ip, prevWin)

		currCount, _ := redisClient.Incr(ctx, keyCurr).Result()
		if currCount == 1 {
			redisClient.Expire(ctx, keyCurr, cfg.Server.RateLimiter.Duration*2)
		}
		prevCount, _ := redisClient.Get(ctx, keyPrev).Int64()

		elapsedWin := now % winSize
		est := float64(currCount) + float64(prevCount)*(1-float64(elapsedWin)/float64(winSize))

		if est > float64(cfg.Server.RateLimiter.Limit) {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{"error": "sliding window exceeded"})
		}

		// ----- ADAPTIVE (contoh sederhana) -----
		// kalau request > 3x limit → penalti (turunkan refillRate sementara)
		if est > float64(cfg.Server.RateLimiter.Limit*3) {
			// turunkan kapasitas refill untuk IP ini
			redisClient.Set(ctx, fmt.Sprintf("penalty:%s", ip), "1", 5*time.Minute)
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{"error": "adaptive block triggered"})
		}

		return c.Next()
	}
}

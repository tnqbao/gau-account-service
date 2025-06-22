package repository

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"
)

const (
	RefreshTokenIDBitmap        = "refresh_token_bitmap"
	RefreshTokenBlacklistBitmap = "refresh_token_blacklist_bitmap"
)

func (r *Repository) AllocateRefreshTokenID(ctx context.Context) (int64, error) {
	id, err := r.cacheDb.BitPos(ctx, RefreshTokenIDBitmap, 0).Result()
	if err != nil || id < 0 {
		return -1, err
	}

	if _, err := r.cacheDb.SetBit(ctx, RefreshTokenIDBitmap, id, 1).Result(); err != nil {
		return -1, err
	}

	// Clear blacklist just in case
	r.cacheDb.SetBit(ctx, RefreshTokenBlacklistBitmap, id, 0)

	return id, nil
}

func (r *Repository) ReleaseAndBlacklistID(ctx context.Context, id int64) error {
	if _, err := r.cacheDb.SetBit(ctx, RefreshTokenIDBitmap, id, 0).Result(); err != nil {
		return err
	}
	if _, err := r.cacheDb.SetBit(ctx, RefreshTokenBlacklistBitmap, id, 1).Result(); err != nil {
		return err
	}
	return nil
}

func (r *Repository) ReleaseID(ctx context.Context, id int64) error {
	if _, err := r.cacheDb.SetBit(ctx, RefreshTokenIDBitmap, id, 0).Result(); err != nil {
		return err
	}
	return nil
}

func (r *Repository) IsRefreshTokenBlacklisted(ctx context.Context, id int64) (bool, error) {
	bit, err := r.cacheDb.GetBit(ctx, RefreshTokenBlacklistBitmap, id).Result()
	if err != nil {
		return false, err
	}
	return bit == 1, nil
}

func (r *Repository) ReleaseAndBlacklistIDWithTTL(ctx context.Context, id int64, ttl time.Duration) error {
	pipe := r.cacheDb.TxPipeline()

	// Release ID from bitmap
	pipe.SetBit(ctx, "blacklist_bitmap", id, 1)

	// Set the ID in the blacklist bitmap
	expireKey := "blacklist_expire:" + strconv.FormatInt(id, 10)
	pipe.Set(ctx, expireKey, 1, ttl)

	_, err := pipe.Exec(ctx)
	return err
}

func (r *Repository) CleanupBlacklistBitmap(ctx context.Context) error {
	const bitmapKey = "blacklist_bitmap"
	const blockSize = 1024
	const maxID = 100_000 // Amount of IDs to scan, adjust as needed

	// 1. Collect all still-active keys
	activeIDs := make(map[int64]struct{})
	var cursor uint64
	for {
		keys, nextCursor, err := r.cacheDb.Scan(ctx, cursor, "blacklist_expire:*", 1000).Result()
		if err != nil {
			return fmt.Errorf("scan error: %w", err)
		}
		for _, key := range keys {
			var id int64
			if _, err := fmt.Sscanf(key, "blacklist_expire:%d", &id); err == nil {
				activeIDs[id] = struct{}{}
			}
		}
		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	// 2. Scan bitmap by chunks
	for start := int64(0); start < maxID; start += blockSize {
		end := start + blockSize
		if end > maxID {
			end = maxID
		}

		// BITFIELD to get blockSize bits starting at offset
		cmds := make([]interface{}, 0, (end-start)*2+1)
		cmds = append(cmds, bitmapKey)
		for i := int64(0); i < end-start; i++ {
			cmds = append(cmds, "GET", "u1", start+i)
		}

		results, err := r.cacheDb.Do(ctx, cmds...).Slice()
		if err != nil {
			log.Printf("BITFIELD failed from %d to %d: %v\n", start, end, err)
			continue
		}

		for i, raw := range results {
			bit, ok := raw.(int64)
			if !ok || bit == 0 {
				continue
			}
			id := start + int64(i)

			// If the ID is not in the active set, clear the bit
			if _, isAlive := activeIDs[id]; !isAlive {
				_, err := r.cacheDb.SetBit(ctx, bitmapKey, id, 0).Result()
				if err != nil {
					log.Printf("Failed to clear bit %d: %v\n", id, err)
				}
			}
		}
	}

	log.Println("Bitmap cleanup complete.")
	return nil
}

func (r *Repository) GetBit(ctx context.Context, key string, offset int64) (int64, error) {
	return r.cacheDb.GetBit(ctx, key, offset).Result()
}

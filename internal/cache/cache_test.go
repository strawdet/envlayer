package cache_test

import (
	"testing"
	"time"

	"github.com/yourorg/envlayer/internal/cache"
)

func TestGet_Miss(t *testing.T) {
	c := cache.New(time.Minute)
	_, ok := c.Get("production")
	if ok {
		t.Fatal("expected cache miss, got hit")
	}
}

func TestSet_AndGet_Hit(t *testing.T) {
	c := cache.New(time.Minute)
	env := map[string]string{"APP_ENV": "production", "DB_HOST": "localhost"}
	c.Set("production", env)

	got, ok := c.Get("production")
	if !ok {
		t.Fatal("expected cache hit, got miss")
	}
	if got["APP_ENV"] != "production" {
		t.Errorf("APP_ENV = %q, want %q", got["APP_ENV"], "production")
	}
}

func TestGet_ReturnsCopy(t *testing.T) {
	c := cache.New(time.Minute)
	c.Set("dev", map[string]string{"KEY": "original"})

	got, _ := c.Get("dev")
	got["KEY"] = "mutated"

	again, _ := c.Get("dev")
	if again["KEY"] != "original" {
		t.Error("cache returned a reference instead of a copy")
	}
}

func TestGet_ExpiredEntry(t *testing.T) {
	c := cache.New(10 * time.Millisecond)
	c.Set("staging", map[string]string{"X": "1"})
	time.Sleep(20 * time.Millisecond)

	_, ok := c.Get("staging")
	if ok {
		t.Fatal("expected expired entry to be a miss")
	}
}

func TestInvalidate(t *testing.T) {
	c := cache.New(time.Minute)
	c.Set("prod", map[string]string{"A": "1"})
	c.Invalidate("prod")

	_, ok := c.Get("prod")
	if ok {
		t.Fatal("expected miss after invalidation")
	}
}

func TestFlush_ClearsAll(t *testing.T) {
	c := cache.New(time.Minute)
	c.Set("a", map[string]string{"K": "v"})
	c.Set("b", map[string]string{"K": "v"})
	c.Flush()

	if c.Len() != 0 {
		t.Errorf("Len = %d after Flush, want 0", c.Len())
	}
}

func TestLen(t *testing.T) {
	c := cache.New(time.Minute)
	if c.Len() != 0 {
		t.Errorf("initial Len = %d, want 0", c.Len())
	}
	c.Set("x", map[string]string{})
	c.Set("y", map[string]string{})
	if c.Len() != 2 {
		t.Errorf("Len = %d, want 2", c.Len())
	}
}

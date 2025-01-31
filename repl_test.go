package main

import (
	"fmt"
	"sync"
	"testing"
	"time"

	pokecache "github.com/Daxin319/Go-Pokedex/internal"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "HELLO   WORLD\tGO",
			expected: []string{"hello", "world", "go"},
		},
		{
			input:    "Trailing   spaces   ",
			expected: []string{"trailing", "spaces"},
		},
		{
			input:    "",
			expected: []string{},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("CleanInput(%q) = %v, want %v", c.input, actual, c.expected)
		}

		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("CleanInput(%q) = %v, want %v", c.input, word, expectedWord)
			}
		}
	}
}

func TestAddGet(t *testing.T) {
	const interval = 5 * time.Second
	cases := []struct {
		key string
		val []byte
	}{
		{
			key: "https://example.com",
			val: []byte("testdata"),
		},
		{
			key: "https://example.com/path",
			val: []byte("moretestdata"),
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			cache := pokecache.NewCache(interval)
			cache.Add(c.key, c.val)
			val, ok := cache.Get(c.key)
			if !ok {
				t.Errorf("expected to find key")
				return
			}
			if string(val) != string(c.val) {
				t.Errorf("expected to find value")
				return
			}
		})
	}
}

func TestReapLoop(t *testing.T) {
	const baseTime = 5 * time.Millisecond
	const waitTime = baseTime + 5*time.Millisecond
	cache := pokecache.NewCache(baseTime)
	cache.Add("https://example.com", []byte("testdata"))

	_, ok := cache.Get("https://example.com")
	if !ok {
		t.Errorf("expected to find key")
		return
	}

	time.Sleep(waitTime)

	_, ok = cache.Get("https://example.com")
	if ok {
		t.Errorf("expected to not find key")
		return
	}
}

func TestAddGet2(t *testing.T) {
	const interval = 5 * time.Second
	cache := pokecache.NewCache(interval)

	key := "testKey"
	value := []byte("testValue")
	cache.Add(key, value)

	retrievedValue, ok := cache.Get(key)
	if !ok {
		t.Errorf("expected to find key %q in cache", key)
	}
	if string(retrievedValue) != string(value) {
		t.Errorf("expected value %q, got %q", string(value), string(retrievedValue))
	}
}

func TestExpiration(t *testing.T) {
	const interval = 5 * time.Millisecond
	const waitTime = interval + 5*time.Millisecond
	cache := pokecache.NewCache(interval)

	key := "testKey"
	value := []byte("testValue")
	cache.Add(key, value)

	_, ok := cache.Get(key)
	if !ok {
		t.Errorf("expected to find key %q in cache", key)
	}

	time.Sleep(waitTime)

	_, ok = cache.Get(key)
	if ok {
		t.Errorf("expected key %q to be expired and removed from cache", key)
	}
}

func TestConcurrency(t *testing.T) {
	const interval = 5 * time.Second
	cache := pokecache.NewCache(interval)
	var wg sync.WaitGroup

	// Simulate concurrent access
	for i := 0; i < 100000; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("key%d", i)
			value := []byte(fmt.Sprintf("value%d", i))
			cache.Add(key, value)

			retrievedValue, ok := cache.Get(key)
			if !ok {
				t.Errorf("expected to find key %q in cache", key)
			}
			if string(retrievedValue) != string(value) {
				t.Errorf("expected value %q, got %q", string(value), string(retrievedValue))
			}
		}(i)
	}

	wg.Wait()
}

func TestNonExistentKey(t *testing.T) {
	const interval = 5 * time.Second
	cache := pokecache.NewCache(interval)

	key := "nonExistentKey"
	_, ok := cache.Get(key)
	if ok {
		t.Errorf("expected key %q to not exist in cache", key)
	}
}

func TestDuplicateKey(t *testing.T) {
	const interval = 5 * time.Second
	cache := pokecache.NewCache(interval)

	key := "testKey"
	value1 := []byte("value1")
	cache.Add(key, value1)

	value2 := []byte("value2")
	cache.Add(key, value2)

	retrievedValue, ok := cache.Get(key)
	if !ok {
		t.Errorf("expected to find key %q in cache", key)
	}
	if string(retrievedValue) != string(value2) {
		t.Errorf("expected value %q, got %q", string(value2), string(retrievedValue))
	}
}

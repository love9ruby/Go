package key

import (
	"math/rand"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateKeys(num int, seed int64, cb func(k string) error) []string {
	if num <= 0 {
		return []string{}
	}

	// local random string
	generator := rand.New(rand.NewSource(seed))
	result := make(chan string, num)
	for i := 0; i < num; i++ {
		go generateKey(generator, result)
	}

	// check in global db by cb
	final := make(chan string, num)
	for i := 0; i < num; i++ {
		key := <-result
		go func(k string) {
			err := cb(key)
			if err == nil {
				final <- key
			} else {
				final <- ""
			}
		}(key)
	}

	// bind result and return in local
	keys := make([]string, 0)
	for i := 0; i < num; i++ {
		k := <-final
		if k != "" {
			keys = append(keys, k)
		}
	}

	return keys
}

func generateKey(randGen *rand.Rand, result chan string) {
	b := make([]byte, 6)
	for i := range b {
		b[i] = charset[randGen.Intn(len(charset))]
	}
	result <- string(b)
}

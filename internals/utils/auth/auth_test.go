package utils_auth_test

import (
	"math/rand"
	utils_auth "multitenant-api-go/internals/utils/auth"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func Test_HashPassword(t *testing.T) {
	t.Parallel()
	t.Run("HashPassword. Generate valid hash", func(t *testing.T) {
		result, err := utils_auth.HashPassword(randSeq(12))
		if err != nil {
			t.Errorf("got error %v", err)
		}
		if result == "" {
			t.Errorf("got empty result '%s'", result)
		}
	})

	t.Run("HashPassword. Should return an error", func(t *testing.T) {
		result, err := utils_auth.HashPassword(randSeq(128))
		if err == nil {
			t.Errorf("expected error, but received nil and result %s", result)
		}
	})
}

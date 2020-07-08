package main

import (
    "os"
    "time"
    "strings"
    "math/rand"
    "crypto/hmac"
    "crypto/sha1"
    "encoding/hex"
)

const (
    letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func contains(arr []string, search string) bool {
    for _, n := range arr {
        if search == n {
            return true
        }
    }
    return false
}

func initRandom() {
    rand.Seed(time.Now().UnixNano())
}

func randStringBytes(n int) string {
    b := make([]byte, n)
    for i := range b {
        b[i] = letterBytes[rand.Intn(len(letterBytes))]
    }
    return string(b)
}

func getEnv(name string, fallback string) string {
	env := os.Getenv(name)
	if len(env) == 0 {
		return fallback
	}

	return env
}

// https://gist.github.com/rjz/b51dc03061dbcff1c521
func verifySignature(secret []byte, signature string, body []byte) bool {
    const signaturePrefix = "sha1="
    const signatureLength = 45 // len(SignaturePrefix) + len(hex(sha1))

    if len(signature) != signatureLength || !strings.HasPrefix(signature, signaturePrefix) {
        return false
    }
	
    actual := make([]byte, 20)
    hex.Decode(actual, []byte(signature[5:]))

    return hmac.Equal(signBody(secret, body), actual)
}

func signBody(secret, body []byte) []byte {
    computed := hmac.New(sha1.New, secret)
    computed.Write(body)
    return []byte(computed.Sum(nil))
}

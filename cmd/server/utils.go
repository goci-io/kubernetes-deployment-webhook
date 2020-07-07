package main

import (
    "time"
    "math/rand"
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

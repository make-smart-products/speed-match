package service

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

var slugRe = regexp.MustCompile(`[^a-z0-9]+`)

func NewID() string {
	return uuid.NewString()
}

func NewToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func Slugify(title string) string {
	s := strings.ToLower(strings.TrimSpace(title))
	s = slugRe.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	if s == "" {
		s = "event"
	}
	return s
}

func UniqueSlug(base string, exists func(string) bool) string {
	if !exists(base) {
		return base
	}
	for i := 2; i < 1000; i++ {
		candidate := fmt.Sprintf("%s-%d", base, i)
		if !exists(candidate) {
			return candidate
		}
	}
	return base + "-" + NewID()[:8]
}

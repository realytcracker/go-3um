package main

import (
	"bytes"
	"math/rand"
	"time"
)

var seed = rand.NewSource(time.Now().UnixNano())

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

//RandomString generates a random alphanumeric string of n length (stolen from stackoverflow)
func RandomString(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, seed.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = seed.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

//SplitHash chops up the user hash into five #xxxxxx character groups for user id colors
//pretty sure this is going to be done client side in the end but i put this here anyway
func SplitHash(s string) []string {
	sub := ""
	subs := []string{}
	n := 6

	runes := bytes.Runes([]byte(s))
	l := len(runes)
	for i, r := range runes {
		sub = sub + string(r)
		if (i+1)%n == 0 {
			subs = append(subs, sub)
			sub = ""
		} else if (i + 1) == l {
			subs = append(subs, sub)
		}
	}

	return subs
}

//IsSetup returns true if an admin exists
func IsSetup() bool {
	var user DBUser
	var err error
	q := dbsession.SelectFrom("users").Where("is_admin = true")
	err = q.One(&user)
	if err == nil {
		return true
	}
	return false
}

//IsAdmin returns true if the userhash is an admin
func IsAdmin(s string) bool {
	var user DBUser
	var err error
	q := dbsession.SelectFrom("users").Where("hash = ? and is_admin = true", s)
	err = q.One(&user)
	if err == nil {
		return true
	}
	return false
}

//IsMod returns true if the userhash is a moderator
func IsMod(s string) bool {
	var user DBUser
	var err error
	q := dbsession.SelectFrom("users").Where("hash = ? and is_mod = true", s)
	err = q.One(&user)
	if err == nil {
		return true
	}
	return false
}

//IsThreadOwner returns true if the userhash matches thread author
func IsThreadOwner(s string, t uint64) bool {
	var user DBUser
	var err error
	q := dbsession.SelectFrom("threads").Where("author = ? and id = ?", s, t)
	err = q.One(&user)
	if err == nil {
		return true
	}
	return false
}

//IsPostOwner returns true if the userhash matches post author
func IsPostOwner(s string, p uint64) bool {
	var user DBUser
	var err error
	q := dbsession.SelectFrom("posts").Where("author = ? and id = ?", s, p)
	err = q.One(&user)
	if err == nil {
		return true
	}
	return false
}

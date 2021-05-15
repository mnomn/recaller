package main

import (
	"fmt"
	"time"
)

type JsonTime struct {
	time.Time
}

func (t JsonTime) MarshalJSON() ([]byte, error) {
	//	stamp := fmt.Sprintf("\"%s\"", t.Format("Mon Jan _2"))
	stamp := "\"\""
	now := time.Now()
	if t.Year() < 2000 {
		stamp = "\"Old times\""
	} else if t.Year() != now.Year() {
		stamp = "\"last year\""
	} else if t.Day() == now.Day()-1 { // Todo: New month.
		stamp = fmt.Sprintf("\"Yesterday %s\"", t.Format("15:04"))
	} else if now.Day() == t.Day() {
		stamp = fmt.Sprintf("\"Today %s\"", t.Format("15:04"))
	} else { // Generic time stamp
		stamp = fmt.Sprintf("%s", t.Format("01-02 15:04"))
	}
	return []byte(stamp), nil
}

type OldPost struct {
	Time   JsonTime
	Input  string
	Output string
}

type OldPosts struct {
	posts []OldPost
}

var oldPostsLists = make(map[string]*OldPosts)

func addToWebLog(in string, out string) {
	ll, ok := oldPostsLists[in]
	if !ok || ll == nil {
		ll = &OldPosts{}
		ll.posts = make([]OldPost, 0)
		// newL := make([]OldPost, 0)
		oldPostsLists[in] = ll
	}
	op := OldPost{JsonTime{time.Now()}, in, out}
	ll.posts = append(ll.posts, op)

	if len(ll.posts) > 20 {
		ll.posts = ll.posts[1:]
	}
}

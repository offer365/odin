package main

import "time"

func main() {
	ticker := time.Tick(1 * time.Minute) // 1分钟
	// expr := cronexpr.MustParse("* * * * *")
	for range ticker {
		// now := time.Date()
		// next := expr.Next(now)
		// time.AfterFunc(next.Sub(now), func() {
		// time.AfterFunc(time.Second, func() {})
	}
}

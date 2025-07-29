package clickhouse_test

//import (
//	"context"
//	"fmt"
//	"math/rand"
//	"testing"
//	"time"
//
//)

//func generateEvents(n int) []event.Event {
//	domains := []string{"www.google.com", "strauhs.dev", "scarydoors.dev", "xkcd.com"}
//	pathnames := []string{"/login", "/blog/1", "/blog/2", "/register", "/about-me"}
//
//	initialTime := time.Now()
//
//	events := make([]event.Event, 0, n)
//	for range n {
//		randomDuration := time.Duration((rand.Intn(100) + 1) * int(time.Millisecond))
//		initialTime = initialTime.Add(randomDuration)
//
//		timestamp := initialTime
//		domain := domains[rand.Intn(len(domains))]
//		pathname := pathnames[rand.Intn(len(pathnames))]
//
//		event := event.Event{
//			Timestamp: timestamp,
//			Domain: domain,
//			Kind: "pageview",
//			Pathname: pathname,
//		}
//
//		events = append(events, event)
//	}
//
//	return events
//}
//
//func BenchmarkInsertEvent(b *testing.B) {
//	ctx := context.Background()
//	if err := clickhouseDB.TruncateTables(ctx, "events"); err != nil {
//		b.Fatalf("truncate events: %s", err)
//	}
//
//	events := generateEvents(10000)
//
//	i := 0
//	for b.Loop() {
//		event := events[i % len(events)]
//		if err := clickhouseDB.InsertEvent(ctx, event); err != nil {
//			b.Errorf("%s", err)
//		}
//		i++
//	}
//	b.ReportMetric(float64(b.N) / b.Elapsed().Seconds(), "events/sec")
//}
//
//func BenchmarkAsyncInsertEvent(b *testing.B) {
//	ctx := context.Background()
//	if err := clickhouseDB.TruncateTables(ctx, "events"); err != nil {
//		b.Fatalf("truncate events: %s", err)
//	}
//
//	events := generateEvents(10000)
//
//	i := 0
//	for b.Loop() {
//		event := events[i % len(events)]
//		if err := clickhouseDB.AsyncInsertEvent(ctx, event); err != nil {
//			b.Errorf("%s", err)
//		}
//		i++
//	}
//	b.ReportMetric(float64(b.N) / b.Elapsed().Seconds(), "events/sec")
//}
//
//func BenchmarkBatchInsertEvent(b *testing.B) {
//	ctx := context.Background()
//
//	batchSizes := [4]int{1000, 10000, 100000, 1000000}
//	for _, size := range batchSizes {
//		b.Run(fmt.Sprintf("BatchSize-%d", size), func (b *testing.B) {
//			if err := clickhouseDB.TruncateTables(ctx, "events"); err != nil {
//				b.Fatalf("truncate events: %s", err)
//			}
//
//			evts := generateEvents(size)
//
//			for b.Loop() {
//				if err := clickhouseDB.BatchInsertEvent(ctx, evts); err != nil {
//					b.Errorf("%s", err)
//				}
//			}
//
//			b.ReportMetric(float64(b.N * size) / b.Elapsed().Seconds(), "events/sec")
//		})
//	}
//}

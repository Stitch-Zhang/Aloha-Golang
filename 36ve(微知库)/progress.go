package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/jedib0t/go-pretty/progress"
	"github.com/tidwall/gjson"
)

var (
	pw  progress.Writer
	cps []courseProgress
	pwg sync.WaitGroup
)

type courseProgress struct {
	track *progress.Tracker
	c     chan int
}

func init() {
	pw = progress.NewWriter()
	pw.SetAutoStop(false)
	pw.SetTrackerLength(25)
	pw.ShowOverallTracker(true)
	pw.ShowTime(true)
	pw.ShowTracker(true)
	pw.ShowValue(true)
	pw.SetMessageWidth(24)
	pw.SetSortBy(progress.SortByPercentDsc)
	pw.SetStyle(progress.StyleDefault)
	pw.SetTrackerPosition(progress.PositionRight)
	pw.SetUpdateFrequency(time.Millisecond * 100)
	pw.Style().Colors = progress.StyleColorsExample
	pw.Style().Options.PercentFormat = "%4.1f%%"
}

func infoLength(resp gjson.Result) (result int) {
	for _, k := range resp.Array() {
		result += len(k.Get("items").Array())
	}
	return
}

func addProgress(name string, resp gjson.Result) {
	//fmt.Println(c)
	counts := infoLength(resp)
	wg.Add(int(counts))
	if len(name) > 20 {
		name = fmt.Sprintf("%s...", name[:20])
	}
	track := &progress.Tracker{
		Message: name,
		Total:   int64(counts),
	}
	cps = append(cps, courseProgress{
		c:     make(chan int, 100),
		track: track,
	})
	//fmt.Println("Added")
	pw.AppendTracker(track)
	return
}

func displayProgress(cp []courseProgress) {
	go pw.Render()
	for _, v := range cp {
		go func(c chan int, track *progress.Tracker) {
			defer wg.Done()
			for !track.IsDone() {
				select {
				case <-c:
					//fmt.Println("Recived")
					track.Increment(1)
				}
			}
		}(v.c, v.track)
	}
	wg.Wait()
}

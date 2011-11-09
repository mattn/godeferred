package main

import . "github.com/mattn/godeferred"
import "syscall"
import "net/http"
import "encoding/xml"

type feed struct {
	Entry []struct {
		Title string
		Link  []struct {
			Rel  string "attr"
			Href string "attr"
		}
		Summary string
	}
}

func main() {
	Deferred().
		Next(func() string {
			return "This return value pass to..."
		}).
		Next(func(v string) {
			println("Here:" + v)
		}).
		Next(func() {
			println("Come on who want to do foot race!")
		}).
		Loop(3, func(i int) {
			println(string("ABC"[i]) + ":Yep!")
		}).
		Next(func() {
			println("Ready Go!")
		}).
		Parallel([]interface{}{
			func() {
				println("A:I'll sleep and go so I'm first to start.")
				syscall.Sleep(1000 * 1000 * 300)
				println("A:What happen!?")
			},
			func() {
				syscall.Sleep(1000 * 1000 * 200)
				println("B:I maybe second.")
			},
			func() {
				syscall.Sleep(1000 * 1000 * 100)
				println("C:I'm fastest!")
			},
		}).
		Next(func() {
			println("Finish!")
		}).
		Next(func() {
			println("Begin to watch feed")
		}).
		HttpGet("http://b.hatena.ne.jp/otsune/atomfeed").
		Next(func(res *http.Response) *feed {
			var f feed
			xml.Unmarshal(res.Body, &f)
			return &f
		}).
		Next(func(f *feed) {
			for _, entry := range f.Entry {
				println(entry.Title + "\n\t" + entry.Link[0].Href)
			}
		}).
		HttpGet("http://b.hatena.ne.otsune/"). // this make error
		Next(func(res *http.Response) {
			println("Don't pass to here.")
		}).
		Error(func(err *error) {
		println("Error occur!")
	})
}

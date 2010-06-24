package main

import . "deferred"
import "syscall"
import "http"
import "xml"

type feed struct {
	Entry []struct {
		Title string
		Link []struct {
			Rel  string "attr"
			Href string "attr"
		}
		Summary string
	}
}

func main() {
	Deferred().
		Next(func() string {
			return "「この戻り値が...」"
		}).
		Next(func(v string) {
			println("ここの引数にくる！:" + v);
		}).
		Next(func() {
			println("かけっこすものよっといで！");
		}).
		Loop(3, func(i int) {
			println(string("ABC"[i]) + ":はい！");
		}).
		Next(func() {
			println("位置についてよーいどん！");
		}).
		Parallel([]interface{} {
			func() {
				println("A:一番手だしちょっと昼寝してから行くか");
				syscall.Sleep(1000*1000*300);
				println("A:ちょwww");
			},
			func() {
				syscall.Sleep(1000*1000*200);
				println("B:たぶん２位かな？");
			},
			func() {
				syscall.Sleep(1000*1000*100);
				println("C:俺いっちばーん！");
			},
		}).
		Next(func() {
			println("しゅーりょー！");
		}).
		Next(func() {
			println("otsune:ネットウォッチでもするか！");
		}).
		HttpGet("http://b.hatena.ne.jp/otsune/atomfeed").
		Next(func(res *http.Response) *feed {
			var f feed;
			err := xml.Unmarshal(res.Body, &f);
			if err != nil {
				println(err.String());
			}
			return &f;
		}).
		Next(func(f *feed) {
			for _, entry := range f.Entry {
				println(entry.Title + "\n\t" + entry.Link[0].Href);
			}
		})
}

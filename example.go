package main

import . "deferred"
import "syscall"

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
		})
}

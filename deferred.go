package deferred

import "reflect"
import "http"
import "os"
import "net"

type deferred struct { ret []reflect.Value; err os.Error }

func (v *deferred) Error(f interface{}) *deferred {
	ff := reflect.NewValue(f).(*reflect.FuncValue)
	ft := ff.Type().(*reflect.FuncType)
	v.ret = ff.Call(v.ret[0:ft.NumIn()])
	return v;
}

func (v *deferred) Next(f interface{}) *deferred {
	if v.err != nil { return v }
	ff := reflect.NewValue(f).(*reflect.FuncValue)
	ft := ff.Type().(*reflect.FuncType)
	v.ret = ff.Call(v.ret[0:ft.NumIn()])
	return v;
}

func (v *deferred) Loop(n int, f interface{}) *deferred {
	if v.err != nil { return v }
	ff := reflect.NewValue(f).(*reflect.FuncValue)
	for i := 0; i < n; i++ { ff.Call([]reflect.Value{reflect.NewValue(i)}) }
	return v;
}

func (v *deferred) Parallel(fa []interface{}) *deferred {
	if v.err != nil { return v }
	wait := make(chan interface{}, len(fa))
	for _, f := range fa {
		ff := reflect.NewValue(f).(*reflect.FuncValue)
		ft := ff.Type().(*reflect.FuncType)
		go func() { wait <- ff.Call(v.ret[0:ft.NumIn()]) }()
	}
	for _ = range fa { <-wait }
	return v;
}

func (v *deferred) HttpGet(url string) *deferred {
	if v.err != nil { return v }
	var r *http.Response;
	var err os.Error;
	if proxy := os.Getenv("HTTP_PROXY"); len(proxy) > 0 {
		proxy_url, _ := http.ParseURL(proxy);
		tcp, _ := net.Dial("tcp", "", proxy_url.Host);
		conn := http.NewClientConn(tcp, nil);
		var req http.Request;
		req.URL, _ = http.ParseURL(url);
		req.Method = "GET";
		err = conn.Write(&req);
		r, err = conn.Read();
	} else {
		r, _, err = http.Get(url);
	}
	return Deferred(r, err);
}

func Deferred(v ...interface{}) *deferred {
	d := &deferred {nil, nil}
	d.ret = make([]reflect.Value, len(v))
	for i := range v { d.ret[i] = reflect.NewValue(v[i]) }
	return d
}

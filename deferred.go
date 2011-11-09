package deferred

import (
	"net/url"
	"reflect"
)
import "net/http"
import "net/http/httputil"
import "os"
import "net"

type deferred struct {
	ret []reflect.Value
	err error
}

func (v *deferred) succeeded() bool {
	if v.err != nil {
		return false
	}
	return true
}

func (v *deferred) check() {
	if len(v.ret) > 0 {
		last := v.ret[len(v.ret)-1]
		switch last.Interface().(type) {
		case error:
			v.err = last.Interface().(error)
		}
	}
}

func (v *deferred) Error(f interface{}) *deferred {
	if v.err == nil {
		return v
	}
	ff := reflect.ValueOf(f)
	ret := ff.Call([]reflect.Value{reflect.ValueOf(&v.err)})
	if len(ret) > 0 {
		v.ret = ret
	}
	return v
}

func (v *deferred) Next(f interface{}) *deferred {
	if !v.succeeded() {
		return v
	}
	ff := reflect.ValueOf(f)
	ft := ff.Type()
	ret := ff.Call(v.ret[0:ft.NumIn()])
	if len(ret) > 0 {
		v.ret = ret
	}
	v.check()
	return v
}

func (v *deferred) Loop(n int, f interface{}) *deferred {
	if !v.succeeded() {
		return v
	}
	ff := reflect.ValueOf(f)
	for i := 0; i < n; i++ {
		ff.Call([]reflect.Value{reflect.ValueOf(i)})
	}
	return v
}

func (v *deferred) Parallel(fa []interface{}) *deferred {
	if !v.succeeded() {
		return v
	}
	wait := make(chan interface{}, len(fa))
	for _, f := range fa {
		ff := reflect.ValueOf(f)
		ft := ff.Type()
		go func() { wait <- ff.Call(v.ret[0:ft.NumIn()]) }()
	}
	for _ = range fa {
		<-wait
	}
	v.check()
	return v
}

func (v *deferred) HttpGet(url_ string) *deferred {
	if !v.succeeded() {
		return v
	}
	var r *http.Response
	var err error
	if proxy := os.Getenv("HTTP_PROXY"); len(proxy) > 0 {
		proxy_url, err := url.Parse(proxy)
		if err != nil {
			v.err = err
			return v
		}
		tcp, err := net.Dial("tcp", proxy_url.Host)
		if err != nil {
			v.err = err
			return v
		}
		conn := httputil.NewClientConn(tcp, nil)
		var req http.Request
		req.URL, err = url.Parse(url_)
		if err != nil {
			v.err = err
			return v
		}
		req.Method = "GET"
		err = conn.Write(&req)
		req.URL, err = url.Parse(url_)
		r, err = conn.Read(&req)
		if r != nil && (r.StatusCode/100) >= 4 {
			v.err = os.NewSyscallError(r.Status, r.StatusCode)
			return v
		}
	} else {
		r, err = http.Get(url_)
	}
	v.check()
	return Deferred(r, err)
}

func Deferred(v ...interface{}) *deferred {
	d := &deferred{nil, nil}
	d.ret = make([]reflect.Value, len(v))
	for i := range v {
		d.ret[i] = reflect.ValueOf(v[i])
	}
	return d
}

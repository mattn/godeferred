package deferred

import "reflect"

type deferred struct { ret []reflect.Value }
func (v *deferred) Next(f interface{}) *deferred {
	ff := reflect.NewValue(f).(*reflect.FuncValue)
	ft := ff.Type().(*reflect.FuncType)
	v.ret = ff.Call(v.ret[0:ft.NumIn()])
	return v;
}

func (v *deferred) Loop(n int, f interface{}) *deferred {
	ff := reflect.NewValue(f).(*reflect.FuncValue)
	for i := 0; i < n; i++ { ff.Call([]reflect.Value{reflect.NewValue(i)}) }
	return v;
}

func (v *deferred) Parallel(fa []interface{}) *deferred {
	wait := make(chan interface{}, len(fa))
	for _, f := range fa {
		ff := reflect.NewValue(f).(*reflect.FuncValue)
		ft := ff.Type().(*reflect.FuncType)
		go func() {
			wait <- ff.Call(v.ret[0:ft.NumIn()])
		}()
	}
	for _ = range fa {
		<-wait
	}
	return v;
}

func Deferred(v ...interface{}) *deferred {
	d := &deferred {nil}
	d.ret = make([]reflect.Value, len(v))
	for i := range v { d.ret[i] = reflect.NewValue(v[i]) }
	return d
}

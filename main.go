package main

import (
	"fmt"
	"reflect"
)

type QueryCallback func() (item interface{}, ok bool)

type LINQer struct {
	Iterate func() QueryCallback
}

type Iterable interface {
	Iterate() QueryCallback
}

func From(d interface{}) LINQer {

	realReflect := reflect.ValueOf(d)
	flyReflect := reflect.New(realReflect.Type()).Elem()
	flyReflect.Set(realReflect)

	switch realReflect.Kind() {
	case reflect.Slice, reflect.Array:
		len := flyReflect.Len()

		return LINQer{
			Iterate: func() QueryCallback {
				index := 0

				return func () (item interface{}, ok bool){
					ok = index < len

					if ok {
						item = flyReflect.Index(index).Interface()
						index++
					}

					return
				}
			},
		}
	default:
		return LINQer{}
	}

}


func (l LINQer) Where(s func(interface{}) bool) LINQer {
	return LINQer{
		Iterate: func() QueryCallback {
			next := l.Iterate()


			return func() (item interface{}, ok bool){
				for item, ok = next(); ok ; item, ok = next() {
					if s(item) {
						return
					}
				}
				return
			}
		},
	}
}

func (l LINQer) Push(d interface{}) LINQer {
	return LINQer{
		Iterate: func() QueryCallback {
			next := l.Iterate()

			items := make([]interface{}, 0)
			for item, ok := next(); ok ; item, ok = next() {
				items = append(items, item)
			}

			items = append(items, d)
			length := len(items)
			index := 0
			return func() (item interface{}, ok bool) {
				if index < length {
					item, ok = items[index], true
					index++
				}
				return
			}
		},
	}
}

func (l LINQer) PutKey(k string, d interface{}) LINQer {
	return LINQer{
		Iterate: func() QueryCallback {
			return func() (item interface{}, ok bool) {
				return
			}
		},
	}
}

func (l LINQer) PutIndex(i uint, d interface{}) LINQer {
	return LINQer{
		Iterate: func() QueryCallback {
			next := l.Iterate()

			items := make([]interface{}, 0)
			copyItems := make([]interface{}, 0)
			for item, ok := next(); ok ; item, ok = next() {
				items = append(items, item)
				copyItems = append(copyItems, item)
			}

			var afterPart []interface{}
			if int(i) > l.Count() { items = append(items, d) } else {

				afterPart = items[:i]

				afterPart = append(afterPart, d)
				afterPart = append(afterPart, copyItems[i:]...)
				}

			length := len(afterPart)
			index := 0
			return func() (item interface{}, ok bool) {
				if index < length {
					item, ok = afterPart[index], true
					index++
				}
				return
			}
		},
	}
}

func slice(val interface{}, i int) []interface{} {
	v := reflect.ValueOf(val.([]interface{}))

	n := reflect.New(v.Type()).Elem()
	n.Set(v)
	return n.Interface().([]interface{})
}

func (l LINQer) First() interface{} {
	item, _ := l.Iterate()()
	return item
}

func (l LINQer) Last() interface{} {
	next := l.Iterate()

	var last interface{}
	for item, ok := next(); ok ; item, ok = next() {
		last = item
	}

	return last
}

func (l LINQer) ForEach(action func(interface{})) {
	next := l.Iterate()

	for item, ok := next(); ok ; item, ok = next(){
		action(item)
	}
}

func (l LINQer) Count() int {
	c := 0
	next := l.Iterate()

	for _, ok := next(); ok ; _, ok = next() {
		c++
	}

	return c
}

func (l LINQer) Results() []interface{} {
	res := make([]interface{}, 0)
	next := l.Iterate()

	for item, ok := next(); ok ; item, ok = next() {
		res = append(res, item)
	}

	return res
}

type Car struct {
	year int
	owner, model string
}

func main() {

	cars := make([]Car, 0)
	cars = append(cars, Car{2000, "emirhan", "m3"})
	cars = append(cars, Car{2006, "emirhan", "e46"})

	From(cars).Where(func(i interface{}) bool {
		return true
	}).PutIndex(0, Car{2004, "Ataman", "e30"}).ForEach(func(i interface{}) {
		fmt.Println(i)
	})

}

func validateQuery(q LINQer) bool {
	next := q.Iterate()

	_, ok := next()
	_, ok2 := next()
	return ok || ok2
}
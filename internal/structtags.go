package internal

import (
	"reflect"
	"sync"

	"github.com/lyraproj/dgo/dgo"
	"github.com/tada/catch"
)

type structTags struct {
	count int
	tags  dgo.Map
}

const tagName = `dgo`

var tagLock sync.RWMutex
var allStructTags = make(map[reflect.Type]*structTags, 10)

func (s *structTags) appendStructMapEntries(rt reflect.Type, entries []dgo.StructMapEntry) []dgo.StructMapEntry {
	cnt := rt.NumField()
	for i := 0; i < cnt; i++ {
		fld := rt.Field(i)
		flt := fld.Type
		if fld.Anonymous {
			if flt.Kind() == reflect.Struct {
				entries = getTags(flt).appendStructMapEntries(flt, entries)
				continue
			}
		}
		var ft dgo.Type
		n := String(fld.Name)
		if s.tags != nil {
			ft, _ = s.tags.Get(n).(dgo.Type)
		}
		if ft == nil {
			ft = TypeFromReflected(fld.Type)
		}
		entries = append(entries, StructMapEntry(n, ft, !ft.Instance(Nil)))
	}
	return entries
}

func (s *structTags) fieldCount() int {
	return s.count
}

func (s *structTags) fieldType(key interface{}) (string, dgo.Type) {
	if ks, ok := Value(key).(dgo.String); ok {
		var tp dgo.Type
		if s.tags != nil {
			tp, _ = s.tags.Get(ks).(dgo.Type)
		}
		return ks.GoString(), tp
	}
	return ``, nil
}

func getTags(rt reflect.Type) *structTags {
	tagLock.RLock()
	st := allStructTags[rt]
	tagLock.RUnlock()
	if st != nil {
		return st
	}

	st = parseTags(rt)
	tagLock.Lock()
	allStructTags[rt] = st
	tagLock.Unlock()
	return st
}

func parseTags(rt reflect.Type) *structTags {
	cnt := 0
	st := &structTags{}
	for i, n := 0, rt.NumField(); i < n; i++ {
		f := rt.Field(i)
		ft := f.Type
		if f.Anonymous && ft.Kind() == reflect.Struct {
			ast := getTags(ft)
			cnt += ast.count
			if ast.tags != nil {
				if st.tags == nil {
					st.tags = ast.tags.Copy(false)
				} else {
					st.tags.PutAll(ast.tags)
				}
			}
			continue
		}
		cnt++

		tag := f.Tag.Get(tagName)
		if tag == `` {
			continue
		}
		var tp dgo.Type
		err := catch.Do(func() {
			tp = AsType(Parse(tag))
		})
		if err != nil {
			panic(catch.Error(`unable to parse dgo tag %s.%s %q: %s`, rt.Name(), f.Name, tag, err.Error()))
		}

		rft := TypeFromReflected(ft)
		if !rft.Assignable(tp) {
			panic(catch.Error(`type %v declared in field tag of %s.%s is not assignable to actual type %v`, tp, rt, f.Name, rft))
		}

		if st.tags == nil {
			st.tags = MutableMap(f.Name, tp)
		} else {
			st.tags.Put(f.Name, tp)
		}
	}
	st.count = cnt
	return st
}

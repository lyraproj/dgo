package loader

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/newtype"
	"github.com/lyraproj/dgo/vf"
)

type (
	// Finder is a function that is called when a loader is asked to deliver a value that it doesn't
	// have. It's the finders responsibility to obtain the value. The finder must return nil when the
	// value could not be found.
	//
	// The key sent to the finder will never be a multi part name. It will reflect a single name that is
	// relative to the loader that the Finder is configured for.
	Finder func(l Loader, key string) interface{}

	// NsCreator is called when a namespace is requested that does not exist. It is the NsCreator's
	// responsibility to create the Loader that represents the new namespace. The creator must return
	// nil when no such namespace can be created.
	//
	// The name sent to the creator will never be a multi part name. It will reflect a single name that is
	// relative to the loader that the Finder is configured for.
	NsCreator func(l Loader, name string) Loader

	// A Loader loads named values on demand. Loaders for nested namespaces can be obtained using the Namespace
	// method.
	//
	// Implementors of Loader must ensure that all methods are safe to use from concurrent go routines.
	//
	// As a convenience, the Loader also implements the Keyed interface. The Get method of that interface will
	// simply dispatch to the Load method.
	Loader interface {
		dgo.Value
		dgo.Keyed

		// AbsoluteName returns the absolute name of this loader, i.e. the absolute name of the parent
		// namespace + '/' + this loaders name or, if this loader has no parent namespace, just this
		// loaders name.
		AbsoluteName() string

		// Load loads value for the given name and returns it, or nil if the value could not be
		// found. A loaded value is cached. The name must not contain the separator character '/'. A
		// load involving a nested name must use Namespace to obtain the correct namespace.
		Load(name string) dgo.Value

		// Name returns this loaders name relative to its parent namespace.
		Name() string

		// Namespace returns the Loader that represents the  given namespace from this loader or nil no such
		// loader exists.
		//
		// The name must not contain the separator character '/'. Nested namespaces must be
		// obtained using multiple calls.
		Namespace(name string) Loader

		// NewChild creates a new loader that is parented by this loader.
		NewChild(finder Finder, nsCreator NsCreator) Loader

		// ParentNamespace returns this loaders parent namespace, or nil, if this loader is in the
		// root namespace.
		ParentNamespace() Loader
	}

	mapLoader struct {
		name     string
		parentNs Loader
		entries  dgo.Map
	}

	loader struct {
		mapLoader
		lock       sync.RWMutex
		namespaces dgo.Map
		finder     Finder
		nsCreator  NsCreator
	}

	childLoader struct {
		Loader
		parent Loader
	}
)

// New returns a new Loader instance
func New(parentNs Loader, name string, entries dgo.Map, finder Finder, nsCreator NsCreator) Loader {
	if entries == nil {
		entries = vf.Map()
	}
	if finder == nil && nsCreator == nil {
		// Immutable map based loader
		return &mapLoader{parentNs: parentNs, name: name, entries: entries.FrozenCopy().(dgo.Map)}
	}

	var namespaces dgo.Map
	if nsCreator == nil {
		namespaces = vf.Map()
	} else {
		namespaces = vf.MutableMap(nil)
	}
	return &loader{
		mapLoader:  mapLoader{parentNs: parentNs, name: name, entries: entries.Copy(finder == nil)},
		namespaces: namespaces,
		finder:     finder,
		nsCreator:  nsCreator}
}

// Type is the basic immutable loader dgo.Type
var Type = newtype.NewNamed(`mapLoader`,
	func(arg dgo.Value) dgo.Value {
		l := &mapLoader{}
		l.init(arg.(dgo.Map))
		return l
	},
	func(v dgo.Value) dgo.Value {
		return v.(*mapLoader).initMap()
	},
	reflect.TypeOf(&mapLoader{}),
	reflect.TypeOf((*Loader)(nil)).Elem())

func (l *mapLoader) init(im dgo.Map) {
	l.name = im.Get(`name`).String()
	l.entries = im.Get(`entries`).(dgo.Map)
}

func (l *mapLoader) initMap() dgo.Map {
	m := vf.MapWithCapacity(5, nil)
	m.Put(`name`, l.name)
	m.Put(`entries`, l.entries)
	return m
}

func (l *mapLoader) String() string {
	return Type.ValueString(l)
}

func (l *mapLoader) Type() dgo.Type {
	return newtype.ExactNamed(Type, l)
}

func (l *mapLoader) Equals(other interface{}) bool {
	if ov, ok := other.(*mapLoader); ok {
		return l.entries.Equals(ov.entries)
	}
	return false
}

func (l *mapLoader) HashCode() int {
	return l.entries.HashCode()
}

func (l *mapLoader) Get(key interface{}) dgo.Value {
	return l.entries.Get(key)
}

func (l *mapLoader) Load(key string) dgo.Value {
	return l.entries.Get(key)
}

func (l *mapLoader) AbsoluteName() string {
	if l.parentNs != nil {
		return l.parentNs.AbsoluteName() + `/` + l.name
	}
	return l.name
}

func (l *mapLoader) Name() string {
	return l.name
}

func (l *mapLoader) Namespace(name string) Loader {
	if name == `` {
		return l
	}
	return nil
}

func (l *mapLoader) NewChild(finder Finder, nsCreator NsCreator) Loader {
	return loaderWithParent(l, finder, nsCreator)
}

func (l *mapLoader) ParentNamespace() Loader {
	return l.parentNs
}

// MutableType is the mutable loader dgo.Type
var MutableType = newtype.NewNamed(`loader`,
	func(args dgo.Value) dgo.Value {
		l := &loader{}
		l.init(args.(dgo.Map))
		return l
	},
	func(v dgo.Value) dgo.Value {
		return v.(*loader).initMap()
	},
	reflect.TypeOf(&loader{}),
	reflect.TypeOf((*Loader)(nil)).Elem())

func (l *loader) init(im dgo.Map) {
	l.mapLoader.init(im)
	l.namespaces = im.Get(`namespaces`).(dgo.Map)
}

func (l *loader) initMap() dgo.Map {
	m := l.mapLoader.initMap()
	m.Put(`namespaces`, l.namespaces)
	return m
}

func (l *loader) add(key string, vi interface{}) {
	value := vf.Value(vi)
	var old dgo.Value
	l.lock.Lock()
	if old = l.entries.Get(key); old == nil {
		l.entries.Put(key, value)
	}
	l.lock.Unlock()
	if !(nil == old || old.Equals(value)) {
		panic(fmt.Errorf(`attempt to override entry %q`, key))
	}
}

func (l *loader) Equals(other interface{}) bool {
	if ov, ok := other.(*loader); ok {
		return l.entries.Equals(ov.entries) && l.namespaces.Equals(ov.namespaces)
	}
	return false
}

func (l *loader) HashCode() int {
	return l.entries.HashCode()*31 + l.namespaces.HashCode()
}

func (l *loader) Get(ki interface{}) dgo.Value {
	if key, ok := vf.Value(ki).(dgo.String); ok {
		return l.Load(key.GoString())
	}
	return nil
}

func (l *loader) Load(k string) dgo.Value {
	key := vf.String(k)
	l.lock.RLock()
	v := l.entries.Get(key)
	l.lock.RUnlock()
	if v != nil || l.finder == nil {
		return v
	}
	if v = vf.Value(l.finder(l, k)); v != nil {
		l.add(k, v)
	}
	return v
}

func (l *loader) Namespace(name string) Loader {
	if name == `` {
		return l
	}

	l.lock.RLock()
	ns, ok := l.namespaces.Get(name).(Loader)
	l.lock.RUnlock()

	if ok || l.nsCreator == nil {
		return ns
	}

	if ns = l.nsCreator(l, name); ns != nil {
		var old dgo.Value

		l.lock.Lock()
		if old = l.namespaces.Get(name); old == nil {
			l.namespaces.Put(name, ns)
		}
		l.lock.Unlock()

		if nil != old {
			// Either the nsCreator did something wrong that resulted in the creation of this
			// namespace or a another one has been created from another go routine.
			if !old.Equals(ns) {
				panic(fmt.Errorf(`namespace %q is already defined`, name))
			}

			// Get rid of the duplicate
			ns = old.(Loader)
		}
	}
	return ns
}

func (l *loader) NewChild(finder Finder, nsCreator NsCreator) Loader {
	return loaderWithParent(l, finder, nsCreator)
}

func (l *loader) String() string {
	return MutableType.ValueString(l)
}

func (l *loader) Type() dgo.Type {
	return newtype.ExactNamed(MutableType, l)
}

// ChildType is the parented loader dgo.Type
var ChildType = newtype.NewNamed(`childLoader`,
	func(args dgo.Value) dgo.Value {
		l := &childLoader{}
		l.init(args.(dgo.Map))
		return l
	},
	func(v dgo.Value) dgo.Value {
		return v.(*childLoader).initMap()
	},
	reflect.TypeOf(&loader{}),
	reflect.TypeOf((*Loader)(nil)).Elem())

func (l *childLoader) init(im dgo.Map) {
	l.Loader = im.Get(`loader`).(Loader)
	l.parent = im.Get(`parent`).(Loader)
}

func (l *childLoader) initMap() dgo.Map {
	m := vf.MapWithCapacity(2, nil)
	m.Put(`loader`, l.Loader)
	m.Put(`parent`, l.parent)
	return m
}

func (l *childLoader) Equals(other interface{}) bool {
	if ov, ok := other.(*childLoader); ok {
		return l.Loader.Equals(ov.Loader) && l.parent.Equals(ov.parent)
	}
	return false
}

func (l *childLoader) Get(key interface{}) dgo.Value {
	v := l.parent.Get(key)
	if v == nil {
		v = l.Loader.Get(key)
	}
	return v
}

func (l *childLoader) HashCode() int {
	return l.Loader.HashCode()*31 + l.parent.HashCode()
}

func (l *childLoader) Namespace(name string) Loader {
	if name == `` {
		return l
	}
	pv := l.parent.Namespace(name)
	v := l.Loader.Namespace(name)
	switch {
	case v == nil:
		v = pv
	case pv == nil:
	default:
		v = &childLoader{Loader: v, parent: pv}
	}
	return v
}

func (l *childLoader) NewChild(finder Finder, nsCreator NsCreator) Loader {
	return loaderWithParent(l, finder, nsCreator)
}

func (l *childLoader) String() string {
	return ChildType.ValueString(l)
}

func (l *childLoader) Type() dgo.Type {
	return newtype.ExactNamed(ChildType, l)
}

func loaderWithParent(parent Loader, finder Finder, nsCreator NsCreator) Loader {
	return &childLoader{Loader: New(parent.ParentNamespace(), parent.Name(), vf.Map(), finder, nsCreator), parent: parent}
}

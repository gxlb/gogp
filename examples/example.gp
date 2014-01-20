//This is an example of using gpg tool for generic-programming
//this is an example of using gpg to define an auto-lock global value with generic type
//it will be realized to real go code by gpg tool through the .gpg file with the same name
package example

import (
	"sync"
)

type AutoLockGbl_<TYPE_NAME> struct {
	val  <VALUE_TYPE>
	lock sync.RWMutex
}

func (me *AutoLockGbl_<TYPE_NAME>) Get() (r <VALUE_TYPE>) {
	me.lock.RLock()
	defer me.lock.RUnlock()
	r = me.val
	return
}

func (me *AutoLockGbl_<TYPE_NAME>) Set(val <VALUE_TYPE>) (r <VALUE_TYPE>) {
	me.lock.Lock()
	defer me.lock.Unlock()
	r = me.val
	me.val = val
	return
}

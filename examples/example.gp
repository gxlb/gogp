//This is an example of using gopg tool for generic-programming
//this is an example of using gopg to define an auto-lock global value with generic type
//it will be realized to real go code by gopg tool through the .gpg file with the same name

package <PACKAGE>

import (
	<LOCK_COMMENT>"sync"
)

//auto locked global value
type AutoLockGbl<TYPE_NAME> struct {
	val  <VALUE_TYPE>
	<LOCK_COMMENT>lock sync.RWMutex
}

//new and init a global value
func New<TYPE_NAME>(val <VALUE_TYPE>) *AutoLockGbl<TYPE_NAME>{
	p := &AutoLockGbl<TYPE_NAME>{}
	p.val = val
	return p
}

//get value, if modify is disable, lock is unneeded
func (me *AutoLockGbl<TYPE_NAME>) Get() (r <VALUE_TYPE>) {
<LOCK_COMMENT>	me.lock.RLock()
<LOCK_COMMENT>	defer me.lock.RUnlock()
	r = me.val
	return
}

//set value, if modify is disable, delete this function
<LOCK_COMMENT>func (me *AutoLockGbl<TYPE_NAME>) Set(val <VALUE_TYPE>) (r <VALUE_TYPE>) {
<LOCK_COMMENT>	me.lock.Lock()
<LOCK_COMMENT>	defer me.lock.Unlock()
<LOCK_COMMENT>	r = me.val
<LOCK_COMMENT>	me.val = val
<LOCK_COMMENT>	return
<LOCK_COMMENT>}

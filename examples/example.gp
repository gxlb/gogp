//This is an example of using gpg tool for generic-programming
//this is an example of using gpg to define an auto-lock global value with generic type
//it will be realized to real go code by gpg tool through the .gpg file with the same name
package example

import (
	<LOCK_COMMENT>"sync"
)

type AutoLockGbl<TYPE_NAME> struct {
	val  <VALUE_TYPE>
	<LOCK_COMMENT>lock sync.RWMutex
}

func New<TYPE_NAME>(val <VALUE_TYPE>) *AutoLockGbl<TYPE_NAME>{
	p := &AutoLockGbl<TYPE_NAME>{}
	p.val = val
	return p
}

func (me *AutoLockGbl<TYPE_NAME>) Get() (r <VALUE_TYPE>) {
<LOCK_COMMENT>	me.lock.RLock()
<LOCK_COMMENT>	defer me.lock.RUnlock()
	r = me.val
	return
}

<LOCK_COMMENT>func (me *AutoLockGbl<TYPE_NAME>) Set(val <VALUE_TYPE>) (r <VALUE_TYPE>) {
<LOCK_COMMENT>	me.lock.Lock()
<LOCK_COMMENT>	defer me.lock.Unlock()
<LOCK_COMMENT>	r = me.val
<LOCK_COMMENT>	me.val = val
<LOCK_COMMENT>	return
<LOCK_COMMENT>}

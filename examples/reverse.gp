//This is an example of using gopg tool for generic-programming
//this is an example of using gopg to define an auto-lock global value with generic type
//it will be realized to real go code by gopg tool through the .gpg file with the same name
<PACKAGE>

import (
	"sync"
)

<GOGP_DUMY_DEFINE_COMMENT>/*
//these defines will never exists in real go files
type <STORE_VALUE> int

<GOGP_DUMY_DEFINE_COMMENT>*/

//auto locked global value
type AutoLockGbl<GBL_NAME_SUFFICE> struct {
	val  <STORE_VALUE>
	lock sync.RWMutex
}

//new and init a global value
func NewGbl<GBL_NAME_SUFFICE>(val <STORE_VALUE>) *AutoLockGbl<GBL_NAME_SUFFICE> {
	p := &AutoLockGbl<GBL_NAME_SUFFICE>{}
	p.val = val
	return p
}

//hello he hehehe

//get value, if modify is disable, lock is unneeded
<LOCK_COMMENT>func (me *AutoLockGblInt) Get() (r TemplateVlue) {
<LOCK_COMMENT>	me.lock.RLock()
<LOCK_COMMENT>	defer me.lock.RUnlock()
<LOCK_COMMENT>	r = me.val
<LOCK_COMMENT>	return
<LOCK_COMMENT>}

//set value, if modify is disable, delete this function
func (me *AutoLockGbl<GBL_NAME_SUFFICE>) Set(val <STORE_VALUE>) (r <STORE_VALUE>) {
	me.lock.Lock()
	defer me.lock.Unlock()
	r = me.val
	me.val = val
	return
}

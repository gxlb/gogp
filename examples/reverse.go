//This is an example of using gopg tool for generic-programming
//this is an example of using gopg to define an auto-lock global value with generic type
//it will be realized to real go code by gopg tool through the .gpg file with the same name
package example

import (
	"sync"
)

//GOGPDummyDefineComment/*
//these defines will never exists in real go files
type GOGPStoreValue int

//GOGPDummyDefineComment*/

//auto locked global value
type AutoLockGblGOPGGblNameSuffice struct {
	val  GOGPStoreValue
	lock sync.RWMutex
}

//new and init a global value
func NewGOGOPGGblNameSuffice(val GOGPStoreValue) *AutoLockGblGOGOPGGblNameSuffice {
	p := &AutoLockGblGOGOPGGblNameSuffice{}
	p.val = val
	return p
}

//hello he hehehe

//get value, if modify is disable, lock is unneeded
//GOGPGLockCommentfunc (me *AutoLockGblInt) Get() (r TemplateVlue) {
//GOGPGLockComment	me.lock.RLock()
//GOGPGLockComment	defer me.lock.RUnlock()
//GOGPGLockComment	r = me.val
//GOGPGLockComment	return
//GOGPGLockComment}

//set value, if modify is disable, delete this function
func (me *AutoLockGblGOGOPGGblNameSuffice) Set(val GOGPStoreValue) (r GOGPStoreValue) {
	me.lock.Lock()
	defer me.lock.Unlock()
	r = me.val
	me.val = val
	return
}

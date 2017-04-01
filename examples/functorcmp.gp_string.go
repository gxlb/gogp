///////////////////////////////////////////////////////////////////
//
// !!!!!!!!!!!! NEVER MODIFY THIS FILE MANUALLY !!!!!!!!!!!!
//
// This file was auto-generated by tool [github.com/vipally/gogp]
// Last update at: [Sat Apr 01 2017 22:48:09]
// Generate from:
//   [github.com/vipally/gogp/examples/gp/functorcmp.gp]
//   [github.com/vipally/gogp/examples/example.gpg] [list_string]
//
// Tool [github.com/vipally/gogp] info:
// CopyRight 2016 @Ally Dale. All rights reserved.
// Author  : Ally Dale(vipally@gmail.com)
// Blog    : http://blog.csdn.net/vipally
// Site    : https://github.com/vipally
// BuildAt :
// Version : 3.0.0.final
//
///////////////////////////////////////////////////////////////////

//this file is used to import by other gp files
//it cannot use independently, simulation C++ stl functors

package examples

//cmp object, zero is Lesser
type CmpString byte

const (
	CmpStringLesser  CmpString = CMPLesser
	CmpStringGreater CmpString = CMPGreater
)

//create cmp object by name
func CreateCmpString(cmpName string) (r CmpString) {
	r = CmpStringLesser.CreateByName(cmpName)
	return
}

//uniformed global function
func (me CmpString) F(left, right string) (ok bool) {
	switch me {
	case CMPLesser:
		ok = me.less(left, right)
	case CMPGreater:
		ok = me.great(left, right)
	}
	return
}

//Lesser object
func (me CmpString) Lesser() CmpString { return CMPLesser }

//Greater object
func (me CmpString) Greater() CmpString { return CMPGreater }

//show as string
func (me CmpString) String() (s string) {
	switch me {
	case CMPLesser:
		s = "Lesser"
	case CMPGreater:
		s = "Greater"
	default:
		s = "error cmp value"
	}
	return
}

//create by bool
func (me CmpString) CreateByBool(bigFirst bool) (r CmpString) {
	if bigFirst {
		r = CMPGreater
	} else {
		r = CMPLesser
	}
	return
}

//create cmp object by name
func (me CmpString) CreateByName(cmpName string) (r CmpString) {
	switch cmpName {
	case "": //default Lesser
		fallthrough
	case "Lesser":
		r = CMPLesser
	case "Greater":
		r = CMPGreater
	default: //unsupport name
		panic(cmpName)
	}
	return
}

//lesser operation
func (me CmpString) less(left, right string) (ok bool) {

	ok = left < right

	return
}

//Greater operation
func (me CmpString) great(left, right string) (ok bool) {

	ok = right < left

	return
}

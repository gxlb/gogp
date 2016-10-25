package examples

//#GOGP_IGNORE_BEGIN//////////////////////////////GOGPCommentDummyGoFile_BEGIN
//
//
///*   //<----This line can be uncommented to disable all this file, and it doesn't effect to the .gp file
//	 //If test or change .gp file required, comment it to modify and cmomile as normal go file
//
//
// This is exactly not a real go code file
// It is used to generate .gp file by gogp tool
// Real go code file will be generated from .gp file
//
//#GOGP_IGNORE_END////////////////////////////////GOGPCommentDummyGoFile

//#GOGP_IFDEF GOGP_Show
import show_bytes "bytes" //
//#GOGP_ENDIF

//#GOGP_IGNORE_BEGIN//ignored import
import dumy_fmt "fmt" //
//#GOGP_IGNORE_END//ignored import

//##GOGP_REQUIRE(github.com/vipally/gx/stl/gp/comparer)

//#GOGP_IGNORE_BEGIN//////////////////////////////GOGPDummyDefine
//
//these defines is used to make sure this dummy go file can be compiled correctlly
//and they will be removed from real go files
//vvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvv

type GOGPValueType int

func (me GOGPValueType) Show() string {
	return dumy_fmt.Sprintf("%d", me)
}

//^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
//#GOGP_IGNORE_END////////////////////////////////GOGPDummyDefine

//stack object
type GOGPStackNamePrefixStack []GOGPValueType

//new object
func NewGOGPStackNamePrefixStack() *GOGPStackNamePrefixStack {
	return &GOGPStackNamePrefixStack{}
}

//push
func (this *GOGPStackNamePrefixStack) Push(v GOGPValueType) {
	*this = append(*this, v)
}

//pop
func (this *GOGPStackNamePrefixStack) Pop() (top GOGPValueType, ok bool) {
	if top, ok = this.Top(); ok {
		*this = (*this)[:this.Depth()-1]
	}
	return
}

//top
func (this *GOGPStackNamePrefixStack) Top() (top GOGPValueType, ok bool) {
	if this.Depth() > 0 {
		top = (*this)[this.Depth()-1]
		ok = true
	}
	return

}

//depth
func (this *GOGPStackNamePrefixStack) Depth() int {
	return len(*this)
}

//#GOGP_IFDEF GOGP_Show
//show
func (this *GOGPStackNamePrefixStack) Show() string {
	var b show_bytes.Buffer
	b.WriteByte('[')
	for _, v := range *this {
		b.WriteString(v.Show())
		b.WriteByte(',')
	}
	if this.Depth() > 0 {
		b.Truncate(b.Len() - 1) //remove last ','
	}
	b.WriteByte(']')
	return b.String()
}

//#GOGP_ENDIF

//#GOGP_IGNORE_BEGIN//////////////////////////////GOGPCommentDummyGoFile
//*/
//#GOGP_IGNORE_END////////////////////////////////GOGPCommentDummyGoFile_END

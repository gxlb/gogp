// MIT License
//
// Copyright (c) 2021 @gxlb
// Url:
//     https://github.com/gxlb
//     https://gitee.com/gxlb
// AUTHORS:
//     Ally Dale <vipally@gamil.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package gogp

type gogpProcessStep int

func (me gogpProcessStep) IsReverse() bool {
	return me >= gogpStepREQUIRE && me <= gogpStepREVERSE
}

func (me gogpProcessStep) String() (s string) {
	switch me {
	case gogpStepREQUIRE:
		s = "Step=[1RequireReplace]"
	case gogpStepREVERSE:
		s = "Step=[2ReverseWork]"
	case gogpStepPRODUCE:
		s = "Step=[3NormalProduce]"
	default:
		s = "Step=Unknown"
	}
	return
}

const (
	gogpStepREQUIRE gogpProcessStep = iota + 1 // require replace in fake go file
	gogpStepREVERSE                            // gen gp file from fake go file
	gogpStepPRODUCE                            // gen go file from gp file
)

// get steps of gogp processor
func getProcessingSteps(removeProductsOnly bool) []gogpProcessStep {
	steps := []gogpProcessStep{gogpStepREVERSE, gogpStepREQUIRE, gogpStepREVERSE, gogpStepPRODUCE} //reverse work first
	if removeProductsOnly {
		steps = []gogpProcessStep{gogpStepPRODUCE, gogpStepREQUIRE, gogpStepREVERSE} //normal work first
	}
	return steps
}

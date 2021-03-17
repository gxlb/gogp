//    CopyRight 2016 @Ally Dale. All rights reserved.
//    Use of this source code is governed by a GPL3.0-style
//    license that can be found in the LICENSE file.
//    Author  : Ally Dale(vipally@gmail.com)
//    Blog    : http://blog.csdn.net/vipally
//    Site    : https://github.com/vipally

//package cpright declaration CopyRight of vipally
package cpright

import (
	"github.com/vipally/cmdline"
)

//<version>
//<buildtime>
//<thiscmd>
//    Use of this source code is governed by a GPL3.0-style.
//    License that can be found in the LICENSE file.
var (
	copyRight = `CopyRight 2021 @Ally Dale. All rights reserved.
		Author  : Ally Dale(vipally@gmail.com)
		Site    : https://github.com/vipally
		Version : <version>
`
)

func init() { //set copyright for vipally
	cmdline.CopyRight(copyRight)
}

//get the CopyRight declaration of vipally
func CopyRight() string {
	return copyRight
}

//add coyright file head for all inside go code files
func AddFileHead(path string) error {
	return nil
}

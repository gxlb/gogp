// package auto runs gogp tool on GoPath when imported.
//
// usage:
//   import (
//       _ "github.com/vipally/gogp/auto" //auto runs gogp tool on GoPath when init()
//   )
package auto

import (
	"fmt"

	"github.com/vipally/gogp"
)

func init() {
	gogp.Silence(true)
	if _, _, _, err := gogp.WorkOnGoPath(); err != nil { //runs gogp tool at GoPath when imported
		fmt.Println(err)
	}
}

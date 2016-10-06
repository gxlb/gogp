// package auto runs gogp tool on GoPath when imported.
//
// usage:
//   import (
//       _ "github.com/vipally/gogp/auto"//auto runs gogp tool on GoPath when init()
//   )
package auto

import (
	"github.com/vipally/gogp"
)

func init() {
	gogp.WorkOnGoPath() //runs gogp tool at GoPath when imported
}

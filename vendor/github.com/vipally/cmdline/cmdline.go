// CopyRight 2016 @Ally Dale. All rights reserved.
// Author  : Ally Dale(vipally@gmail.com)
// Blog    : http://blog.csdn.net/vipally
// Site    : https://github.com/vipally

/*
   package cmdline extends std.flag to support more useful features.

       Ally Dale(vipally@gmail.com) modify from std.flag version 1.7
       Main change list from std.flag:
       1.  Add LogicName and Required field for flag, and modify the flag define interface
       2.  Add Summary and Details and Version info to usage page
       3.  Add labels <thiscmd> <buildtime> <version> for getting runtime info in usage page
       4.  Add interface GetUsage() string
       5.  Modify Parse() logic
       6.  Add no-name flag support
       7.  Add "/flag" style support, named flags can lead with "/" or "-" or "--"
       8.  Fix "-flag = x" or "-flag= x" or "-flag =x" cause panic bug
       9.  Add synonyms support for with-name flags
       10. Format usage page lines head with proper num of space

   Usage as follow:

       import (
           "fmt"
           "github.com/vipally/cmdline"
       )
       func main() {
           var (
               host string
               v4   = false
               ttl  = 0
           )
           cmdline.Version("1.0.2")
           cmdline.Summary("<thiscmd> is an example of cmdline package usage.")
           cmdline.Details(`Version   :<version>
               BulidTime :<buildtime>
               <thiscmd> is an example usage of github.com/vipally/cmdline package.`)
           cmdline.CopyRight("no copyright defined")

           //no-name flag and required ones
           cmdline.StringVar(&host, "", "host", "", true, "host ip or name")
           host2 := cmdline.String("", "host2", "", false, "second host ip or name")

           cmdline.BoolVar(&v4, "4", "v4", v4, false, "ipv4")

           //synonym with the same variables
           cmdline.IntVar(&ttl, "t", "ttl", ttl, false, "ttl")
           cmdline.IntVar(&ttl, "ttl", "synonym of -t", ttl, true, "this is synonym of -t")

           //define a synonym with method AnotherName
           c := cmdline.Int("c", "count", 0, false, "count")
           cmdline.AnotherName("count", "c")

           cmdline.Parse()

           fmt.Printf("host:%s host2:%s v4:%t ttl:%d count:%d\n", host, *host2, v4, ttl, *c)
           fmt.Println(cmdline.GetUsage())

           //cmd example: ping -t=20 /4 127.0.0.1 --count =4 localhost -ttl= 5
           //output:
           //
           //host:127.0.0.1 host2:localhost v4:true ttl:5 count:4
           //Usage of ([ping] Build [Sep 30 2016 21:54:21]):
           //  Summary:
           //    ping is an example of cmdline package usage.
           //
           //  Usage:
           //    ping [-4=<v4>] [-c|count=<count>] [-t|ttl=<ttl>] <host> [<host2>
           //  -4=<v4>
           //      ipv4
           //  -c|count=<count>  int
           //      count
           //  -t|ttl=<ttl>  int
           //      ttl
           //  <host>  required  string
           //      host ip or name
           //  <host2>  string
           //      second host ip or name
           //
           //  CopyRight:
           //    no copyright defined
           //
           //  Details:
           //    Version   :1.0.2
           //    BulidTime :[Sep 30 2016 21:54:21]
           //    ping is an example usage of github.com/vipally/cmdline package.
       }
*/
package cmdline

import (
	"os"
	"path"
	"path/filepath"
	"strings"
)

var (
	thisCmd    = getCmd(os.Args[0])
	workDir, _ = os.Getwd()
)

//format path strings
func formatPath(s string) string {
	return path.Clean(s)
}

//get command
func getCmd(arg0 string) string {
	ext := filepath.Ext(arg0)
	return strings.TrimSuffix(filepath.Base(arg0), ext)
}

//WorkDir get current working path
func WorkDir() string {
	return workDir
}

//adapt for os.Exit
func Exit(code int) {
	os.Exit(code)
}

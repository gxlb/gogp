//Copyright @Ally 2014. All rights reserved.
//Author:    vipally@gmail.com
//Version:   1.0.0
//Build at:  [Jan 19 2014  15:28:01]
//Blog site: http://blog.sina.com.cn/ally2014

//Tool gpg is used to generate Generic-Programming code
package main

import (
	"fmt"
	"os"
	"path"
	"strings"
)

const (
	g_version  = "1.0.1"
	g_help_op  = "-h"
	g_help_op2 = "-help"
)

//replace <T> with val that define in *.gpg

func main() {
	file_path := ""
	exit_code := 0
	args := len(os.Args)
	if args == 2 {
		file_path = os.Args[1]
		if file_path == g_help_op || file_path == g_help_op {
			exit_code = -1
		}
	} else if args == 1 {
		var err error
		if file_path, err = os.Getwd(); err != nil {
			fmt.Println(err)
			exit_code = -2
		}
	} else {
		exit_code = -1
	}

	if exit_code == 0 {
		file_path = strings.Replace(file_path, "\\", "/", -1)
		f, l, e := gen_gp_code_by_path(file_path)
		if e != nil {
			fmt.Println(e)
			exit_code = 1
		}
		fmt.Printf("%d code file(s) have been generated from %d %s file(s)\nat path [%s]\n",
			l, f, g_gpg_ext, file_path)
	} else {
		usage()
	}

	os.Exit(exit_code)
	return
}

func usage() {
	cmd := strings.Replace(os.Args[0], "\\", "/", -1)
	d := path.Dir(cmd) + "/"
	app := strings.TrimPrefix(cmd, d)
	fmt.Printf("Tool [%s] is used to generate Generic-Programming code", app)
	fmt.Printf("\nTool help info as belows:\n%s\n", copy_right(""))
	fmt.Printf("\nusage: %s <help> <path>", app)
	fmt.Print(
		`
  <path> String. Specify path to generate GP code, if not exists, means current path.
  <help>="-h" or "-help". Means show this help info.
`)
}

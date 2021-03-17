// +build windows

package cmdline

func isFlagLeadByte(c byte) bool {
	return c == '-' || c == '/'
}

func isFlagLead(s string) bool {
	return s == "-" || s == "--" || s == "/"
}

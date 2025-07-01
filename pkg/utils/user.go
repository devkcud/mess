package utils

import "os/user"

var RootUser = func() string {
	u, err := user.LookupId("0")
	if err != nil {
		return ""
	}
	return u.Username
}()

var CurrentUser = func() string {
	u, err := user.Current()
	if err != nil {
		return ""
	}
	return u.Username
}()

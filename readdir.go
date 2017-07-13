package sftpd

import (
	"fmt"
	"time"
)

func readdirLongName(fi *NamedAttr) string {
	return fmt.Sprintf("%10s %3d %-8s %-8s %8d %12s %s",
		fi.Mode.String(),
		1, // links
		fi.User, fi.Group,
		fi.Size,
		readdirTimeFormat(fi.MTime),
		fi.Name,
	)
}

func readdirTimeFormat(t time.Time) string {
	// We return timestamps in UTC, should we offer a customisation point for users?
	if t.Year() == time.Now().Year() {
		return t.Format("Jan _2 15:04")
	}
	return t.Format("Jan _2  2006")
}

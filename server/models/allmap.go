package models

import "github.com/yushaona/gmessage/server/job"

func init() {
	job.AddMap(10, &Grow{})
}

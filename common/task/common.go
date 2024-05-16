package task

import "github.com/mrst2000/my-ray/common"

// Close returns a func() that closes v.
func Close(v interface{}) func() error {
	return func() error {
		return common.Close(v)
	}
}

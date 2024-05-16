package reverse_test

import (
	"testing"

	"github.com/mrst2000/my-ray/app/reverse"
	"github.com/mrst2000/my-ray/common"
)

func TestStaticPickerEmpty(t *testing.T) {
	picker, err := reverse.NewStaticMuxPicker()
	common.Must(err)
	worker, err := picker.PickAvailable()
	if err == nil {
		t.Error("expected error, but nil")
	}
	if worker != nil {
		t.Error("expected nil worker, but not nil")
	}
}

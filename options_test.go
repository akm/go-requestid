package requestid

import "testing"

func TestOptions(t *testing.T) {
	t.Run("LogAttr", func(t *testing.T) {
		options := &Options{}
		LogAttr("test-attr")(options)
		if options.logAttr != "test-attr" {
			t.Errorf("LogAttr() = %v; want test-attr", options.logAttr)
		}
	})
}

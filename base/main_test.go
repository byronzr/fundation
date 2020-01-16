package base

import "testing"

func Test_main(t *testing.T) {
	Dump(map[string]string{"test": "test_content"}, map[string]string{"test2": "test_content2"})
}

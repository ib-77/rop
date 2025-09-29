package main

import (
	"github.com/ib-77/rop/pkg/rop/solo"
)

/*
func TestMain(t *testing.M) {
	setupAll()
	code := t.Run()
	tearDownAll()
	os.Exit(code)
}

func setupAll() {

}

func tearDownAll() {

}


func Test_Escape(t *testing.T) {
	t.Parallel()

} */

type testNoEscape struct {
	value [130000]byte // on escape
}

type testEscape struct {
	value [140000]byte // on escape
}

// go run -gcflags="-m" test/others/escape.go
func main() {
	a := solo.Validate(testNoEscape{}, func(in testNoEscape) bool {
		return false
	}, "error")

	_ = a

	//b := solo.Validate(testEscape{}, func(in testEscape) bool {
	//	return false
	//}, "error")
	//
	//_ = b
}

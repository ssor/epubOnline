package tests

import (
    "testing"

    . "github.com/smartystreets/goconvey/convey"
)

func TestDefault(t *testing.T) {
    Convey("TestDefault ", t, func() {
        So(0, ShouldNotEqual, 0)
    })
}
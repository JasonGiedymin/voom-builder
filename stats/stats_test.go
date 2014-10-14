package stats

import (
    "testing"

    "github.com/JasonGiedymin/voom-builder/common"
)

func TestSuperStats(t *testing.T) {
    s := NewStats()
    s.Reserve(100)
    s.Withdraw(50, nil)
    if s.Claims() != 50 {
        t.Error("Count should have been 50")
    }

    s.Withdraw(50, &common.WorkError{"some error"}) //withdraw with error
}

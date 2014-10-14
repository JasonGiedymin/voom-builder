package stats

import (
    "github.com/JasonGiedymin/voom-builder/common"

    "github.com/rcrowley/go-metrics"
)

type ApiStats struct {
    calls metrics.Counter // count of api calls
}

func (s *ApiStats) Mark() {
    s.calls.Inc(1)
}

func (s *ApiStats) Calls() int64 {
    return s.calls.Count()
}

// == Stats ==
type SupervisorStats struct {
    workersInProgress metrics.Counter

    claimsCompleted metrics.Counter // count of errors encountered during a claim
    claimCompletion metrics.Meter   // measure the time between completions

    // errors can also be service restarts, not necessarily related to the work
    claimsErrors          metrics.Counter // count of errors encountered during a claim
    claimCompletionErrors metrics.Meter   // measure the time between failures
}

func (s *SupervisorStats) Success(shares int64) {
    s.claimsCompleted.Inc(shares)
    s.claimCompletion.Mark(1)
}

// We tie the completion of a claim with it's withdrawl.
// Errors are marked properly.
func (s *SupervisorStats) Failure(shares int64, err *common.WorkError) {
    s.claimsErrors.Inc(1)
    s.claimCompletionErrors.Mark(1)
}

func (s *SupervisorStats) SuccessCount() int64 {
    return s.claimsCompleted.Count()
}

func (s *SupervisorStats) Errors() int64 {
    return s.claimsErrors.Count()
}

func (s *SupervisorStats) Snapshot() metrics.Meter {
    return s.claimCompletion.Snapshot()
}

func NewSupervisorStats() *SupervisorStats {
    return &SupervisorStats{
        metrics.NewCounter(),
        metrics.NewCounter(),
        metrics.NewMeter(),
        metrics.NewCounter(),
        metrics.NewMeter(),
    }
}

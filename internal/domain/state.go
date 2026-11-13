package domain

const (
	TakeStateValidating = "validating"
	TakeStateAdmitted   = "admitted"

	AnalysisStateQueued    = "queued"
	AnalysisStateDetecting = "detecting"
	AnalysisStateReady     = "ready"
	AnalysisStateFailed    = "failed"

	ConformStateQueued     = "queued"
	ConformStateRendering  = "rendering"
	ConformStateCommitting = "committing"
	ConformStateCompleted  = "completed"
	ConformStateFailed     = "failed"
)

func CanAdmitTake(next string) bool {
	return next == TakeStateAdmitted
}

func CanAdvanceAnalysis(current, next string) bool {
	switch current {
	case AnalysisStateQueued:
		return next == AnalysisStateDetecting || next == AnalysisStateFailed
	case AnalysisStateDetecting:
		return next == AnalysisStateReady || next == AnalysisStateFailed
	default:
		return false
	}
}

func AnalysisTerminal(state string) bool {
	return state == AnalysisStateReady || state == AnalysisStateFailed
}

func CanAdvanceConform(current, next string) bool {
	switch current {
	case ConformStateQueued:
		return next == ConformStateRendering || next == ConformStateFailed
	case ConformStateRendering:
		return next == ConformStateCommitting || next == ConformStateFailed
	case ConformStateCommitting:
		return next == ConformStateCompleted || next == ConformStateFailed
	default:
		// completed and failed are terminal: a delivered bundle is never
		// rewritten and a failed job keeps its stable error.
		return false
	}
}

// CanRecoverConform reports whether a conform interrupted in the given state
// may still be driven through its lifecycle. Terminal states are never
// recovered: completed jobs stay delivered and failed jobs stay failed.
func CanRecoverConform(state string) bool {
	switch state {
	case ConformStateQueued,
		ConformStateRendering,
		ConformStateCommitting:
		return true
	default:
		return false
	}
}

func ConformTerminal(state string) bool {
	return state == ConformStateCompleted || state == ConformStateFailed
}

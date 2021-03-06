package engine

import (
	"github.com/wailsapp/wails"
	"gopkg.in/freeeve/uci.v1"
)

// UciEngine is a communication layer above an UCI engine.
type UciEngine struct {
	eng     *uci.Engine
	Runtime *wails.Runtime
}

// WailsInit initialize the structure for Wails
func (engine *UciEngine) WailsInit(r *wails.Runtime) error {
	engine.Runtime = r
	return nil
}

// NewUciEngine builds a new UciEngine instance.
func NewUciEngine() *UciEngine {
	return &UciEngine{}
}

// GetUserEnginePath make user select his engine file, with a file chooser
// and returns the path as string
// [output] string: the choosen file, or empty string if user cancelled
func (engine *UciEngine) GetUserEnginePath() string {
	path := engine.Runtime.Dialog.SelectFile()
	return path
}

// PlayPosition attempts to get the best move from the given position
// positionFen string: the position in Forsyth-Edwards Notation
// [output] string: the best move (ie: a7a8q or c3d4) if no error
// otherwise the error code:
// #EngineNotSet: not engine configured
// #ComputationError: either the given position is illegal or a misc. computation error occurred
func (engine *UciEngine) PlayPosition(positionFen string, depth int) string {
	if engine.eng == nil {
		return "#EngineNotSet"
	}

	engine.eng.SetFEN(positionFen)
	resultOpts := uci.HighestDepthOnly | uci.IncludeUpperbounds | uci.IncludeLowerbounds
	results, error := engine.eng.GoDepth(depth, resultOpts)

	if error != nil {
		return "#ComputationError"
	}

	return results.BestMove
}

// LoadEngineWithPathProviden loads an engine from a given path, without a file chooser.
// path string: the absolute path
// [output] string: the error type if any
// #ConfigEngineErr: could not set the required engine
func (engine *UciEngine) LoadEngineWithPathProviden(path string) string {
	newEngine, err := uci.NewEngine(path)
	if err != nil {
		return "#ConfigEngineErr"
	}

	newEngine.SetOptions(uci.Options{
		Hash:    128,
		Ponder:  false,
		OwnBook: true,
	})

	engine.eng = newEngine
	return ""
}

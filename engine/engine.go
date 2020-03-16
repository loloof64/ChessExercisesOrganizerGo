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

// Load let the user select the engine thanks to a dialog
// [output] string: the error type if any
// #ConfigEngineErr: could not set the required engine
func (engine *UciEngine) Load() string {
	path := engine.Runtime.Dialog.SelectFile()

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

// PlayPosition attempts to get the best move from the given position
// positionFen string: the position in Forsyth-Edwards Notation
// [output] string: the best move (ie: a7a8q or c3d4) if no error
// otherwise the error code:
// #EngineNotSet: not engine configured
// #ComputationError: either the given position is illegal or a misc. computation error occurred
func (engine *UciEngine) PlayPosition(positionFen string) string {
	if engine.eng == nil {
		return "#EngineNotSet"
	}

	engine.eng.SetFEN(positionFen)
	resultOpts := uci.HighestDepthOnly | uci.IncludeUpperbounds | uci.IncludeLowerbounds
	results, error := engine.eng.GoDepth(10, resultOpts)

	if error != nil {
		return "#ComputationError"
	}

	return results.BestMove
}

// SetEnginePathManually loads an engine from a given path, without a file chooser.
// path string: the absolute path
// [output] string: the error type if any
// #ConfigEngineErr: could not set the required engine
func (engine *UciEngine) SetEnginePathManually(path string) string {
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

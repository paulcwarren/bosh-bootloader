package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
)

type logger interface {
	Println(message string)
}

type StateBootstrap struct {
	logger     logger
	bblVersion string
}

func NewStateBootstrap(logger logger, bblVersion string) StateBootstrap {
	return StateBootstrap{
		logger:     logger,
		bblVersion: bblVersion,
	}
}

func (b StateBootstrap) GetState(dir string) (State, error) {
	state := State{}

	_, err := os.Stat(dir)
	if err != nil {
		return state, err
	}

	file, err := os.Open(filepath.Join(dir, StateFileName))
	if err != nil {
		if os.IsNotExist(err) {
			return state, nil
		}
		return state, err
	}

	err = json.NewDecoder(file).Decode(&state)
	if err != nil {
		return state, err
	}

	emptyState := State{}
	if reflect.DeepEqual(state, emptyState) {
		state = State{
			Version:    STATE_SCHEMA,
			BBLVersion: b.bblVersion,
		}
	} else {
		if state.BBLVersion == "" {
			state.BBLVersion = b.getBBLVersion(state.Version)
		}
	}

	if state.Version < 3 {
		return state, errors.New("Existing bbl environment is incompatible with bbl v3. Create a new environment with v3 to continue.")
	}

	if state.Version > STATE_SCHEMA {
		return state, fmt.Errorf("Existing bbl environment was created with a newer version of bbl. Please upgrade to bbl v%s.\n", state.BBLVersion)
	}

	if state.Version < STATE_SCHEMA {
		b.logger.Println(fmt.Sprintf("Warning: Current schema version (%d) is newer than existing bbl environment schema (%d). Some things may not work as expected until you bbl up again.", STATE_SCHEMA, state.Version))
	}

	return state, nil
}

// Get the earliest bbl version compatible with the given bbl state version.
func (b StateBootstrap) getBBLVersion(stateSchema int) string {
	stateToBBLVersion := map[int]string{
		3:  "3.0.0",
		5:  "4.0.0",
		6:  "4.0.0",
		7:  "4.0.0",
		8:  "4.0.0",
		9:  "4.4.0",
		10: "4.6.0",
		11: "5.1.0",
		12: "5.1.0",
		13: "5.4.0",
	}
	bblVersion, ok := stateToBBLVersion[stateSchema]
	if ok {
		return bblVersion
	}
	return "dev"
}

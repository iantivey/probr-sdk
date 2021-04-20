// Package probeengine contains the types and functions responsible for managing tests and test execution.  This is the primary
// entry point to the core of the application and should be utilised by the probr library to create, execute and report
// on tests.
package probeengine

import (
	"errors"
	"log"
	"sync"

	"github.com/citihub/probr-sdk/audit"
)

// ProbeStatus type describes the status of the test, e.g. Pending, Running, CompleteSuccess, CompleteFail and Error
type ProbeStatus int

//ProbeStatus enumeration for the ProbeStatus type.
const (
	Pending ProbeStatus = iota
	Running
	CompleteSuccess
	CompleteFail
	Error
	Excluded
)

func (s ProbeStatus) String() string {
	return [...]string{"Pending", "Running", "CompleteSuccess", "CompleteFail", "Error", "Excluded"}[s]
}

// Group type describes the group to which the test belongs, e.g. kubernetes, clouddriver, probeengine, etc.
type Group int

// Group type enumeration
const (
	Kubernetes Group = iota
	CloudDriver
	CoreEngine
)

// ProbeStore maintains a collection of probes to be run and their status.  FailedProbes is an explicit
// collection of failed probes.
type ProbeStore struct {
	Name         string
	Probes       map[string]*GodogProbe
	FailedProbes map[ProbeStatus]*GodogProbe
	Lock         sync.RWMutex
}

// NewProbeStore creates a new object to store GodogProbes
func NewProbeStore(name string) *ProbeStore {
	return &ProbeStore{
		Name:   name,
		Probes: make(map[string]*GodogProbe),
	}
}

// AddProbe provided GodogProbe to the ProbeStore.
func (ps *ProbeStore) AddProbe(probe Probe) {
	ps.Lock.Lock()
	defer ps.Lock.Unlock()
	status := Pending

	ps.Probes[probe.Name()] = &GodogProbe{
		Name:                probe.Name(),
		Group:               ps.Name,
		ProbeInitializer:    probe.ProbeInitialize,
		ScenarioInitializer: probe.ScenarioInitialize,
		FeaturePath:         probe.Path(),
		Status:              &status,
	}

	audit.State.GetProbeLog(probe.Name()).Result = status.String()
	audit.State.LogProbeMeta(probe.Name(), "group", ps.Name)
}

// GetProbe returns the test identified by the given name.
func (ps *ProbeStore) GetProbe(name string) (*GodogProbe, error) {
	ps.Lock.Lock()
	defer ps.Lock.Unlock()

	//get the test from the store
	p, exists := ps.Probes[name]

	if !exists {
		return nil, errors.New("test with name '" + name + "' not found")
	}
	return p, nil
}

// ExecProbe executes the test identified by the specified name.
func (ps *ProbeStore) ExecProbe(name string) (int, error) {
	p, err := ps.GetProbe(name)
	log.Printf("Executing Probe: %v", p)
	if err != nil {
		return 1, err // Failure
	}
	if p.Status.String() != Excluded.String() {
		return ps.RunProbe(p) // Return test results
	}
	return 0, nil // Succeed if test is excluded
}

// ExecAllProbes executes all tests that are present in the ProbeStore.
func (ps *ProbeStore) ExecAllProbes() (int, error) {
	status := 0
	var err error

	for name := range ps.Probes {
		st, err := ps.ExecProbe(name)
		audit.State.ProbeComplete(name)
		if err != nil {
			//log but continue with remaining probe
			log.Printf("[ERROR] error executing probe: %v", err)
		}
		if st > status {
			status = st
		}
	}
	return status, err
}

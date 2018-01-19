package minipipeline

import (
	"fmt"
	"log"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

type Set struct {
	sync.RWMutex
	pipelines []*Pipeline
}

type Pipeline struct {
	sync.RWMutex
	Name  string
	Steps []*Step
	Err   error
}

type Step struct {
	sync.RWMutex
	Name       string
	StartedAt  time.Time
	FinishedAt time.Time
	Err        error

	progressDone uint32
	progressAll  uint32
}

func (s *Step) ProgressDone(delta uint32) {
	s.Lock()
	defer s.Unlock()

	s.progressDone += delta
}

func (s *Step) ProgressAll(delta uint32) {
	s.Lock()
	defer s.Unlock()

	s.progressAll += delta
}

func (s *Set) Snapshot() []*Pipeline {
	s.RLock()
	defer s.RUnlock()

	pipelines := make([]*Pipeline, len(s.pipelines))
	for i, p := range s.pipelines {
		pipelines[i] = p.Copy()
	}
	return pipelines
}

func (s *Set) Pipeline(name string) *Pipeline {
	s.Lock()
	defer s.Unlock()

	p := &Pipeline{Name: name}
	s.pipelines = append(s.pipelines, p)

	return p
}

func (p *Pipeline) Step(name string, fn func() error) error {
	p.Lock()

	if p.Err != nil {
		p.Unlock()
		return p.Err
	}

	step := &Step{Name: name}
	p.Steps = append(p.Steps, step)

	log.Printf("%s {", name)

	step.StartedAt = time.Now()
	p.Unlock()
	err := fn()
	p.Lock()
	step.FinishedAt = time.Now()

	log.Printf("} // %s", name)

	step.Err = err
	p.Err = err

	p.Unlock()
	return err
}

func (p *Pipeline) Current() *Step {
	p.RLock()
	defer p.RUnlock()

	if len(p.Steps) == 0 {
		return nil
	}

	step := p.Steps[len(p.Steps)-1]
	if step.FinishedAt.IsZero() {
		// step is running
		return step
	}

	return nil
}

func (p *Pipeline) Copy() *Pipeline {
	p.RLock()
	defer p.RUnlock()

	var copy = *p

	steps := make([]*Step, len(p.Steps))
	for i, s := range p.Steps {
		s.RLock()
		var copy = *s
		s.RUnlock()
		steps[i] = &copy
	}

	copy.Steps = steps

	return &copy
}

type ProgressGroup struct {
	step *Step
	*errgroup.Group
}

func (s *Step) ProgressGroup() ProgressGroup {
	return ProgressGroup{
		step:  s,
		Group: new(errgroup.Group),
	}
}

func (s *Step) String() string {
	if s.FinishedAt.IsZero() {
		// running
		s.Lock()
		defer s.Unlock()

		if s.progressDone != 0 || s.progressAll != 0 {
			return fmt.Sprintf("● %d/%d", s.progressDone, s.progressAll)
		} else {
			return "●"
		}
	} else if s.Err != nil {
		return fmt.Sprintf("✗ %s", s.Err)
	} else {
		return fmt.Sprintf("✓ %s", s.FinishedAt.Sub(s.StartedAt))
	}
}

func (pg *ProgressGroup) Go(fn func() error) {
	pg.step.ProgressAll(1)
	pg.Group.Go(func() error {
		defer pg.step.ProgressDone(1)
		return fn()
	})
}

// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package resource

import "github.com/drone/runner-go/manifest"

var (
	_ manifest.Resource          = (*Pipeline)(nil)
	_ manifest.TriggeredResource = (*Pipeline)(nil)
	_ manifest.DependantResource = (*Pipeline)(nil)
	_ manifest.PlatformResource  = (*Pipeline)(nil)
)

// Defines the Resource Kind and Type.
const (
	Kind = "pipeline"
	Type = "exec"
)

type (
	// Pipeline is a pipeline resource that executes pipelines
	// on the host machine without any virtualization.
	Pipeline struct {
		Version   string              `json:"version,omitempty"`
		Kind      string              `json:"kind,omitempty"`
		Type      string              `json:"type,omitempty"`
		Name      string              `json:"name,omitempty"`
		Deps      []string            `json:"depends_on,omitempty"`
		Clone     manifest.Clone      `json:"clone,omitempty"`
		Platform  manifest.Platform   `json:"platform,omitempty"`
		Trigger   manifest.Conditions `json:"conditions,omitempty"`
		Workspace manifest.Workspace  `json:"workspace,omitempty"`

		Steps []*Step `json:"steps,omitempty"`
	}

	// Step defines a Pipeline step.
	Step struct {
		Name        string                        `json:"name,omitempty"`
		Shell       string                        `json:"shell,omitempty"`
		DependsOn   []string                      `json:"depends_on,omitempty" yaml:"depends_on"`
		Detach      bool                          `json:"detach,omitempty"`
		Environment map[string]*manifest.Variable `json:"environment,omitempty"`
		Failure     string                        `json:"failure,omitempty"`
		Commands    []string                      `json:"commands,omitempty"`
		When        manifest.Conditions           `json:"when,omitempty"`
	}
)

// GetVersion returns the resource version.
func (p *Pipeline) GetVersion() string { return p.Version }

// GetKind returns the resource kind.
func (p *Pipeline) GetKind() string { return p.Kind }

// GetType returns the resource type.
func (p *Pipeline) GetType() string { return p.Type }

// GetName returns the resource name.
func (p *Pipeline) GetName() string { return p.Name }

// GetDependsOn returns the resource dependencies.
func (p *Pipeline) GetDependsOn() []string { return p.Deps }

// GetTrigger returns the resource triggers.
func (p *Pipeline) GetTrigger() manifest.Conditions { return p.Trigger }

// GetPlatform returns the resource platform.
func (p *Pipeline) GetPlatform() manifest.Platform { return p.Platform }

// GetStep returns the named step. If no step exists with the
// given name, a nil value is returned.
func (p *Pipeline) GetStep(name string) *Step {
	for _, step := range p.Steps {
		if step.Name == name {
			return step
		}
	}
	return nil
}

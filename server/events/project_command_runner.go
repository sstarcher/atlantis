// Copyright 2017 HootSuite Media Inc.
//
// Licensed under the Apache License, Version 2.0 (the License);
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an AS IS BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// Modified hereafter by contributors to runatlantis/atlantis.

package events

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/runatlantis/atlantis/server/events/models"
	"github.com/runatlantis/atlantis/server/events/runtime"
	"github.com/runatlantis/atlantis/server/events/webhooks"
	"github.com/runatlantis/atlantis/server/events/yaml/raw"
	"github.com/runatlantis/atlantis/server/events/yaml/valid"
	"github.com/runatlantis/atlantis/server/logging"
)

//go:generate pegomock generate -m --use-experimental-model-gen --package mocks -o mocks/mock_lock_url_generator.go LockURLGenerator

// LockURLGenerator generates urls to locks.
type LockURLGenerator interface {
	// GenerateLockURL returns the full URL to the lock at lockID.
	GenerateLockURL(lockID string) string
}

//go:generate pegomock generate -m --use-experimental-model-gen --package mocks -o mocks/mock_step_runner.go StepRunner

// StepRunner runs steps. Steps are individual pieces of execution like
// `terraform plan`.
type StepRunner interface {
	// Run runs the step.
	Run(ctx models.ProjectCommandContext, extraArgs []string, path string) (string, error)
}

//go:generate pegomock generate -m --use-experimental-model-gen --package mocks -o mocks/mock_webhooks_sender.go WebhooksSender

// WebhooksSender sends webhook.
type WebhooksSender interface {
	// Send sends the webhook.
	Send(log *logging.SimpleLogger, res webhooks.ApplyResult) error
}

// PlanSuccess is the result of a successful plan.
type PlanSuccess struct {
	// TerraformOutput is the output from Terraform of running plan.
	TerraformOutput string
	// LockURL is the full URL to the lock held by this plan.
	LockURL string
	// RePlanCmd is the command that users should run to re-plan this project.
	RePlanCmd string
	// ApplyCmd is the command that users should run to apply this plan.
	ApplyCmd string
}

//go:generate pegomock generate -m --use-experimental-model-gen --package mocks -o mocks/mock_project_command_runner.go ProjectCommandRunner

// ProjectCommandRunner runs project commands. A project command is a command
// for a specific TF project.
type ProjectCommandRunner interface {
	// Plan runs terraform plan for the project described by ctx.
	Plan(ctx models.ProjectCommandContext) ProjectResult
	// Apply runs terraform apply for the project described by ctx.
	Apply(ctx models.ProjectCommandContext) ProjectResult
}

// DefaultProjectCommandRunner implements ProjectCommandRunner.
type DefaultProjectCommandRunner struct {
	Locker                  ProjectLocker
	LockURLGenerator        LockURLGenerator
	InitStepRunner          StepRunner
	PlanStepRunner          StepRunner
	ApplyStepRunner         StepRunner
	RunStepRunner           StepRunner
	PullApprovedChecker     runtime.PullApprovedChecker
	WorkingDir              WorkingDir
	Webhooks                WebhooksSender
	WorkingDirLocker        WorkingDirLocker
	RequireApprovalOverride bool
}

// Plan runs terraform plan for the project described by ctx.
func (p *DefaultProjectCommandRunner) Plan(ctx models.ProjectCommandContext) ProjectResult {
	planSuccess, failure, err := p.doPlan(ctx)
	return ProjectResult{
		PlanSuccess: planSuccess,
		Error:       err,
		Failure:     failure,
		RepoRelDir:  ctx.RepoRelDir,
		Workspace:   ctx.Workspace,
		ProjectName: ctx.GetProjectName(),
	}
}

// Apply runs terraform apply for the project described by ctx.
func (p *DefaultProjectCommandRunner) Apply(ctx models.ProjectCommandContext) ProjectResult {
	applyOut, failure, err := p.doApply(ctx)
	return ProjectResult{
		Failure:      failure,
		Error:        err,
		ApplySuccess: applyOut,
		RepoRelDir:   ctx.RepoRelDir,
		Workspace:    ctx.Workspace,
		ProjectName:  ctx.GetProjectName(),
	}
}

func (p *DefaultProjectCommandRunner) doPlan(ctx models.ProjectCommandContext) (*PlanSuccess, string, error) {
	// Acquire Atlantis lock for this repo/dir/workspace.
	lockAttempt, err := p.Locker.TryLock(ctx.Log, ctx.Pull, ctx.User, ctx.Workspace, models.NewProject(ctx.BaseRepo.FullName, ctx.RepoRelDir))
	if err != nil {
		return nil, "", errors.Wrap(err, "acquiring lock")
	}
	if !lockAttempt.LockAcquired {
		return nil, lockAttempt.LockFailureReason, nil
	}
	ctx.Log.Debug("acquired lock for project")

	// Acquire internal lock for the directory we're going to operate in.
	unlockFn, err := p.WorkingDirLocker.TryLock(ctx.BaseRepo.FullName, ctx.Pull.Num, ctx.Workspace)
	if err != nil {
		return nil, "", err
	}
	defer unlockFn()

	// Clone is idempotent so okay to run even if the repo was already cloned.
	repoDir, cloneErr := p.WorkingDir.Clone(ctx.Log, ctx.BaseRepo, ctx.HeadRepo, ctx.Pull, ctx.RebaseRepo, ctx.Workspace)
	if cloneErr != nil {
		if unlockErr := lockAttempt.UnlockFn(); unlockErr != nil {
			ctx.Log.Err("error unlocking state after plan error: %v", unlockErr)
		}
		return nil, "", cloneErr
	}
	projAbsPath := filepath.Join(repoDir, ctx.RepoRelDir)

	// Use default stage unless another workflow is defined in config
	stage := p.defaultPlanStage()
	if ctx.ProjectConfig != nil && ctx.ProjectConfig.Workflow != nil {
		ctx.Log.Debug("project configured to use workflow %q", *ctx.ProjectConfig.Workflow)
		configuredStage := ctx.GlobalConfig.GetPlanStage(*ctx.ProjectConfig.Workflow)
		if configuredStage != nil {
			ctx.Log.Debug("project will use the configured stage for that workflow")
			stage = *configuredStage
		}
	}
	outputs, err := p.runSteps(stage.Steps, ctx, projAbsPath)
	if err != nil {
		if unlockErr := lockAttempt.UnlockFn(); unlockErr != nil {
			ctx.Log.Err("error unlocking state after plan error: %v", unlockErr)
		}
		return nil, "", fmt.Errorf("%s\n%s", err, strings.Join(outputs, "\n"))
	}

	return &PlanSuccess{
		LockURL:         p.LockURLGenerator.GenerateLockURL(lockAttempt.LockKey),
		TerraformOutput: strings.Join(outputs, "\n"),
		RePlanCmd:       ctx.RePlanCmd,
		ApplyCmd:        ctx.ApplyCmd,
	}, "", nil
}

func (p *DefaultProjectCommandRunner) runSteps(steps []valid.Step, ctx models.ProjectCommandContext, absPath string) ([]string, error) {
	var outputs []string
	for _, step := range steps {
		var out string
		var err error
		switch step.StepName {
		case "init":
			out, err = p.InitStepRunner.Run(ctx, step.ExtraArgs, absPath)
		case "plan":
			out, err = p.PlanStepRunner.Run(ctx, step.ExtraArgs, absPath)
		case "apply":
			out, err = p.ApplyStepRunner.Run(ctx, step.ExtraArgs, absPath)
		case "run":
			out, err = p.RunStepRunner.Run(ctx, step.RunCommand, absPath)
		}

		if out != "" {
			outputs = append(outputs, out)
		}
		if err != nil {
			return outputs, err
		}
	}
	return outputs, nil
}

func (p *DefaultProjectCommandRunner) doApply(ctx models.ProjectCommandContext) (applyOut string, failure string, err error) {
	repoDir, err := p.WorkingDir.GetWorkingDir(ctx.BaseRepo, ctx.Pull, ctx.Workspace)
	if err != nil {
		if os.IsNotExist(err) {
			return "", "", errors.New("project has not been cloned–did you run plan?")
		}
		return "", "", err
	}
	absPath := filepath.Join(repoDir, ctx.RepoRelDir)

	var applyRequirements []string
	if ctx.ProjectConfig != nil {
		applyRequirements = ctx.ProjectConfig.ApplyRequirements
	}
	// todo: this class shouldn't know about the server-side approval requirement.
	// Instead the project_command_builder should figure this out and store this information in the ctx. # refactor
	if p.RequireApprovalOverride {
		applyRequirements = []string{raw.ApprovedApplyRequirement}
	}
	for _, req := range applyRequirements {
		switch req {
		case raw.ApprovedApplyRequirement:
			approved, err := p.PullApprovedChecker.PullIsApproved(ctx.BaseRepo, ctx.Pull) // nolint: vetshadow
			if err != nil {
				return "", "", errors.Wrap(err, "checking if pull request was approved")
			}
			if !approved {
				return "", "Pull request must be approved before running apply.", nil
			}
		}
	}
	// Acquire internal lock for the directory we're going to operate in.
	unlockFn, err := p.WorkingDirLocker.TryLock(ctx.BaseRepo.FullName, ctx.Pull.Num, ctx.Workspace)
	if err != nil {
		return "", "", err
	}
	defer unlockFn()

	// Use default stage unless another workflow is defined in config
	stage := p.defaultApplyStage()
	if ctx.ProjectConfig != nil && ctx.ProjectConfig.Workflow != nil {
		configuredStage := ctx.GlobalConfig.GetApplyStage(*ctx.ProjectConfig.Workflow)
		if configuredStage != nil {
			stage = *configuredStage
		}
	}
	outputs, err := p.runSteps(stage.Steps, ctx, absPath)
	p.Webhooks.Send(ctx.Log, webhooks.ApplyResult{ // nolint: errcheck
		Workspace: ctx.Workspace,
		User:      ctx.User,
		Repo:      ctx.BaseRepo,
		Pull:      ctx.Pull,
		Success:   err == nil,
	})
	if err != nil {
		return "", "", fmt.Errorf("%s\n%s", err, strings.Join(outputs, "\n"))
	}
	return strings.Join(outputs, "\n"), "", nil
}

func (p DefaultProjectCommandRunner) defaultPlanStage() valid.Stage {
	return valid.Stage{
		Steps: []valid.Step{
			{
				StepName: "init",
			},
			{
				StepName: "plan",
			},
		},
	}
}

func (p DefaultProjectCommandRunner) defaultApplyStage() valid.Stage {
	return valid.Stage{
		Steps: []valid.Step{
			{
				StepName: "apply",
			},
		},
	}
}

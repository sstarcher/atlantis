// Automatically generated by pegomock. DO NOT EDIT!
// Source: github.com/runatlantis/atlantis/server/events (interfaces: WorkingDir)

package mocks

import (
	"reflect"

	pegomock "github.com/petergtz/pegomock"
	models "github.com/runatlantis/atlantis/server/events/models"
	logging "github.com/runatlantis/atlantis/server/logging"
)

type MockWorkingDir struct {
	fail func(message string, callerSkip ...int)
}

func NewMockWorkingDir() *MockWorkingDir {
	return &MockWorkingDir{fail: pegomock.GlobalFailHandler}
}

func (mock *MockWorkingDir) Clone(log *logging.SimpleLogger, baseRepo models.Repo, headRepo models.Repo, p models.PullRequest, rebase bool, workspace string) (string, error) {
	params := []pegomock.Param{log, baseRepo, headRepo, p, rebase, workspace}
	result := pegomock.GetGenericMockFrom(mock).Invoke("Clone", params, []reflect.Type{reflect.TypeOf((*string)(nil)).Elem(), reflect.TypeOf((*error)(nil)).Elem()})
	var ret0 string
	var ret1 error
	if len(result) != 0 {
		if result[0] != nil {
			ret0 = result[0].(string)
		}
		if result[1] != nil {
			ret1 = result[1].(error)
		}
	}
	return ret0, ret1
}

func (mock *MockWorkingDir) GetWorkingDir(r models.Repo, p models.PullRequest, workspace string) (string, error) {
	params := []pegomock.Param{r, p, workspace}
	result := pegomock.GetGenericMockFrom(mock).Invoke("GetWorkingDir", params, []reflect.Type{reflect.TypeOf((*string)(nil)).Elem(), reflect.TypeOf((*error)(nil)).Elem()})
	var ret0 string
	var ret1 error
	if len(result) != 0 {
		if result[0] != nil {
			ret0 = result[0].(string)
		}
		if result[1] != nil {
			ret1 = result[1].(error)
		}
	}
	return ret0, ret1
}

func (mock *MockWorkingDir) GetPullDir(r models.Repo, p models.PullRequest) (string, error) {
	params := []pegomock.Param{r, p}
	result := pegomock.GetGenericMockFrom(mock).Invoke("GetPullDir", params, []reflect.Type{reflect.TypeOf((*string)(nil)).Elem(), reflect.TypeOf((*error)(nil)).Elem()})
	var ret0 string
	var ret1 error
	if len(result) != 0 {
		if result[0] != nil {
			ret0 = result[0].(string)
		}
		if result[1] != nil {
			ret1 = result[1].(error)
		}
	}
	return ret0, ret1
}

func (mock *MockWorkingDir) Delete(r models.Repo, p models.PullRequest) error {
	params := []pegomock.Param{r, p}
	result := pegomock.GetGenericMockFrom(mock).Invoke("Delete", params, []reflect.Type{reflect.TypeOf((*error)(nil)).Elem()})
	var ret0 error
	if len(result) != 0 {
		if result[0] != nil {
			ret0 = result[0].(error)
		}
	}
	return ret0
}

func (mock *MockWorkingDir) DeleteForWorkspace(r models.Repo, p models.PullRequest, workspace string) error {
	params := []pegomock.Param{r, p, workspace}
	result := pegomock.GetGenericMockFrom(mock).Invoke("DeleteForWorkspace", params, []reflect.Type{reflect.TypeOf((*error)(nil)).Elem()})
	var ret0 error
	if len(result) != 0 {
		if result[0] != nil {
			ret0 = result[0].(error)
		}
	}
	return ret0
}

func (mock *MockWorkingDir) VerifyWasCalledOnce() *VerifierWorkingDir {
	return &VerifierWorkingDir{mock, pegomock.Times(1), nil}
}

func (mock *MockWorkingDir) VerifyWasCalled(invocationCountMatcher pegomock.Matcher) *VerifierWorkingDir {
	return &VerifierWorkingDir{mock, invocationCountMatcher, nil}
}

func (mock *MockWorkingDir) VerifyWasCalledInOrder(invocationCountMatcher pegomock.Matcher, inOrderContext *pegomock.InOrderContext) *VerifierWorkingDir {
	return &VerifierWorkingDir{mock, invocationCountMatcher, inOrderContext}
}

type VerifierWorkingDir struct {
	mock                   *MockWorkingDir
	invocationCountMatcher pegomock.Matcher
	inOrderContext         *pegomock.InOrderContext
}

func (verifier *VerifierWorkingDir) Clone(log *logging.SimpleLogger, baseRepo models.Repo, headRepo models.Repo, p models.PullRequest, workspace string) *WorkingDir_Clone_OngoingVerification {
	params := []pegomock.Param{log, baseRepo, headRepo, p, workspace}
	methodInvocations := pegomock.GetGenericMockFrom(verifier.mock).Verify(verifier.inOrderContext, verifier.invocationCountMatcher, "Clone", params)
	return &WorkingDir_Clone_OngoingVerification{mock: verifier.mock, methodInvocations: methodInvocations}
}

type WorkingDir_Clone_OngoingVerification struct {
	mock              *MockWorkingDir
	methodInvocations []pegomock.MethodInvocation
}

func (c *WorkingDir_Clone_OngoingVerification) GetCapturedArguments() (*logging.SimpleLogger, models.Repo, models.Repo, models.PullRequest, string) {
	log, baseRepo, headRepo, p, workspace := c.GetAllCapturedArguments()
	return log[len(log)-1], baseRepo[len(baseRepo)-1], headRepo[len(headRepo)-1], p[len(p)-1], workspace[len(workspace)-1]
}

func (c *WorkingDir_Clone_OngoingVerification) GetAllCapturedArguments() (_param0 []*logging.SimpleLogger, _param1 []models.Repo, _param2 []models.Repo, _param3 []models.PullRequest, _param4 []string) {
	params := pegomock.GetGenericMockFrom(c.mock).GetInvocationParams(c.methodInvocations)
	if len(params) > 0 {
		_param0 = make([]*logging.SimpleLogger, len(params[0]))
		for u, param := range params[0] {
			_param0[u] = param.(*logging.SimpleLogger)
		}
		_param1 = make([]models.Repo, len(params[1]))
		for u, param := range params[1] {
			_param1[u] = param.(models.Repo)
		}
		_param2 = make([]models.Repo, len(params[2]))
		for u, param := range params[2] {
			_param2[u] = param.(models.Repo)
		}
		_param3 = make([]models.PullRequest, len(params[3]))
		for u, param := range params[3] {
			_param3[u] = param.(models.PullRequest)
		}
		_param4 = make([]string, len(params[4]))
		for u, param := range params[4] {
			_param4[u] = param.(string)
		}
	}
	return
}

func (verifier *VerifierWorkingDir) GetWorkingDir(r models.Repo, p models.PullRequest, workspace string) *WorkingDir_GetWorkingDir_OngoingVerification {
	params := []pegomock.Param{r, p, workspace}
	methodInvocations := pegomock.GetGenericMockFrom(verifier.mock).Verify(verifier.inOrderContext, verifier.invocationCountMatcher, "GetWorkingDir", params)
	return &WorkingDir_GetWorkingDir_OngoingVerification{mock: verifier.mock, methodInvocations: methodInvocations}
}

type WorkingDir_GetWorkingDir_OngoingVerification struct {
	mock              *MockWorkingDir
	methodInvocations []pegomock.MethodInvocation
}

func (c *WorkingDir_GetWorkingDir_OngoingVerification) GetCapturedArguments() (models.Repo, models.PullRequest, string) {
	r, p, workspace := c.GetAllCapturedArguments()
	return r[len(r)-1], p[len(p)-1], workspace[len(workspace)-1]
}

func (c *WorkingDir_GetWorkingDir_OngoingVerification) GetAllCapturedArguments() (_param0 []models.Repo, _param1 []models.PullRequest, _param2 []string) {
	params := pegomock.GetGenericMockFrom(c.mock).GetInvocationParams(c.methodInvocations)
	if len(params) > 0 {
		_param0 = make([]models.Repo, len(params[0]))
		for u, param := range params[0] {
			_param0[u] = param.(models.Repo)
		}
		_param1 = make([]models.PullRequest, len(params[1]))
		for u, param := range params[1] {
			_param1[u] = param.(models.PullRequest)
		}
		_param2 = make([]string, len(params[2]))
		for u, param := range params[2] {
			_param2[u] = param.(string)
		}
	}
	return
}

func (verifier *VerifierWorkingDir) GetPullDir(r models.Repo, p models.PullRequest) *WorkingDir_GetPullDir_OngoingVerification {
	params := []pegomock.Param{r, p}
	methodInvocations := pegomock.GetGenericMockFrom(verifier.mock).Verify(verifier.inOrderContext, verifier.invocationCountMatcher, "GetPullDir", params)
	return &WorkingDir_GetPullDir_OngoingVerification{mock: verifier.mock, methodInvocations: methodInvocations}
}

type WorkingDir_GetPullDir_OngoingVerification struct {
	mock              *MockWorkingDir
	methodInvocations []pegomock.MethodInvocation
}

func (c *WorkingDir_GetPullDir_OngoingVerification) GetCapturedArguments() (models.Repo, models.PullRequest) {
	r, p := c.GetAllCapturedArguments()
	return r[len(r)-1], p[len(p)-1]
}

func (c *WorkingDir_GetPullDir_OngoingVerification) GetAllCapturedArguments() (_param0 []models.Repo, _param1 []models.PullRequest) {
	params := pegomock.GetGenericMockFrom(c.mock).GetInvocationParams(c.methodInvocations)
	if len(params) > 0 {
		_param0 = make([]models.Repo, len(params[0]))
		for u, param := range params[0] {
			_param0[u] = param.(models.Repo)
		}
		_param1 = make([]models.PullRequest, len(params[1]))
		for u, param := range params[1] {
			_param1[u] = param.(models.PullRequest)
		}
	}
	return
}

func (verifier *VerifierWorkingDir) Delete(r models.Repo, p models.PullRequest) *WorkingDir_Delete_OngoingVerification {
	params := []pegomock.Param{r, p}
	methodInvocations := pegomock.GetGenericMockFrom(verifier.mock).Verify(verifier.inOrderContext, verifier.invocationCountMatcher, "Delete", params)
	return &WorkingDir_Delete_OngoingVerification{mock: verifier.mock, methodInvocations: methodInvocations}
}

type WorkingDir_Delete_OngoingVerification struct {
	mock              *MockWorkingDir
	methodInvocations []pegomock.MethodInvocation
}

func (c *WorkingDir_Delete_OngoingVerification) GetCapturedArguments() (models.Repo, models.PullRequest) {
	r, p := c.GetAllCapturedArguments()
	return r[len(r)-1], p[len(p)-1]
}

func (c *WorkingDir_Delete_OngoingVerification) GetAllCapturedArguments() (_param0 []models.Repo, _param1 []models.PullRequest) {
	params := pegomock.GetGenericMockFrom(c.mock).GetInvocationParams(c.methodInvocations)
	if len(params) > 0 {
		_param0 = make([]models.Repo, len(params[0]))
		for u, param := range params[0] {
			_param0[u] = param.(models.Repo)
		}
		_param1 = make([]models.PullRequest, len(params[1]))
		for u, param := range params[1] {
			_param1[u] = param.(models.PullRequest)
		}
	}
	return
}

func (verifier *VerifierWorkingDir) DeleteForWorkspace(r models.Repo, p models.PullRequest, workspace string) *WorkingDir_DeleteForWorkspace_OngoingVerification {
	params := []pegomock.Param{r, p, workspace}
	methodInvocations := pegomock.GetGenericMockFrom(verifier.mock).Verify(verifier.inOrderContext, verifier.invocationCountMatcher, "DeleteForWorkspace", params)
	return &WorkingDir_DeleteForWorkspace_OngoingVerification{mock: verifier.mock, methodInvocations: methodInvocations}
}

type WorkingDir_DeleteForWorkspace_OngoingVerification struct {
	mock              *MockWorkingDir
	methodInvocations []pegomock.MethodInvocation
}

func (c *WorkingDir_DeleteForWorkspace_OngoingVerification) GetCapturedArguments() (models.Repo, models.PullRequest, string) {
	r, p, workspace := c.GetAllCapturedArguments()
	return r[len(r)-1], p[len(p)-1], workspace[len(workspace)-1]
}

func (c *WorkingDir_DeleteForWorkspace_OngoingVerification) GetAllCapturedArguments() (_param0 []models.Repo, _param1 []models.PullRequest, _param2 []string) {
	params := pegomock.GetGenericMockFrom(c.mock).GetInvocationParams(c.methodInvocations)
	if len(params) > 0 {
		_param0 = make([]models.Repo, len(params[0]))
		for u, param := range params[0] {
			_param0[u] = param.(models.Repo)
		}
		_param1 = make([]models.PullRequest, len(params[1]))
		for u, param := range params[1] {
			_param1[u] = param.(models.PullRequest)
		}
		_param2 = make([]string, len(params[2]))
		for u, param := range params[2] {
			_param2[u] = param.(string)
		}
	}
	return
}

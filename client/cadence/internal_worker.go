package cadence

// All code in this file is private to the package.

import (
	"github.com/uber-common/bark"

	m "code.uber.internal/devexp/minions-client-go.git/.gen/go/minions"
	"code.uber.internal/devexp/minions-client-go.git/common/metrics"
	"code.uber.internal/go-common.git/x/log"
)

type (
	// WorkflowWorker wraps the code for hosting workflow types.
	// And worker is mapped 1:1 with task list. If the user want's to poll multiple
	// task list names they might have to manage 'n' workers for 'n' task lists.
	workflowWorker struct {
		executionParameters WorkerExecutionParameters
		workflowDefFactory  workflowDefinitionFactory
		workflowService     m.TChanWorkflowService
		poller              taskPoller // taskPoller to poll the tasks.
		worker              *baseWorker
		identity            string
		logger              bark.Logger
	}

	// activityRegistry collection of activity implementations
	activityRegistry map[string]Activity

	// ActivityWorker wraps the code for hosting activity types.
	// TODO: Worker doing heartbeating automatically while activity task is running
	activityWorker struct {
		executionParameters WorkerExecutionParameters
		activityRegistry    activityRegistry
		workflowService     m.TChanWorkflowService
		poller              *activityTaskPoller
		worker              *baseWorker
		identity            string
		logger              bark.Logger
	}

	// Worker overrides.
	workerOverrides struct {
		workflowTaskHander  workflowTaskHandler
		activityTaskHandler activityTaskHandler
	}
)

// NewWorkflowWorker returns an instance of the workflow worker.
func newWorkflowWorker(params WorkerExecutionParameters, factory workflowDefinitionFactory,
	service m.TChanWorkflowService, logger bark.Logger,
	reporter metrics.Reporter, pressurePoints map[string]map[string]string) *workflowWorker {
	return newWorkflowWorkerInternal(params, factory, service, logger, reporter, pressurePoints, nil)
}

func newWorkflowWorkerInternal(params WorkerExecutionParameters, factory workflowDefinitionFactory,
	service m.TChanWorkflowService, logger bark.Logger, reporter metrics.Reporter,
	pressurePoints map[string]map[string]string, overrides *workerOverrides) *workflowWorker {
	// Get an identity.
	identity := params.Identity
	if identity == "" {
		identity = getWorkerIdentity(params.TaskListName)
	}

	// Get a workflow task handler.
	var taskHandler workflowTaskHandler
	if overrides != nil && overrides.workflowTaskHander != nil {
		taskHandler = overrides.workflowTaskHander
	} else {
		taskHandler = newWorkflowTaskHandler(params.TaskListName, identity, factory, logger, reporter, pressurePoints)
	}

	poller := newWorkflowTaskPoller(
		service,
		params.TaskListName,
		identity,
		taskHandler,
		logger,
		reporter)
	worker := newBaseWorker(baseWorkerOptions{
		routineCount:    params.ConcurrentPollRoutineSize,
		taskPoller:      poller,
		workflowService: service,
		identity:        identity})

	return &workflowWorker{
		executionParameters: params,
		workflowDefFactory:  factory,
		workflowService:     service,
		poller:              poller,
		worker:              worker,
		identity:            identity,
	}
}

// Start the worker.
func (ww *workflowWorker) Start() error {
	ww.worker.Start()
	return nil // TODO: propagate error
}

// Shutdown the worker.
func (ww *workflowWorker) Stop() {
	ww.worker.Stop()
}

func newActivityWorkerInternal(executionParameters WorkerExecutionParameters, activities []Activity,
	service m.TChanWorkflowService, logger bark.Logger, reporter metrics.Reporter, overrides *workerOverrides) *activityWorker {
	// Get an identity.
	identity := executionParameters.Identity
	if identity == "" {
		identity = getWorkerIdentity(executionParameters.TaskListName)
	}

	if logger == nil {
		logger = log.WithFields(log.Fields{tagTaskListName: executionParameters.TaskListName})
	}

	// Get a activity task handler.
	var taskHandler activityTaskHandler
	if overrides != nil && overrides.activityTaskHandler != nil {
		taskHandler = overrides.activityTaskHandler
	} else {
		taskHandler = newActivityTaskHandler(executionParameters.TaskListName, executionParameters.Identity,
			activities, service, logger, reporter)
	}
	poller := newActivityTaskPoller(
		service,
		executionParameters.TaskListName,
		identity,
		taskHandler,
		reporter,
		logger)
	worker := newBaseWorker(baseWorkerOptions{
		routineCount:    executionParameters.ConcurrentPollRoutineSize,
		taskPoller:      poller,
		workflowService: service,
		identity:        identity})

	return &activityWorker{
		executionParameters: executionParameters,
		activityRegistry:    make(map[string]Activity),
		workflowService:     service,
		worker:              worker,
		poller:              poller,
		identity:            identity,
	}
}

// Start the worker.
func (aw *activityWorker) Start() error {
	aw.worker.Start()
	return nil // TODO: propagate errors
}

// Shutdown the worker.
func (aw *activityWorker) Stop() {
	aw.worker.Stop()
}

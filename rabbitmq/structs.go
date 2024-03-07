package rabbitmq

type Saga struct {
	SagaName        string
	SagaDescription string
	ChannelName     string
	CorrelationId   string
	SagaSteps       []SagaStep
}

type SagaStep struct {
	StepName             string
	StepDescription      string
	StepSequence         int
	FunctionName         string
	CompensatingFunction string
	QueueName            string
	ReplyQueueName       string
	ActionType           string
	DataObject           []byte
	Results              []string
}

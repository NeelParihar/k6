package execution

// Controller implementations are used to control the k6 execution of a test or
// test suite, either locally or in a distributed environment.
type Controller interface {
	// TODO: split apart into `Once()`, `SetData(), `GetData()`?
	GetOrCreateData(id string, callback func() ([]byte, error)) ([]byte, error)

	Wait(eventId string) func() error
	Signal(eventId string) error
}

package serviceexecutor

import "context"

type Service interface {
	Run() error
	Shutdown(ctx context.Context) error
}

type Setupable interface {
	Setup() error
}

type Hooks struct {
	OnServiceRunFinished func(err error)
	OnServiceShutdownFinished func(err error)
}

type Multi struct {
	Services []Service
	Hooks Hooks
}

func (m *Multi) Run() error {
	for _, s := range m.Services {
		go func() {
			err := s.Run()
			m.Hooks.OnServiceRunFinished(err)
		}()
	}
	return nil
}

func (m *Multi) Shutdown(ctx context.Context) error {
	for i :=len(m.Services);i>=0;i-- {
		err := m.Services[i].Shutdown(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}
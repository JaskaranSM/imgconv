package manager

import (
	"imgconv/storage"
	"imgconv/utils"
	"sync"

	"github.com/google/uuid"
)

type EventType int

const (
	OnConversionStartedEvent   EventType = 0
	OnConversionCompletedEvent EventType = 1
	OnConversionFailedEvent    EventType = 2
)

const (
	TypeStatusProcessing string = "processing"
	TypeStatusCompleted  string = "completed"
	TypeStatusFailed     string = "error"
)

type AddConversionOptions struct {
	InputBytes []byte
	Format     string
	Filename   string
}

type ConversionStatus struct {
	Id       string
	Status   string
	Error    error
	Format   string
	Filename string
}

func NewConversionStatus() *ConversionStatus {
	return &ConversionStatus{}
}

type Conversion struct {
	ConversionStatus
	InputBytes []byte
}

func NewConversion() *Conversion {
	return &Conversion{}
}

type ConversionManager interface {
	AddConversion(*AddConversionOptions) (string, error)
	GetConversionStatusById(string) *ConversionStatus
	AddConversionListener(ConversionListener)
	GetOutputImageBytesById(string) ([]byte, error)
}

type Manager struct {
	mut         sync.Mutex
	conversions []*Conversion
	listeners   []ConversionListener
	storage     storage.StorageRepo
}

func (m *Manager) AddConversionListener(listener ConversionListener) {
	m.mut.Lock()
	defer m.mut.Unlock()
	m.listeners = append(m.listeners, listener)
}

func (m *Manager) getConversionById(id string) *Conversion {
	m.mut.Lock()
	defer m.mut.Unlock()
	for _, c := range m.conversions {
		if c.Id == id {
			return c
		}
	}
	return nil
}

func (m *Manager) updateConversionStatusById(id string, status string, err error) {
	conversion := m.getConversionById(id)
	conversion.Status = status
	conversion.Error = err
}

func (m *Manager) GetConversionStatusById(id string) *ConversionStatus {
	status := NewConversionStatus()
	conversion := m.getConversionById(id)
	status.Error = conversion.Error
	status.Filename = conversion.Filename
	status.Format = conversion.Format
	status.Id = conversion.Id
	status.Status = conversion.Status
	return status
}

func (m *Manager) notifyListeners(event EventType, id string, err error) {
	for _, l := range m.listeners {
		switch event {
		case OnConversionStartedEvent:
			l.OnConversionStarted(id)
		case OnConversionCompletedEvent:
			m.updateConversionStatusById(id, TypeStatusCompleted, nil)
			l.OnConversionCompleted(id)
		case OnConversionFailedEvent:
			m.updateConversionStatusById(id, TypeStatusFailed, err)
			l.OnConversionFailed(id, err)
		}
	}
}

func (m *Manager) startConversionTask(id string, inputBytes []byte) {
	m.notifyListeners(OnConversionStartedEvent, id, nil)
	outputBytes, err := utils.ConvertImageToPngBytes(inputBytes)
	if err != nil {
		m.notifyListeners(OnConversionFailedEvent, id, err)
		return
	}
	err = m.storage.SaveImageBytesId(id, outputBytes)
	if err != nil {
		m.notifyListeners(OnConversionFailedEvent, id, err)
		return
	}
	m.notifyListeners(OnConversionCompletedEvent, id, nil)
}

func (m *Manager) AddConversion(opts *AddConversionOptions) (string, error) {
	id := uuid.New().String()
	conversion := NewConversion()
	conversion.Id = id
	conversion.InputBytes = opts.InputBytes
	conversion.Filename = opts.Filename
	conversion.Format = opts.Format
	conversion.Status = TypeStatusProcessing
	m.mut.Lock()
	m.conversions = append(m.conversions, conversion)
	m.mut.Unlock()
	go m.startConversionTask(id, conversion.InputBytes)
	return id, nil
}

func (m *Manager) GetOutputImageBytesById(id string) ([]byte, error) {
	return m.storage.GetImageBytesById(id)
}

func NewConversionManager(storage storage.StorageRepo) *Manager {
	return &Manager{
		storage: storage,
	}
}

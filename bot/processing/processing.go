package processing

import "sync"

type ProcessingState struct {
	ForwardedMessages   map[int64][]string
	Emails              map[int64][]string
	IsProcessing        map[int64]bool
	InaccessibleClosed  map[int64][]string
	InaccessibleNotInDB map[int64][]string
	mutex               sync.Mutex // Защита от конкурентного доступа
}

func NewProcessingState() *ProcessingState {
	return &ProcessingState{
		ForwardedMessages:   make(map[int64][]string),
		Emails:              make(map[int64][]string),
		IsProcessing:        make(map[int64]bool),
		InaccessibleClosed:  make(map[int64][]string),
		InaccessibleNotInDB: make(map[int64][]string),
	}
}

func (ps *ProcessingState) AddToClosedAccount(userID int64, name string) {
	ps.InaccessibleClosed[userID] = append(ps.InaccessibleClosed[userID], name)
}

func (ps *ProcessingState) AddToNotInDB(userID int64, name string) {
	ps.InaccessibleNotInDB[userID] = append(ps.InaccessibleNotInDB[userID], name)
}

func (ps *ProcessingState) ClearUserData(userID int64) {
	ps.mutex.Lock()
	defer ps.mutex.Unlock()
	delete(ps.ForwardedMessages, userID)
	delete(ps.Emails, userID)
	delete(ps.InaccessibleClosed, userID)
	delete(ps.InaccessibleNotInDB, userID)
}

package velocitylimit

import "sync"

type Store struct {
	sync.RWMutex
	accounts map[string]*Account
	txns     map[string]bool
}

func NewStore() *Store {
	return &Store{
		accounts: make(map[string]*Account),
		txns:     make(map[string]bool),
	}
}

func (s *Store) GetAccount(custID string) *Account {
	acc := s.getAccountFromStore(custID)
	if acc != nil {
		return acc
	}
	return s.addAccountToStore(custID)
}

func (s *Store) getAccountFromStore(custID string) *Account {
	s.RLock()
	defer s.RUnlock()
	if acc, ok := s.accounts[custID]; ok {
		return acc
	}
	return nil
}

func (s *Store) addAccountToStore(custID string) *Account {
	s.Lock()
	defer s.Unlock()
	acc := NewAccount(custID)
	s.accounts[custID] = acc
	return acc
}

func (s *Store) AddTxn(id, custID string) {
	s.Lock()
	defer s.Unlock()
	s.txns[id+custID] = true
}

func (s *Store) IsDupTxn(id, custID string) bool {
	s.RLock()
	defer s.RUnlock()
	_, ok := s.txns[id+custID]
	return ok
}

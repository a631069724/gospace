package manager

import (
	"errors"
	"fmt"
	"log"
	"sync"
)

type keeper struct {
	name     string
	accounts map[string]*Account
	sync.Mutex
}

func NewKeeper(name string) *keeper {
	return &keeper{
		name:     name,
		accounts: make(map[string]*Account),
	}
}

func (k *keeper) Name() string {
	return k.name
}

func (k *keeper) DelAccount(name string) {
	k.Lock()
	defer k.Unlock()
	delete(k.accounts, name)
}

func (k *keeper) AddAccount(acc *Account) error {
	k.Lock()
	defer k.Unlock()
	if acc == nil {
		return errors.New("account is nil")
	}
	name := acc.Name()
	if name == "" {
		return errors.New("account name is nil")
	}
	if _, ok := k.accounts[name]; ok {
		return errors.New(fmt.Sprintf("account %s has exist", name))
	}
	k.accounts[name] = acc
	return nil
}

func (k *keeper) GetAccounts() map[string]*Account {
	return k.accounts
}

func (k *keeper) GetAccount(account string) *Account {
	k.Lock()
	defer k.Unlock()
	return k.accounts[account]
}

func (k *keeper) AccountRun() {
	k.Lock()
	defer k.Unlock()
	accounts := k.GetAccounts()
	for _, acc := range accounts {
		Wg.Add(1)
		go func() {
			defer Wg.Done()
			if err := acc.Run(); err != nil {
				log.Println("Account Running error: ", err.Error())
			}
		}()

	}
}

func (k *keeper) AccountExit() {
	k.Lock()
	defer k.Unlock()
	accounts := k.GetAccounts()
	for _, acc := range accounts {
		acc.Exit()
	}
}

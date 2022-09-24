package inmemory

import (
	"context"
	"dk-go-gophkeeper/internal/client/grpcclient"
	"dk-go-gophkeeper/internal/client/storage/modelstroage"
	"fmt"
	"golang.org/x/sync/errgroup"
	"sync"
)

type Storage struct {
	mu              sync.Mutex
	userID          string
	bankCardDB      map[string]modelstroage.BankCard
	loginPasswordDB map[string]modelstroage.LoginAndPassword
	textBinaryDB    map[string]modelstroage.TextOrBinary
	clientGRPC      grpcclient.GRPCClient
}

func InitStorage() *Storage {
	bankCardDB := make(map[string]modelstroage.BankCard)
	loginPasswordDB := make(map[string]modelstroage.LoginAndPassword)
	textBinaryDB := make(map[string]modelstroage.TextOrBinary)
	st := Storage{
		bankCardDB:      bankCardDB,
		loginPasswordDB: loginPasswordDB,
		textBinaryDB:    textBinaryDB,
	}
	return &st
}

func (s *Storage) AddBankCard(identifier, number, holder, cvv, meta string) error {
	newBankCardEntry := modelstroage.BankCard{
		Identifier: identifier,
		Number:     number,
		Holder:     holder,
		Cvv:        cvv,
		Meta:       meta,
	}
	_, ok := s.bankCardDB[identifier]
	if ok {
		return fmt.Errorf("entry of type 'Bank Card' with ID %s already exists", identifier)
	}
	s.bankCardDB[identifier] = newBankCardEntry
	return nil
}

func (s *Storage) AddLoginPassword(identifier, login, password, meta string) error {
	newLoginPasswordEntry := modelstroage.LoginAndPassword{
		Identifier: identifier,
		Login:      login,
		Password:   password,
		Meta:       meta,
	}
	_, ok := s.loginPasswordDB[identifier]
	if ok {
		return fmt.Errorf("entry of type 'Login And Passowrd' with ID %s already exists", identifier)
	}
	s.loginPasswordDB[identifier] = newLoginPasswordEntry
	return nil
}

func (s *Storage) AddTextBinary(identifier, entry, meta string) error {
	newTextBinaryEntry := modelstroage.TextOrBinary{
		Identifier: identifier,
		Entry:      entry,
		Meta:       meta,
	}
	_, ok := s.textBinaryDB[identifier]
	if ok {
		return fmt.Errorf("entry of type 'Text Or Binary' with ID %s already exists", identifier)
	}
	s.textBinaryDB[identifier] = newTextBinaryEntry
	return nil
}

func (s *Storage) Sync(ctx context.Context) error {
	grp, ctx := errgroup.WithContext(ctx)
	funcs := []func() error{s.dumpTextsBinaries, s.dumpLoginsPasswords, s.dumpBankCards}
	for _, fn := range funcs {
		grp.Go(fn)
	}
	if err := grp.Wait(); err != nil {
		return err
	}
	err := s.sendTextsBinaries(ctx)
	if err != nil {
		return err
	}
	err = s.sendLoginsPasswords(ctx)
	if err != nil {
		return err
	}
	err = s.sendBankCards(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) dumpBankCards() error {
	cloudDataBankCards, err := s.clientGRPC.GetBankCards()
	if err != nil {
		return err
	}
	for identifier, value := range cloudDataBankCards {
		// overwrite any local data with cloud data
		s.bankCardDB[identifier] = value
	}
	return nil
}

func (s *Storage) sendBankCards(ctx context.Context) error {
	grp, ctx := errgroup.WithContext(ctx)
	for _, value := range s.bankCardDB {
		grp.Go(func() error {
			err := s.clientGRPC.SendBankCard(value)
			if err != nil {
				return err
			}
			return nil
		})
	}
	if err := grp.Wait(); err != nil {
		return err
	}
	return nil
}

func (s *Storage) dumpLoginsPasswords() error {
	cloudDataLoginsPasswords, err := s.clientGRPC.GetLoginsPasswords()
	if err != nil {
		return err
	}
	for identifier, value := range cloudDataLoginsPasswords {
		// overwrite any local data with cloud data
		s.loginPasswordDB[identifier] = value
	}
	return nil
}

func (s *Storage) sendLoginsPasswords(ctx context.Context) error {
	grp, ctx := errgroup.WithContext(ctx)
	for _, value := range s.loginPasswordDB {
		grp.Go(func() error {
			err := s.clientGRPC.SendLoginPassword(value)
			if err != nil {
				return err
			}
			return nil
		})
	}
	if err := grp.Wait(); err != nil {
		return err
	}
	return nil
}

func (s *Storage) dumpTextsBinaries() error {
	cloudDataTextsBinaries, err := s.clientGRPC.GetTextsBinaries()
	if err != nil {
		return err
	}
	for identifier, value := range cloudDataTextsBinaries {
		// overwrite any local data with cloud data
		s.textBinaryDB[identifier] = value
	}
	return nil
}

func (s *Storage) sendTextsBinaries(ctx context.Context) error {
	grp, ctx := errgroup.WithContext(ctx)
	for _, value := range s.textBinaryDB {
		grp.Go(func() error {
			err := s.clientGRPC.SendTextBinary(value)
			if err != nil {
				return err
			}
			return nil
		})
	}
	if err := grp.Wait(); err != nil {
		return err
	}
	return nil
}

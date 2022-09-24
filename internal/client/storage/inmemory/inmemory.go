package inmemory

import (
	"context"
	"dk-go-gophkeeper/internal/client/grpcclient"
	"dk-go-gophkeeper/internal/client/storage"
	"dk-go-gophkeeper/internal/client/storage/modelstorage"
	"fmt"
	"golang.org/x/sync/errgroup"
	"log"
)

var (
	_ storage.DataStorage = (*Storage)(nil)
)

type Storage struct {
	bankCardDB      map[string]modelstorage.BankCard
	loginPasswordDB map[string]modelstorage.LoginAndPassword
	textBinaryDB    map[string]modelstorage.TextOrBinary
	clientGRPC      grpcclient.GRPCClient
	logger          *log.Logger
}

func InitStorage(logger *log.Logger) *Storage {
	bankCardDB := make(map[string]modelstorage.BankCard)
	loginPasswordDB := make(map[string]modelstorage.LoginAndPassword)
	textBinaryDB := make(map[string]modelstorage.TextOrBinary)
	st := Storage{
		bankCardDB:      bankCardDB,
		loginPasswordDB: loginPasswordDB,
		textBinaryDB:    textBinaryDB,
		logger:          logger,
	}
	return &st
}

func (s *Storage) Remove(identifier string) (string, error) {
	status := ""
	_, ok := s.bankCardDB[identifier]
	if ok {
		s.logger.Print("Removed entry from bank card storage:", identifier)
		status += fmt.Sprintf("removed entry %s from bank card storage\n", identifier)
		delete(s.bankCardDB, identifier)
	}
	_, ok = s.loginPasswordDB[identifier]
	if ok {
		s.logger.Print("Removed entry from login/password storage:", identifier)
		status += fmt.Sprintf("removed entry %s from login/password storage\n", identifier)
		delete(s.loginPasswordDB, identifier)
	}
	_, ok = s.textBinaryDB[identifier]
	if ok {
		s.logger.Print("Removed entry from text/binary storage:", identifier)
		status += fmt.Sprintf("removed entry %s from text/binary storage\n", identifier)
		delete(s.textBinaryDB, identifier)
	}
	if status == "" {
		s.logger.Print("Not found in storage:", identifier)
		return status, fmt.Errorf("identifier %s was not found in storage", identifier)
	}
	return status, nil
}

func (s *Storage) AddBankCard(identifier, number, holder, cvv, meta string) error {
	newBankCardEntry := modelstorage.BankCard{
		Identifier: identifier,
		Number:     number,
		Holder:     holder,
		Cvv:        cvv,
		Meta:       meta,
	}
	_, ok := s.bankCardDB[identifier]
	if ok {
		s.logger.Print("Unique violation in bank card storage for ID:", identifier)
		return fmt.Errorf("entry of type 'Bank Card' with ID %s already exists", identifier)
	}
	s.bankCardDB[identifier] = newBankCardEntry
	s.logger.Print("Added to bank card storage:", newBankCardEntry)
	//err := s.clientGRPC.SendBankCard(newBankCardEntry)
	//if err != nil {
	//	delete(s.bankCardDB, identifier)
	//	s.logger.Print("Could not upload bank card entry:", err.Error())
	//}
	return nil
}

func (s *Storage) AddLoginPassword(identifier, login, password, meta string) error {
	newLoginPasswordEntry := modelstorage.LoginAndPassword{
		Identifier: identifier,
		Login:      login,
		Password:   password,
		Meta:       meta,
	}
	_, ok := s.loginPasswordDB[identifier]
	if ok {
		s.logger.Print("Unique violation in login/password storage for ID:", identifier)
		return fmt.Errorf("entry of type 'Login And Passowrd' with ID %s already exists", identifier)
	}
	s.loginPasswordDB[identifier] = newLoginPasswordEntry
	s.logger.Print("Added to login/password storage:", newLoginPasswordEntry)
	//err := s.clientGRPC.SendLoginPassword(newLoginPasswordEntry)
	//if err != nil {
	//	delete(s.loginPasswordDB, identifier)
	//	s.logger.Print("Could not upload login.password entry:", err.Error())
	//}
	return nil
}

func (s *Storage) AddTextBinary(identifier, entry, meta string) error {
	newTextBinaryEntry := modelstorage.TextOrBinary{
		Identifier: identifier,
		Entry:      entry,
		Meta:       meta,
	}
	_, ok := s.textBinaryDB[identifier]
	if ok {
		s.logger.Print("Unique violation in text/binary storage for ID:", identifier)
		return fmt.Errorf("entry of type 'Text Or Binary' with ID %s already exists", identifier)
	}
	s.textBinaryDB[identifier] = newTextBinaryEntry
	s.logger.Print("Added to text/binary storage:", newTextBinaryEntry)
	//err := s.clientGRPC.SendTextBinary(newTextBinaryEntry)
	//if err != nil {
	//	delete(s.textBinaryDB, identifier)
	//	s.logger.Print("Could not upload text/binary entry:", err.Error())
	//}
	return nil
}

func (s *Storage) Sync(ctx context.Context) error {
	s.logger.Print("Attempting sync")
	grp, ctx := errgroup.WithContext(ctx)
	funcs := []func() error{s.dumpTextsBinaries, s.dumpLoginsPasswords, s.dumpBankCards}
	for _, fn := range funcs {
		grp.Go(fn)
	}
	if err := grp.Wait(); err != nil {
		return err
	}
	s.logger.Print("Sync performed successfully")
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

func (s *Storage) ShowAllData() string {
	data := ""
	for _, value := range s.textBinaryDB {
		data += fmt.Sprintf("%#v", value) + "\n"
	}
	for _, value := range s.loginPasswordDB {
		data += fmt.Sprintf("%#v", value) + "\n"
	}
	for _, value := range s.bankCardDB {
		data += fmt.Sprintf("%#v", value) + "\n"
	}
	return data
}

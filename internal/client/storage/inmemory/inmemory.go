// Package inmemory provides local client data storing functionality.
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

// check for interface compliance
var (
	_ storage.DataStorage = (*Storage)(nil)
)

// Storage defines atrributes and methods of a Storage instance.
type Storage struct {
	bankCardDB      map[string]modelstorage.BankCard
	loginPasswordDB map[string]modelstorage.LoginAndPassword
	textBinaryDB    map[string]modelstorage.TextOrBinary
	clientGRPC      grpcclient.GRPCClient
	logger          *log.Logger
}

// InitStorage initializes a Storage instance.
func InitStorage(logger *log.Logger, client grpcclient.GRPCClient) *Storage {
	logger.Print("Attempting to initialize storage")
	bankCardDB := make(map[string]modelstorage.BankCard)
	loginPasswordDB := make(map[string]modelstorage.LoginAndPassword)
	textBinaryDB := make(map[string]modelstorage.TextOrBinary)
	st := Storage{
		bankCardDB:      bankCardDB,
		loginPasswordDB: loginPasswordDB,
		textBinaryDB:    textBinaryDB,
		logger:          logger,
		clientGRPC:      client,
	}
	return &st
}

// Remove deletes data from local storage and sends delete requests to the server.
func (s *Storage) Remove(identifier, db string) error {
	var err error
	switch db {
	case "bankCard":
		_, ok := s.bankCardDB[identifier]
		if ok {
			s.logger.Print("Removing entry from bank card storage:", identifier)
			_, err_ := s.clientGRPC.RemoveBankCard(identifier)
			if err_ != nil {
				s.logger.Print("Could not remove bank card entry:", err_.Error())
				return err_
			}
			delete(s.bankCardDB, identifier)
		} else {
			err = fmt.Errorf("entry ID %s in %s storage does not exist", identifier, db)
		}
	case "loginPassword":
		_, ok := s.loginPasswordDB[identifier]
		if ok {
			s.logger.Print("Removing entry from login/password storage:", identifier)
			_, err_ := s.clientGRPC.RemoveLoginPassword(identifier)
			if err_ != nil {
				s.logger.Print("Could not remove login/password entry:", err_.Error())
				return err_
			}
			delete(s.loginPasswordDB, identifier)
		} else {
			err = fmt.Errorf("entry ID %s in %s storage does not exist", identifier, db)
		}
	case "textBinary":
		_, ok := s.textBinaryDB[identifier]
		if ok {
			s.logger.Print("Removing entry from text/binary storage:", identifier)
			_, err_ := s.clientGRPC.RemoveTextBinary(identifier)
			if err_ != nil {
				s.logger.Print("Could not remove text/binary entry:", err_.Error())
				return err_
			}
			delete(s.textBinaryDB, identifier)
		} else {
			err = fmt.Errorf("entry ID %s in %s storage does not exist", identifier, db)
		}
	default:
		err = fmt.Errorf("invalid db %s", db)
	}
	return err
}

// Login sends a login request to the server and cleans local DB upon successful response.
func (s *Storage) Login(login, password string) error {
	newLoginRegisterEntry := modelstorage.RegisterLogin{
		Login:    login,
		Password: password,
	}
	_, err := s.clientGRPC.Login(newLoginRegisterEntry)
	if err != nil {
		s.logger.Print("Could not perform login request:", err.Error())
		return err
	}
	s.CleanDB()
	return nil
}

// Register sends a register request to the server and cleans local DB upon successful response.
func (s *Storage) Register(login, password string) error {
	newLoginRegisterEntry := modelstorage.RegisterLogin{
		Login:    login,
		Password: password,
	}
	_, err := s.clientGRPC.Register(newLoginRegisterEntry)
	if err != nil {
		s.logger.Print("Could not perform register request:", err.Error())
		return err
	}
	s.CleanDB()
	return nil
}

// AddBankCard adds a new bank card entry to the local client storage and sends it to the server.
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
	_, err := s.clientGRPC.SendBankCard(newBankCardEntry)
	if err != nil {
		delete(s.bankCardDB, identifier)
		s.logger.Print("Could not upload bank card entry:", err.Error())
		return err
	}
	return nil
}

// AddLoginPassword adds a new login/password entry to the local client storage and sends it to the server.
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
	_, err := s.clientGRPC.SendLoginPassword(newLoginPasswordEntry)
	if err != nil {
		delete(s.loginPasswordDB, identifier)
		s.logger.Print("Could not upload login.password entry:", err.Error())
		return err
	}
	return nil
}

// AddTextBinary adds a new text/binary entry to the local client storage and sends it to the server.
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
	_, err := s.clientGRPC.SendTextBinary(newTextBinaryEntry)
	if err != nil {
		delete(s.textBinaryDB, identifier)
		s.logger.Print("Could not upload text/binary entry:", err.Error())
		return err
	}
	return nil
}

// Sync performs retrieval of all data from server overwriting local storage.
func (s *Storage) Sync() error {
	s.logger.Print("Attempting sync")
	grp, _ := errgroup.WithContext(context.Background())
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

// dumpBankCards retrieves bank card entries from server.
func (s *Storage) dumpBankCards() error {
	cloudDataBankCards, _, err := s.clientGRPC.GetBankCards()
	if err != nil {
		return err
	}
	for identifier, value := range cloudDataBankCards {
		// overwrite any local data with cloud data
		s.bankCardDB[identifier] = value
	}
	return nil
}

// dumpLoginsPasswords retrieves login/password entries from server.
func (s *Storage) dumpLoginsPasswords() error {
	cloudDataLoginsPasswords, _, err := s.clientGRPC.GetLoginsPasswords()
	if err != nil {
		return err
	}
	for identifier, value := range cloudDataLoginsPasswords {
		// overwrite any local data with cloud data
		s.loginPasswordDB[identifier] = value
	}
	return nil
}

// dumpTextsBinaries retrieves text/binary entries from server.
func (s *Storage) dumpTextsBinaries() error {
	cloudDataTextsBinaries, _, err := s.clientGRPC.GetTextsBinaries()
	if err != nil {
		return err
	}
	for identifier, value := range cloudDataTextsBinaries {
		// overwrite any local data with cloud data
		s.textBinaryDB[identifier] = value
	}
	return nil
}

// Get retrieves a data piece from local storage.
func (s *Storage) Get(identifier, db string) (string, error) {
	var err error
	var data string
	switch db {
	case "bankCard":
		value, ok := s.bankCardDB[identifier]
		if ok {
			data = fmt.Sprintf("%#v", value) + "\n"
		} else {
			err = fmt.Errorf("entry ID %s in %s storage does not exist", identifier, db)
		}
	case "loginPassword":
		value, ok := s.loginPasswordDB[identifier]
		if ok {
			data = fmt.Sprintf("%#v", value) + "\n"
		} else {
			err = fmt.Errorf("entry ID %s in %s storage does not exist", identifier, db)
		}
	case "textBinary":
		value, ok := s.textBinaryDB[identifier]
		if ok {
			data = fmt.Sprintf("%#v", value) + "\n"
		} else {
			err = fmt.Errorf("entry ID %s in %s storage does not exist", identifier, db)
		}
	default:
		err = fmt.Errorf("invalid db %s", db)
	}
	return data, err
}

// CleanDB re-initializes a local DB.
func (s *Storage) CleanDB() {
	bankCardDB := make(map[string]modelstorage.BankCard)
	loginPasswordDB := make(map[string]modelstorage.LoginAndPassword)
	textBinaryDB := make(map[string]modelstorage.TextOrBinary)
	s.bankCardDB = bankCardDB
	s.loginPasswordDB = loginPasswordDB
	s.textBinaryDB = textBinaryDB
}

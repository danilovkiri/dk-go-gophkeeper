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

func (s *Storage) Remove(identifier, db string) error {
	var err error
	switch db {
	case "bankCard":
		_, ok := s.bankCardDB[identifier]
		if ok {
			s.logger.Print("Removing entry from bank card storage:", identifier)
			_, err := s.clientGRPC.RemoveBankCard(identifier)
			if err != nil {
				s.logger.Print("Could not remove bank card entry:", err.Error())
				return err
			}
			delete(s.bankCardDB, identifier)
		} else {
			err = fmt.Errorf("entry ID %s in %s storage does not exist", identifier, db)
		}
	case "loginPassword":
		_, ok := s.loginPasswordDB[identifier]
		if ok {
			s.logger.Print("Removing entry from login/password storage:", identifier)
			_, err := s.clientGRPC.RemoveLoginPassword(identifier)
			if err != nil {
				s.logger.Print("Could not remove login/password entry:", err.Error())
				return err
			}
			delete(s.loginPasswordDB, identifier)
		} else {
			err = fmt.Errorf("entry ID %s in %s storage does not exist", identifier, db)
		}
	case "textBinary":
		_, ok := s.textBinaryDB[identifier]
		if ok {
			s.logger.Print("Removing entry from text/binary storage:", identifier)
			_, err := s.clientGRPC.RemoveTextBinary(identifier)
			if err != nil {
				s.logger.Print("Could not remove text/binary entry:", err.Error())
				return err
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

func (s *Storage) LoginRegister(login, password string) error {
	newLoginRegisterEntry := modelstorage.RegisterLogin{
		Login:    login,
		Password: password,
	}
	_, err := s.clientGRPC.LoginRegister(newLoginRegisterEntry)
	if err != nil {
		s.logger.Print("Could not perform login/register request:", err.Error())
		return err
	}
	s.CleanDB()
	return nil
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
	_, err := s.clientGRPC.SendBankCard(newBankCardEntry)
	if err != nil {
		delete(s.bankCardDB, identifier)
		s.logger.Print("Could not upload bank card entry:", err.Error())
		return err
	}
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
	_, err := s.clientGRPC.SendLoginPassword(newLoginPasswordEntry)
	if err != nil {
		delete(s.loginPasswordDB, identifier)
		s.logger.Print("Could not upload login.password entry:", err.Error())
		return err
	}
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
	_, err := s.clientGRPC.SendTextBinary(newTextBinaryEntry)
	if err != nil {
		delete(s.textBinaryDB, identifier)
		s.logger.Print("Could not upload text/binary entry:", err.Error())
		return err
	}
	return nil
}

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

func (s *Storage) CleanDB() {
	bankCardDB := make(map[string]modelstorage.BankCard)
	loginPasswordDB := make(map[string]modelstorage.LoginAndPassword)
	textBinaryDB := make(map[string]modelstorage.TextOrBinary)
	s.bankCardDB = bankCardDB
	s.loginPasswordDB = loginPasswordDB
	s.textBinaryDB = textBinaryDB
}

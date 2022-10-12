// Package inmemory provides local client data storing functionality.
package inmemory

import (
	"context"
	"dk-go-gophkeeper/internal/client/grpcclient"
	"dk-go-gophkeeper/internal/client/storage"
	"dk-go-gophkeeper/internal/client/storage/modelstorage"
	"dk-go-gophkeeper/internal/config"
	"errors"
	"fmt"

	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
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
	logger          *zerolog.Logger
	cfg             *config.Config
}

// InitStorage initializes a Storage instance.
func InitStorage(logger *zerolog.Logger, client grpcclient.GRPCClient, cfg *config.Config) *Storage {
	logger.Info().Msg("Attempting to initialize storage")
	bankCardDB := make(map[string]modelstorage.BankCard)
	loginPasswordDB := make(map[string]modelstorage.LoginAndPassword)
	textBinaryDB := make(map[string]modelstorage.TextOrBinary)
	st := Storage{
		bankCardDB:      bankCardDB,
		loginPasswordDB: loginPasswordDB,
		textBinaryDB:    textBinaryDB,
		logger:          logger,
		clientGRPC:      client,
		cfg:             cfg,
	}
	return &st
}

// Remove deletes data from local storage and sends delete requests to the server.
func (s *Storage) Remove(identifier, db string) error {
	var err error
	if identifier == "" {
		return fmt.Errorf("identifier cannot be empty in db %s", db)
	}
	switch db {
	case s.cfg.BankCardDB:
		_, ok := s.bankCardDB[identifier]
		if ok {
			s.logger.Info().Msgf("Removing entry from bank card storage: %s", identifier)
			_, err_ := s.clientGRPC.RemoveBankCard(identifier)
			if err_ != nil {
				s.logger.Error().Err(err_).Msg("Could not remove bank card entry")
				return err_
			}
			delete(s.bankCardDB, identifier)
		} else {
			err = fmt.Errorf("entry ID %s in %s storage does not exist", identifier, db)
		}
	case s.cfg.LoginPasswordDB:
		_, ok := s.loginPasswordDB[identifier]
		if ok {
			s.logger.Info().Msgf("Removing entry from login/password storage: %s", identifier)
			_, err_ := s.clientGRPC.RemoveLoginPassword(identifier)
			if err_ != nil {
				s.logger.Error().Err(err_).Msg("Could not remove login/password entry")
				return err_
			}
			delete(s.loginPasswordDB, identifier)
		} else {
			err = fmt.Errorf("entry ID %s in %s storage does not exist", identifier, db)
		}
	case s.cfg.TextBinaryDB:
		_, ok := s.textBinaryDB[identifier]
		if ok {
			s.logger.Info().Msgf("Removing entry from text/binary storage: %s", identifier)
			_, err_ := s.clientGRPC.RemoveTextBinary(identifier)
			if err_ != nil {
				s.logger.Error().Err(err_).Msg("Could not remove text/binary entry")
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
	if login == "" || password == "" {
		return errors.New("Login/Password fields cannot be empty")
	}
	newLoginRegisterEntry := modelstorage.RegisterLogin{
		Login:    login,
		Password: password,
	}
	_, err := s.clientGRPC.Login(newLoginRegisterEntry)
	if err != nil {
		s.logger.Error().Err(err).Msg("Could not perform login request")
		return err
	}
	s.CleanDB()
	return nil
}

// Register sends a register request to the server and cleans local DB upon successful response.
func (s *Storage) Register(login, password string) error {
	if login == "" || password == "" {
		return errors.New("Login/Password fields cannot be empty")
	}
	newLoginRegisterEntry := modelstorage.RegisterLogin{
		Login:    login,
		Password: password,
	}
	_, err := s.clientGRPC.Register(newLoginRegisterEntry)
	if err != nil {
		s.logger.Error().Err(err).Msg("Could not perform register request")
		return err
	}
	s.CleanDB()
	return nil
}

// AddBankCard adds a new bank card entry to the local client storage and sends it to the server.
func (s *Storage) AddBankCard(identifier, number, holder, cvv, meta string) error {
	if identifier == "" {
		return errors.New("identifier cannot be empty")
	}
	newBankCardEntry := modelstorage.BankCard{
		Identifier: identifier,
		Number:     number,
		Holder:     holder,
		Cvv:        cvv,
		Meta:       meta,
	}
	_, ok := s.bankCardDB[identifier]
	if ok {
		s.logger.Error().Msgf("Unique violation in bank card storage for ID: %s", identifier)
		return fmt.Errorf("entry of type 'Bank Card' with ID %s already exists", identifier)
	}
	s.bankCardDB[identifier] = newBankCardEntry
	s.logger.Info().Msgf("Added to bank card storage: %v", newBankCardEntry)
	_, err := s.clientGRPC.SendBankCard(newBankCardEntry)
	if err != nil {
		delete(s.bankCardDB, identifier)
		s.logger.Error().Err(err).Msg("Could not upload bank card entry")
		return err
	}
	return nil
}

// AddLoginPassword adds a new login/password entry to the local client storage and sends it to the server.
func (s *Storage) AddLoginPassword(identifier, login, password, meta string) error {
	if identifier == "" {
		return errors.New("identifier cannot be empty")
	}
	newLoginPasswordEntry := modelstorage.LoginAndPassword{
		Identifier: identifier,
		Login:      login,
		Password:   password,
		Meta:       meta,
	}
	_, ok := s.loginPasswordDB[identifier]
	if ok {
		s.logger.Error().Msgf("Unique violation in login/password storage for ID: %s", identifier)
		return fmt.Errorf("entry of type 'Login And Password' with ID %s already exists", identifier)
	}
	s.loginPasswordDB[identifier] = newLoginPasswordEntry
	s.logger.Info().Msgf("Added to login/password storage: %v", newLoginPasswordEntry)
	_, err := s.clientGRPC.SendLoginPassword(newLoginPasswordEntry)
	if err != nil {
		delete(s.loginPasswordDB, identifier)
		s.logger.Error().Err(err).Msg("Could not upload login/password entry")
		return err
	}
	return nil
}

// AddTextBinary adds a new text/binary entry to the local client storage and sends it to the server.
func (s *Storage) AddTextBinary(identifier, entry, meta string) error {
	if identifier == "" {
		return errors.New("identifier cannot be empty")
	}
	newTextBinaryEntry := modelstorage.TextOrBinary{
		Identifier: identifier,
		Entry:      entry,
		Meta:       meta,
	}
	_, ok := s.textBinaryDB[identifier]
	if ok {
		s.logger.Error().Msgf("Unique violation in text/binary storage for ID: %s", identifier)
		return fmt.Errorf("entry of type 'Text Or Binary' with ID %s already exists", identifier)
	}
	s.textBinaryDB[identifier] = newTextBinaryEntry
	s.logger.Info().Msgf("Added to text/binary storage: %v", newTextBinaryEntry)
	_, err := s.clientGRPC.SendTextBinary(newTextBinaryEntry)
	if err != nil {
		delete(s.textBinaryDB, identifier)
		s.logger.Error().Err(err).Msg("Could not upload text/binary entry")
		return err
	}
	return nil
}

// Sync performs retrieval of all data from server overwriting local storage.
func (s *Storage) Sync() error {
	s.logger.Info().Msg("Attempting sync")
	grp, _ := errgroup.WithContext(context.Background())
	funcs := []func() error{s.dumpTextsBinaries, s.dumpLoginsPasswords, s.dumpBankCards}
	for _, fn := range funcs {
		grp.Go(fn)
	}
	if err := grp.Wait(); err != nil {
		return err
	}
	s.logger.Info().Msg("Sync performed successfully")
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
	if identifier == "" {
		return "", fmt.Errorf("identifier cannot be empty in db %s", db)
	}
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

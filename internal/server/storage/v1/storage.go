package storage

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"database/sql"
	"dk-go-gophkeeper/internal/config"
	"dk-go-gophkeeper/internal/server/storage"
	storageErrors "dk-go-gophkeeper/internal/server/storage/errors"
	"dk-go-gophkeeper/internal/server/storage/modelstorage"
	"errors"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/lib/pq"
	"log"
	"sync"
	"time"
)

var (
	_ storage.DataStorage = (*Storage)(nil)
)

type Storage struct {
	mu     sync.Mutex
	cfg    *config.Config
	DB     *sql.DB
	logger *log.Logger
	ch     chan modelstorage.Removal
}

func InitStorage(ctx context.Context, logger *log.Logger, cfg *config.Config, wg *sync.WaitGroup) *Storage {
	logger.Print("Attempting to initialize storage")
	db, err := sql.Open("pgx", cfg.DatabaseDSN)
	if err != nil {
		logger.Fatal(err)
	}
	st := Storage{
		cfg:    cfg,
		logger: logger,
		DB:     db,
	}
	err = st.createTables(ctx)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Print("PSQL DB connection was established")

	const flushPartsAmount = 10
	const flushPartsInterval = time.Second * 10

	wg.Add(1)
	go func() {
		defer wg.Done()
		t := time.NewTicker(flushPartsInterval)
		parts := make([]modelstorage.Removal, 0, flushPartsAmount)
		for {
			select {
			case <-ctx.Done():
				if len(parts) > 0 {
					logger.Print("Deleting data due to context cancellation", parts)
					err := st.Flush(ctx, parts)
					if err != nil {
						logger.Fatal(err)
					}
				}
				close(st.ch)
				err := st.DB.Close()
				if err != nil {
					logger.Fatal(err)
				}
				logger.Print("PSQL DB connection closed successfully")
				return
			case <-t.C:
				if len(parts) > 0 {
					logger.Print("Deleting data due to timeout", parts)
					err := st.Flush(ctx, parts)
					if err != nil {
						logger.Fatal(err)
					}
					parts = make([]modelstorage.Removal, 0, flushPartsAmount)
				}
			case part, ok := <-st.ch:
				if !ok {
					return
				}
				parts = append(parts, part)
				if len(parts) >= flushPartsAmount {
					logger.Print("Deleting data due to exceeding capacity", parts)
					err := st.Flush(ctx, parts)
					if err != nil {
						logger.Fatal(err)
					}
					parts = make([]modelstorage.Removal, 0, flushPartsAmount)
				}
			}
		}
	}()
	return &st
}

func (s *Storage) GetBankCardData(ctx context.Context, userID, identifier string) ([]modelstorage.BankCardStorageEntry, error) {
	// prepare SELECT statement
	selectStmt, err := s.DB.PrepareContext(ctx, "SELECT * FROM logins_passwords WHERE user_id = $1 AND identifier = $2")
	defer selectStmt.Close()
	if err != nil {
		return nil, &storageErrors.StatementPSQLError{Err: err}
	}
	chanOk := make(chan []modelstorage.BankCardStorageEntry)
	chanEr := make(chan error)
	go func() {
		s.mu.Lock()
		defer s.mu.Unlock()
		rows, err := selectStmt.QueryContext(ctx, userID, identifier)
		if err != nil {
			chanEr <- &storageErrors.ExecutionPSQLError{Err: err}
			return
		}
		defer rows.Close()
		var queryOutput []modelstorage.BankCardStorageEntry
		for rows.Next() {
			var queryOutputRow modelstorage.BankCardStorageEntry
			err = rows.Scan(&queryOutputRow.ID, &queryOutputRow.UserID, &queryOutputRow.Identifier, &queryOutputRow.Number, &queryOutputRow.Holder, &queryOutputRow.CVV, &queryOutputRow.Meta)
			if err != nil {
				chanEr <- &storageErrors.ScanningPSQLError{Err: err}
				return
			}
			queryOutput = append(queryOutput, queryOutputRow)
		}
		err = rows.Err()
		if err != nil {
			chanEr <- &storageErrors.ScanningPSQLError{Err: err}
		}
		chanOk <- queryOutput
	}()
	select {
	case <-ctx.Done():
		s.logger.Print("getting bank card failed")
		return nil, &storageErrors.ContextTimeoutExceededError{Err: ctx.Err()}
	case methodErr := <-chanEr:
		s.logger.Print("getting bank card failed")
		return nil, methodErr
	case query := <-chanOk:
		s.logger.Print("getting bank card done")
		return query, nil
	}
}

func (s *Storage) GetLoginPasswordData(ctx context.Context, userID, identifier string) ([]modelstorage.LoginPasswordStorageEntry, error) {
	// prepare SELECT statement
	selectStmt, err := s.DB.PrepareContext(ctx, "SELECT * FROM logins_passwords WHERE user_id = $1 AND identifier = $2")
	defer selectStmt.Close()
	if err != nil {
		return nil, &storageErrors.StatementPSQLError{Err: err}
	}
	chanOk := make(chan []modelstorage.LoginPasswordStorageEntry)
	chanEr := make(chan error)
	go func() {
		s.mu.Lock()
		defer s.mu.Unlock()
		rows, err := selectStmt.QueryContext(ctx, userID, identifier)
		if err != nil {
			chanEr <- &storageErrors.ExecutionPSQLError{Err: err}
			return
		}
		defer rows.Close()
		var queryOutput []modelstorage.LoginPasswordStorageEntry
		for rows.Next() {
			var queryOutputRow modelstorage.LoginPasswordStorageEntry
			err = rows.Scan(&queryOutputRow.ID, &queryOutputRow.UserID, &queryOutputRow.Identifier, &queryOutputRow.Login, &queryOutputRow.Password, &queryOutputRow.Meta)
			if err != nil {
				chanEr <- &storageErrors.ScanningPSQLError{Err: err}
				return
			}
			queryOutput = append(queryOutput, queryOutputRow)
		}
		err = rows.Err()
		if err != nil {
			chanEr <- &storageErrors.ScanningPSQLError{Err: err}
		}
		chanOk <- queryOutput
	}()
	select {
	case <-ctx.Done():
		s.logger.Print("getting login/password failed")
		return nil, &storageErrors.ContextTimeoutExceededError{Err: ctx.Err()}
	case methodErr := <-chanEr:
		s.logger.Print("getting login/password failed")
		return nil, methodErr
	case query := <-chanOk:
		s.logger.Print("getting login/password done")
		return query, nil
	}
}

func (s *Storage) GetTextBinaryData(ctx context.Context, userID, identifier string) ([]modelstorage.TextBinaryStorageEntry, error) {
	// prepare SELECT statement
	selectStmt, err := s.DB.PrepareContext(ctx, "SELECT * FROM texts_binaries WHERE user_id = $1 AND identifier = $2")
	defer selectStmt.Close()
	if err != nil {
		return nil, &storageErrors.StatementPSQLError{Err: err}
	}
	chanOk := make(chan []modelstorage.TextBinaryStorageEntry)
	chanEr := make(chan error)
	go func() {
		s.mu.Lock()
		defer s.mu.Unlock()
		rows, err := selectStmt.QueryContext(ctx, userID, identifier)
		if err != nil {
			chanEr <- &storageErrors.ExecutionPSQLError{Err: err}
			return
		}
		defer rows.Close()
		var queryOutput []modelstorage.TextBinaryStorageEntry
		for rows.Next() {
			var queryOutputRow modelstorage.TextBinaryStorageEntry
			err = rows.Scan(&queryOutputRow.ID, &queryOutputRow.UserID, &queryOutputRow.Identifier, &queryOutputRow.Entry, &queryOutputRow.Meta)
			if err != nil {
				chanEr <- &storageErrors.ScanningPSQLError{Err: err}
				return
			}
			queryOutput = append(queryOutput, queryOutputRow)
		}
		err = rows.Err()
		if err != nil {
			chanEr <- &storageErrors.ScanningPSQLError{Err: err}
		}
		chanOk <- queryOutput
	}()
	select {
	case <-ctx.Done():
		s.logger.Print("getting text/binary failed")
		return nil, &storageErrors.ContextTimeoutExceededError{Err: ctx.Err()}
	case methodErr := <-chanEr:
		s.logger.Print("getting text/binary failed")
		return nil, methodErr
	case query := <-chanOk:
		s.logger.Print("getting text/binary done")
		return query, nil
	}
}

func (s *Storage) SetBankCardData(ctx context.Context, userID, identifier, number, holder, cvv, meta string) error {
	selectStmt, err := s.DB.PrepareContext(ctx, "SELECT * FROM bank_cards WHERE user_id = $1 AND identifier = $2")
	if err != nil {
		return &storageErrors.StatementPSQLError{Err: err}
	}
	newDataStmt, err := s.DB.PrepareContext(ctx, "INSERT INTO bank_cards (user_id, identifier, card_number, card_holder, card_cvv, card_meta) VALUES ($1, $2, $3, $4, $5, $6)")
	if err != nil {
		return &storageErrors.StatementPSQLError{Err: err}
	}
	defer selectStmt.Close()
	defer newDataStmt.Close()
	chanOk := make(chan bool)
	chanEr := make(chan error)
	go func() {
		var queryOutput modelstorage.BankCardStorageEntry
		err := selectStmt.QueryRowContext(ctx, userID, identifier).Scan(&queryOutput.ID, &queryOutput.UserID, &queryOutput.Identifier, &queryOutput.Number, &queryOutput.Holder, &queryOutput.CVV, &queryOutput.Meta)
		if err != nil {
			switch {
			case errors.Is(err, sql.ErrNoRows):
				_, err = newDataStmt.ExecContext(ctx, userID, identifier, number, holder, cvv, meta)
				if err != nil {
					chanEr <- &storageErrors.ExecutionPSQLError{Err: err}
					return
				}
				chanOk <- true
			default:
				chanEr <- &storageErrors.ExecutionPSQLError{Err: err}
				return
			}
		}
		chanEr <- &storageErrors.AlreadyExistsError{Err: err}
		return
	}()

	select {
	case <-ctx.Done():
		s.logger.Printf("adding new bank card failed for ID %s", identifier)
		return &storageErrors.ContextTimeoutExceededError{Err: ctx.Err()}
	case methodErr := <-chanEr:
		s.logger.Printf("adding new bank card failed for ID %s", identifier)
		return methodErr
	case <-chanOk:
		s.logger.Printf("adding new bank card done for ID %s", identifier)
		return nil
	}
}

func (s *Storage) SetLoginPasswordData(ctx context.Context, userID, identifier, login, password, meta string) error {
	selectStmt, err := s.DB.PrepareContext(ctx, "SELECT * FROM logins_passwords WHERE user_id = $1 AND identifier = $2")
	if err != nil {
		return &storageErrors.StatementPSQLError{Err: err}
	}
	newDataStmt, err := s.DB.PrepareContext(ctx, "INSERT INTO logins_passwords (user_id, identifier, login, password, cred_meta) VALUES ($1, $2, $3, $4, $5)")
	if err != nil {
		return &storageErrors.StatementPSQLError{Err: err}
	}
	defer selectStmt.Close()
	defer newDataStmt.Close()
	chanOk := make(chan bool)
	chanEr := make(chan error)
	go func() {
		var queryOutput modelstorage.LoginPasswordStorageEntry
		err := selectStmt.QueryRowContext(ctx, userID, identifier).Scan(&queryOutput.ID, &queryOutput.UserID, &queryOutput.Identifier, &queryOutput.Login, &queryOutput.Password, &queryOutput.Meta)
		if err != nil {
			switch {
			case errors.Is(err, sql.ErrNoRows):
				_, err = newDataStmt.ExecContext(ctx, userID, identifier, login, password, meta)
				if err != nil {
					chanEr <- &storageErrors.ExecutionPSQLError{Err: err}
					return
				}
				chanOk <- true
			default:
				chanEr <- &storageErrors.ExecutionPSQLError{Err: err}
				return
			}
		}
		chanEr <- &storageErrors.AlreadyExistsError{Err: err}
		return
	}()

	select {
	case <-ctx.Done():
		s.logger.Printf("adding new login/password failed for ID %s", identifier)
		return &storageErrors.ContextTimeoutExceededError{Err: ctx.Err()}
	case methodErr := <-chanEr:
		s.logger.Printf("adding new login/password failed for ID %s", identifier)
		return methodErr
	case <-chanOk:
		s.logger.Printf("adding new login/password done for ID %s", identifier)
		return nil
	}
}

func (s *Storage) SetTextBinaryData(ctx context.Context, userID, identifier, entry, meta string) error {
	selectStmt, err := s.DB.PrepareContext(ctx, "SELECT * FROM texts_binaries WHERE user_id = $1 AND identifier = $2")
	if err != nil {
		return &storageErrors.StatementPSQLError{Err: err}
	}
	newDataStmt, err := s.DB.PrepareContext(ctx, "INSERT INTO texts_binaries (user_id, identifier, text_entry, text_meta) VALUES ($1, $2, $3, $4)")
	if err != nil {
		return &storageErrors.StatementPSQLError{Err: err}
	}
	defer selectStmt.Close()
	defer newDataStmt.Close()
	chanOk := make(chan bool)
	chanEr := make(chan error)
	go func() {
		var queryOutput modelstorage.TextBinaryStorageEntry
		err := selectStmt.QueryRowContext(ctx, userID, identifier).Scan(&queryOutput.ID, &queryOutput.UserID, &queryOutput.Identifier, &queryOutput.Entry, &queryOutput.Meta)
		if err != nil {
			switch {
			case errors.Is(err, sql.ErrNoRows):
				_, err = newDataStmt.ExecContext(ctx, userID, identifier, entry, meta)
				if err != nil {
					chanEr <- &storageErrors.ExecutionPSQLError{Err: err}
					return
				}
				chanOk <- true
			default:
				chanEr <- &storageErrors.ExecutionPSQLError{Err: err}
				return
			}
		}
		chanEr <- &storageErrors.AlreadyExistsError{Err: err}
		return
	}()

	select {
	case <-ctx.Done():
		s.logger.Printf("adding new text/binary failed for ID %s", identifier)
		return &storageErrors.ContextTimeoutExceededError{Err: ctx.Err()}
	case methodErr := <-chanEr:
		s.logger.Printf("adding new text/binary failed for ID %s", identifier)
		return methodErr
	case <-chanOk:
		s.logger.Printf("adding new text/binary done for ID %s", identifier)
		return nil
	}
}

func (s *Storage) AddNewUser(ctx context.Context, login, password, userID string) error {
	newUserStmt, err := s.DB.PrepareContext(ctx, "INSERT INTO users (user_id, login, password, registered_at) VALUES ($1, $2, $3, $4)")
	if err != nil {
		return &storageErrors.StatementPSQLError{Err: err}
	}
	defer newUserStmt.Close()
	chanOk := make(chan bool)
	chanEr := make(chan error)
	go func() {
		s.mu.Lock()
		defer s.mu.Unlock()
		_, err := newUserStmt.ExecContext(ctx, userID, login, password, time.Now().Format(time.RFC3339))
		if err != nil {
			if err, ok := err.(*pgconn.PgError); ok && err.Code == pgerrcode.UniqueViolation {
				chanEr <- &storageErrors.AlreadyExistsError{Err: err, ID: login}
				return
			}
			chanEr <- &storageErrors.ExecutionPSQLError{Err: err}
			return
		}
		chanOk <- true
	}()

	select {
	case <-ctx.Done():
		s.logger.Printf("adding new user failed for %s", login)
		return &storageErrors.ContextTimeoutExceededError{Err: ctx.Err()}
	case methodErr := <-chanEr:
		s.logger.Printf("adding new user failed for %s", login)
		return methodErr
	case <-chanOk:
		s.logger.Printf("adding new user done for %s", login)
		return nil
	}
}

func (s *Storage) CheckUser(ctx context.Context, login, password string) (string, error) {
	selectStmt, err := s.DB.PrepareContext(ctx, "SELECT * FROM users WHERE login = $1")
	if err != nil {
		return "", &storageErrors.StatementPSQLError{Err: err}
	}
	defer selectStmt.Close()
	chanOk := make(chan string)
	chanEr := make(chan error)
	go func() {
		s.mu.Lock()
		defer s.mu.Unlock()
		var queryOutput modelstorage.UserStorageEntry
		err := selectStmt.QueryRowContext(ctx, login).Scan(&queryOutput.ID, &queryOutput.UserID, &queryOutput.Login, &queryOutput.Password, &queryOutput.RegisteredAt)
		if err != nil {
			switch {
			case errors.Is(err, sql.ErrNoRows):
				chanEr <- &storageErrors.NotFoundError{Err: err}
				return
			default:
				chanEr <- err
				return
			}
		}
		passwordHash := sha256.Sum256([]byte(password))
		expectedPasswordHash := sha256.Sum256([]byte(queryOutput.Password))
		passwordMatch := subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1
		if !passwordMatch {
			chanEr <- &storageErrors.NotFoundError{Err: nil}
		}
		chanOk <- queryOutput.UserID
	}()

	select {
	case <-ctx.Done():
		s.logger.Print("user authentication failed")
		return "", &storageErrors.ContextTimeoutExceededError{Err: ctx.Err()}
	case methodErr := <-chanEr:
		s.logger.Print("user authentication failed")
		return "", methodErr
	case userID := <-chanOk:
		s.logger.Print("user authentication done")
		return userID, nil
	}
}

func (s *Storage) SendToQueue(item modelstorage.Removal) {
	s.ch <- item
}

func (s *Storage) DeleteBatch(ctx context.Context, identifiers []string, userID, db string) error {
	// prepare DELETE statement
	var deleteStmt *sql.Stmt
	var err error
	switch db {
	case "bankCard":
		deleteStmt, err = s.DB.PrepareContext(ctx, "DELETE FROM bank_cards WHERE user_id = $1 AND identifier = ANY($2)")
	case "loginPassword":
		deleteStmt, err = s.DB.PrepareContext(ctx, "DELETE FROM logins_passwords WHERE user_id = $1 AND identifier = ANY($2)")
	case "textBinary":
		deleteStmt, err = s.DB.PrepareContext(ctx, "DELETE FROM texts_binaries WHERE user_id = $1 AND identifier = ANY($2)")
	default:
		return &storageErrors.WrongDBError{
			Err: errors.New("wrong DB identifier"),
			ID:  db,
		}
	}
	defer deleteStmt.Close()
	if err != nil {
		return &storageErrors.StatementPSQLError{Err: err}
	}

	//begin transaction
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return &storageErrors.ExecutionPSQLError{Err: err}
	}
	defer tx.Rollback()
	txDeleteStmt := tx.StmtContext(ctx, deleteStmt)
	// create channels for listening to the go routine result
	deleteDone := make(chan bool)
	deleteError := make(chan error)
	go func() {
		_, err = txDeleteStmt.ExecContext(
			ctx,
			userID,
			pq.Array(identifiers),
		)
		if err != nil {
			deleteError <- &storageErrors.ExecutionPSQLError{Err: err}
		}
		deleteDone <- true
	}()

	// wait for the first channel to retrieve a value
	select {
	case <-ctx.Done():
		log.Println("Deleting data:", ctx.Err())
		return &storageErrors.ContextTimeoutExceededError{Err: ctx.Err()}
	case dltError := <-deleteError:
		log.Println("Deleting data:", dltError.Error())
		return dltError
	case <-deleteDone:
		log.Println("Deleting data:", identifiers)
		return tx.Commit()
	}
}

func (s *Storage) Flush(ctx context.Context, batch []modelstorage.Removal) error {
	uniqueMapBankCards := make(map[string][]string)
	uniqueMapLoginsPasswords := make(map[string][]string)
	uniqueMapTextsBinaries := make(map[string][]string)
	for _, b := range batch {
		switch b.Db {
		case "bankCard":
			if _, exist := uniqueMapBankCards[b.UserID]; !exist {
				uniqueMapBankCards[b.UserID] = []string{b.Identifier}
			} else {
				uniqueMapBankCards[b.UserID] = append(uniqueMapBankCards[b.UserID], b.Identifier)
			}
		case "loginPassword":
			if _, exist := uniqueMapLoginsPasswords[b.UserID]; !exist {
				uniqueMapLoginsPasswords[b.UserID] = []string{b.Identifier}
			} else {
				uniqueMapLoginsPasswords[b.UserID] = append(uniqueMapLoginsPasswords[b.UserID], b.Identifier)
			}
		case "textBinary":
			if _, exist := uniqueMapTextsBinaries[b.UserID]; !exist {
				uniqueMapTextsBinaries[b.UserID] = []string{b.Identifier}
			} else {
				uniqueMapTextsBinaries[b.UserID] = append(uniqueMapTextsBinaries[b.UserID], b.Identifier)
			}
		}
	}
	for userID, identifiers := range uniqueMapBankCards {
		err := s.DeleteBatch(ctx, identifiers, userID, "bankCard")
		if err != nil {
			return err
		}
	}
	for userID, identifiers := range uniqueMapLoginsPasswords {
		err := s.DeleteBatch(ctx, identifiers, userID, "loginPassword")
		if err != nil {
			return err
		}
	}
	for userID, identifiers := range uniqueMapTextsBinaries {
		err := s.DeleteBatch(ctx, identifiers, userID, "textBinary")
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Storage) createTables(ctx context.Context) error {
	var queries []string
	query := `CREATE TABLE IF NOT EXISTS users (
		id				BIGSERIAL   	NOT NULL UNIQUE,
		user_id       	TEXT        	NOT NULL UNIQUE,
		login         	TEXT        	NOT NULL UNIQUE,
		password      	TEXT        	NOT NULL,
		registered_at 	TIMESTAMPTZ 	NOT NULL  
	);`
	queries = append(queries, query)
	query = `CREATE TABLE IF NOT EXISTS logins_passwords (
		id           	BIGSERIAL      	NOT NULL UNIQUE,
		user_id      	TEXT           	NOT NULL,
		identifier      TEXT           	NOT NULL,
		login		  	TEXT           	NOT NULL,
		password	  	TEXT 		   	NOT NULL,
		cred_meta		TEXT 
	);`
	queries = append(queries, query)
	query = `CREATE TABLE IF NOT EXISTS texts_binaries (
		id           	BIGSERIAL      	NOT NULL UNIQUE,
		user_id      	TEXT           	NOT NULL,
		identifier      TEXT           	NOT NULL,
		text_entry  	TEXT           	NOT NULL,
		text_meta		TEXT
	);`
	queries = append(queries, query)
	query = `CREATE TABLE IF NOT EXISTS bank_cards (
		id           	BIGSERIAL      	NOT NULL UNIQUE,
		user_id      	TEXT           	NOT NULL,
		identifier      TEXT           	NOT NULL,
		card_number  	TEXT           	NOT NULL,
		card_holder  	TEXT 		   	NOT NULL,
		card_cvv		TEXT			NOT NULL,
		card_meta		TEXT
	);`
	queries = append(queries, query)
	for _, subquery := range queries {
		_, err := s.DB.ExecContext(ctx, subquery)
		if err != nil {
			return err
		}
	}
	return nil
}

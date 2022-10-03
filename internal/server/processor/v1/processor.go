package processor

import (
	"context"
	"dk-go-gophkeeper/internal/server/cipher"
	"dk-go-gophkeeper/internal/server/modeldto"
	"dk-go-gophkeeper/internal/server/processor"
	"dk-go-gophkeeper/internal/server/storage"
	"dk-go-gophkeeper/internal/server/storage/modelstorage"
	"log"
)

var (
	_ processor.Processor = (*Processor)(nil)
)

type Processor struct {
	storage storage.DataStorage
	cipher  cipher.Cipher
	logger  *log.Logger
}

func InitService(st storage.DataStorage, cp cipher.Cipher, logger *log.Logger) *Processor {
	logger.Print("Attempting to initialize processor")
	serviceProcessor := &Processor{
		storage: st,
		cipher:  cp,
		logger:  logger,
	}
	return serviceProcessor
}

func (proc *Processor) GetUserID(accessToken string) (string, error) {
	userID, err := proc.cipher.ValidateToken(accessToken)
	return userID, err
}

func (proc *Processor) GetTokenForUser(userID string) string {
	token := proc.cipher.Encode(userID)
	return token
}

func (proc *Processor) AddNewUser(ctx context.Context, login, password string) (string, error) {
	accessToken, userID := proc.cipher.NewToken()
	err := proc.storage.AddNewUser(ctx, proc.cipher.Encode(login), proc.cipher.Encode(password), userID)
	if err != nil {
		return "", err
	}
	return accessToken, nil
}

func (proc *Processor) LoginUser(ctx context.Context, login, password string) (string, error) {
	userID, err := proc.storage.CheckUser(ctx, proc.cipher.Encode(login), proc.cipher.Encode(password))
	if err != nil {
		return "", err
	}
	accessToken := proc.cipher.Encode(userID)
	return accessToken, nil
}

func (proc *Processor) GetBankCardData(ctx context.Context, userID, identifier string) ([]modeldto.BankCard, error) {
	bankCards, err := proc.storage.GetBankCardData(ctx, userID, proc.cipher.Encode(identifier))
	if err != nil {
		return nil, err
	}
	var responseBankCards []modeldto.BankCard
	for _, bankCard := range bankCards {
		decodedIdentifier, err := proc.cipher.Decode(bankCard.Identifier)
		if err != nil {
			return nil, err
		}
		decodedNumber, err := proc.cipher.Decode(bankCard.Number)
		if err != nil {
			return nil, err
		}
		decodedHolder, err := proc.cipher.Decode(bankCard.Holder)
		if err != nil {
			return nil, err
		}
		decodedCVV, err := proc.cipher.Decode(bankCard.CVV)
		if err != nil {
			return nil, err
		}
		decodedMeta, err := proc.cipher.Decode(bankCard.Meta)
		if err != nil {
			return nil, err
		}
		responseBankCard := modeldto.BankCard{
			Identifier: decodedIdentifier,
			Number:     decodedNumber,
			Holder:     decodedHolder,
			CVV:        decodedCVV,
			Meta:       decodedMeta,
		}
		responseBankCards = append(responseBankCards, responseBankCard)
	}
	return responseBankCards, nil
}

func (proc *Processor) GetLoginPasswordData(ctx context.Context, userID, identifier string) ([]modeldto.LoginPassword, error) {
	loginsPasswords, err := proc.storage.GetLoginPasswordData(ctx, userID, proc.cipher.Encode(identifier))
	if err != nil {
		return nil, err
	}
	var responseLoginsPasswords []modeldto.LoginPassword
	for _, loginPassword := range loginsPasswords {
		decodedIdentifier, err := proc.cipher.Decode(loginPassword.Identifier)
		if err != nil {
			return nil, err
		}
		decodedLogin, err := proc.cipher.Decode(loginPassword.Login)
		if err != nil {
			return nil, err
		}
		decodedPassword, err := proc.cipher.Decode(loginPassword.Password)
		if err != nil {
			return nil, err
		}
		decodedMeta, err := proc.cipher.Decode(loginPassword.Meta)
		if err != nil {
			return nil, err
		}
		responseLoginPassword := modeldto.LoginPassword{
			Identifier: decodedIdentifier,
			Login:      decodedLogin,
			Password:   decodedPassword,
			Meta:       decodedMeta,
		}
		responseLoginsPasswords = append(responseLoginsPasswords, responseLoginPassword)
	}
	return responseLoginsPasswords, nil
}

func (proc *Processor) GetTextBinaryData(ctx context.Context, userID, identifier string) ([]modeldto.TextBinary, error) {
	textsBinaries, err := proc.storage.GetTextBinaryData(ctx, userID, proc.cipher.Encode(identifier))
	if err != nil {
		return nil, err
	}
	var responseTextsBinaries []modeldto.TextBinary
	for _, textBinary := range textsBinaries {
		decodedIdentifier, err := proc.cipher.Decode(textBinary.Identifier)
		if err != nil {
			return nil, err
		}
		decodedEntry, err := proc.cipher.Decode(textBinary.Entry)
		if err != nil {
			return nil, err
		}
		decodedMeta, err := proc.cipher.Decode(textBinary.Meta)
		if err != nil {
			return nil, err
		}
		responsetextBinary := modeldto.TextBinary{
			Identifier: decodedIdentifier,
			Entry:      decodedEntry,
			Meta:       decodedMeta,
		}
		responseTextsBinaries = append(responseTextsBinaries, responsetextBinary)
	}
	return responseTextsBinaries, nil
}

func (proc *Processor) SetBankCardData(ctx context.Context, userID, identifier, number, holder, cvv, meta string) error {
	encodedIndentifier := proc.cipher.Encode(identifier)
	encodedNumber := proc.cipher.Encode(number)
	encodedHolder := proc.cipher.Encode(holder)
	encodedCvv := proc.cipher.Encode(cvv)
	encodedMeta := proc.cipher.Encode(meta)
	err := proc.storage.SetBankCardData(ctx, userID, encodedIndentifier, encodedNumber, encodedHolder, encodedCvv, encodedMeta)
	return err
}

func (proc *Processor) SetLoginPasswordData(ctx context.Context, userID, identifier, login, password, meta string) error {
	encodedIndentifier := proc.cipher.Encode(identifier)
	encodedLogin := proc.cipher.Encode(login)
	encodedPassword := proc.cipher.Encode(password)
	encodedMeta := proc.cipher.Encode(meta)
	err := proc.storage.SetLoginPasswordData(ctx, userID, encodedIndentifier, encodedLogin, encodedPassword, encodedMeta)
	return err
}

func (proc *Processor) SetTextBinaryData(ctx context.Context, userID, identifier, entry, meta string) error {
	encodedIndentifier := proc.cipher.Encode(identifier)
	encodedEntry := proc.cipher.Encode(entry)
	encodedMeta := proc.cipher.Encode(meta)
	err := proc.storage.SetTextBinaryData(ctx, userID, encodedIndentifier, encodedEntry, encodedMeta)
	return err
}

func (proc *Processor) Delete(userID, identifier, db string) {
	encodedIndentifier := proc.cipher.Encode(identifier)
	item := modelstorage.Removal{
		UserID:     userID,
		Identifier: encodedIndentifier,
		Db:         db,
	}
	proc.storage.SendToQueue(item)
}

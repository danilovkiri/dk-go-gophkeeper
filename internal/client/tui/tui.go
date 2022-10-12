// Package tui provides TUI functionality.
package tui

import (
	"context"
	"dk-go-gophkeeper/internal/client/storage"
	"dk-go-gophkeeper/internal/client/tui/modeltui"
	"dk-go-gophkeeper/internal/config"
	"fmt"
	"log"
	"strings"
	"unicode"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog"
)

// build parameters to be used with ldflags

var (
	buildVersion = "NA"
	buildDate    = "NA"
	buildCommit  = "NA"
)

// tview page identifiers

const (
	pageStoreLoginPassword = "store_login_password"
	pageStoreTextBinary    = "store_text_binary"
	pageStoreBankCard      = "store_bank_card"
	pageRemove             = "remove"
	pageRegister           = "register"
	pageLogin              = "login"
	pageGetData            = "get_data"
	pageResult             = "result"
	pageMenu               = "menu"
)

// tview form field lengths

const (
	identifierLength     = 20
	metaLength           = 50
	bankCardNumberLength = 16
	bankCardHolderLength = 20
	bankCardCVVLength    = 3
	loginLength          = 20
	passwordLength       = 20
	textEntryLength      = 50
)

// shared static attributes

var flex = tview.NewFlex()
var pages = tview.NewPages()

// start screen static attributes

var header = tview.NewTextView().SetText(fmt.Sprintf("GophKeeper: build %s, date %s, commit %s", buildVersion, buildDate, buildCommit)).SetTextAlign(1).SetTextColor(tcell.ColorGreen)
var footer = tview.NewTextView().SetText("Kirill Danilov, 2022, https://github.com/danilovkiri/").SetTextAlign(1).SetTextColor(tcell.ColorGreen)
var statusHeader = tview.NewTextView().SetText("Last operation status:").SetTextAlign(1).SetTextColor(tcell.ColorGreen)
var buttonSync = tview.NewButton("Sync")
var buttonQuit = tview.NewButton("Quit")
var buttonLogin = tview.NewButton("Login")
var buttonRegister = tview.NewButton("Register")
var menu = tview.NewFlex().
	AddItem(buttonSync, 0, 1, false).
	AddItem(tview.NewBox(), 0, 1, false).
	AddItem(buttonQuit, 0, 1, false).
	AddItem(tview.NewBox(), 0, 1, false).
	AddItem(buttonLogin, 0, 1, false).
	AddItem(tview.NewBox(), 0, 1, false).
	AddItem(buttonRegister, 0, 1, false)
var buttonStoreLoginPassword = tview.NewButton("Add login/password item")
var buttonStoreTextBinary = tview.NewButton("Add text/binary item")
var buttonStoreBankCard = tview.NewButton("Add bank card item")
var buttonGetData = tview.NewButton("Get item")
var buttonRemove = tview.NewButton("Remove item")
var buttonBackToMainScreen = tview.NewButton("Back to menu")
var input = tview.NewFlex().SetDirection(tview.FlexRow).
	AddItem(buttonStoreLoginPassword, 0, 10, false).
	AddItem(tview.NewBox(), 0, 2, false).
	AddItem(buttonStoreTextBinary, 0, 10, false).
	AddItem(tview.NewBox(), 0, 2, false).
	AddItem(buttonStoreBankCard, 0, 10, false).
	AddItem(tview.NewBox(), 0, 2, false).
	AddItem(buttonGetData, 0, 10, false).
	AddItem(tview.NewBox(), 0, 2, false).
	AddItem(buttonRemove, 0, 10, false)
var body = tview.NewFlex().AddItem(input, 0, 1, false)

// App defines attributes and methods of an App instance.
type App struct {
	App                    *tview.Application
	storage                storage.DataStorage
	cancel                 context.CancelFunc
	registerLoginDetails   modeltui.RegisterLogin
	registerForm           *tview.Form
	loginForm              *tview.Form
	bankCards              []modeltui.BankCard
	storeBankCardForm      *tview.Form
	textsOrBinaries        []modeltui.TextOrBinary
	storeTextOrBinaryForm  *tview.Form
	loginsAndPasswords     []modeltui.LoginAndPassword
	storeLoginPasswordForm *tview.Form
	removeForm             *tview.Form
	retrieveDataPieceForm  *tview.Form
	loginStatus            *tview.TextView
	operationStatus        *tview.TextView
	result                 *tview.TextView
	logger                 *zerolog.Logger
	cfg                    *config.Config
}

// addRetrieveDataPieceForm defines form behavior and its contents.
func (a *App) addRetrieveDataPieceForm() *tview.Form {
	query := modeltui.Get{}
	a.retrieveDataPieceForm.AddInputField("Identifier", "", identifierLength, nil, func(id string) {
		if strings.ReplaceAll(id, " ", "") == "" {
			a.operationStatus.SetText("Identifier cannot be empty")
			pages.SwitchToPage("menu")
		} else {
			query.Identifier = id
		}
	})
	a.retrieveDataPieceForm.AddDropDown("DB type", []string{a.cfg.BankCardDB, a.cfg.LoginPasswordDB, a.cfg.TextBinaryDB}, 0, func(db string, idx int) {
		query.Db = db
	})
	a.retrieveDataPieceForm.AddButton("Get", func() {
		result, err := a.storage.Get(query.Identifier, query.Db)
		if err != nil {
			a.operationStatus.SetText(err.Error())
			pages.SwitchToPage("menu")
		} else {
			a.operationStatus.SetText("Getting data: OK")
			a.result.SetText(result)
			pages.SwitchToPage("result")
		}
	})
	a.retrieveDataPieceForm.AddButton("Cancel", func() {
		pages.SwitchToPage("menu")
	})
	return a.retrieveDataPieceForm
}

// addRemovalForm defines form behavior and its contents.
func (a *App) addRemovalForm() *tview.Form {
	removal := modeltui.Removal{}
	a.removeForm.AddInputField("Identifier", "", identifierLength, nil, func(id string) {
		if strings.ReplaceAll(id, " ", "") == "" {
			a.operationStatus.SetText("Identifier cannot be empty")
			pages.SwitchToPage("menu")
		} else {
			removal.Identifier = id
		}
	})
	a.removeForm.AddDropDown("DB type", []string{a.cfg.BankCardDB, a.cfg.LoginPasswordDB, a.cfg.TextBinaryDB}, 0, func(db string, idx int) {
		removal.Db = db
	})
	a.removeForm.AddButton("Remove", func() {
		err := a.storage.Remove(removal.Identifier, removal.Db)
		if err != nil {
			a.operationStatus.SetText(err.Error())
		} else {
			a.operationStatus.SetText("Removal: OK")
		}
		pages.SwitchToPage("menu")
	})
	a.removeForm.AddButton("Cancel", func() {
		pages.SwitchToPage("menu")
	})
	return a.removeForm
}

// addLoginPasswordForm defines form behavior and its contents.
func (a *App) addLoginPasswordForm() *tview.Form {
	loginAndPassword := modeltui.LoginAndPassword{}
	a.storeLoginPasswordForm.AddInputField("Identifier", "", identifierLength, nil, func(id string) {
		if strings.ReplaceAll(id, " ", "") == "" {
			a.operationStatus.SetText("Identifier cannot be empty")
			pages.SwitchToPage("menu")
		} else {
			loginAndPassword.Identifier = id
		}
	})
	a.storeLoginPasswordForm.AddInputField("Login", "", loginLength, nil, func(login string) {
		if strings.ReplaceAll(login, " ", "") == "" {
			a.operationStatus.SetText("Login cannot be empty")
			pages.SwitchToPage("menu")
		} else {
			loginAndPassword.Login = login
		}
	})
	a.storeLoginPasswordForm.AddInputField("Password", "", passwordLength, nil, func(password string) {
		if strings.ReplaceAll(password, " ", "") == "" {
			a.operationStatus.SetText("Password cannot be empty")
			pages.SwitchToPage("menu")
		} else {
			loginAndPassword.Password = password
		}
	})
	a.storeLoginPasswordForm.AddInputField("Meta", "", metaLength, nil, func(meta string) {
		loginAndPassword.Meta = meta
	})
	a.storeLoginPasswordForm.AddButton("Submit", func() {
		err := a.storage.AddLoginPassword(loginAndPassword.Identifier, loginAndPassword.Login, loginAndPassword.Password, loginAndPassword.Meta)
		if err != nil {
			a.operationStatus.SetText(err.Error())
		} else {
			a.operationStatus.SetText("Adding login/password: OK")
		}
		pages.SwitchToPage("menu")
	})
	a.storeLoginPasswordForm.AddButton("Cancel", func() {
		pages.SwitchToPage("menu")
	})
	return a.storeLoginPasswordForm
}

// addTextOrBinaryForm defines form behavior and its contents.
func (a *App) addTextOrBinaryForm() *tview.Form {
	textOrBinary := modeltui.TextOrBinary{}
	a.storeTextOrBinaryForm.AddInputField("Identifier", "", identifierLength, nil, func(id string) {
		if strings.ReplaceAll(id, " ", "") == "" {
			a.operationStatus.SetText("Identifier cannot be empty")
			pages.SwitchToPage("menu")
		} else {
			textOrBinary.Identifier = id
		}
	})
	a.storeTextOrBinaryForm.AddInputField("Input", "", textEntryLength, nil, func(entry string) {
		textOrBinary.Entry = entry
	})
	a.storeTextOrBinaryForm.AddInputField("Meta", "", metaLength, nil, func(meta string) {
		textOrBinary.Meta = meta
	})
	a.storeTextOrBinaryForm.AddButton("Submit", func() {
		err := a.storage.AddTextBinary(textOrBinary.Identifier, textOrBinary.Entry, textOrBinary.Meta)
		if err != nil {
			a.operationStatus.SetText(err.Error())
		} else {
			a.operationStatus.SetText("Adding text/binary: OK")
		}
		pages.SwitchToPage("menu")
	})
	a.storeTextOrBinaryForm.AddButton("Cancel", func() {
		pages.SwitchToPage("menu")
	})
	return a.storeTextOrBinaryForm
}

// addBankCardForm defines form behavior and its contents.
func (a *App) addBankCardForm() *tview.Form {
	bankCard := modeltui.BankCard{}
	a.storeBankCardForm.AddInputField("Identifier", "", identifierLength, nil, func(id string) {
		if strings.ReplaceAll(id, " ", "") == "" {
			a.operationStatus.SetText("Identifier cannot be empty")
			pages.SwitchToPage("menu")
		} else {
			bankCard.Identifier = id
		}
	})
	a.storeBankCardForm.AddInputField("Number", "", bankCardNumberLength, nil, func(number string) {
		if len(number) != 16 && !isInt(number) {
			a.operationStatus.SetText("Bank card number must be a 16-digit code")
			pages.SwitchToPage("menu")
		} else {
			bankCard.Number = number
		}
	})
	a.storeBankCardForm.AddInputField("Holder", "", bankCardHolderLength, nil, func(holder string) {
		if strings.ReplaceAll(holder, " ", "") == "" {
			a.operationStatus.SetText("Bank card holder cannot be empty")
			pages.SwitchToPage("menu")
		} else {
			bankCard.Holder = holder
		}
	})
	a.storeBankCardForm.AddInputField("CVV", "", bankCardCVVLength, nil, func(cvv string) {
		if len(cvv) != 3 && !isInt(cvv) {
			a.operationStatus.SetText("Bank card CVV must be a 3-digit code")
			pages.SwitchToPage("menu")
		} else {
			bankCard.Number = cvv
		}
	})
	a.storeBankCardForm.AddInputField("Meta", "", metaLength, nil, func(meta string) {
		bankCard.Meta = meta
	})
	a.storeBankCardForm.AddButton("Submit", func() {
		err := a.storage.AddBankCard(bankCard.Identifier, bankCard.Number, bankCard.Holder, bankCard.Cvv, bankCard.Meta)
		if err != nil {
			a.operationStatus.SetText(err.Error())
		} else {
			a.operationStatus.SetText("Adding bank card: OK")
		}
		pages.SwitchToPage("menu")
	})
	a.storeBankCardForm.AddButton("Cancel", func() {
		pages.SwitchToPage("menu")
	})
	return a.storeBankCardForm
}

// addRegisterForm defines form behavior and its contents.
func (a *App) addRegisterForm() *tview.Form {
	a.registerForm.AddInputField("Login", "", loginLength, nil, func(login string) {
		if strings.ReplaceAll(login, " ", "") == "" {
			a.operationStatus.SetText("Login cannot be empty")
			pages.SwitchToPage("menu")
		} else {
			a.registerLoginDetails.Login = login
		}
	})
	a.registerForm.AddInputField("Password", "", passwordLength, nil, func(password string) {
		if strings.ReplaceAll(password, " ", "") == "" {
			a.operationStatus.SetText("Password cannot be empty")
			pages.SwitchToPage("menu")
		} else {
			a.registerLoginDetails.Password = password
		}
	})
	a.registerForm.AddButton("Submit", func() {
		err := a.storage.Register(a.registerLoginDetails.Login, a.registerLoginDetails.Password)
		if err != nil {
			a.operationStatus.SetText(err.Error())
		} else {
			a.operationStatus.SetText("Register: OK")
			a.loginStatus.SetText(fmt.Sprintf("Logged in as: %s", a.registerLoginDetails.Login))
		}
		pages.SwitchToPage("menu")
	})
	a.registerForm.AddButton("Cancel", func() {
		pages.SwitchToPage("menu")
	})
	return a.registerForm
}

// addLoginForm defines form behavior and its contents.
func (a *App) addLoginForm() *tview.Form {
	a.loginForm.AddInputField("Login", "", loginLength, nil, func(login string) {
		if strings.ReplaceAll(login, " ", "") == "" {
			a.operationStatus.SetText("Login cannot be empty")
			pages.SwitchToPage("menu")
		} else {
			a.registerLoginDetails.Login = login
		}
	})
	a.loginForm.AddInputField("Password", "", passwordLength, nil, func(password string) {
		if strings.ReplaceAll(password, " ", "") == "" {
			a.operationStatus.SetText("Password cannot be empty")
			pages.SwitchToPage("menu")
		} else {
			a.registerLoginDetails.Password = password
		}
	})
	a.loginForm.AddButton("Submit", func() {
		err := a.storage.Login(a.registerLoginDetails.Login, a.registerLoginDetails.Password)
		if err != nil {
			a.operationStatus.SetText(err.Error())
		} else {
			a.operationStatus.SetText("Login: OK")
			a.loginStatus.SetText(fmt.Sprintf("Logged in as: %s", a.registerLoginDetails.Login))
		}
		pages.SwitchToPage("menu")
	})
	a.loginForm.AddButton("Cancel", func() {
		pages.SwitchToPage("menu")
	})
	return a.loginForm
}

// InitTUI initializes a TUI instance and defines non-static attributes.
func InitTUI(cancel context.CancelFunc, storage storage.DataStorage, logger *zerolog.Logger, cfg *config.Config) App {
	logger.Print("Attempting to initialize TUI")
	var app = tview.NewApplication()
	application := App{
		App:                    app,
		storage:                storage,
		cancel:                 cancel,
		registerLoginDetails:   modeltui.RegisterLogin{},
		registerForm:           tview.NewForm(),
		loginForm:              tview.NewForm(),
		bankCards:              make([]modeltui.BankCard, 0),
		storeBankCardForm:      tview.NewForm(),
		textsOrBinaries:        make([]modeltui.TextOrBinary, 0),
		storeTextOrBinaryForm:  tview.NewForm(),
		loginsAndPasswords:     make([]modeltui.LoginAndPassword, 0),
		storeLoginPasswordForm: tview.NewForm(),
		removeForm:             tview.NewForm(),
		retrieveDataPieceForm:  tview.NewForm(),
		loginStatus:            tview.NewTextView().SetText("Logged in as: NA").SetTextAlign(1).SetScrollable(true).SetTextColor(tcell.ColorRed),
		operationStatus:        tview.NewTextView().SetText("Nothing to report yet").SetTextAlign(1).SetScrollable(true).SetTextColor(tcell.ColorRed),
		result:                 tview.NewTextView().SetText("Nothing was requested yet").SetTextAlign(1).SetScrollable(true),
		logger:                 logger,
		cfg:                    cfg,
	}
	return application
}

// Run starts the TUI instance, sets on-the-run behaviour.
func (a *App) Run() {
	defer a.cancel()
	buttonStoreLoginPassword.SetSelectedFunc(func() {
		a.storeLoginPasswordForm.Clear(true)
		a.addLoginPasswordForm()
		pages.SwitchToPage(pageStoreLoginPassword)
	})
	buttonStoreTextBinary.SetSelectedFunc(func() {
		a.storeTextOrBinaryForm.Clear(true)
		a.addTextOrBinaryForm()
		pages.SwitchToPage(pageStoreTextBinary)
	})
	buttonStoreBankCard.SetSelectedFunc(func() {
		a.storeBankCardForm.Clear(true)
		a.addBankCardForm()
		pages.SwitchToPage(pageStoreBankCard)
	})
	buttonRemove.SetSelectedFunc(func() {
		a.removeForm.Clear(true)
		a.addRemovalForm()
		pages.SwitchToPage(pageRemove)
	})
	buttonRegister.SetSelectedFunc(func() {
		a.registerForm.Clear(true)
		a.addRegisterForm()
		pages.SwitchToPage(pageRegister)
	})
	buttonLogin.SetSelectedFunc(func() {
		a.loginForm.Clear(true)
		a.addLoginForm()
		pages.SwitchToPage(pageLogin)
	})
	buttonQuit.SetSelectedFunc(func() {
		a.App.Stop()
		a.cancel()
	})
	buttonSync.SetSelectedFunc(func() {
		err := a.storage.Sync()
		if err != nil {
			a.operationStatus.SetText(err.Error())
		} else {
			a.operationStatus.SetText("Syncing OK")
		}
	})
	buttonGetData.SetSelectedFunc(func() {
		a.retrieveDataPieceForm.Clear(true)
		a.addRetrieveDataPieceForm()
		pages.SwitchToPage(pageGetData)
	})
	buttonBackToMainScreen.SetSelectedFunc(func() {
		pages.SwitchToPage(pageMenu)
	})

	resultView := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(a.result, 0, 9, false).
		AddItem(buttonBackToMainScreen, 0, 1, false)

	flex.SetDirection(tview.FlexRow).AddItem(header, 0, 1, false).
		AddItem(a.loginStatus, 0, 2, false).
		AddItem(menu, 0, 1, false).
		AddItem(tview.NewBox(), 0, 2, false).
		AddItem(body, 0, 20, false).
		AddItem(tview.NewBox(), 0, 2, false).
		AddItem(statusHeader, 0, 1, false).
		AddItem(a.operationStatus, 0, 5, false).
		AddItem(footer, 0, 1, false)

	pages.AddPage(pageMenu, flex, true, true)
	pages.AddPage(pageStoreLoginPassword, a.storeLoginPasswordForm, true, false)
	pages.AddPage(pageStoreTextBinary, a.storeTextOrBinaryForm, true, false)
	pages.AddPage(pageStoreBankCard, a.storeBankCardForm, true, false)
	pages.AddPage(pageRegister, a.registerForm, true, false)
	pages.AddPage(pageLogin, a.loginForm, true, false)
	pages.AddPage(pageRemove, a.removeForm, true, false)
	pages.AddPage(pageGetData, a.retrieveDataPieceForm, true, false)
	pages.AddPage(pageResult, resultView, true, false)

	a.logger.Info().Msg("Starting the TUI")
	if err := a.App.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		log.Panic(err)
	}
	a.logger.Info().Msg("TUI closed, Run() function returned")
}

// isInt checks that any rune inside a string is a digit
func isInt(s string) bool {
	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}

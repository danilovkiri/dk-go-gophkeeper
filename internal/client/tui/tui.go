package tui

import (
	"dk-go-gophkeeper/internal/client/tui/modeltui"
	"fmt"
	"github.com/rivo/tview"
	"log"
)

var (
	buildVersion = "NA"
	buildDate    = "NA"
	buildCommit  = "NA"
)

var flex = tview.NewFlex()
var pages = tview.NewPages()

// start screen attributes

var header = tview.NewTextView().SetText(fmt.Sprintf("GophKeeper: build %s, date %s, commit %s", buildVersion, buildDate, buildCommit)).SetTextAlign(1)
var footer = tview.NewTextView().SetText(fmt.Sprint("Kirill Danilov, 2022, https://github.com/danilovkiri/")).SetTextAlign(1)
var statusHeader = tview.NewTextView().SetText(fmt.Sprint("Operation status:")).SetTextAlign(1)
var buttonSync = tview.NewButton("Sync")
var buttonQuit = tview.NewButton("Quit")
var buttonRegisterLogin = tview.NewButton("Register/Login")
var menu = tview.NewFlex().
	AddItem(buttonSync, 0, 1, false).
	AddItem(buttonQuit, 0, 1, false).
	AddItem(buttonRegisterLogin, 0, 1, false)
var buttonStoreLoginPassword = tview.NewButton("login/password")
var buttonStoreTextBinary = tview.NewButton("text/binary")
var buttonStoreBankCard = tview.NewButton("bank card")
var buttonGetAllData = tview.NewButton("get all data")
var buttonRemove = tview.NewButton("remove")
var buttonBackToMainScreen = tview.NewButton("back to menu")
var input = tview.NewFlex().SetDirection(tview.FlexRow).
	AddItem(buttonStoreLoginPassword, 0, 1, false).
	AddItem(buttonStoreTextBinary, 0, 1, false).
	AddItem(buttonStoreBankCard, 0, 1, false).
	AddItem(buttonGetAllData, 0, 1, false).
	AddItem(buttonRemove, 0, 1, false)
var body = tview.NewFlex().AddItem(input, 0, 1, false)

type App struct {
	App                    *tview.Application
	registerLoginDetails   modeltui.RegisterLogin
	registerLoginForm      *tview.Form
	bankCards              []modeltui.BankCard
	storeBankCardForm      *tview.Form
	textsOrBinaries        []modeltui.TextOrBinary
	storeTextOrBinaryForm  *tview.Form
	loginsAndPasswords     []modeltui.LoginAndPassword
	storeLoginPasswordForm *tview.Form
	removals               []modeltui.Removal
	removeForm             *tview.Form
	loginStatus            *tview.TextView
	operationStatus        *tview.TextView
	result                 *tview.TextView
}

func (a *App) addRemovalForm() *tview.Form {
	removal := modeltui.Removal{}
	a.removeForm.AddInputField("Identifier", "", 20, nil, func(id string) {
		removal.Identifier = id
	})
	a.removeForm.AddButton("Save", func() {
		a.removals = append(a.removals, removal)
		pages.SwitchToPage("menu")
	})
	a.removeForm.AddButton("Exit", func() {
		pages.SwitchToPage("menu")
	})
	return a.removeForm
}

func (a *App) addLoginPasswordForm() *tview.Form {
	loginAndPassword := modeltui.LoginAndPassword{}
	a.storeLoginPasswordForm.AddInputField("Identifier", "", 20, nil, func(id string) {
		loginAndPassword.Identifier = id
	})
	a.storeLoginPasswordForm.AddInputField("Login", "", 20, nil, func(login string) {
		loginAndPassword.Login = login
	})
	a.storeLoginPasswordForm.AddInputField("Password", "", 20, nil, func(password string) {
		loginAndPassword.Password = password
	})
	a.storeLoginPasswordForm.AddInputField("Meta", "", 20, nil, func(meta string) {
		loginAndPassword.Meta = meta
	})
	a.storeLoginPasswordForm.AddButton("Save", func() {
		a.loginsAndPasswords = append(a.loginsAndPasswords, loginAndPassword)
		pages.SwitchToPage("menu")
	})
	a.storeLoginPasswordForm.AddButton("Exit", func() {
		pages.SwitchToPage("menu")
	})
	return a.storeLoginPasswordForm
}

func (a *App) addTextOrBinaryForm() *tview.Form {
	textOrBinary := modeltui.TextOrBinary{}
	a.storeTextOrBinaryForm.AddInputField("Identifier", "", 20, nil, func(id string) {
		textOrBinary.Identifier = id
	})
	a.storeTextOrBinaryForm.AddInputField("Input", "", 20, nil, func(entry string) {
		textOrBinary.Entry = entry
	})
	a.storeTextOrBinaryForm.AddInputField("Meta", "", 20, nil, func(meta string) {
		textOrBinary.Meta = meta
	})
	a.storeTextOrBinaryForm.AddButton("Save", func() {
		a.textsOrBinaries = append(a.textsOrBinaries, textOrBinary)
		pages.SwitchToPage("menu")
	})
	a.storeTextOrBinaryForm.AddButton("Exit", func() {
		pages.SwitchToPage("menu")
	})
	return a.storeTextOrBinaryForm
}

func (a *App) addBankCardForm() *tview.Form {
	bankCard := modeltui.BankCard{}
	a.storeBankCardForm.AddInputField("Identifier", "", 20, nil, func(id string) {
		bankCard.Identifier = id
	})
	a.storeBankCardForm.AddInputField("Number", "", 20, nil, func(number string) {
		bankCard.Number = number
	})
	a.storeBankCardForm.AddInputField("Holder", "", 20, nil, func(holder string) {
		bankCard.Holder = holder
	})
	a.storeBankCardForm.AddInputField("CVV", "", 20, nil, func(cvv string) {
		bankCard.Cvv = cvv
	})
	a.storeBankCardForm.AddInputField("Meta", "", 20, nil, func(meta string) {
		bankCard.Meta = meta
	})
	a.storeBankCardForm.AddButton("Save", func() {
		a.bankCards = append(a.bankCards, bankCard)
		pages.SwitchToPage("menu")
	})
	a.storeBankCardForm.AddButton("Exit", func() {
		pages.SwitchToPage("menu")
	})
	return a.storeBankCardForm
}

func (a *App) addRegisterLoginForm() *tview.Form {
	a.registerLoginForm.AddInputField("Login", "", 20, nil, func(login string) {
		a.registerLoginDetails.Login = login
	})
	a.registerLoginForm.AddInputField("Password", "", 20, nil, func(password string) {
		a.registerLoginDetails.Password = password
	})
	a.registerLoginForm.AddButton("OK", func() {
		pages.SwitchToPage("menu")
	})
	a.registerLoginForm.AddButton("Exit", func() {
		pages.SwitchToPage("menu")
	})
	return a.registerLoginForm
}

func (a *App) Sync() error {
	return nil
}

func (a *App) GetAllData() error {
	return nil
}

func InitTUI() App {
	var app = tview.NewApplication()
	application := App{
		App:                    app,
		registerLoginDetails:   modeltui.RegisterLogin{},
		registerLoginForm:      tview.NewForm(),
		bankCards:              make([]modeltui.BankCard, 0),
		storeBankCardForm:      tview.NewForm(),
		textsOrBinaries:        make([]modeltui.TextOrBinary, 0),
		storeTextOrBinaryForm:  tview.NewForm(),
		loginsAndPasswords:     make([]modeltui.LoginAndPassword, 0),
		storeLoginPasswordForm: tview.NewForm(),
		removals:               make([]modeltui.Removal, 0),
		removeForm:             tview.NewForm(),
		loginStatus:            tview.NewTextView().SetText("Logged in: NA").SetTextAlign(1).SetScrollable(true),
		operationStatus:        tview.NewTextView().SetText("Nothing to report yet").SetTextAlign(1).SetScrollable(true),
		result:                 tview.NewTextView().SetText("Nothing was requested yet").SetTextAlign(1).SetScrollable(true),
	}
	return application
}

func (a *App) Run() {
	buttonStoreLoginPassword.SetSelectedFunc(func() {
		a.storeLoginPasswordForm.Clear(true)
		a.addLoginPasswordForm()
		pages.SwitchToPage("store_login_password")
	})
	buttonStoreTextBinary.SetSelectedFunc(func() {
		a.storeTextOrBinaryForm.Clear(true)
		a.addTextOrBinaryForm()
		pages.SwitchToPage("store_text_binary")
	})
	buttonStoreBankCard.SetSelectedFunc(func() {
		a.storeBankCardForm.Clear(true)
		a.addBankCardForm()
		pages.SwitchToPage("store_bank_card")
	})
	buttonRemove.SetSelectedFunc(func() {
		a.removeForm.Clear(true)
		a.addRemovalForm()
		pages.SwitchToPage("remove")
	})
	buttonRegisterLogin.SetSelectedFunc(func() {
		a.registerLoginForm.Clear(true)
		a.addRegisterLoginForm()
		pages.SwitchToPage("register_login")
	})
	buttonQuit.SetSelectedFunc(func() {
		a.App.Stop()
	})
	buttonSync.SetSelectedFunc(func() {
		err := a.Sync()
		if err != nil {
			a.operationStatus.SetText(err.Error())
		} else {
			a.operationStatus.SetText("OK")
		}
	})
	buttonGetAllData.SetSelectedFunc(func() {
		err := a.GetAllData()
		pages.SwitchToPage("result")
		if err != nil {
			a.result.SetText(err.Error())
		} else {
			a.result.SetText("OK")
		}
	})
	buttonBackToMainScreen.SetSelectedFunc(func() {
		pages.SwitchToPage("menu")
	})

	resultView := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(a.result, 0, 9, false).
		AddItem(buttonBackToMainScreen, 0, 1, false)

	flex.SetDirection(tview.FlexRow).AddItem(header, 0, 1, false).
		AddItem(a.loginStatus, 0, 1, false).
		AddItem(menu, 0, 1, false).
		AddItem(body, 0, 20, false).
		AddItem(statusHeader, 0, 1, false).
		AddItem(a.operationStatus, 0, 5, false).
		AddItem(footer, 0, 1, false)

	pages.AddPage("menu", flex, true, true)
	pages.AddPage("store_login_password", a.storeLoginPasswordForm, true, false)
	pages.AddPage("store_text_binary", a.storeTextOrBinaryForm, true, false)
	pages.AddPage("store_bank_card", a.storeBankCardForm, true, false)
	pages.AddPage("register_login", a.registerLoginForm, true, false)
	pages.AddPage("remove", a.removeForm, true, false)
	pages.AddPage("result", resultView, true, false)

	if err := a.App.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		log.Panic(err)
	}
}

package main

import (
	"log"
	"net/url"
	"os"
	"strconv"

	"github.com/mxschmitt/playwright-go"
)

type student struct {
	Name  string
	Email string
}

func MustGetenv(name string) string {
	val := os.Getenv(name)
	if val == "" {
		log.Fatalf("Environment variable %s is not set!", name)
	}
	return val
}

var Username = MustGetenv("ADMIN_USERNAME")
var Password = MustGetenv("ADMIN_PASSWORD")

const MaxPerPage = "50"
const LoginPage = "https://www.homeosapiens.eu/login/"
const StudentsPage = "https://www.homeosapiens.eu/dashboard/students/"

func initPage() (playwright.Page, error) {
	pw, err := playwright.Run()
	if err != nil {
		return nil, err
	}

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: new(false),
	})
	if err != nil {
		return nil, err
	}

	return browser.NewPage()
}

func signIn(page playwright.Page) error {
	if _, err := page.Goto(LoginPage); err != nil {
		return err
	}

	if err := page.Locator("input#email").Fill(Username); err != nil {
		return err
	}

	if err := page.Locator("input#password").Fill(Password); err != nil {
		return err
	}

	if err := page.Locator("button[type=submit]").Click(); err != nil {
		return err
	}

	return page.WaitForURL("**/dashboard/")
}

func gotoStudentsPage(page playwright.Page, pageNumber int) error {
	qs := url.Values{
		"page":     {strconv.Itoa(pageNumber)},
		"per_page": {MaxPerPage},
	}

	uri := StudentsPage + "?" + qs.Encode()
	if _, err := page.Goto(uri); err != nil {
		return err
	}

	return page.Locator("table tbody tr[data-row-key]").WaitFor()
}

func getPageCount(page playwright.Page) (int, error) {
	lastPage, err := page.Locator(".ant-pagination-item").Last().InnerText()
	if err != nil {
		return 0, err
	}

	asInt, err := strconv.ParseInt(lastPage, 10, 64)
	return int(asInt), err
}

func getAllStudents(page playwright.Page) ([]student, error) {

}

func main() {
	page, err := initPage()
	if err != nil {
		log.Fatal(err)
	}

	if err := signIn(page); err != nil {
		log.Fatal(err)
	}

	gotoStudentsPage(page, 1)

	pageCount, err := getPageCount(page)
	if err != nil {
		log.Fatal(err)
	}

}

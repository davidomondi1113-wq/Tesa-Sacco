package routes

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"text/template"
	"time"

	_ "modernc.org/sqlite"
)

type RegisterData struct {
	AccountType string `json:"institution_type"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	ConfirmPass string `json:"confirm_pass"`
	Type        string `json:"type"`
	Years       string `json:"years"`
	License     string `json:"license"`
	Kra         string `json:"kra"`
	Phone       string `json:"phone"`
}

type Loan struct {
	Id      string `json:"id"`
	Date    string `json:"date"`
	Amount  string `json:"amount"`
	Purpose string `json:"purpose"`
	Status  string `json:"status"`
}

type User struct {
	Id              string `json:"id"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	InstitutionType string `json:"institution_type"`
	Loans           []Loan `json:"loans"`
}

type DashboardData struct {
	User         User
	TotalLoans   int
	PendingLoans int
	TotalAmount  string
	AccountStatus string
}



var userProfile User
var db *sql.DB
var sessions = make(map[string]User)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	user := getSessionUser(r)
	tmpl := template.Must(template.ParseFiles("template/index.html"))
	tmpl.Execute(w, user)
}

func AboutHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("template/about.html"))
	tmpl.Execute(w, nil)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// Handle GET requests - show registration form
	if r.Method == http.MethodGet {
		tmpl := template.Must(template.ParseFiles("template/register.html"))
		tmpl.Execute(w, nil)
		return
	}

	// Handle POST requests - process registration
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}

	data := RegisterData{
		AccountType: r.FormValue("account-type"),
		Email:       r.FormValue("email"),
		Password:    r.FormValue("pass"),
		ConfirmPass: r.FormValue("confirm"),
		Type:        r.FormValue("type"),
		Years:       r.FormValue("years"),
		License:     r.FormValue("license"),
		Kra:         r.FormValue("kra"),
		Phone:       r.FormValue("phone"),
	}

	// Basic validation
	if data.Email == "" || data.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	// Password confirmation validation
	if data.Password != data.ConfirmPass {
		http.Error(w, "Passwords do not match", http.StatusBadRequest)
		return
	}

	// Check if user already exists
	if userExists(data.Email) {
		http.Error(w, "Email already exists", http.StatusConflict)
		return
	}

	// Create new user
	newUser := User{
		Id:              generateUserID(),
		Email:           data.Email,
		Password:        data.Password,
		InstitutionType: data.AccountType,
		Loans:           []Loan{},
	}

	if err := addUser(newUser); err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	// Redirect to login page after successful registration
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("template/login.html"))
	tmpl.Execute(w, nil)
}

func AuthHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	instType := r.FormValue("inst-type")

	// Find user by email
	user := findUserByEmail(email)
	if user == nil {
		http.Error(w, "User not found: "+email, http.StatusUnauthorized)
		return
	}

	// For demo purposes, accept any password for existing users
	// In production, implement proper password hashing
	if instType == user.InstitutionType {
		userProfile = *user
		setSessionUser(w, userProfile)
		switch instType {
		case "business":
			http.Redirect(w, r, "/business", http.StatusSeeOther)
		case "microfinance":
			http.Redirect(w, r, "/mfi", http.StatusSeeOther)
		default:
			http.Error(w, "Unknown institution type", http.StatusBadRequest)
		}
	} else {
		http.Error(w, "Invalid institution type. Expected: "+user.InstitutionType+", Got: "+instType, http.StatusUnauthorized)
	}
}

func BusinessHandler(w http.ResponseWriter, r *http.Request) {
	user := getSessionUser(r)
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	
	// Calculate dashboard statistics
	totalLoans := len(user.Loans)
	pendingLoans := 0
	totalAmount := 0
	
	for _, loan := range user.Loans {
		if loan.Status == "pending" {
			pendingLoans++
		}
		if amount, err := strconv.Atoi(loan.Amount); err == nil {
			totalAmount += amount
		}
	}
	
	accountStatus := "New Member"
	if totalLoans > 0 {
		accountStatus = "Active Member"
	}
	
	dashboardData := DashboardData{
		User:         *user,
		TotalLoans:   totalLoans,
		PendingLoans: pendingLoans,
		TotalAmount:  fmt.Sprintf("%d", totalAmount),
		AccountStatus: accountStatus,
	}
	
	tmpl := template.Must(template.ParseFiles("template/business_dashboard.html"))
	tmpl.Execute(w, dashboardData)
}

func MfiHandler(w http.ResponseWriter, r *http.Request) {
	user := getSessionUser(r)
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	tmpl := template.Must(template.ParseFiles("template/mfi_dashboard.html"))
	tmpl.Execute(w, user)
}

func ErrorHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("template/error.html"))
	tmpl.Execute(w, nil)
}

func MfiReportHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("template/mfi/records.html"))
	tmpl.Execute(w, userProfile)
}

func MfiReportDownloadHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("template/mfi/download.html"))
	tmpl.Execute(w, userProfile)
}

// business pages routes
func BusinessActiveLoansHandler(w http.ResponseWriter, r *http.Request) {
	user := getSessionUser(r)
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	tmpl := template.Must(template.ParseFiles("template/business/active_loans.html"))
	tmpl.Execute(w, user)
}

func BusinessProfileHandler(w http.ResponseWriter, r *http.Request) {
	user := getSessionUser(r)
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	tmpl := template.Must(template.ParseFiles("template/business/business_profile.html"))
	tmpl.Execute(w, user)
}

func BusinessLoanApplicationHandler(w http.ResponseWriter, r *http.Request) {
	user := getSessionUser(r)
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	tmpl := template.Must(template.ParseFiles("template/business/loan_applications.html"))
	tmpl.Execute(w, user)
}

func BusinessTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	user := getSessionUser(r)
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	tmpl := template.Must(template.ParseFiles("template/business/transactions.html"))
	tmpl.Execute(w, user)
}

func LoanApplicationHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("template/business/loan_applications.html"))
	tmpl.Execute(w, userProfile)
}

func MsiListHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("template/mfi/mfis.html"))
	tmpl.Execute(w, userProfile)
}

func ApplyLoanHandler(w http.ResponseWriter, r *http.Request) {
	user := getSessionUser(r)
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		// Handle loan application submission
		amount := r.FormValue("amount")
		purpose := r.FormValue("purpose")
		
		if amount == "" || purpose == "" {
			http.Error(w, "Amount and purpose are required", http.StatusBadRequest)
			return
		}
		
		// Create new loan application
		newLoan := Loan{
			Id:      "L" + generateUserID(),
			Date:    time.Now().Format("2006-01-02"),
			Amount:  amount,
			Purpose: purpose,
			Status:  "pending",
		}
		
		// Add to user's loans and update session
		user.Loans = append(user.Loans, newLoan)
		sessionId := getCookieValue(r, "session_id")
		sessions[sessionId] = *user
		
		http.Redirect(w, r, "/loan-application?success=1", http.StatusSeeOther)
		return
	}

	tmpl := template.Must(template.ParseFiles("template/mfi/apply_loan.html"))
	tmpl.Execute(w, user)
}

func getCookieValue(r *http.Request, name string) string {
	cookie, err := r.Cookie(name)
	if err != nil {
		return ""
	}
	return cookie.Value
}

// SQLite database functions
func init() {
	initDB()
}

func initDB() {
	var err error
	db, err = sql.Open("sqlite", "users.db")
	if err != nil {
		log.Fatal(err)
	}

	createTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		institution_type TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`

	if _, err := db.Exec(createTable); err != nil {
		log.Fatal(err)
	}
}

func userExists(email string) bool {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", email).Scan(&count)
	return err == nil && count > 0
}

func findUserByEmail(email string) *User {
	var user User

	err := db.QueryRow(`
		SELECT id, email, password, institution_type 
		FROM users WHERE email = ?
	`, email).Scan(
		&user.Id, &user.Email, &user.Password, &user.InstitutionType,
	)

	if err != nil {
		return nil
	}

	// Add sample loans for demo - in production, load from database
	if user.InstitutionType == "business" {
		// No pre-existing loans - users start fresh
		user.Loans = []Loan{}
	}

	return &user
}

func addUser(user User) error {
	_, err := db.Exec(`
		INSERT INTO users (email, password, institution_type)
		VALUES (?, ?, ?)
	`, user.Email, user.Password, user.InstitutionType)
	return err
}

const (
	businessShortCode = "174379"
	passKey           = "bfb279f9aa9bdbcf158e97dd71a467cd2e0c893059b10f78e6b72ada1ed2c919"
)

func generateUserID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano()%10000)
}

func generateTimestamp() string {
	return time.Now().Format("20060102150405")
}

func base64Encode(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

func FetchAccessToken() (string, error) {
	url := "https://sandbox.safaricom.co.ke/oauth/v1/generate?grant_type=client_credentials"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("Authorization", "Basic eUFHd3U0UnJoQ2tSNjBWUndYQUdHdlJVampxOHd0b2dBc1ZaaUdJbGhhRlVmd3dCOjdQbkxiVnpqeXBJWGUycWNOOGRpbXpHeXFSR1VRODVteUVzZ3RoZWp6Z0hTMHVqOWpvSTFUWVZRR2UyRWFBenE=")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return "", err
	}

	accessToken, ok := response["access_token"].(string)
	if !ok {
		return "", fmt.Errorf("access token not found")
	}

	return accessToken, nil
}

func makePayment(accessToken, amount, phoneNumber string) error {
	timestamp := generateTimestamp()
	password := base64Encode(businessShortCode + passKey + timestamp)

	payload := map[string]interface{}{
		"BusinessShortCode": businessShortCode,
		"Password":          password,
		"Timestamp":         timestamp,
		"TransactionType":   "CustomerPayBillOnline",
		"Amount":            amount,
		"PartyA":            phoneNumber,
		"PartyB":            businessShortCode,
		"PhoneNumber":       phoneNumber,
		"CallBackURL":       "https://mydomain.com/path",
		"AccountReference":  "TESAA",
		"TransactionDesc":   "Payment",
	}

	payloadBytes, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", "https://sandbox.safaricom.co.ke/mpesa/stkpush/v1/processrequest", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func MakePaymentHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, `{"status":"error","message":"Invalid request method"}`)
		return
	}

	r.ParseForm()
	amount := r.FormValue("amount")
	phoneNumber := r.FormValue("phoneNumber")

	if amount == "" || phoneNumber == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"status":"error","message":"Amount and phone number required"}`)
		return
	}

	accessToken, err := FetchAccessToken()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"status":"error","message":"Failed to get access token: %s"}`, err.Error())
		return
	}

	if err := makePayment(accessToken, amount, phoneNumber); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"status":"error","message":"Payment failed: %s"}`, err.Error())
		return
	}

	fmt.Fprintf(w, `{"status":"success","message":"Payment initiated successfully"}`)
}

type PendingLoan struct {
	Id       string `json:"id"`
	Borrower string `json:"borrower"`
	Amount   string `json:"amount"`
	Purpose  string `json:"purpose"`
	Date     string `json:"date"`
	Status   string `json:"status"`
}

func LoanManagementHandler(w http.ResponseWriter, r *http.Request) {
	user := getSessionUser(r)
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		// Handle loan approval/rejection
		loanId := r.FormValue("loan_id")
		action := r.FormValue("action")
		borrowerEmail := r.FormValue("borrower_email")
		
		// Update loan status (in production, update database)
		for sessionId, sessionUser := range sessions {
			if sessionUser.Email == borrowerEmail {
				for i, loan := range sessionUser.Loans {
					if loan.Id == loanId {
						sessionUser.Loans[i].Status = action
						sessions[sessionId] = sessionUser
						break
					}
				}
				break
			}
		}
		
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":"success","message":"Loan %s successfully", "action":"%s"}`, action, action)
		return
	}

	// Get all pending loans from all business users
	var pendingLoans []PendingLoan
	for _, sessionUser := range sessions {
		if sessionUser.InstitutionType == "business" {
			for _, loan := range sessionUser.Loans {
				if loan.Status == "pending" {
					pendingLoans = append(pendingLoans, PendingLoan{
						Id:       loan.Id,
						Borrower: sessionUser.Email,
						Amount:   loan.Amount,
						Purpose:  loan.Purpose,
						Date:     loan.Date,
						Status:   loan.Status,
					})
				}
			}
		}
	}

	tmpl := template.Must(template.ParseFiles("template/mfi/loan_management.html"))
	tmpl.Execute(w, map[string]interface{}{
		"User":         user,
		"PendingLoans": pendingLoans,
	})
}

func MfiDirectoryHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("template/mfi/mfi_directory.html"))
	tmpl.Execute(w, userProfile)
}

func ApplicationProcessingHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("template/mfi/application_processing.html"))
	tmpl.Execute(w, userProfile)
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("template/search.html")
	if err != nil {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func SearchApiHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, `{"status":"error","message":"Invalid request method"}`)
		return
	}

	query := r.URL.Query().Get("q")
	if query == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"status":"error","message":"Search query required"}`)
		return
	}

	// Mock search results
	results := []map[string]interface{}{
		{"type": "loan", "id": "L001", "borrower": "John Doe", "amount": "5000", "status": "Active"},
		{"type": "user", "id": "U001", "name": "Jane Smith", "email": "jane@example.com", "category": "business"},
		{"type": "institution", "id": "I001", "name": "Equity Bank", "category": "microfinance", "location": "Nairobi"},
	}

	response := map[string]interface{}{
		"status": "success",
		"query": query,
		"results": results,
	}

	jsonResponse, _ := json.Marshal(response)
	fmt.Fprintf(w, string(jsonResponse))
}

func getSessionUser(r *http.Request) *User {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return nil
	}
	user, exists := sessions[cookie.Value]
	if !exists {
		return nil
	}
	return &user
}

func setSessionUser(w http.ResponseWriter, user User) {
	sessionID := generateUserID()
	sessions[sessionID] = user
	http.SetCookie(w, &http.Cookie{
		Name:  "session_id",
		Value: sessionID,
		Path:  "/",
	})
}

func MfiProfileHandler(w http.ResponseWriter, r *http.Request) {
	user := getSessionUser(r)
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	tmpl := template.Must(template.ParseFiles("template/mfi/mfi_profile.html"))
	tmpl.Execute(w, user)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err == nil {
		delete(sessions, cookie.Value)
	}
	http.SetCookie(w, &http.Cookie{
		Name:   "session_id",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}





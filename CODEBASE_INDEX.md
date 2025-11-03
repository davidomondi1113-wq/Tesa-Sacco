# TESAA Microfinance Platform - Codebase Index

## Overview
TESAA is a comprehensive microfinance platform built with Go (Golang) that connects lenders and borrowers, facilitating loan applications, approvals, and management. The platform serves both business entities and microfinance institutions with dedicated dashboards and functionality.

## Technology Stack
- **Backend**: Go 1.24+ with net/http package
- **Database**: SQLite with modernc.org/sqlite driver
- **Frontend**: HTML5, CSS3, Bootstrap 5.3, JavaScript
- **Payment Integration**: M-Pesa API (Safaricom)
- **Template Engine**: Go's built-in html/template

## Project Structure

```
Tesa-Sacco/
├── main.go                     # Application entry point
├── go.mod                      # Go module dependencies
├── go.sum                      # Dependency checksums
├── users.db                    # SQLite database (auto-generated)
├── tesaa.exe                   # Compiled executable
├── README.md                   # Project documentation
├── PROJECT_DOCUMENTATION.md    # Comprehensive project docs
├── LICENSE                     # MIT License
├── CODEBASE_INDEX.md          # This index file
├── routes/
│   └── routes.go              # HTTP handlers and business logic
├── template/                  # HTML templates
│   ├── index.html             # Home page
│   ├── login.html             # Login form
│   ├── register.html          # Registration form
│   ├── business_dashboard.html # Business main dashboard
│   ├── mfi_dashboard.html     # MFI main dashboard
│   ├── about.html             # About page
│   ├── error.html             # Error page
│   ├── payment.html           # Payment interface
│   ├── search.html            # Search functionality
│   ├── business/              # Business-specific templates
│   │   ├── active_loans.html  # Active loans management
│   │   ├── apply_loan.html    # Loan application form
│   │   ├── business_profile.html # Profile management
│   │   ├── loan_applications.html # Application tracking
│   │   ├── login.html         # Business login
│   │   └── transactions.html  # Transaction history
│   └── mfi/                   # MFI-specific templates
│       ├── application_processing.html # Process applications
│       ├── apply_loan.html    # MFI loan interface
│       ├── download.html      # Report downloads
│       ├── loan_management.html # Loan portfolio management
│       ├── mfi_directory.html # Institution directory
│       ├── mfis.html          # MFI listings
│       └── records.html       # Reports and analytics
├── assets/                    # Static assets
│   ├── businessowner.jpg      # Business owner images
│   ├── equity.png             # Bank logos
│   ├── faulu.png              # MFI logos
│   ├── mpesa.jpg              # Payment method images
│   └── [other image assets]
├── api/
│   └── database/
│       ├── configure.sh       # Database setup script
│       └── users.json         # JSON user data (legacy)
└── auth/
    └── auth.html              # Authentication template
```

## Core Components

### 1. Main Application (main.go)
**Purpose**: Application entry point and routing configuration
**Key Features**:
- HTTP server setup on port 8080
- Route mapping to handler functions
- Static asset serving
- Template parsing

**Routes Defined**:
```go
/ - HomeHandler (Landing page)
/about - AboutHandler
/register - RegisterHandler (GET/POST)
/login - LoginHandler
/auth - AuthHandler (Authentication)
/business - BusinessHandler (Business dashboard)
/mfi - MfiHandler (MFI dashboard)
/active-loans - BusinessActiveLoansHandler
/business-profile - BusinessProfileHandler
/loan-application - BusinessLoanApplicationHandler
/business-transactions - BusinessTransactionsHandler
/records - MfiReportHandler
/download - MfiReportDownloadHandler
/applyloan - ApplyLoanHandler
/mfis - MsiListHandler
/make-payment - MakePaymentHandler
/loan-management - LoanManagementHandler
/mfi-directory - MfiDirectoryHandler
/application-processing - ApplicationProcessingHandler
/search - SearchHandler
/search-api - SearchApiHandler
```

### 2. Business Logic (routes/routes.go)
**Purpose**: Core application logic and HTTP handlers
**Key Components**:

#### Data Structures
```go
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
```

#### Database Functions
- `initDB()`: Initialize SQLite database and create tables
- `userExists(email)`: Check if user exists
- `findUserByEmail(email)`: Retrieve user by email
- `addUser(user)`: Add new user to database

#### Payment Integration
- M-Pesa STK Push implementation
- OAuth token management
- Transaction processing
- Payment verification

### 3. Database Schema
**SQLite Database**: `users.db`
```sql
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    institution_type TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### 4. Frontend Templates

#### Landing Page (template/index.html)
- Responsive Bootstrap design
- Hero carousel with platform benefits
- Feature highlights and statistics
- Success stories section
- Call-to-action buttons

#### Registration System (template/register.html)
- Dual registration for businesses and MFIs
- Form validation (client and server-side)
- Business information collection
- Institution verification fields

#### Business Dashboard (template/business_dashboard.html)
- Sidebar navigation
- Loan portfolio overview
- Transaction history
- Payment scheduling
- Profile management access

#### MFI Dashboard (template/mfi_dashboard.html)
- Administrative interface
- Loan application processing
- Portfolio analytics
- Report generation
- Institution directory

## Key Features

### 1. User Management
- **Registration**: Separate flows for businesses and MFIs
- **Authentication**: Email/password with institution type selection
- **Session Management**: User profile storage and routing
- **Profile Management**: Business and institutional information updates

### 2. Loan Management
- **Application Process**: Streamlined loan application forms
- **Status Tracking**: Real-time application and loan status
- **Portfolio Management**: Active loan monitoring and analytics
- **Payment Processing**: M-Pesa integration for loan payments
- **Repayment Scheduling**: Automated payment reminders and tracking

### 3. Dashboard Systems
#### Business Dashboard Features:
- Active loan overview with progress indicators
- Payment history and transaction records
- Loan application status tracking
- Profile and business information management
- Payment processing with M-Pesa integration

#### MFI Dashboard Features:
- Loan application review and processing
- Portfolio performance analytics
- Report generation and export
- Institution directory and networking
- Risk assessment and monitoring tools

### 4. Payment Integration
- **M-Pesa STK Push**: Automated payment requests
- **Transaction Verification**: Real-time payment confirmation
- **Payment History**: Comprehensive transaction records
- **Receipt Generation**: Automated payment receipts

### 5. Search and Analytics
- **Advanced Search**: Filter loans, users, and institutions
- **Reporting System**: Financial and operational reports
- **Export Capabilities**: PDF, Excel, CSV formats
- **Performance Metrics**: Portfolio analytics and KPIs

## Security Features
- **Input Validation**: Form data sanitization and validation
- **Session Management**: Secure user session handling
- **Database Security**: Parameterized queries to prevent SQL injection
- **Authentication**: Email/password verification system
- **Data Protection**: Secure credential handling for payment APIs

## API Integrations
### M-Pesa Payment API
- **Endpoint**: Safaricom M-Pesa STK Push API
- **Authentication**: OAuth 2.0 token-based
- **Features**: Payment initiation, status checking, transaction verification
- **Security**: Encrypted credential handling and secure callbacks

## Configuration Files
- **go.mod**: Go module dependencies and version management
- **go.sum**: Dependency checksums for security verification
- **configure.sh**: Database setup and configuration script

## Development Guidelines
### Code Organization
- **Separation of Concerns**: Clear separation between routing, business logic, and data access
- **Template Structure**: Organized template hierarchy for maintainability
- **Asset Management**: Centralized static asset serving
- **Error Handling**: Comprehensive error handling and user feedback

### Best Practices Implemented
- **Responsive Design**: Mobile-first Bootstrap implementation
- **Form Validation**: Client and server-side validation
- **Security**: Input sanitization and secure data handling
- **Performance**: Efficient database queries and template rendering
- **Maintainability**: Clear code structure and documentation

## Deployment Information
- **Server Requirements**: Go 1.22+, SQLite support
- **Port Configuration**: Default port 8080
- **Database**: Auto-generated SQLite database
- **Static Assets**: Served from /assets/ directory
- **Build Process**: `go build -o tesaa main.go`

## Future Enhancement Opportunities
1. **Mobile Application**: Native iOS/Android apps
2. **Advanced Analytics**: Machine learning for credit scoring
3. **Blockchain Integration**: Enhanced security and transparency
4. **Multi-language Support**: Localization for different regions
5. **Advanced Security**: Multi-factor authentication, encryption
6. **API Development**: RESTful APIs for third-party integrations
7. **Cloud Deployment**: Scalable cloud infrastructure
8. **Real-time Features**: WebSocket integration for live updates

## Documentation References
- **README.md**: Installation and usage instructions
- **PROJECT_DOCUMENTATION.md**: Comprehensive project documentation
- **LICENSE**: MIT License terms and conditions

## Contact Information
- **Developer**: SMARTEDGE COMPUTERS & PRINTING HUB
- **Email**: davidomondi1113@gmail.com
- **Project**: TESAA Microfinance Platform
- **Version**: 1.4.0 (Current)

---

*This index provides a comprehensive overview of the TESAA microfinance platform codebase. For detailed implementation information, refer to the individual source files and project documentation.*
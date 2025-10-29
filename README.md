# TESAA - Microfinance Platform

## Overview

TESAA is a comprehensive microfinance platform that connects lenders and borrowers, facilitating loan applications, approvals, and management. The platform serves both business entities and microfinance institutions with dedicated dashboards and functionality.

## Features

- **User Registration & Authentication**: Secure registration for businesses and microfinance institutions
- **Dual Dashboard System**: Separate interfaces for businesses and MFIs
- **Loan Management**: Application, tracking, and management of loans
- **SQLite Database**: Reliable data persistence with pure Go SQLite driver
- **Responsive Design**: Bootstrap-powered responsive web interface

## Technology Stack

- **Backend**: Go (Golang) 1.24+
- **Database**: SQLite with modernc.org/sqlite driver
- **Frontend**: HTML5, CSS3, Bootstrap 5.3, JavaScript
- **Template Engine**: Go's built-in html/template

## Prerequisites

- Go 1.22 or higher
- Git (for cloning the repository)
- Web browser (Chrome, Firefox, Safari, Edge)

## Installation

### 1. Clone the Repository

```bash
git clone https://github.com/your-username/tesa-sacco-main.git
cd tesa-sacco-main/tesa-sacco-main
```

### 2. Install Dependencies

```bash
go mod tidy
```

### 3. Verify Installation

```bash
go version
```

Ensure Go version is 1.22 or higher.

## Configuration

### Database Setup

The application automatically creates and configures the SQLite database on first run. No manual database setup is required.



## Running the Application

### Development Mode

```bash
go run main.go
```

### Production Build

```bash
go build -o tesaa main.go
./tesaa
```

The server will start on `http://localhost:8080`

## Usage Guide

### 1. Accessing the Application

Open your web browser and navigate to:
```
http://localhost:8080
```

### 2. User Registration

#### For Businesses:
1. Click "Register" in the navigation menu
2. Select "Business" from the account type dropdown
3. Fill in the required information:
   - Full Name/Organization Name
   - Email Address
   - Password and Confirmation
   - Business Type
   - Years in Operation
   - Business License Number
   - KRA PIN
   - Phone Number
4. Click "Sign up"
5. You will be redirected to the login page upon successful registration

#### For Microfinance Institutions:
1. Click "Register" in the navigation menu
2. Select "Microfinance Institution" from the account type dropdown
3. Complete the registration form with institutional details
4. Submit the form

### 3. User Login

1. Navigate to the login page
2. Select your account type (Business or Microfinance Institution)
3. Enter your email and password
4. Click "Login"
5. You will be redirected to your respective dashboard

### 4. Business Dashboard Features

After logging in as a business user, you can access:

- **Active Loans**: View current loan status and details
- **Loan Applications**: Submit new loan applications
- **Business Profile**: Manage your business information
- **Transactions**: View payment history and transaction records

### 5. Microfinance Institution Dashboard Features

MFI users have access to:

- **Loan Management**: Review and process loan applications
- **Reports**: Generate and download financial reports
- **MFI Directory**: View other participating institutions
- **Application Processing**: Approve or reject loan requests

### 6. Navigation Structure

```
/                    - Home page with platform overview
/register           - User registration form
/login              - User authentication
/business           - Business dashboard
/mfi                - MFI dashboard
/active-loans       - Business active loans page
/business-profile   - Business profile management
/loan-application   - Loan application form
/business-transactions - Transaction history
/records            - MFI reports page
/download           - Report download functionality
/mfis               - MFI directory
/applyloan          - Loan application interface
```

## File Structure

```
tesa-sacco-main/
├── main.go                 # Application entry point
├── go.mod                  # Go module dependencies
├── go.sum                  # Dependency checksums
├── users.db               # SQLite database (auto-generated)
├── routes/
│   └── routes.go          # HTTP handlers and business logic
├── template/              # HTML templates
│   ├── index.html         # Home page
│   ├── login.html         # Login form
│   ├── register.html      # Registration form
│   ├── business_dashboard.html
│   ├── mfi_dashboard.html
│   ├── business/          # Business-specific templates
│   └── mfi/              # MFI-specific templates
└── assets/               # Static assets (images, etc.)
```

## API Endpoints

### Authentication Endpoints
- `GET /register` - Display registration form
- `POST /register` - Process user registration
- `GET /login` - Display login form
- `POST /auth` - Authenticate user

### Business Endpoints
- `GET /business` - Business dashboard
- `GET /active-loans` - Active loans page
- `GET /business-profile` - Profile management
- `GET /loan-application` - Loan application form
- `GET /business-transactions` - Transaction history

### MFI Endpoints
- `GET /mfi` - MFI dashboard
- `GET /records` - Reports page
- `GET /download` - Download reports
- `GET /mfis` - MFI directory
- `GET /applyloan` - Loan application interface



## Database Schema

### Users Table
```sql
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    institution_type TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

## Troubleshooting

### Common Issues

1. **Port 8080 already in use**
   ```bash
   # Kill existing processes
   netstat -ano | findstr :8080
   taskkill /PID <process_id> /F
   ```

2. **Database connection errors**
   - Ensure the application has write permissions in the directory
   - Check if `users.db` file is created successfully

3. **Template not found errors**
   - Verify all template files exist in the `template/` directory
   - Check file paths are correct in the route handlers

### Development Tips

- Use `go run main.go` for development with automatic recompilation
- Check console output for detailed error messages
- Verify form field names match the expected parameters in handlers
- Test with different user types (business vs MFI) to ensure proper routing

## Security Considerations

- Passwords are stored in plain text (implement hashing for production)
- Add HTTPS support for production deployment
- Implement proper session management
- Add CSRF protection for forms
- Validate and sanitize all user inputs

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## License

Copyright 2025 TESAA. All rights reserved.

## Support

For technical support or questions, contact:
- **Developer**: SMARTEDGE COMPUTERS & PRINTING HUB
- **Email**: davidomondi1113@gmail.com

## Version History

- **v1.0.0** - Initial release with basic functionality
- **v1.1.0** - Added SQLite database integration
- **v1.2.0** - Implemented M-Pesa payment processing
- **v1.3.0** - Enhanced user interface and responsive design
- **v1.4.0** - Removed M-Pesa payment processing, focused on loan display functionality
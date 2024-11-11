# Snippetbox (Let's Go)

> A simple web app for storing and sharing code snippets, built with Go.

## Features

- **User Accounts**: Register, login, and manage your snippets.
- **Snippet Management**: Create, view, edit, and delete code snippets.
- **Simple UI**: Clean, responsive design that works on mobile and desktop.
- **Database**: Uses SQLite to store snippets and user data.

## Getting Started

### Requirements

- Go 1.22.2+
- SQLite3

### Installation

Clone the repo
```bash
git clone https://github.com/Danvs60/snippetbox.git
cd snippetbox
```
Set up the database Run the following command to create the SQLite database using the provided schema:
```bash
sqlite3 snippetbox.db < schema.sql
```
Run the app Start the application with:
```bash
go run ./cmd/web
```
The app should now be accessible at http://localhost:4000.

### Running Tests

To run tests, use the following command:
```bash
go test ./...
```
This will run all tests in the project. Ensure that any necessary test data or test database setup is completed before running tests.

## Project Structure

- cmd/: Main application entry
- internal/: Contains handlers, models, and templates for the app
- ui/: Static files like CSS and JavaScript
  
## Usage

- Register or log in to create snippets.
- View, edit, or delete your snippets from the dashboard.

<div align="center">

🚀 GoCourse: Teacher Management REST API 🚀

A high-performance, lightweight REST API built in Go using only the standard library. No external routers, just pure Go.

</div>

<p align="center">
<img src="https://www.google.com/search?q=https://img.shields.io/badge/Go-1.18%252B-blue.svg%3Fstyle%3Dfor-the-badge%26logo%3Dgo" alt="Go Version">
<img src="https://www.google.com/search?q=https://img.shields.io/badge/License-MIT-yellow.svg%3Fstyle%3Dfor-the-badge" alt="License: MIT">
<img src="https://www.google.com/search?q=https://img.shields.io/badge/PRs-Welcome-brightgreen.svg%3Fstyle%3Dfor-the-badge" alt="PRs Welcome">
</p>

Welcome to GoCourse! This project is a lean, powerful RESTful API for managing teacher records. It's built from the ground up using Go's native net/http package to demonstrate a clean, dependency-minimal approach to building web services.

This API handles all core CRUD (Create, Read) operations, connects to a MySQL/MariaDB backend, and uses custom middleware for security (CORS).

✅ Core Features

Zero Dependencies: No external routers (like Gorilla Mux or Chi). Just the standard library!

Database Integration: Connects to MySQL/MariaDB using the database/sql package.

RESTful Endpoints: Clean, predictable API design for POST and GET operations.

Custom Middleware: Includes a from-scratch CORS middleware handler.

Environment-Ready: Uses .env files for secure credential management.

JSON Handling: Efficiently decodes and encodes JSON request/response bodies.

🛠️ Tech Stack

<p align="center">
<img src="https://www.google.com/search?q=https://img.shields.io/badge/Go-00ADD8%3Fstyle%3Dfor-the-badge%26logo%3Dgo%26logoColor%3Dwhite" alt="Go">
<img src="https://img.shields.io/badge/MySQL-4479A1?style=for-the-badge&logo=mysql&logoColor=white" alt="MySQL">
<img src="https://www.google.com/search?q=https://img.shields.io/badge/.env-ECD53F%3Fstyle%3Dfor-the-badge%26logo%3Ddotenv%26logoColor%3Dblack" alt="DotEnv">
</p>

🏃‍♂️ Getting Started

Follow these steps to get the project running on your local machine.

1. 🔌 Prerequisites

Go (version 1.18 or higher)

MySQL or MariaDB

Git

2. 📂 Clone the Repository

git clone [https://github.com/ProNeverDies/gocourse.git](https://github.com/ProNeverDies/gocourse.git)
cd gocourse


3. 📦 Install Dependencies

This project has minimal dependencies. go mod tidy will fetch them for you.

go mod tidy


4. 🗄️ Set Up the Database

You need to create your database and the teachers table.

-- 1. Log in to your database (e.g., mysql -u root -p)
-- 2. Create the database
CREATE DATABASE gocourse_db;
USE gocourse_db;

-- 3. Create the teachers table
CREATE TABLE IF NOT EXISTS teachers (
    id INT AUTO_INCREMENT PRIMARY KEY,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    class VARCHAR(50),
    subject VARCHAR(100)
);


5. 🔑 Configure Environment

Create a .env file in the root of the project and add your database credentials.

DB_USER=your_db_username
DB_PASSWORD=your_db_password
HOST=127.0.0.1
DB_PORT=3306
DB_NAME=gocourse_db


6. 🚀 Run the Application

go run ./cmd/web/main.go


The server will start on http://localhost:3000.

⚡️ API Endpoints

1. Create Teacher(s)

Creates one or more new teachers in the database.

<p>
<img src="https://www.google.com/search?q=https://img.shields.io/badge/POST-success%3Fstyle%3Dfor-the-badge" alt="POST">
<code style="font-size: 1.1em; padding: 5px; border-radius: 5px;">/teachers/</code>
</p>

Request Body (JSON):

[
    {
        "first_name": "Alice",
        "last_name": "Brown",
        "email": "alice@example.com",
        "class": "6A",
        "subject": "World History"
    },
    {
        "first_name": "Bob",
        "last_name": "White",
        "email": "bob@example.com",
        "class": "9B",
        "subject": "Quant Mechanics"
    }
]


Success Response (201 Created):

{
    "status": "success",
    "count": 2,
    "data": [
        {
            "id": 1,
            "first_name": "Alice",
            "last_name": "Brown",
            "email": "alice@example.com",
            "class": "6A",
            "subject": "World History"
        },
        {
            "id": 2,
            "first_name": "Bob",
            "last_name": "White",
            "email": "bob@example.com",
            "class": "9B",
            "subject": "Quant Mechanics"
        }
    ]
}


2. Get All Teachers

Retrieves a list of all teachers. Can be filtered by query parameters.

<p>
<img src="https://www.google.com/search?q=https://img.shields.io/badge/GET-blue%3Fstyle%3Dfor-the-badge" alt="GET">
<code style="font-size: 1.1em; padding: 5px; border-radius: 5px;">/teachers/</code>
</p>

Query Parameters (Optional):

first_name (string): Filters by first name (e.g., /teachers/?first_name=Alice)

last_name (string): Filters by last name (e.g., /teachers/?last_name=Brown)

Success Response (200 OK):

{
    "status": "success",
    "count": 1,
    "data": [
        {
            "id": 1,
            "first_name": "Alice",
            "last_name": "Brown",
            "email": "alice@example.com",
            "class": "6A",
            "subject": "World History"
        }
    ]
}


3. Get Teacher by ID

Retrieves a single teacher by their unique ID.

<p>
<img src="https://www.google.com/search?q=https://img.shields.io/badge/GET-blue%3Fstyle%3Dfor-the-badge" alt="GET">
<code style="font-size: 1.1em; padding: 5px; border-radius: 5px;">/teachers/{id}</code>
</p>

Example URL:
/teachers/1

Success Response (200 OK):

{
    "id": 1,
    "first_name": "Alice",
    "last_name": "Brown",
    "email": "alice@example.com",
    "class": "6A",
    "subject": "World History"
}


Error Response (404 Not Found):

Teacher not found


📂 Project Structure

gocourse/
├── cmd/web/            # Main application entry point
│   └── main.go
├── internal/
│   ├── handlers/       # HTTP request handlers (teacher_handler.go)
│   ├── middlewares/    # Custom middleware (Cors.go)
│   ├── models/         # Struct definitions (models.Teacher)
│   └── repository/     # Database logic
│       └── sqlconnect/
│           └── sqlconnect.go
├── .env                # Environment variables (private)
├── go.mod              # Go module dependencies
├── go.sum
└── README.md


🤝 Contributing

Pull requests are welcome! For major changes, please open an issue first to discuss what you would like to change.

📜 License

This project is licensed under the MIT License - see the LICENSE file for details.

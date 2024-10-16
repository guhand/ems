# Employee Management System

## Overview

The **Employee Management System** is a RESTful API service built with Go. It provides endpoints for managing employees, departments, roles, and permissions. This system allows administrators to create, update, and delete employee records, manage department assignments, and handle role-based access control (RBAC). Additional features include leave management, permission management, notice board, and dashboard functionalities.

## Features

- **Employee Management**: CRUD operations for employee records.
- **Department Management**: CRUD operations for department records.
- **Role-Based Access Control (RBAC)**: Manage role-based access and permissions.
- **Employee-Department Association**: Link employees with departments.
- **Dashboard Data**: View summary and statistics about employees and departments.
- **Leave Management**: Handle employee leave requests and approvals.
- **Permission Management**: Manage employee-specific permissions.
- **Notice Board Management**: Post and manage notices visible to employees.
- **Pagination and Filtering**: Support for pagination and filtering of large datasets.

## Technologies Used

- **Go**: 1.21.0
- **Gin**: HTTP web framework for Go.
- **GORM**: ORM library for Go.
- **SQLite**: Database used in the project.

## Getting Started

### Prerequisites

- Go 1.21.0 installed.
- SQLite installed[Optional].

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/guhand/ems.git
   ```

2. Change the directory:

   ```bash
   cd employee-management-system
   ```

3. Install the necessary Go modules:

   ```bash
   go mod tidy
   ```

4. Run the service:

   ```bash
   go run main.go
   ```

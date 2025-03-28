# 🌐 DataForge: JSON Data Management CLI

## 📖 Overview

DataForge is a powerful and user-friendly Command Line Interface (CLI) tool designed to simplify JSON data file management. Whether you're working with language translations, configuration files, or any key-value JSON data, DataForge provides an intuitive interface to interact with your JSON files.

## ✨ Features

- 📂 **Load JSON Files**: Easily load your JSON files
- 🔍 **Smart Value Search**: Search through values with precise matching
- 📊 **Duplicate Detection**: Identify and list duplicate values across keys
- ➕ **Add New Entries**: Seamlessly add new keys
- ➖ **Remove Entries**: Quickly remove unnecessary keys
- 🎨 **Interactive UI**: Colorful terminal interface with emojis for a delightful user experience

## 🚀 Prerequisites

- Go 1.16+
- Dependencies:
  - github.com/fatih/color
  - github.com/manifoldco/promptui

## 🛠 Installation

1. Clone the repository
```bash
git clone https://github.com/your-username/dataforge.git
cd dataforge
```

2. Install dependencies
```bash
go mod init dataforge
go get github.com/fatih/color
go get github.com/manifoldco/promptui
go mod tidy
```

3. Build the application
```bash
go build
```

## 💻 Usage

Run the application:
```bash
./dataforge
```

### Workflow

1. Select "Load JSON File" and provide the path to your JSON file
2. Choose from options:
   - Search values
   - Find duplicate values
   - Add new keys
   - Remove existing keys
   - Exit the application

## 🌍 Example JSON Structures

### Language Translation
```json
{
  "welcome_message": "Welcome",
  "login_button": "Log In",
  "signup_button": "Sign Up"
}
```

### Configuration Example
```json
{
  "database_host": "localhost",
  "max_connections": "100",
  "cache_enabled": "true"
}
```

## 🆕 New Features

### Value Search
- Search through JSON values with case-insensitive matching
- Quickly find entries containing specific text

### Duplicate Value Detection
- Identify keys with identical values
- Helps in finding redundant or mistakenly duplicated entries

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📝 License

Distributed under the MIT License. See `LICENSE` for more information.

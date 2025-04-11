# Birthday Script

A simple project to manage and display birthdays using a Go-based backend and a Python script for generating `.ics` files for calendar integration.

---

## Features

- **Go Backend**:

  - Displays a list of birthdays in a web interface.
  - Allows adding new birthdays dynamically.
  - Stores birthdays in a JSON file (`birthdays.json`).

- **Python Script**:
  - Reads birthdays from a CSV file.
  - Generates a `.ics` file with yearly recurring events for calendar integration.

---

## Requirements

### Backend (Go)

- Go 1.16 or higher
- A JSON file (`birthdays.json`) to store birthdays.

### Python Script

- Python 3.8 or higher
- `vobject` library for `.ics` file generation.

---

## Installation

### Backend (Go)

1. Clone the repository:
   ```bash
   git clone https://github.com/HanniOfHyrule/birthday_list
   cd Birthday_script
   ```

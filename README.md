# AetherGuard 🛡️

AetherGuard is a lightweight, high-performance multi-language security monitoring system designed for modern infrastructure. It provides real-time traffic ingestion, AI-driven threat analysis, and a sleek management dashboard.

## Features

*   **Multi-Language Architecture**: Core engine in Go, AI module in Python.
*   **High Performance**: Go-based ingestion capable of handling high request volumes.
*   **AI Threat Detection**: Pattern matching and score-based analysis for common exploits.
*   **Live Dashboard**: Modern, dark-themed UI for real-time system monitoring.
*   **RESTful API**: Simple integration with existing logging systems.

## Modules & Ports

| Module | Description | Port |
| :--- | :--- | :--- |
| `core-system` | Go Ingestion Engine | 8080 |
| `ai-analyzer` | Python Analysis API | 5000 |
| `dashboard` | HTML/JS Management UI | Static |

## Configuration

`config.json`

```json
{
  "system": {
    "name": "AetherGuard Core",
    "version": "1.0.0",
    "port": 8080
  },
  "ai_settings": {
    "sensitivity": 0.8,
    "log_history": true
  }
}
```

## Installation

1.  Clone the repository to your local machine.
2.  Install Python dependencies: `pip install flask flask-cors`
3.  Ensure Go is installed for the core engine.

## Build & Run

### 1. Start Core Engine
```bash
cd core-system
go run main.go
```

### 2. Start AI Analyzer
```bash
cd ai-analyzer
python app.py
```

### 3. Open Dashboard
Simply open `dashboard/index.html` in any modern web browser.

## Author

*   Nrleryx

## Update
*   Improved AI analysis logic
*   Enhanced Dashboard UI
*   Optimized Go core performance
*   Added config.json support

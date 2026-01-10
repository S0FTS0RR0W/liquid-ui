# Liquid UI

A web-based dashboard for monitoring and controlling liquid coolers, built on top of [liquidctl](https://github.com/liquidctl/liquidctl).

## Features

- **Device Discovery**: Automatically lists compatible devices connected to your system.
- **Status Monitoring**: Real-time updates for liquid temperature, fan RPM, and pump RPM.
- **Profile Management**: Create, save, and apply fan curve profiles.

## Prerequisites

- **[liquidctl](https://github.com/liquidctl/liquidctl)**: Must be installed and accessible in your system PATH.
- **Go**: For running the backend server.
- **Node.js & npm**: For the SvelteKit frontend.

## Installation

1.  **Clone the repository**

    ```bash
    git clone https://github.com/S0FTS0RR0W/liquid-ui.git
    cd liquid-ui
    ```

2.  **Start the Backend**

    Navigate to the server directory and run the Go application:

    ```bash
    cd backend/cmd/server
    # Ensure dependencies are installed
    go mod tidy
    go run .
    ```

    The backend listens on `http://localhost:8765`.

3.  **Start the Frontend**

    In a new terminal, navigate to the frontend directory:

    ```bash
    cd frontend
    npm install
    npm run dev
    ```

    The dashboard will be available at `http://localhost:5173`.# liquid-ui

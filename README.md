# wallet

A REST-API built handling a simple wallet application.

## Overview

The wallet service is built using Go and PostgreSQL, providing functionality to handle user wallets and balances.

## Prerequisites

To run this project locally, ensure you have the following installed:

- Docker
- Docker Compose

## Setup Instructions

1. Clone the repository:

   ```bash
   git clone https://github.com/sophiacaspar/wallet.git
   ```

2. Navigate to the project directory.
3. Build and run the application using Docker Compose:

```bash
  docker-compose up --build
```

## Usage

Once the application is up and running, you can access the wallet service API:

Base URL: http://localhost:4000
Endpoints: - GET: http://localhost:4000/v1/wallet/{id} - PUT: http://localhost:4000/v1/wallet/subtract/{id} - PUT: http://localhost:4000/v1/wallet/add/{id}

## Testing

To run tests, use the following command:

```bash
docker-compose run app go test ./...
```

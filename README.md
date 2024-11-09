# URL Shortener
[![Go Version](https://img.shields.io/badge/Go-1.17-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](https://opensource.org/licenses/MIT)
[![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)](https://github.com/yourusername/url-shortener/actions)

A simple URL shortener service built with Go. This service allows users to shorten long URLs and retrieve them later by using a unique short code.

## Table of Contents
- [Features](#features)
- [Getting Started](#getting-started)
- [Installation](#installation)
- [Configuration](#configuration)
- [Usage](#usage)
- [API Endpoints](#api-endpoints)
- [Contributing](#contributing)
- [License](#license)

## Features
- Shorten long URLs into unique short codes
- Retrieve original URLs by short code
- Track usage statistics (number of times a short URL has been accessed)
- Configurable expiration times for short URLs
- Easy deployment with Docker

## Getting Started
To get a local copy of this project up and running, follow these steps.

### Prerequisites
- [Go](https://golang.org/doc/install) 1.17 or higher
- [Docker](https://docs.docker.com/get-docker/) (optional, for containerized deployment)

## Installation
Clone the repository:
```bash
git clone https://github.com/yourusername/url-shortener.git
cd url-shortener
```
## Configuration
In docker-compose:
```bash
version: '3.8'

services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DATABASE_URL=postgres://postgres:password@db:5432/urlshortener?sslmode=disable
      - SERVER_PORT=8080
      - LOG_LEVEL=development
      - RATE_LIMIT=100
      - BASE_URL=http://localhost:8080
    depends_on:
      - db

  db:
    image: postgres:13
    platform: linux/amd64
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: password
      POSTGRES_DB: urlshortener
    volumes:
      # first my local dir, second default for everyone
      #- D:/Base/Downloads/example_db:/var/lib/postgresql/data
      - db-data:/var/lib/postgresql/data
      - ./db:/docker-entrypoint-initdb.d
    ports:
      - "5432:5432"

volumes:
  db-data:
```
You need to compile the database and api
##API Endpoints
-POST /shorten - Shorten a new URL.
-- Body: JSON with the original URL, e.g., { "url": "https://example.com" }
-- Response: JSON with the short code.
-GET /{shortCode} - Redirects to the original URL associated with {shortCode}.
-GET /{shortCode}/stats - Retrieves usage statistics for a specific short URL.

## Usage
Once the service is running, you can use it via HTTP requests.
![image](https://github.com/user-attachments/assets/abdebdab-60f2-47c8-b1a8-158c707e57ea)

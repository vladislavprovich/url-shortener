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
Clone the repository:
```bash
DB_DRIVER=postgres
DATABASE_URL=postgres://user:password@localhost:5432/shortener_db
SERVER_PORT=8080
RATE_LIMIT=10

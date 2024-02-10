# Cadana Exchange Challenge

## Overview

The Cadana Exchange Challenge involves building a client that retrieves exchange rate data from multiple services, selecting results from the quickest response, and discarding slow requests. This document provides a brief description of the challenge's goal and outlines the architecture used.

## Architecture

![](https://res.cloudinary.com/tutcan/image/upload/v1707426613/general/Cadana_Architecture.png)

### Goal

Build a client that efficiently retrieves exchange rate data by prioritizing responses from the fastest services, discarding slower ones.

### Tools used

- Go (programming language)
- Docker (containerization platform)
- Docker Compose (orchestration tool for Docker)
- Redis (in-memory data structure store)

## Setup

### Prerequisites

Ensure that Docker is installed on your machine before proceeding.

### Instructions

1. Clone the repository to your local machine

   ```
   git clone https://github.com/josephakayesi/cadana
   ```

2. Navigate to the project root.

   ```
   cd cadana
   ```

3. Start the exchange servers and client servers using Docker Compose.

   ```
   docker compose up
   ```

4. Access the people module and start the application that consumes data from the client server and performs manipulations on the Person objects.

   ```
   cd people
   go run cmd/main.go
   ```

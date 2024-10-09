> [!WARNING]  
> WORK IN VERY MUCH PROGRESS


# Micropub to Git-Powered Blogs

## Overview

This project provides a **Micropub-compliant API endpoint** that allows users to post content to Git-powered static blogs like **Hugo** and **Jekyll**. It supports creating, updating, and deleting blog posts via Micropub clients, with content automatically committed and pushed to a Git repository. The service uses **PocketBase** for authentication and indexing, and is built entirely in **Go**.

## Features

- **Micropub Compliance**
  - Supports `create`, `update`, and `delete` actions.
  - Handles `x-www-form-urlencoded` and JSON requests.
  - Implements `q=config` query for client configuration.
  - Separate **Media Endpoint** for asynchronous media uploads.

- **Content Types**
  - Text posts
  - Images (via URL or direct upload)
  - Links
  - Status updates
  - Replies (`in-reply-to` support)

- **Authentication**
  - OAuth2 authentication using PocketBase.
  - Access tokens validated per Micropub spec.

- **Git Integration**
  - Auto-commit and push to Git repositories (initially GitHub).
  - Option to create new branches or pull requests.
  - Frontmatter auto-generation for Hugo/Jekyll.

- **Draft Support**
  - Mark posts as drafts based on Micropub client input.
  - Draft status reflected in frontmatter (`draft: true`).

- **Web Interface**
  - Configure Git repository settings.
  - Manage OAuth2 settings.
  - View and manage blog entries indexed in PocketBase.

- **Indexing**
  - Background crawling of the Git repository.
  - PocketBase stores metadata for efficient retrieval.

- **Error Handling**
  - Appropriate HTTP status codes (`400`, `401`, `403`).
  - Detailed JSON error responses.

## Getting Started

### Prerequisites

- **Go** (latest version recommended)
- **PocketBase** (embedded, no additional setup required)
- **Git** (installed and configured)
- **Micropub Client** (e.g., Micropublish)

### Installation

1. **Clone the Repository**

   ```bash
   git clone https://github.com/harperreed/micropub-to-git.git
   cd micropub-to-git
   ```

2. **Install Dependencies**

   Ensure you have Go modules enabled:

   ```bash
   go mod download
   ```

3. **Build the Application**

   ```bash
   go build -o micropub-server main.go
   ```

4. **Run the Application**

   ```bash
   ./micropub-server
   ```

## Configuration

### PocketBase Setup

PocketBase is used for authentication and indexing.

- **OAuth2 Configuration**
  - Users authenticate via OAuth2 to obtain access tokens.
  - Configure OAuth2 clients and scopes as needed.

- **Database**
  - PocketBase uses SQLite; no additional setup required.
  - Data is stored in `pb_data` directory by default.

### Git Repository Settings

Configure your Git repository settings via the web interface or a configuration file.

- **Repository URL**
- **Branch (e.g., `main`, `master`)**
- **Commit Options**
  - Direct commit and push
  - Create a new branch
  - Submit as a pull request

### Web Interface

Access the web interface at `http://localhost:PORT` to manage settings.

## Usage

### Micropub Client Configuration

1. **Set Endpoint URL**

   Point your Micropub client to the server endpoint:

   ```
   http://localhost:PORT/micropub
   ```

2. **Authenticate**

   - Use OAuth2 to authenticate and obtain an access token.
   - The client should handle token storage per Micropub spec.

3. **Create, Update, Delete Posts**

   - Use your Micropub client to create new posts, update existing ones, or delete posts.
   - Supports text, images, links, status updates, and replies.

### Media Uploads

Use the separate Media Endpoint for uploading files.

- **Media Endpoint URL**

  ```
  http://localhost:PORT/media
  ```

- **Upload Process**
  - Upload media files (images, etc.) to the media endpoint.
  - Receive a URL for the uploaded file.
  - Reference the media URL in your Micropub post.

## Testing

This project uses **Test-Driven Development (TDD)** and includes both unit tests and integration tests.

- **Running Tests**

  To run all tests, including unit tests and integration tests:

  ```bash
  go test ./...
  ```

- **Test Coverage**

  To run tests with coverage:

  ```bash
  go test -coverprofile=coverage.out ./...
  go tool cover -html=coverage.out
  ```

- **Integration Tests**

  Integration tests for server initialization and middleware functionality are located in `cmd/server/main_test.go`.

- **Test Configuration**

  The `internal/config/config.go` file includes support for test configurations. Use `LoadTestConfig()` for test setups.

Ensure all tests pass before deploying or updating the application.

## Contributing

Contributions are welcome! Please follow these guidelines:

1. **Fork the Repository**

   ```bash
   git fork https://github.com/harperreed/micropub-to-git.git
   ```

2. **Create a Feature Branch**

   ```bash
   git checkout -b feature/your-feature-name
   ```

3. **Commit Your Changes**

   ```bash
   git commit -m "Description of your changes"
   ```

4. **Push to Your Fork**

   ```bash
   git push origin feature/your-feature-name
   ```

5. **Submit a Pull Request**

   - Describe the changes and why they're needed.
   - Ensure all tests pass and the code adheres to the project's coding standards.

## License

This project is licensed under the **MIT License**.

## Acknowledgments

- **Micropub Specification**: [W3C Micropub Recommendation](https://www.w3.org/TR/micropub/)
- **PocketBase**: [PocketBase.io](https://pocketbase.io/)
- **Go Language**: [golang.org](https://golang.org/)

## Contact

For questions or support, please open an issue on the [GitHub repository](https://github.com/harperreed/micropub-to-git/issues).

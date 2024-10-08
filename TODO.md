# TODO List: Micropub to Git-Powered Blogs Project

## Features

### 1. Micropub Endpoint
- [x] Write tests and stubs for `create` post action
  - [x] Support `x-www-form-urlencoded` requests
  - [x] Support JSON requests
- [x] Implement basic functionality for creating posts
- [x] Implement Git integration for creating posts
- [x] Write tests and stubs for `update` post action
- [x] Write tests and stubs for `delete` post action
- [x] Implement full functionality for updating and deleting posts

### 2. OAuth2 Authentication (PocketBase)
- [x] Write tests and stubs for OAuth2 authentication flow
  - [x] Test user registration process
  - [x] Test token generation and validation
  - [x] Test token refresh mechanism
  - [x] Test token revocation
- [x] Implement real functionality after passing tests
  - [x] Set up PocketBase OAuth2 server
  - [x] Implement user registration and login
  - [x] Implement token generation and validation
  - [x] Implement token refresh mechanism
  - [x] Implement token revocation
- [x] Integrate OAuth2 authentication with Micropub endpoint

### 3. Git Integration
- [x] Write functionality for committing posts to Git
- [x] Write functionality for pushing posts to Git repository
- [x] Write functionality for creating new branches
- [x] Write tests and stubs for creating pull requests
- [x] Implement real functionality for creating pull requests after passing tests

### 4. Frontmatter Generation
- [x] Implement basic frontmatter generation
- [x] Write tests and stubs for advanced frontmatter generation
  - [x] Support for metadata (tags, categories)
  - [x] Mark posts as drafts
- [x] Implement real functionality after passing tests

### 5. Media Endpoint
- [x] Write tests and stubs for media endpoint
- [x] Implement real functionality after passing tests

### 6. Post Metadata and PocketBase Indexing
- [x] Write tests and stubs for storing blog entry metadata in PocketBase
- [ ] Write tests and stubs for indexing Git repository blog entries
- [ ] Implement real functionality after passing tests

### 7. Post Draft Support
- [ ] Write tests and stubs for draft support based on client input
- [ ] Implement real functionality after passing tests

### 8. Web Interface for Settings
- [ ] Write tests and stubs for web interface features:
  - [ ] Git repository configuration (URL, branch, pull request options)
  - [ ] OAuth2 settings management
- [ ] Implement real functionality after passing tests

### 9. Config Query Support (`q=config`)
- [ ] Write tests and stubs for `q=config` query response
  - [ ] Include media endpoint URL and supported content types in the response
- [ ] Implement real functionality after passing tests

### 10. Error Handling
- [ ] Write tests and stubs for error handling with appropriate HTTP status codes:
  - [ ] `400` for invalid requests
  - [ ] `401` for authentication errors
  - [ ] `403` for forbidden actions
  - [ ] Descriptive `error` field in JSON response
- [ ] Implement real functionality after passing tests

### 11. Background Crawling for Blog Index
- [ ] Write tests and stubs for background crawling of the Git repository
- [ ] Write tests and stubs for updating PocketBase index with blog metadata
- [ ] Implement real functionality after passing tests

### 12. Post Syndication (Future Enhancement)
- [ ] Consider adding support for syndication targets (`q=syndicate-to`) in future versions

---

## Testing and Quality Assurance

### Test-Driven Development (TDD)
- [ ] Use TDD for all features
  - [ ] Write tests and stubs for each feature
  - [ ] Mock interfaces before implementing real functionality
  - [ ] Ensure all tests pass before moving to real implementation

### Automated Testing
- [ ] Set up automated testing for continuous integration
  - [ ] Configure GitHub Actions or similar CI/CD tool
  - [ ] Set up test runners for Go tests
  - [ ] Implement code coverage reporting

### Integration Testing
- [ ] Develop integration tests for key system components
  - [ ] Test Micropub endpoint with various clients
  - [ ] Test Git integration with different repository hosts
  - [ ] Test OAuth2 flow with PocketBase

### Performance Testing
- [ ] Implement performance benchmarks
  - [ ] Test response times for Micropub endpoint
  - [ ] Measure Git operations performance
  - [ ] Evaluate PocketBase query performance

### Security Testing
- [ ] Conduct security audits
  - [ ] Perform OAuth2 implementation security review
  - [ ] Test for common web vulnerabilities (XSS, CSRF, etc.)
  - [ ] Review secure handling of tokens and sensitive data

---

## Miscellaneous
- [x] Set up project in Go
- [x] Use PocketBase as the backend for managing the blog index and OAuth2
- [x] Ensure the system is scalable for multi-blog support in future iterations


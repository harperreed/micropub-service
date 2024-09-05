# TODO List: Micropub to Git-Powered Blogs Project

## Features

### 1. Micropub Endpoint
- [ ] Write tests and stubs for `create` post action
  - [ ] Support `x-www-form-urlencoded` requests
  - [ ] Support JSON requests
- [ ] Write tests and stubs for `update` post action
- [ ] Write tests and stubs for `delete` post action
- [ ] Implement real functionality after passing tests

### 2. OAuth2 Authentication (PocketBase)
- [ ] Write tests and stubs for OAuth2 authentication flow
- [ ] Implement real functionality after passing tests

### 3. Git Integration
- [ ] Write tests and stubs for committing posts to Git
- [ ] Write tests and stubs for pushing posts to Git repository
- [ ] Write tests and stubs for creating new branches or pull requests
- [ ] Implement real functionality after passing tests

### 4. Frontmatter Generation
- [ ] Write tests and stubs for frontmatter generation
  - [ ] Support for metadata (tags, categories, dates)
  - [ ] Mark posts as drafts
- [ ] Implement real functionality after passing tests

### 5. Media Endpoint
- [ ] Write tests and stubs for media endpoint
- [ ] Implement real functionality after passing tests

### 6. Post Metadata and PocketBase Indexing
- [ ] Write tests and stubs for storing blog entry metadata in PocketBase
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

## Testing and TDD

- [ ] Use Test-Driven Development (TDD) for all features.
  - [ ] Write tests and stubs for each feature.
  - [ ] Mock interfaces before implementing real functionality.
  - [ ] Ensure all tests pass before moving to real implementation.
- [ ] Set up automated testing for continuous integration.

---

## Miscellaneous
- [ ] Set up project in Go
- [ ] Use PocketBase as the backend for managing the blog index and OAuth2
- [ ] Ensure the system is scalable for multi-blog support in future iterations


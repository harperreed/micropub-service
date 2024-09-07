### **Product Requirements Document: Micropub to Git-powered Blogs using PocketBase**

---

#### **Product Overview:**
The goal is to create a fully Micropub-compliant API that posts content to Git-powered static blogs (e.g., Hugo, Jekyll). The service will use PocketBase as the backend to manage an index of blog entries and handle OAuth authentication, with the content committed directly to a Git repository.

#### **Primary User Story:**
As a user, I want to post content (text, images, links, replies) to my Git-powered blog via a Micropub client, and have the content automatically committed and pushed to the Git repository with generated frontmatter.

---

### **Features**

1. **Micropub Compliance**
   The service will follow the [W3C Micropub Recommendation (23 May 2017)](https://www.w3.org/TR/micropub/), ensuring:
   - **Create**: Users can create blog posts via form-encoded or JSON requests.
   - **Update**: Users can update existing posts with new data.
   - **Delete**: Users can delete or undelete blog posts.
   - **q=config**: Clients can query the endpoint for configuration details like syndication targets and media endpoints.

2. **Content Types Supported:**
   - Text posts
   - Images (via URLs or direct uploads)
   - Links
   - Status updates
   - Replies (with `in-reply-to` property)

3. **PocketBase Integration:**
   - **Data Storage**: PocketBase will handle the index of all blog entries, storing metadata for posts, including titles, tags, categories, dates, and any other relevant fields.
   - **OAuth Integration**: PocketBase will manage user authentication using OAuth2. Users will authenticate via OAuth2 and obtain an access token to use with the Micropub API.

4. **Git Integration**
   The Micropub server will:
   - Automatically commit and push content to a configured Git repository (initially supporting GitHub, with broader support for any Git-based platform).
   - Provide the option to create new branches or submit pull requests instead of pushing to the main branch.
   - Generate appropriate frontmatter based on the content posted (e.g., for Hugo/Jekyll).

5. **Media Uploads via Media Endpoint**
   A separate **Media Endpoint** will handle media file uploads asynchronously. Users will upload images or other media files, and the system will return a URL that can be used in the main Micropub post request.

6. **Post Draft Support**
   If the client marks content as a draft, it will be reflected in the frontmatter of the Markdown file (e.g., `draft: true`). Drafts will not be published until ready.

7. **Post Metadata**
   - The service will automatically generate frontmatter metadata (tags, categories, published date) and store it in the index.
   - PocketBase will store this metadata, and Micropub clients will be able to send custom fields as needed.
   - The service will ignore unrecognized properties to allow for extensibility.

8. **Configuration Queries (`q=config`)**
   The Micropub API will support the `q=config` query, allowing clients to discover:
   - The available Media Endpoint URL
   - Supported content types
   - Syndication targets (if future support is added)

9. **Error Handling**
   - The service will return appropriate HTTP status codes for errors:
     - `400` for bad requests (e.g., invalid content types)
     - `401` for authentication errors (missing or invalid OAuth tokens)
     - `403` for forbidden actions
   - Error responses will include a descriptive `error` field in the JSON body, as required by the spec.

10. **Crawling & Indexing Blog Entries**
    - Upon hitting `q=config`, the service will trigger a background crawl of the Git repository to index all blog entries in PocketBase.
    - The index will store metadata for each post (file paths, tags, categories, etc.) and will allow quick retrieval for post updates or deletions.

11. **Web Interface for Settings Management**
    - A simple web interface will be built to allow users to:
      - Configure their Git repository (URL, branch, pull request options).
      - Manage OAuth settings.
      - View and manage the index of blog entries stored in PocketBase.

12. **Scalability**
   - The initial implementation will support a single blog per user, but the system will be designed to allow for multi-blog support in the future.

---

### **Technical Stack**

- **Backend**: Golang (to handle the Micropub endpoint, PocketBase, and Git interactions)
- **Database**: PocketBase (SQLite-based, lightweight, with built-in OAuth support)
- **Git Integration**: GitHub (initially) with plans for expanding to support any Git-based platform (via Go's `git` libraries or APIs)
- **Frontend**: Simple web UI built in Go or JavaScript for managing settings.

---

### **Micropub Specification Support**

- **Authentication (OAuth2)**:
   PocketBase will handle OAuth2, and users will authenticate via OAuth2 and receive access tokens, which will be passed in the HTTP Authorization header or form body as per the Micropub spec.

- **Form-encoded and JSON Support**:
   The service will support both `x-www-form-urlencoded` and JSON-based content creation requests, fully compliant with the Micropub specification.

- **Post Actions**:
   The service will handle `create`, `update`, and `delete` actions:
   - **Create**: Users can create new blog posts.
   - **Update**: Posts can be updated, with properties like `content`, `category`, or `tags` being modified.
   - **Delete**: Posts can be deleted or undeleted using `action=delete` or `action=undelete`.

- **Media Endpoint**:
   The service will implement a separate media endpoint for handling file uploads. Clients will upload files asynchronously, and the endpoint will return a URL for the uploaded file, which can then be referenced in the main post.

- **Error Handling**:
   All errors will return appropriate HTTP status codes and a JSON error body, including detailed descriptions of the error, as per the Micropub spec.

- **Syndication Targets** (future support):
   Support for `q=syndicate-to` will be considered for future iterations, allowing users to push content to other platforms.

---

### **Development Milestones**

1. **Micropub Endpoint**:
   - Basic endpoint that supports content creation and form-encoded requests.
   - Git integration for auto-committing and pushing to GitHub.

2. **OAuth Integration via PocketBase**:
   - Set up PocketBase for OAuth2 authentication and user management.

3. **Media Endpoint**:
   - Build a separate media endpoint for handling image uploads asynchronously.

4. **Post Actions**:
   - Add support for updating and deleting posts.

5. **Web Interface**:
   - Create a basic settings management UI for configuring Git and OAuth settings.

6. **PocketBase Index**:
   - Set up PocketBase to index blog entries and handle queries efficiently.

7. **Crawling & Indexing**:
   - Implement background crawling of the Git repository on `q=config` to keep the PocketBase index up to date.

---

### **Acceptance Criteria**

- Full compliance with the Micropub specification.
- Ability to create, update, and delete posts.
- OAuth2 authentication integrated with PocketBase.
- Git commits and pushes for all new content.
- Working media endpoint for file uploads.
- Simple web interface for settings management.
- Background crawling and indexing of blog entries.

---

This PRD outlines the complete scope of the project, ensuring that all features are fully compliant with the Micropub specification and providing room for future scalability and improvements (like multi-blog support and syndication targets).

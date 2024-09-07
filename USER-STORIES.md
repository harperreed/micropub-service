### **User Stories for Micropub to Git-Powered Blogs**

1. **Micropub Post Creation**
   **As a user**, I want to create blog posts (text, images, links) via a Micropub client, so that I can publish content directly to my Git-powered blog (e.g., Hugo or Jekyll).

2. **OAuth Authentication**
   **As a user**, I want to authenticate with OAuth2, so that my Micropub client can securely post content to my blog.

3. **Git Integration (Commit & Push)**
   **As a user**, I want the service to automatically commit and push my blog posts to my Git repository, so that my blog is updated without manual intervention.

4. **Frontmatter Auto-generation**
   **As a user**, I want the service to automatically generate frontmatter for my blog posts, so that they follow the correct format for my static site generator (Hugo/Jekyll).

5. **Post Drafts**
   **As a user**, I want to mark posts as drafts via my Micropub client, so that they are not published until I’m ready.

6. **Media Upload via Media Endpoint**
   **As a user**, I want to upload media files (e.g., images) via a separate media endpoint, so that I can easily embed images in my blog posts.

7. **Post Updates**
   **As a user**, I want to update existing blog posts, so that I can make changes to content like text, categories, or images without creating a new post.

8. **Post Deletions**
   **As a user**, I want to delete or undelete posts, so that I can remove content from my blog when needed.

9. **Config Query**
   **As a user**, I want my Micropub client to query the server for configuration (`q=config`), so that I can discover the available media endpoint and other server capabilities.

10. **Error Feedback**
   **As a user**, I want to receive clear error messages when something goes wrong (e.g., invalid content, authentication issues), so that I can resolve the problem quickly.

11. **Web Interface for Settings**
   **As a user**, I want a web interface to configure my Git repository settings (URL, branch, pull request options), so that I can easily manage where my content is pushed.

12. **Background Crawling for Blog Index**
   **As a user**, I want the service to index my existing blog posts, so that I can update or delete content already in the Git repository.

13. **PocketBase Integration for Post Index**
   **As a developer**, I want to use PocketBase to manage the index of blog entries, so that I can store metadata (like post titles, tags, categories) and retrieve them efficiently for updates or deletions.

14. **OAuth2 Token Management via PocketBase**
   **As a user**, I want PocketBase to manage my OAuth2 authentication, so that I can securely log in and authenticate my Micropub client.

15. **Commit as Branch or Pull Request**
   **As a user**, I want the option to commit my blog posts to a new branch or create a pull request, so that I can review changes before they are merged into my main blog.

16. **Real-time Post Indexing on Config Query**
   **As a user**, I want the blog post index to be updated when I query the configuration (`q=config`), so that I know my index reflects the latest posts from the Git repository.

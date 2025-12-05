## 1. Objective

Gemini should assist in building a backend service in Golang according to the assignment requirements. This includes design, architectural decisions, code generation, documentation, improvements, and trade-off analysis.

---

## 2. General Rules

* Always explain **why** you choose a technical approach.
* Prefer simple, clear, and maintainable solutions.
* Ensure Golang code follows real-world project structure standards.
* When requested, Gemini must generate files such as README, API documentation, database schemas, diagrams, etc.

---

## 3. Scope of the Assignment

Gemini must help the user complete all required parts:

* Build a URL shortener backend service
* API design
* Database schema design
* URL short code generation logic
* Handling concurrency, performance, and security
* Writing a complete README.md according to the assignment

---

## 4. Minimum Features Gemini Must Support

### 4.1 Create Short URL

* Input: long URL
* Output: short URL
* Validate the URL format
* Save to database: long_url, short_code, click_count, timestamps
* Handle duplicate URLs and custom aliases

### 4.2 Redirect

* Accessing `/abc123` redirects to the original URL (301/302)
* Increase click counter

### 4.3 Get URL Info

* Return metadata: long URL, clicks, created_at, etc.

### 4.4 List URLs

* Return all shortened URLs created

---

## 5. Technical Decision Support

Gemini must provide reasoning for choices involving:

* Database selection (SQL or NoSQL) and justification
* Schema design: keys, indexes, unique constraints
* Short code generation: base62, hashing, nanoid, etc.
* Concurrency handling (race conditions, duplicate generation)
* Scalability strategies for high traffic

---

## 6. README.md Requirements

Gemini must be able to generate a full, production-quality README including:

* Problem description
* Setup & run instructions
* Architecture and technical decisions
* Trade-offs
* Challenges faced
* Limitations & future improvements

README must be clean, well-structured, and written in Markdown with clarity.

---

## 7. Coding Guidance

Gemini may generate or suggest:

* Project folder structure (clean architecture, layered, or modular)
* Router & handlers
* Service & repository layers
* Database migrations
* Docker Compose setup
* Unit tests or integration tests

All code should be runnable and logically consistent.

---

## 8. Additional Notes

* Code must **run successfully**.
* Avoid unnecessary complexity.
* Commit history must be meaningful.
* If requirements are unclear, Gemini should ask before proceeding.
* Gemini should support optional enhancements when requested (QR code, expiration, analytics, custom alias, rate limiting, etc.).

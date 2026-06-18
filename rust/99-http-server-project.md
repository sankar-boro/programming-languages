# Chapter 99 — Final Project: HTTP Server From Scratch

> *Build a real HTTP/1.1 server using only the Rust standard library — no frameworks, no external crates. Just TCP sockets, manual parsing, and everything you've learned in this book.*

---

## Project Goal

Build a functioning HTTP/1.1 server that:
- Accepts TCP connections
- Parses raw HTTP requests
- Routes GET requests to handlers
- Sends valid HTTP responses with proper headers
- Handles concurrent connections with a thread pool
- Serves static files and dynamic responses

**Constraint**: `std` only. No `tokio`, no `hyper`, no `actix`. Just `std::net`, `std::io`, and `std::thread`.

---

## Project Structure

```
http-server/
├── Cargo.toml
└── src/
    ├── main.rs          ← server startup and configuration
    ├── server.rs        ← TCP listener + thread pool
    ├── request.rs       ← HTTP request parsing
    ├── response.rs      ← HTTP response building
    ├── router.rs        ← route matching and dispatch
    └── threadpool.rs    ← worker thread pool
```

```toml
# Cargo.toml
[package]
name = "http-server"
version = "0.1.0"
edition = "2021"

# No external dependencies — std only
```

---

## Part 1: The Thread Pool

```rust
// src/threadpool.rs
use std::sync::{Arc, Mutex};
use std::thread;

type Job = Box<dyn FnOnce() + Send + 'static>;

enum Message {
    NewJob(Job),
    Shutdown,
}

struct Worker {
    id: usize,
    thread: Option<thread::JoinHandle<()>>,
}

impl Worker {
    fn new(id: usize, receiver: Arc<Mutex<std::sync::mpsc::Receiver<Message>>>) -> Worker {
        let thread = thread::spawn(move || {
            loop {
                let message = receiver.lock().unwrap().recv().unwrap();
                match message {
                    Message::NewJob(job) => {
                        job();
                    }
                    Message::Shutdown => {
                        break;
                    }
                }
            }
        });

        Worker { id, thread: Some(thread) }
    }
}

pub struct ThreadPool {
    workers: Vec<Worker>,
    sender: std::sync::mpsc::Sender<Message>,
}

impl ThreadPool {
    pub fn new(size: usize) -> ThreadPool {
        assert!(size > 0, "Thread pool size must be greater than zero");

        let (sender, receiver) = std::sync::mpsc::channel();
        let receiver = Arc::new(Mutex::new(receiver));

        let mut workers = Vec::with_capacity(size);
        for id in 0..size {
            workers.push(Worker::new(id, Arc::clone(&receiver)));
        }

        ThreadPool { workers, sender }
    }

    pub fn execute<F>(&self, f: F)
    where
        F: FnOnce() + Send + 'static,
    {
        let job = Box::new(f);
        self.sender.send(Message::NewJob(job)).unwrap();
    }
}

impl Drop for ThreadPool {
    fn drop(&mut self) {
        // Send shutdown signal to all workers
        for _ in &self.workers {
            self.sender.send(Message::Shutdown).unwrap();
        }

        // Wait for all workers to finish
        for worker in &mut self.workers {
            if let Some(thread) = worker.thread.take() {
                thread.join().unwrap();
            }
        }
    }
}
```

---

## Part 2: HTTP Request Parsing

```rust
// src/request.rs
use std::collections::HashMap;
use std::io::{BufRead, BufReader, Read};
use std::net::TcpStream;

#[derive(Debug, Clone, PartialEq)]
pub enum Method {
    Get,
    Post,
    Put,
    Delete,
    Head,
    Options,
    Unknown(String),
}

impl Method {
    fn parse(s: &str) -> Method {
        match s {
            "GET" => Method::Get,
            "POST" => Method::Post,
            "PUT" => Method::Put,
            "DELETE" => Method::Delete,
            "HEAD" => Method::Head,
            "OPTIONS" => Method::Options,
            other => Method::Unknown(other.to_string()),
        }
    }
}

#[derive(Debug)]
pub struct Request {
    pub method: Method,
    pub path: String,
    pub query: HashMap<String, String>,
    pub version: String,
    pub headers: HashMap<String, String>,
    pub body: Vec<u8>,
}

#[derive(Debug)]
pub enum ParseError {
    EmptyRequest,
    InvalidRequestLine(String),
    InvalidHeader(String),
    IoError(std::io::Error),
}

impl std::fmt::Display for ParseError {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match self {
            ParseError::EmptyRequest => write!(f, "Empty request"),
            ParseError::InvalidRequestLine(s) => write!(f, "Invalid request line: {}", s),
            ParseError::InvalidHeader(s) => write!(f, "Invalid header: {}", s),
            ParseError::IoError(e) => write!(f, "IO error: {}", e),
        }
    }
}

impl From<std::io::Error> for ParseError {
    fn from(e: std::io::Error) -> Self {
        ParseError::IoError(e)
    }
}

impl Request {
    pub fn parse(stream: &TcpStream) -> Result<Request, ParseError> {
        let mut reader = BufReader::new(stream);
        let mut headers = HashMap::new();
        let mut lines = Vec::new();

        // Read headers (until blank line)
        loop {
            let mut line = String::new();
            let n = reader.read_line(&mut line)?;
            if n == 0 {
                return Err(ParseError::EmptyRequest);
            }
            let line = line.trim_end_matches('\n').trim_end_matches('\r').to_string();
            if line.is_empty() {
                break;
            }
            lines.push(line);
        }

        if lines.is_empty() {
            return Err(ParseError::EmptyRequest);
        }

        // Parse request line: "GET /path HTTP/1.1"
        let request_line = &lines[0];
        let parts: Vec<&str> = request_line.splitn(3, ' ').collect();
        if parts.len() != 3 {
            return Err(ParseError::InvalidRequestLine(request_line.clone()));
        }

        let method = Method::parse(parts[0]);
        let raw_path = parts[1].to_string();
        let version = parts[2].to_string();

        // Parse path and query string
        let (path, query) = parse_path_and_query(&raw_path);

        // Parse headers: "Key: Value"
        for line in &lines[1..] {
            let colon_pos = line.find(':').ok_or_else(|| {
                ParseError::InvalidHeader(line.clone())
            })?;
            let key = line[..colon_pos].trim().to_lowercase();
            let value = line[colon_pos + 1..].trim().to_string();
            headers.insert(key, value);
        }

        // Read body if Content-Length is set
        let body = if let Some(len_str) = headers.get("content-length") {
            let len: usize = len_str.parse().unwrap_or(0);
            let mut body = vec![0u8; len];
            reader.read_exact(&mut body)?;
            body
        } else {
            Vec::new()
        };

        Ok(Request { method, path, query, version, headers, body })
    }
}

fn parse_path_and_query(raw: &str) -> (String, HashMap<String, String>) {
    let mut query_map = HashMap::new();

    if let Some(q_pos) = raw.find('?') {
        let path = raw[..q_pos].to_string();
        let query_str = &raw[q_pos + 1..];

        for pair in query_str.split('&') {
            if let Some(eq_pos) = pair.find('=') {
                let key = url_decode(&pair[..eq_pos]);
                let value = url_decode(&pair[eq_pos + 1..]);
                query_map.insert(key, value);
            }
        }

        (path, query_map)
    } else {
        (raw.to_string(), query_map)
    }
}

fn url_decode(s: &str) -> String {
    let mut result = String::new();
    let mut chars = s.chars().peekable();

    while let Some(c) = chars.next() {
        if c == '%' {
            let hex: String = chars.by_ref().take(2).collect();
            if let Ok(byte) = u8::from_str_radix(&hex, 16) {
                result.push(byte as char);
            } else {
                result.push('%');
                result.push_str(&hex);
            }
        } else if c == '+' {
            result.push(' ');
        } else {
            result.push(c);
        }
    }

    result
}
```

---

## Part 3: HTTP Response Builder

```rust
// src/response.rs
use std::collections::HashMap;
use std::io::Write;
use std::net::TcpStream;

#[derive(Debug)]
pub struct Response {
    pub status_code: u16,
    pub status_text: &'static str,
    pub headers: HashMap<String, String>,
    pub body: Vec<u8>,
}

impl Response {
    pub fn new(status_code: u16) -> Self {
        Response {
            status_code,
            status_text: status_text(status_code),
            headers: HashMap::new(),
            body: Vec::new(),
        }
    }

    pub fn ok() -> Self { Response::new(200) }
    pub fn not_found() -> Self { Response::new(404) }
    pub fn method_not_allowed() -> Self { Response::new(405) }
    pub fn internal_error() -> Self { Response::new(500) }
    pub fn bad_request() -> Self { Response::new(400) }

    pub fn header(mut self, key: &str, value: &str) -> Self {
        self.headers.insert(key.to_string(), value.to_string());
        self
    }

    pub fn body_str(mut self, content: &str) -> Self {
        self.body = content.as_bytes().to_vec();
        self.headers.insert("Content-Length".to_string(), self.body.len().to_string());
        self
    }

    pub fn body_bytes(mut self, content: Vec<u8>) -> Self {
        let len = content.len();
        self.body = content;
        self.headers.insert("Content-Length".to_string(), len.to_string());
        self
    }

    pub fn content_type(self, ct: &str) -> Self {
        self.header("Content-Type", ct)
    }

    pub fn html(self, content: &str) -> Self {
        self.content_type("text/html; charset=utf-8").body_str(content)
    }

    pub fn json(self, content: &str) -> Self {
        self.content_type("application/json").body_str(content)
    }

    pub fn text(self, content: &str) -> Self {
        self.content_type("text/plain; charset=utf-8").body_str(content)
    }

    pub fn send(&self, stream: &mut TcpStream) -> std::io::Result<()> {
        // Status line
        let status_line = format!("HTTP/1.1 {} {}\r\n", self.status_code, self.status_text);
        stream.write_all(status_line.as_bytes())?;

        // Headers
        for (key, value) in &self.headers {
            let header_line = format!("{}: {}\r\n", key, value);
            stream.write_all(header_line.as_bytes())?;
        }

        // Default headers
        stream.write_all(b"Connection: close\r\n")?;
        stream.write_all(b"Server: rust-http/0.1\r\n")?;

        // Blank line separating headers from body
        stream.write_all(b"\r\n")?;

        // Body
        if !self.body.is_empty() {
            stream.write_all(&self.body)?;
        }

        stream.flush()?;
        Ok(())
    }
}

fn status_text(code: u16) -> &'static str {
    match code {
        200 => "OK",
        201 => "Created",
        204 => "No Content",
        301 => "Moved Permanently",
        302 => "Found",
        304 => "Not Modified",
        400 => "Bad Request",
        401 => "Unauthorized",
        403 => "Forbidden",
        404 => "Not Found",
        405 => "Method Not Allowed",
        500 => "Internal Server Error",
        503 => "Service Unavailable",
        _ => "Unknown",
    }
}
```

---

## Part 4: The Router

```rust
// src/router.rs
use crate::request::{Method, Request};
use crate::response::Response;
use std::collections::HashMap;

pub type Handler = Box<dyn Fn(&Request) -> Response + Send + Sync>;

pub struct Router {
    routes: Vec<Route>,
}

struct Route {
    method: Method,
    pattern: String,
    handler: Handler,
}

impl Router {
    pub fn new() -> Self {
        Router { routes: Vec::new() }
    }

    pub fn get<F>(&mut self, pattern: &str, handler: F) -> &mut Self
    where
        F: Fn(&Request) -> Response + Send + Sync + 'static,
    {
        self.routes.push(Route {
            method: Method::Get,
            pattern: pattern.to_string(),
            handler: Box::new(handler),
        });
        self
    }

    pub fn post<F>(&mut self, pattern: &str, handler: F) -> &mut Self
    where
        F: Fn(&Request) -> Response + Send + Sync + 'static,
    {
        self.routes.push(Route {
            method: Method::Post,
            pattern: pattern.to_string(),
            handler: Box::new(handler),
        });
        self
    }

    pub fn handle(&self, request: &Request) -> Response {
        for route in &self.routes {
            if route.method == request.method {
                if let Some(_params) = match_pattern(&route.pattern, &request.path) {
                    return (route.handler)(request);
                }
            }
        }

        // No route matched
        Response::not_found()
            .html("<h1>404 Not Found</h1><p>No route matches this path.</p>")
    }
}

fn match_pattern(pattern: &str, path: &str) -> Option<HashMap<String, String>> {
    let pattern_parts: Vec<&str> = pattern.split('/').collect();
    let path_parts: Vec<&str> = path.split('/').collect();

    if pattern_parts.len() != path_parts.len() {
        return None;
    }

    let mut params = HashMap::new();

    for (p, v) in pattern_parts.iter().zip(path_parts.iter()) {
        if p.starts_with(':') {
            // Dynamic segment: :id, :name, etc.
            params.insert(p[1..].to_string(), v.to_string());
        } else if p != v {
            return None;
        }
    }

    Some(params)
}
```

---

## Part 5: The Server

```rust
// src/server.rs
use crate::request::Request;
use crate::response::Response;
use crate::router::Router;
use crate::threadpool::ThreadPool;
use std::net::{TcpListener, TcpStream};
use std::sync::Arc;
use std::time::Duration;

pub struct Server {
    addr: String,
    router: Arc<Router>,
    workers: usize,
}

impl Server {
    pub fn new(addr: &str) -> Self {
        Server {
            addr: addr.to_string(),
            router: Arc::new(Router::new()),
            workers: 4,
        }
    }

    pub fn with_router(mut self, router: Router) -> Self {
        self.router = Arc::new(router);
        self
    }

    pub fn workers(mut self, n: usize) -> Self {
        self.workers = n;
        self
    }

    pub fn run(self) {
        let listener = TcpListener::bind(&self.addr).expect("Failed to bind address");
        let pool = ThreadPool::new(self.workers);

        println!("Server listening on http://{}", self.addr);

        for stream in listener.incoming() {
            match stream {
                Ok(stream) => {
                    let router = Arc::clone(&self.router);
                    pool.execute(move || {
                        handle_connection(stream, &router);
                    });
                }
                Err(e) => {
                    eprintln!("Connection error: {}", e);
                }
            }
        }
    }
}

fn handle_connection(mut stream: TcpStream, router: &Router) {
    // Set read timeout to prevent hanging connections
    stream.set_read_timeout(Some(Duration::from_secs(30))).ok();

    let response = match Request::parse(&stream) {
        Ok(request) => {
            println!("{} {} {}", request.version, 
                     format!("{:?}", request.method), request.path);
            router.handle(&request)
        }
        Err(e) => {
            eprintln!("Request parse error: {}", e);
            Response::bad_request().text(&format!("Bad Request: {}", e))
        }
    };

    if let Err(e) = response.send(&mut stream) {
        eprintln!("Failed to send response: {}", e);
    }
}
```

---

## Part 6: Static File Server

```rust
// In main.rs or a separate file_server.rs

use crate::request::Request;
use crate::response::Response;
use std::path::{Path, PathBuf};

pub fn serve_static(root: &str, req: &Request) -> Response {
    let root = Path::new(root);
    let req_path = req.path.trim_start_matches('/');

    // Security: prevent path traversal
    let full_path = root.join(req_path);
    if !is_safe_path(&root, &full_path) {
        return Response::new(403).text("Forbidden");
    }

    // Default to index.html for directory requests
    let file_path = if full_path.is_dir() {
        full_path.join("index.html")
    } else {
        full_path
    };

    match std::fs::read(&file_path) {
        Ok(content) => {
            let content_type = mime_type(file_path.extension()
                .and_then(|e| e.to_str())
                .unwrap_or(""));
            Response::ok()
                .content_type(content_type)
                .body_bytes(content)
        }
        Err(_) => Response::not_found()
            .html("<h1>404 Not Found</h1>"),
    }
}

fn is_safe_path(root: &Path, path: &Path) -> bool {
    // Ensure resolved path is within root — prevents ../../../etc/passwd
    let root = root.canonicalize().unwrap_or_else(|_| root.to_path_buf());
    let path = match path.canonicalize() {
        Ok(p) => p,
        Err(_) => return false,
    };
    path.starts_with(root)
}

fn mime_type(extension: &str) -> &'static str {
    match extension {
        "html" | "htm" => "text/html; charset=utf-8",
        "css"          => "text/css",
        "js"           => "application/javascript",
        "json"         => "application/json",
        "png"          => "image/png",
        "jpg" | "jpeg" => "image/jpeg",
        "gif"          => "image/gif",
        "svg"          => "image/svg+xml",
        "ico"          => "image/x-icon",
        "txt"          => "text/plain; charset=utf-8",
        "pdf"          => "application/pdf",
        _              => "application/octet-stream",
    }
}
```

---

## Part 7: Main — Putting It All Together

```rust
// src/main.rs
mod threadpool;
mod request;
mod response;
mod router;
mod server;

use response::Response;
use router::Router;
use server::Server;

fn main() {
    let mut router = Router::new();

    // Home page
    router.get("/", |_req| {
        Response::ok().html(r#"
            <!DOCTYPE html>
            <html>
            <head><title>Rust HTTP Server</title></head>
            <body>
                <h1>Hello from Rust!</h1>
                <p>Built with std only — no frameworks.</p>
                <ul>
                    <li><a href="/about">About</a></li>
                    <li><a href="/api/time">Current Time API</a></li>
                    <li><a href="/api/echo?message=hello">Echo API</a></li>
                </ul>
            </body>
            </html>
        "#)
    });

    // About page
    router.get("/about", |_req| {
        Response::ok().html(r#"
            <h1>About This Server</h1>
            <p>An HTTP/1.1 server written in Rust using only std.</p>
            <p>Features: TCP, manual HTTP parsing, thread pool, routing.</p>
        "#)
    });

    // Time API — returns current time as JSON
    router.get("/api/time", |_req| {
        use std::time::{SystemTime, UNIX_EPOCH};
        let ts = SystemTime::now()
            .duration_since(UNIX_EPOCH)
            .map(|d| d.as_secs())
            .unwrap_or(0);
        Response::ok().json(&format!(r#"{{"timestamp": {}}}"#, ts))
    });

    // Echo API — returns the query parameter
    router.get("/api/echo", |req| {
        let message = req.query.get("message")
            .cloned()
            .unwrap_or_else(|| "(no message)".to_string());
        Response::ok().json(&format!(r#"{{"echo": "{}"}}"#, message))
    });

    // Health check
    router.get("/health", |_req| {
        Response::ok().json(r#"{"status": "ok"}"#)
    });

    // Headers inspector
    router.get("/api/headers", |req| {
        let headers_json: String = req.headers
            .iter()
            .map(|(k, v)| format!(r#"  "{}": "{}""#, k, v))
            .collect::<Vec<_>>()
            .join(",\n");
        Response::ok().json(&format!("{{\n{}\n}}", headers_json))
    });

    // 404 fallback is handled by the router automatically

    Server::new("127.0.0.1:7878")
        .with_router(router)
        .workers(4)
        .run();
}
```

---

## Running the Server

```bash
# Start the server
cargo run

# In another terminal:
curl http://localhost:7878/
curl http://localhost:7878/api/time
curl http://localhost:7878/api/echo?message=hello
curl http://localhost:7878/health
curl -v http://localhost:7878/api/headers

# Test 404
curl -v http://localhost:7878/doesnt-exist

# Load test with wrk (if installed)
wrk -t4 -c100 -d10s http://localhost:7878/health
```

---

## What This Server Demonstrates

| Concept | Where Used |
|---------|-----------|
| Ownership | Request parsing takes TcpStream, router handles Request |
| Borrowing | `&Request` passed to handlers, `&Router` shared across threads |
| Lifetimes | `Handler` type with `'static` bound for thread safety |
| Generics | `ThreadPool::execute<F: FnOnce() + Send + 'static>` |
| Traits | `Write` for TcpStream, `Send + Sync` for handlers |
| Enums | `Method`, `ParseError`, `Message` (threadpool) |
| HashMap | Headers, query params, route params |
| Arc | Shared router across thread pool |
| Mutex | Shared channel receiver across workers |
| Drop | ThreadPool graceful shutdown |
| Closures | Route handlers captured by router |
| Pattern matching | Method matching, route matching |
| Error handling | `Result` + `?` in request parsing |
| Iterators | Header parsing, route matching |

---

## Extension Challenges

1. **HTTPS**: Add TLS support using `native-tls` or `rustls` crate
2. **Keep-Alive**: Support persistent connections with `Connection: keep-alive`
3. **Chunked Transfer**: Implement chunked encoding for streaming responses
4. **Middleware**: Add a middleware system (logging, auth, compression)
5. **Path Parameters**: Extract `:id` from routes and pass to handlers
6. **POST handling**: Implement a POST endpoint that echoes the body
7. **Rate limiting**: Implement per-IP rate limiting with `Arc<Mutex<HashMap>>`
8. **Graceful shutdown**: Handle `Ctrl+C` with `std::sync::atomic::AtomicBool`

---

## Congratulations!

You've built a real HTTP server from scratch in Rust, without any frameworks. Along the way, you've applied:

- **Ownership and borrowing**: Thread-safe data sharing
- **Lifetimes**: Ensuring handlers outlive the router
- **Generics and traits**: A flexible, type-safe handler interface
- **Enums**: Clean error and method representation
- **Iterators**: Parsing and transformation pipelines
- **Concurrency**: A production-grade thread pool with graceful shutdown
- **Error handling**: Propagating errors from I/O without panicking

This is Rust. Fast. Safe. No magic.

---

*End of The Rust Book*

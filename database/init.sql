-- Web Crawler Database Schema
-- Supports all test task requirements with proper indexing

-- URLs table - stores target URLs for crawling
CREATE TABLE urls (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    url VARCHAR(2048) NOT NULL,
    status ENUM('queued', 'running', 'completed', 'error') DEFAULT 'queued',
    error_message TEXT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    -- Indexes for performance
    INDEX idx_status (status),
    INDEX idx_created_at (created_at),
    UNIQUE KEY unique_url (url(255)) -- Prevent duplicate URLs
);

-- Crawl results - stores extracted data from each URL
CREATE TABLE crawl_results (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    url_id BIGINT NOT NULL,
    
    -- Required crawl data
    html_version VARCHAR(50) NULL,
    page_title VARCHAR(500) NULL,
    h1_count INT DEFAULT 0,
    h2_count INT DEFAULT 0,
    h3_count INT DEFAULT 0,
    h4_count INT DEFAULT 0,
    h5_count INT DEFAULT 0,
    h6_count INT DEFAULT 0,
    internal_links_count INT DEFAULT 0,
    external_links_count INT DEFAULT 0,
    inaccessible_links_count INT DEFAULT 0,
    has_login_form BOOLEAN DEFAULT FALSE,
    
    -- Metadata
    crawled_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    crawl_duration_ms INT NULL,
    
    FOREIGN KEY (url_id) REFERENCES urls(id) ON DELETE CASCADE,
    INDEX idx_url_id (url_id),
    INDEX idx_crawled_at (crawled_at)
);

-- Links found during crawling - for detailed analysis
CREATE TABLE found_links (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    url_id BIGINT NOT NULL,
    link_url VARCHAR(2048) NOT NULL,
    link_text VARCHAR(500) NULL,
    is_internal BOOLEAN NOT NULL,
    is_accessible BOOLEAN NULL, -- NULL = not checked yet
    status_code INT NULL,
    error_message TEXT NULL,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (url_id) REFERENCES urls(id) ON DELETE CASCADE,
    INDEX idx_url_id (url_id),
    INDEX idx_is_internal (is_internal),
    INDEX idx_is_accessible (is_accessible),
    INDEX idx_status_code (status_code)
);

-- API tokens for authentication
CREATE TABLE api_tokens (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    token_hash VARCHAR(255) NOT NULL,
    name VARCHAR(100) NOT NULL DEFAULT 'Default Token',
    is_active BOOLEAN DEFAULT TRUE,
    expires_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_used_at TIMESTAMP NULL,
    
    UNIQUE KEY unique_token_hash (token_hash),
    INDEX idx_is_active (is_active),
    INDEX idx_expires_at (expires_at)
);

-- Insert default API token for development
-- Token: "dev-token-12345" -> SHA256 hash
INSERT INTO api_tokens (token_hash, name) VALUES 
('a665a45920422f9d417e4867efdc4fb8a04a1f3fff1fa07e998e86f7f7a27ae3', 'Development Token');
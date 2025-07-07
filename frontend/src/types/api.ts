// URL Status Types
export type URLStatus = 'queued' | 'running' | 'completed' | 'error'

// API Response Types
export interface APIResponse<T = any> {
  success: boolean
  data?: T
  error?: {
    code: string
    message: string
    details?: string
  }
  meta?: {
    page: number
    page_size: number
    total: number
    total_pages: number
  }
}

// URL Types
export interface URL {
  id: number
  url: string
  status: URLStatus
  error_message?: string
  created_at: string
  updated_at: string
  crawl_result?: CrawlResult
}

// Crawl Result Types
export interface CrawlResult {
  id: number
  html_version?: string
  page_title?: string
  heading_counts: Record<string, number>
  internal_links_count: number
  external_links_count: number
  inaccessible_links_count: number
  has_login_form: boolean
  crawled_at: string
  crawl_duration_ms?: number
  total_links: number
}

// Found Link Types
export interface FoundLink {
  id: number
  link_url: string
  link_text?: string
  is_internal: boolean
  is_accessible?: boolean
  status_code?: number
  error_message?: string
  is_broken: boolean
  status_category: string
  created_at: string
}

// URL Detail Response (with links)
export interface URLDetail extends URL {
  found_links: FoundLink[]
}

// Request Types
export interface AddURLRequest {
  url: string
}

export interface PaginationParams {
  page?: number
  page_size?: number
  search?: string
  status?: URLStatus
  sort_by?: 'id' | 'url' | 'status' | 'created_at' | 'updated_at'
  sort_dir?: 'asc' | 'desc'
}

// Queue Status Types
export interface QueueStatus {
  is_running: boolean
  queue_length: number
  queue_size: number
}

export interface CrawlStatusResponse {
  queue_manager: QueueStatus
  database_stats: {
    queued_count: number
    running_count: number
    completed_count: number
    error_count: number
  }
}
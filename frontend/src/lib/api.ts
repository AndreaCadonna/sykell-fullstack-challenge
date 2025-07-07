import axios from 'axios'
import type { 
  APIResponse, 
  URL, 
  URLDetail, 
  AddURLRequest, 
  PaginationParams,
  CrawlStatusResponse 
} from '@/types/api'

// Create axios instance with base configuration
const api = axios.create({
  baseURL: import.meta.env.REACT_APP_API_URL || 'http://localhost:8080',
  headers: {
    'Content-Type': 'application/json',
  },
})

// Add auth token to all requests
api.interceptors.request.use((config) => {
  const token = 'dev-token-12345' // Development token
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// Response interceptor for error handling
api.interceptors.response.use(
  (response) => response,
  (error) => {
    console.error('API Error:', error)
    throw error
  }
)

// API Functions
export const urlsApi = {
  // Get paginated list of URLs
  getURLs: async (params: PaginationParams = {}): Promise<APIResponse<URL[]>> => {
    const response = await api.get('/api/urls', { params })
    return response.data
  },

  // Get single URL with basic info
  getURL: async (id: number): Promise<APIResponse<URL>> => {
    const response = await api.get(`/api/urls/${id}`)
    return response.data
  },

  // Get URL with detailed info including links
  getURLDetails: async (id: number): Promise<APIResponse<URLDetail>> => {
    const response = await api.get(`/api/urls/${id}/details`)
    return response.data
  },

  // Add new URL
  addURL: async (data: AddURLRequest): Promise<APIResponse<URL>> => {
    const response = await api.post('/api/urls', data)
    return response.data
  },

  // Delete URL
  deleteURL: async (id: number): Promise<APIResponse<{ message: string }>> => {
    const response = await api.delete(`/api/urls/${id}`)
    return response.data
  },

  // Bulk delete URLs
  bulkDeleteURLs: async (ids: number[]): Promise<APIResponse<{ deleted_count: number }>> => {
    const response = await api.delete('/api/urls/bulk', { data: { ids } })
    return response.data
  },
}

export const crawlApi = {
  // Start crawling a URL
  startCrawl: async (id: number): Promise<APIResponse<any>> => {
    const response = await api.post(`/api/urls/${id}/crawl`)
    return response.data
  },

  // Get crawl status for a URL
  getCrawlStatus: async (id: number): Promise<APIResponse<any>> => {
    const response = await api.get(`/api/urls/${id}/crawl/status`)
    return response.data
  },

  // Start bulk crawl
  startBulkCrawl: async (urlIds: number[]): Promise<APIResponse<any>> => {
    const response = await api.post('/api/crawls/bulk', { url_ids: urlIds })
    return response.data
  },

  // Get queue status
  getQueueStatus: async (): Promise<APIResponse<CrawlStatusResponse>> => {
    const response = await api.get('/api/crawls/queue/status')
    return response.data
  },
}

export const healthApi = {
  // Health check
  getHealth: async (): Promise<any> => {
    const response = await api.get('/health')
    return response.data
  },
}
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { urlsApi, crawlApi } from '@/lib/api'
import type { PaginationParams, AddURLRequest } from '@/types/api'

// Query Keys
export const queryKeys = {
  urls: ['urls'] as const,
  urlsList: (params: PaginationParams) => [...queryKeys.urls, 'list', params] as const,
  urlDetail: (id: number) => [...queryKeys.urls, 'detail', id] as const,
  queueStatus: ['queue', 'status'] as const,
}

// URLs List Hook with real-time updates
export function useURLs(params: PaginationParams = {}) {
  return useQuery({
    queryKey: queryKeys.urlsList(params),
    queryFn: () => urlsApi.getURLs(params),
    refetchInterval: 3000, // Refetch every 3 seconds for real-time updates
    staleTime: 1000, // Consider data stale after 1 second
  })
}

// Single URL Detail Hook
export function useURLDetail(id: number) {
  return useQuery({
    queryKey: queryKeys.urlDetail(id),
    queryFn: () => urlsApi.getURLDetails(id),
    refetchInterval: 5000, // Less frequent updates for detail view
    enabled: !!id, // Only run if ID is provided
  })
}

// Queue Status Hook
export function useQueueStatus() {
  return useQuery({
    queryKey: queryKeys.queueStatus,
    queryFn: () => crawlApi.getQueueStatus(),
    refetchInterval: 2000, // Quick updates for queue status
  })
}

// Add URL Mutation
export function useAddURL() {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: (data: AddURLRequest) => urlsApi.addURL(data),
    onSuccess: () => {
      // Invalidate URLs list to refetch with new data
      queryClient.invalidateQueries({ queryKey: queryKeys.urls })
    },
  })
}

// Delete URL Mutation
export function useDeleteURL() {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: (id: number) => urlsApi.deleteURL(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.urls })
    },
  })
}

// Bulk Delete URLs Mutation
export function useBulkDeleteURLs() {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: (ids: number[]) => urlsApi.bulkDeleteURLs(ids),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.urls })
    },
  })
}

// Start Crawl Mutation
export function useStartCrawl() {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: (id: number) => crawlApi.startCrawl(id),
    onSuccess: (_, id) => {
      // Invalidate the specific URL and the list
      queryClient.invalidateQueries({ queryKey: queryKeys.urlDetail(id) })
      queryClient.invalidateQueries({ queryKey: queryKeys.urls })
      queryClient.invalidateQueries({ queryKey: queryKeys.queueStatus })
    },
  })
}

// Bulk Crawl Mutation
export function useBulkCrawl() {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: (urlIds: number[]) => crawlApi.startBulkCrawl(urlIds),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.urls })
      queryClient.invalidateQueries({ queryKey: queryKeys.queueStatus })
    },
  })
}
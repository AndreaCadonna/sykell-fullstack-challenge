# Frontend Documentation

## Overview

The frontend is a modern React application built with TypeScript, featuring a responsive design that works seamlessly across desktop and mobile devices. It provides a comprehensive interface for managing and monitoring web crawling operations.

## Tech Stack

- **React 19** - Latest React with modern hooks and features
- **TypeScript** - Full type safety throughout the application
- **Vite** - Fast build tool with hot module replacement
- **Tailwind CSS v4** - Utility-first CSS framework with working hot reload
- **shadcn/ui** - High-quality, accessible UI components
- **React Query (@tanstack/react-query)** - Server state management with real-time updates
- **React Router** - Client-side routing
- **Recharts** - Data visualization for charts
- **Axios** - HTTP client for API communication
- **Lucide React** - Beautiful, customizable icons
- **Sonner** - Toast notifications

## Architecture

### Project Structure

```
frontend/
├── src/
│   ├── components/
│   │   ├── ui/              # shadcn/ui components
│   │   ├── layout/          # Layout components
│   │   ├── dashboard/       # Dashboard-specific components
│   │   ├── charts/          # Data visualization components
│   │   └── url-details/     # URL details page components
│   ├── pages/              # Main route components
│   │   ├── Dashboard.tsx   # Main dashboard page
│   │   └── URLDetails.tsx  # URL details page
│   ├── hooks/              # Custom React hooks
│   │   └── useURLs.ts      # React Query hooks for API calls
│   ├── lib/                # Utility libraries
│   │   ├── api.ts          # API client configuration
│   │   └── utils.ts        # Helper functions
│   ├── types/              # TypeScript type definitions
│   │   └── api.ts          # API response types
│   └── test/               # Test utilities and setup
├── public/                 # Static assets
├── Dockerfile             # Multi-stage container build
├── package.json           # Dependencies and scripts
├── tailwind.config.js     # Tailwind CSS configuration
├── vite.config.ts         # Vite build configuration
└── tsconfig.json          # TypeScript configuration
```

### Component Architecture

#### **Layout Components**
- **Layout.tsx** - Main application layout with header, navigation, and footer
- **Header** - Responsive navigation with queue status indicators
- **Footer** - Application information and statistics

#### **Dashboard Components**
- **Dashboard.tsx** - Main dashboard page with URL management
- **AddURLDialog.tsx** - Modal dialog for adding new URLs
- **URLsTable.tsx** - Responsive table/card view for URL listing
- **MobileURLCard.tsx** - Mobile-optimized card layout
- **StatusBadge.tsx** - Visual status indicators with icons
- **PaginationControls.tsx** - Advanced pagination with page size controls

#### **URL Details Components**
- **URLDetails.tsx** - Comprehensive URL analysis page
- **LinkDistributionChart.tsx** - Pie chart showing internal vs external links
- **HeadingsChart.tsx** - Bar chart displaying heading structure (H1-H6)
- **BrokenLinksTable.tsx** - Filterable table of discovered links

#### **Shared UI Components**
- **shadcn/ui components** - Button, Card, Table, Dialog, Input, Select, etc.
- **Responsive design** - All components work on mobile and desktop

## Key Features

### 🎯 **URL Management**
- **Add URLs** - Modal dialog with validation and error handling
- **URL List** - Responsive table/card view with real-time updates
- **Search & Filter** - Search by URL, filter by status (queued, running, completed, error)
- **Sorting** - Sort by URL, status, creation date with visual indicators
- **Bulk Operations** - Select multiple URLs for bulk actions (crawl, delete)

### 📱 **Mobile-First Responsive Design**
- **Adaptive Layout** - Automatically switches between table and card views
- **Touch-Friendly** - Large touch targets and gesture-friendly interface
- **Mobile Filters** - Collapsible filter panel for mobile devices
- **Smart Truncation** - Content adapts to screen size automatically
- **Responsive Charts** - Data visualization works on all screen sizes

### 📊 **Data Visualization**
- **Link Distribution Chart** - Interactive pie chart with tooltips
- **Headings Structure Chart** - Bar chart showing H1-H6 distribution
- **Real-time Statistics** - Live queue status and URL counts
- **Progress Indicators** - Visual feedback for crawl progress

### 🔄 **Real-Time Updates**
- **Auto-refresh** - Data updates every 3 seconds for URLs, 2 seconds for queue
- **Live Status** - Real-time crawl status updates
- **Queue Monitoring** - Live queue length and processing status
- **Optimistic Updates** - Immediate UI feedback for user actions

### 🎨 **Professional UI/UX**
- **Loading States** - Skeleton screens and spinners
- **Error Handling** - Graceful error recovery with retry options
- **Toast Notifications** - Success/error feedback for all actions
- **Empty States** - Contextual messages for empty data
- **Accessibility** - ARIA labels, keyboard navigation, color contrast

## API Integration

### React Query Implementation

The application uses React Query for efficient server state management:

```typescript
// Real-time URL list with automatic updates
export function useURLs(params: PaginationParams = {}) {
  return useQuery({
    queryKey: queryKeys.urlsList(params),
    queryFn: () => urlsApi.getURLs(params),
    refetchInterval: 3000, // Refresh every 3 seconds
    staleTime: 1000,
  })
}

// Queue status monitoring
export function useQueueStatus() {
  return useQuery({
    queryKey: queryKeys.queueStatus,
    queryFn: () => crawlApi.getQueueStatus(),
    refetchInterval: 2000, // Quick updates for queue
  })
}
```

### API Service Layer

Centralized API client with authentication:

```typescript
// Axios instance with auth token
const api = axios.create({
  baseURL: 'http://localhost:8080',
  headers: { 'Content-Type': 'application/json' },
})

// Auto-inject dev token
api.interceptors.request.use((config) => {
  config.headers.Authorization = 'Bearer dev-token-12345'
  return config
})
```

### Mutation Handling

Optimistic updates with error recovery:

```typescript
export function useAddURL() {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: (data: AddURLRequest) => urlsApi.addURL(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.urls })
    },
  })
}
```

## Responsive Design Implementation

### Breakpoint Strategy

The application uses Tailwind's responsive design system:

```typescript
// Mobile-first approach
- Mobile: Default styles (< 768px)
- Tablet: md: prefix (≥ 768px)
- Desktop: lg: prefix (≥ 1024px)
- Large: xl: prefix (≥ 1280px)
```

### Adaptive Components

#### **Table vs Cards**
```typescript
// Automatic layout switching
const [isMobile, setIsMobile] = useState(window.innerWidth < 768)

// Render different layouts based on screen size
{isMobile ? (
  <MobileURLCards />
) : (
  <DesktopTable />
)}
```

#### **Mobile Filter Panel**
```typescript
// Collapsible filters for mobile
const [showFilters, setShowFilters] = useState(false)

// Mobile filter toggle
<Button onClick={() => setShowFilters(!showFilters)}>
  <Filter className="h-4 w-4 mr-2" />
  Filters
  {hasActiveFilters && <Badge>!</Badge>}
</Button>
```

#### **Responsive Navigation**
```typescript
// Adaptive header content
<h1 className="text-lg md:text-xl font-semibold">
  <span className="hidden sm:inline">Web Crawler Dashboard</span>
  <span className="sm:hidden">Crawler</span>
</h1>
```

## State Management

### Server State (React Query)
- **URLs data** - Paginated lists with filters and sorting
- **Queue status** - Real-time monitoring
- **URL details** - Individual URL information with charts
- **Mutations** - Add, delete, crawl operations

### Local State (React hooks)
- **UI state** - Modal visibility, filter panels, selections
- **Form state** - Form inputs and validation
- **Responsive state** - Screen size detection

### Example State Flow
```typescript
// Dashboard state management
const [params, setParams] = useState<PaginationParams>({
  page: 1,
  page_size: 20,
  search: '',
  status: undefined,
  sort_by: 'created_at',
  sort_dir: 'desc',
})

// Real-time data fetching
const { data: urlsResponse, isLoading } = useURLs(params)

// UI state
const [selectedIds, setSelectedIds] = useState<number[]>([])
const [showFilters, setShowFilters] = useState(false)
```

## Performance Optimizations

### Build Optimizations
- **Vite bundling** - Fast builds with tree shaking
- **Code splitting** - Route-based lazy loading
- **Asset optimization** - Optimized images and fonts
- **TypeScript compilation** - Efficient type checking

### Runtime Optimizations
- **React Query caching** - Intelligent data caching and background updates
- **Memo optimization** - React.memo for expensive components
- **Debounced search** - Search input optimization
- **Efficient re-renders** - Optimized dependency arrays

### Network Optimizations
- **Request deduplication** - React Query prevents duplicate requests
- **Background updates** - Stale-while-revalidate pattern
- **Error boundaries** - Graceful error handling
- **Retry logic** - Automatic retry for failed requests

## Development Features

### Hot Reload Setup
```typescript
// Vite configuration with hot reload
export default defineConfig({
  plugins: [react(), tailwindcss()],
  server: {
    host: "0.0.0.0",
    port: 3000,
    watch: { usePolling: true },
  },
})
```

### TypeScript Configuration
```json
{
  "compilerOptions": {
    "strict": true,
    "baseUrl": ".",
    "paths": { "@/*": ["./src/*"] }
  }
}
```

### Development Tools
- **React DevTools** - Component inspection
- **React Query DevTools** - State debugging
- **Vite DevTools** - Build analysis
- **TypeScript** - Real-time type checking

## Testing Strategy

### Testing Setup
- **Vitest** - Fast unit testing framework
- **React Testing Library** - Component testing utilities
- **Jest DOM** - DOM testing matchers
- **User Event** - User interaction simulation

### Test Categories
- **Component Tests** - Individual component functionality
- **Integration Tests** - Component interaction flows
- **API Tests** - Mocked API response handling
- **Accessibility Tests** - ARIA and keyboard navigation

### Example Test Structure
```typescript
describe('Dashboard Component', () => {
  test('renders URL list correctly', async () => {
    render(<Dashboard />)
    await waitFor(() => {
      expect(screen.getByText('URL Dashboard')).toBeInTheDocument()
    })
  })
})
```

## Deployment Configuration

### Docker Multi-Stage Build
```dockerfile
# Development stage with hot reload
FROM node:20-alpine AS development
ENV CHOKIDAR_USEPOLLING=true
CMD ["npm", "run", "dev", "--", "--host"]

# Production stage with nginx
FROM nginx:alpine AS production
COPY --from=builder /app/dist /usr/share/nginx/html
```

### Environment Configuration
```bash
# Development
REACT_APP_API_URL=http://localhost:8080

# Production  
REACT_APP_API_URL=${API_URL}
```

## Accessibility Features

### WCAG Compliance
- **Color contrast** - AA compliant color schemes
- **Keyboard navigation** - Full keyboard support
- **Screen readers** - ARIA labels and descriptions
- **Focus management** - Clear focus indicators

### Implementation Examples
```typescript
// ARIA labels for interactive elements
<Button aria-label="Start crawling URL">
  <Play className="h-4 w-4" />
</Button>

// Screen reader support
<span className="sr-only">Loading URLs</span>

// Keyboard navigation
<Table role="table" aria-label="URLs list">
```

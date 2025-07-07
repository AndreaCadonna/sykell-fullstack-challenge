import { useState } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Badge } from "@/components/ui/badge";
import { AddURLDialog } from "@/components/dashboard/AddURLDialog";
import { URLsTable } from "@/components/dashboard/URLsTable";
import { PaginationControls } from "@/components/dashboard/PaginationControls";
import { useURLs, useQueueStatus } from "@/hooks/useURLs";
import { Search, RefreshCw, Plus, Filter, X } from "lucide-react";
import type { URLStatus, PaginationParams } from "@/types/api";

export function Dashboard() {
  const [params, setParams] = useState<PaginationParams>({
    page: 1,
    page_size: 20,
    search: "",
    status: undefined,
    sort_by: "created_at",
    sort_dir: "desc",
  });
  const [showFilters, setShowFilters] = useState(false);

  const { data: urlsResponse, isLoading, error, refetch } = useURLs(params);
  const { data: queueStatus } = useQueueStatus();

  const handleSearchChange = (search: string) => {
    setParams((prev) => ({ ...prev, search, page: 1 }));
  };

  const handleStatusFilter = (status: URLStatus | "all") => {
    setParams((prev) => ({
      ...prev,
      status: status === "all" ? undefined : status,
      page: 1,
    }));
  };

  const clearFilters = () => {
    setParams((prev) => ({
      ...prev,
      search: "",
      status: undefined,
      page: 1,
    }));
    setShowFilters(false);
  };

  const hasActiveFilters = params.search || params.status;

  const urls = urlsResponse?.data || [];
  const meta = urlsResponse?.meta;

  return (
    <div className="space-y-4 md:space-y-6">
      {/* Header */}
      <div className="flex flex-col space-y-4 md:flex-row md:items-center md:justify-between md:space-y-0">
        <div>
          <h1 className="text-2xl md:text-3xl font-bold text-gray-900">
            URL Dashboard
          </h1>
          <p className="text-gray-600 mt-1 text-sm md:text-base">
            Manage and monitor website crawling operations
          </p>
        </div>
        <AddURLDialog>
          <Button size="default" className="w-full md:w-auto">
            <Plus className="h-4 w-4 mr-2" />
            Add URL
          </Button>
        </AddURLDialog>
      </div>

      {/* Stats Cards */}
      {queueStatus?.data && (
        <div className="grid grid-cols-2 lg:grid-cols-4 gap-3 md:gap-4">
          <Card>
            <CardHeader className="pb-2">
              <CardTitle className="text-xs md:text-sm font-medium text-gray-600">
                Queued
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="text-xl md:text-2xl font-bold text-blue-600">
                {queueStatus.data.database_stats.queued_count}
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="pb-2">
              <CardTitle className="text-xs md:text-sm font-medium text-gray-600">
                Running
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="text-xl md:text-2xl font-bold text-orange-600">
                {queueStatus.data.database_stats.running_count}
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="pb-2">
              <CardTitle className="text-xs md:text-sm font-medium text-gray-600">
                Completed
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="text-xl md:text-2xl font-bold text-green-600">
                {queueStatus.data.database_stats.completed_count}
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="pb-2">
              <CardTitle className="text-xs md:text-sm font-medium text-gray-600">
                Errors
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="text-xl md:text-2xl font-bold text-red-600">
                {queueStatus.data.database_stats.error_count}
              </div>
            </CardContent>
          </Card>
        </div>
      )}

      {/* Filters and Search */}
      <Card>
        <CardHeader>
          <div className="flex flex-col space-y-4 md:flex-row md:items-center md:justify-between md:space-y-0">
            <CardTitle className="text-lg">URLs</CardTitle>

            {/* Mobile Filter Toggle */}
            <div className="flex items-center gap-2 md:hidden">
              <Button
                variant="outline"
                size="sm"
                onClick={() => setShowFilters(!showFilters)}
                className="flex-1"
              >
                <Filter className="h-4 w-4 mr-2" />
                Filters
                {hasActiveFilters && (
                  <Badge
                    variant="secondary"
                    className="ml-2 h-5 w-5 rounded-full p-0 flex items-center justify-center text-xs"
                  >
                    !
                  </Badge>
                )}
              </Button>
              <Button variant="outline" size="sm" onClick={() => refetch()}>
                <RefreshCw className="h-4 w-4" />
              </Button>
            </div>

            {/* Desktop Filters */}
            <div className="hidden md:flex md:items-center md:gap-2">
              {/* Search */}
              <div className="relative">
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-gray-400" />
                <Input
                  placeholder="Search URLs..."
                  value={params.search}
                  onChange={(e) => handleSearchChange(e.target.value)}
                  className="pl-10 w-64"
                />
              </div>

              {/* Status Filter */}
              <Select
                value={params.status || "all"}
                onValueChange={handleStatusFilter}
              >
                <SelectTrigger className="w-32">
                  <SelectValue placeholder="Status" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">All Status</SelectItem>
                  <SelectItem value="queued">Queued</SelectItem>
                  <SelectItem value="running">Running</SelectItem>
                  <SelectItem value="completed">Completed</SelectItem>
                  <SelectItem value="error">Error</SelectItem>
                </SelectContent>
              </Select>

              {/* Clear Filters */}
              {hasActiveFilters && (
                <Button variant="outline" size="sm" onClick={clearFilters}>
                  <X className="h-4 w-4 mr-1" />
                  Clear
                </Button>
              )}

              {/* Refresh Button */}
              <Button variant="outline" onClick={() => refetch()}>
                <RefreshCw className="h-4 w-4" />
              </Button>
            </div>
          </div>

          {/* Mobile Filters Dropdown */}
          {showFilters && (
            <div className="md:hidden space-y-3 pt-4 border-t">
              {/* Search */}
              <div className="relative">
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-gray-400" />
                <Input
                  placeholder="Search URLs..."
                  value={params.search}
                  onChange={(e) => handleSearchChange(e.target.value)}
                  className="pl-10"
                />
              </div>

              {/* Status Filter */}
              <Select
                value={params.status || "all"}
                onValueChange={handleStatusFilter}
              >
                <SelectTrigger>
                  <SelectValue placeholder="Filter by status" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">All Status</SelectItem>
                  <SelectItem value="queued">Queued</SelectItem>
                  <SelectItem value="running">Running</SelectItem>
                  <SelectItem value="completed">Completed</SelectItem>
                  <SelectItem value="error">Error</SelectItem>
                </SelectContent>
              </Select>

              {/* Clear Filters */}
              {hasActiveFilters && (
                <Button
                  variant="outline"
                  onClick={clearFilters}
                  className="w-full"
                >
                  <X className="h-4 w-4 mr-2" />
                  Clear Filters
                </Button>
              )}
            </div>
          )}
        </CardHeader>

        <CardContent>
          {/* Loading State */}
          {isLoading && (
            <div className="flex justify-center py-8">
              <RefreshCw className="h-6 w-6 animate-spin text-gray-400" />
            </div>
          )}

          {/* Error State */}
          {error && (
            <div className="text-center py-8">
              <p className="text-red-600 mb-2">Failed to load URLs</p>
              <Button variant="outline" onClick={() => refetch()}>
                Try Again
              </Button>
            </div>
          )}

          {/* Empty State */}
          {!isLoading && !error && urls.length === 0 && (
            <div className="text-center py-8">
              <p className="text-gray-500 mb-4">
                {params.search || params.status
                  ? "No URLs match your filters"
                  : "No URLs added yet"}
              </p>
              {!params.search && !params.status && (
                <AddURLDialog>
                  <Button>Add Your First URL</Button>
                </AddURLDialog>
              )}
            </div>
          )}

          {/* URLs Table */}
          {!isLoading && !error && urls.length > 0 && (
            <div className="space-y-4">
              <URLsTable
                urls={urls}
                params={params}
                onParamsChange={setParams}
                isLoading={isLoading}
              />

              {/* Pagination */}
              {meta && (
                <PaginationControls
                  params={params}
                  meta={meta}
                  onParamsChange={setParams}
                />
              )}
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}

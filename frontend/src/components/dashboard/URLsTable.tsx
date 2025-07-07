import { useState } from "react";
import { Link } from "react-router-dom";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
import { Badge } from "@/components/ui/badge";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { StatusBadge } from "./StatusBadge";
import {
  useStartCrawl,
  useDeleteURL,
  useBulkDeleteURLs,
  useBulkCrawl,
} from "@/hooks/useURLs";
import { formatDate, truncateUrl } from "@/lib/utils";
import {
  ArrowUpDown,
  ArrowUp,
  ArrowDown,
  MoreHorizontal,
  Play,
  Trash2,
  Eye,
  Loader2,
} from "lucide-react";
import { toast } from "sonner";
import type { URL, PaginationParams } from "@/types/api";

interface URLsTableProps {
  urls: URL[];
  params: PaginationParams;
  onParamsChange: (params: PaginationParams) => void;
  isLoading?: boolean;
}

export function URLsTable({
  urls,
  params,
  onParamsChange,
  isLoading,
}: URLsTableProps) {
  const [selectedIds, setSelectedIds] = useState<number[]>([]);

  const startCrawl = useStartCrawl();
  const deleteURL = useDeleteURL();
  const bulkDeleteURLs = useBulkDeleteURLs();
  const bulkCrawl = useBulkCrawl();

  // Handle sorting
  const handleSort = (column: string) => {
    const newSortDir =
      params.sort_by === column && params.sort_dir === "asc" ? "desc" : "asc";

    onParamsChange({
      ...params,
      sort_by: column as any,
      sort_dir: newSortDir,
    });
  };

  // Handle row selection
  const handleSelectAll = (checked: boolean) => {
    setSelectedIds(checked ? urls.map((url) => url.id) : []);
  };

  const handleSelectRow = (id: number, checked: boolean) => {
    setSelectedIds((prev) =>
      checked ? [...prev, id] : prev.filter((selectedId) => selectedId !== id)
    );
  };

  // Handle actions
  const handleStartCrawl = async (id: number) => {
    try {
      await startCrawl.mutateAsync(id);
      toast.success("Crawl started successfully");
    } catch (error: any) {
      toast.error(
        error.response?.data?.error?.message || "Failed to start crawl"
      );
    }
  };

  const handleDelete = async (id: number) => {
    try {
      await deleteURL.mutateAsync(id);
      toast.success("URL deleted successfully");
      setSelectedIds((prev) => prev.filter((selectedId) => selectedId !== id));
    } catch (error: any) {
      toast.error(
        error.response?.data?.error?.message || "Failed to delete URL"
      );
    }
  };

  const handleBulkDelete = async () => {
    if (selectedIds.length === 0) return;

    try {
      await bulkDeleteURLs.mutateAsync(selectedIds);
      toast.success(`${selectedIds.length} URLs deleted successfully`);
      setSelectedIds([]);
    } catch (error: any) {
      toast.error(
        error.response?.data?.error?.message || "Failed to delete URLs"
      );
    }
  };

  const handleBulkCrawl = async () => {
    if (selectedIds.length === 0) return;

    try {
      await bulkCrawl.mutateAsync(selectedIds);
      toast.success(`Started crawling ${selectedIds.length} URLs`);
      setSelectedIds([]);
    } catch (error: any) {
      toast.error(
        error.response?.data?.error?.message || "Failed to start bulk crawl"
      );
    }
  };

  // Render sort icon
  const renderSortIcon = (column: string) => {
    if (params.sort_by !== column) {
      return <ArrowUpDown className="h-4 w-4" />;
    }
    return params.sort_dir === "asc" ? (
      <ArrowUp className="h-4 w-4" />
    ) : (
      <ArrowDown className="h-4 w-4" />
    );
  };

  const allSelected = urls.length > 0 && selectedIds.length === urls.length;
  const someSelected = selectedIds.length > 0;

  return (
    <div className="space-y-4">
      {/* Bulk Actions */}
      {someSelected && (
        <div className="flex items-center gap-2 p-3 bg-blue-50 rounded-lg border">
          <span className="text-sm font-medium">
            {selectedIds.length} URL{selectedIds.length > 1 ? "s" : ""} selected
          </span>
          <Button
            size="sm"
            onClick={handleBulkCrawl}
            disabled={bulkCrawl.isPending}
          >
            {bulkCrawl.isPending && (
              <Loader2 className="h-4 w-4 mr-1 animate-spin" />
            )}
            <Play className="h-4 w-4 mr-1" />
            Start Crawling
          </Button>
          <Button
            size="sm"
            variant="destructive"
            onClick={handleBulkDelete}
            disabled={bulkDeleteURLs.isPending}
          >
            {bulkDeleteURLs.isPending && (
              <Loader2 className="h-4 w-4 mr-1 animate-spin" />
            )}
            <Trash2 className="h-4 w-4 mr-1" />
            Delete
          </Button>
          <Button
            size="sm"
            variant="outline"
            onClick={() => setSelectedIds([])}
          >
            Clear Selection
          </Button>
        </div>
      )}

      {/* Table */}
      <div className="border rounded-lg">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead className="w-12">
                <Checkbox
                  checked={allSelected}
                  onCheckedChange={handleSelectAll}
                  aria-label="Select all URLs"
                />
              </TableHead>

              <TableHead>
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={() => handleSort("url")}
                  className="h-auto p-0 font-medium"
                >
                  URL
                  {renderSortIcon("url")}
                </Button>
              </TableHead>

              <TableHead>
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={() => handleSort("status")}
                  className="h-auto p-0 font-medium"
                >
                  Status
                  {renderSortIcon("status")}
                </Button>
              </TableHead>

              <TableHead>Title</TableHead>

              <TableHead>Links</TableHead>

              <TableHead>
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={() => handleSort("created_at")}
                  className="h-auto p-0 font-medium"
                >
                  Created
                  {renderSortIcon("created_at")}
                </Button>
              </TableHead>

              <TableHead className="w-12"></TableHead>
            </TableRow>
          </TableHeader>

          <TableBody>
            {urls.map((url) => (
              <TableRow
                key={url.id}
                className={selectedIds.includes(url.id) ? "bg-blue-50" : ""}
              >
                <TableCell>
                  <Checkbox
                    checked={selectedIds.includes(url.id)}
                    onCheckedChange={(checked) =>
                      handleSelectRow(url.id, !!checked)
                    }
                    aria-label={`Select ${url.url}`}
                  />
                </TableCell>

                <TableCell className="font-medium">
                  <Link
                    to={`/urls/${url.id}`}
                    className="text-blue-600 hover:text-blue-800 hover:underline"
                    title={url.url}
                  >
                    {truncateUrl(url.url, 60)}
                  </Link>
                </TableCell>

                <TableCell>
                  <StatusBadge status={url.status} />
                </TableCell>

                <TableCell>
                  {url.crawl_result?.page_title ? (
                    <span
                      className="text-gray-900"
                      title={url.crawl_result.page_title}
                    >
                      {truncateUrl(url.crawl_result.page_title, 40)}
                    </span>
                  ) : (
                    <span className="text-gray-400 italic">
                      {url.status === "completed"
                        ? "No title found"
                        : "Not crawled yet"}
                    </span>
                  )}
                </TableCell>

                <TableCell>
                  {url.crawl_result ? (
                    <div className="flex items-center gap-2">
                      <Badge variant="outline" className="text-xs">
                        {url.crawl_result.internal_links_count} internal
                      </Badge>
                      <Badge variant="outline" className="text-xs">
                        {url.crawl_result.external_links_count} external
                      </Badge>
                    </div>
                  ) : (
                    <span className="text-gray-400 text-sm">—</span>
                  )}
                </TableCell>

                <TableCell className="text-sm text-gray-600">
                  {formatDate(url.created_at)}
                </TableCell>

                <TableCell>
                  <DropdownMenu>
                    <DropdownMenuTrigger asChild>
                      <Button variant="ghost" size="sm">
                        <MoreHorizontal className="h-4 w-4" />
                        <span className="sr-only">Open menu</span>
                      </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent align="end">
                      <DropdownMenuItem asChild>
                        <Link to={`/urls/${url.id}`}>
                          <Eye className="h-4 w-4 mr-2" />
                          View Details
                        </Link>
                      </DropdownMenuItem>

                      {(url.status === "queued" || url.status === "error") && (
                        <DropdownMenuItem
                          onClick={() => handleStartCrawl(url.id)}
                          disabled={startCrawl.isPending}
                        >
                          <Play className="h-4 w-4 mr-2" />
                          {url.status === "error"
                            ? "Retry Crawl"
                            : "Start Crawl"}
                        </DropdownMenuItem>
                      )}

                      <DropdownMenuSeparator />

                      <DropdownMenuItem
                        onClick={() => handleDelete(url.id)}
                        disabled={deleteURL.isPending}
                        className="text-red-600"
                      >
                        <Trash2 className="h-4 w-4 mr-2" />
                        Delete
                      </DropdownMenuItem>
                    </DropdownMenuContent>
                  </DropdownMenu>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>

        {urls.length === 0 && !isLoading && (
          <div className="text-center py-8 text-gray-500">No URLs found</div>
        )}
      </div>
    </div>
  );
}

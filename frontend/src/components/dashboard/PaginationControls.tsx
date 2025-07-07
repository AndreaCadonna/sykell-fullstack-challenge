import {
  Pagination,
  PaginationContent,
  PaginationEllipsis,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from "@/components/ui/pagination";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import type { PaginationParams } from "@/types/api";

interface PaginationControlsProps {
  params: PaginationParams;
  meta: {
    page: number;
    page_size: number;
    total: number;
    total_pages: number;
  };
  onParamsChange: (params: PaginationParams) => void;
}

export function PaginationControls({
  params,
  meta,
  onParamsChange,
}: PaginationControlsProps) {
  const handlePageChange = (page: number) => {
    onParamsChange({ ...params, page });
  };

  const handlePageSizeChange = (pageSize: string) => {
    onParamsChange({
      ...params,
      page_size: parseInt(pageSize),
      page: 1, // Reset to first page when changing page size
    });
  };

  // Generate page numbers to show
  const generatePageNumbers = () => {
    const totalPages = meta.total_pages;
    const currentPage = meta.page;
    const pages: (number | "ellipsis")[] = [];

    if (totalPages <= 7) {
      // Show all pages if 7 or fewer
      for (let i = 1; i <= totalPages; i++) {
        pages.push(i);
      }
    } else {
      // Always show first page
      pages.push(1);

      if (currentPage <= 4) {
        // Near the beginning
        for (let i = 2; i <= 5; i++) {
          pages.push(i);
        }
        pages.push("ellipsis");
        pages.push(totalPages);
      } else if (currentPage >= totalPages - 3) {
        // Near the end
        pages.push("ellipsis");
        for (let i = totalPages - 4; i <= totalPages; i++) {
          pages.push(i);
        }
      } else {
        // In the middle
        pages.push("ellipsis");
        for (let i = currentPage - 1; i <= currentPage + 1; i++) {
          pages.push(i);
        }
        pages.push("ellipsis");
        pages.push(totalPages);
      }
    }

    return pages;
  };

  const pageNumbers = generatePageNumbers();

  return (
    <div className="flex flex-col sm:flex-row items-center justify-between gap-4 pt-4">
      {/* Results info */}
      <div className="text-sm text-gray-600">
        Showing {(meta.page - 1) * meta.page_size + 1} to{" "}
        {Math.min(meta.page * meta.page_size, meta.total)} of {meta.total}{" "}
        results
      </div>

      <div className="flex items-center gap-4">
        {/* Page size selector */}
        <div className="flex items-center gap-2">
          <span className="text-sm text-gray-600">Rows per page:</span>
          <Select
            value={params.page_size?.toString()}
            onValueChange={handlePageSizeChange}
          >
            <SelectTrigger className="w-20">
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="10">10</SelectItem>
              <SelectItem value="20">20</SelectItem>
              <SelectItem value="50">50</SelectItem>
              <SelectItem value="100">100</SelectItem>
            </SelectContent>
          </Select>
        </div>

        {/* Pagination */}
        {meta.total_pages > 1 && (
          <Pagination>
            <PaginationContent>
              <PaginationItem>
                <PaginationPrevious
                  href="#"
                  onClick={(e) => {
                    e.preventDefault();
                    if (meta.page > 1) handlePageChange(meta.page - 1);
                  }}
                  className={
                    meta.page <= 1 ? "pointer-events-none opacity-50" : ""
                  }
                />
              </PaginationItem>

              {pageNumbers.map((pageNumber, index) => (
                <PaginationItem key={index}>
                  {pageNumber === "ellipsis" ? (
                    <PaginationEllipsis />
                  ) : (
                    <PaginationLink
                      href="#"
                      onClick={(e) => {
                        e.preventDefault();
                        handlePageChange(pageNumber);
                      }}
                      isActive={pageNumber === meta.page}
                    >
                      {pageNumber}
                    </PaginationLink>
                  )}
                </PaginationItem>
              ))}

              <PaginationItem>
                <PaginationNext
                  href="#"
                  onClick={(e) => {
                    e.preventDefault();
                    if (meta.page < meta.total_pages)
                      handlePageChange(meta.page + 1);
                  }}
                  className={
                    meta.page >= meta.total_pages
                      ? "pointer-events-none opacity-50"
                      : ""
                  }
                />
              </PaginationItem>
            </PaginationContent>
          </Pagination>
        )}
      </div>
    </div>
  );
}

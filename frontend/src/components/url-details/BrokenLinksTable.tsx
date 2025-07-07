import { useState } from "react";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { ExternalLink, Search, AlertTriangle } from "lucide-react";
import { truncateUrl } from "@/lib/utils";
import type { FoundLink } from "@/types/api";

interface BrokenLinksTableProps {
  links: FoundLink[];
}

export function BrokenLinksTable({ links }: BrokenLinksTableProps) {
  const [search, setSearch] = useState("");
  const [statusFilter, setStatusFilter] = useState<string>("all");
  const [typeFilter, setTypeFilter] = useState<string>("all");

  // Filter links
  const filteredLinks = links.filter((link) => {
    const matchesSearch =
      !search ||
      link.link_url.toLowerCase().includes(search.toLowerCase()) ||
      link.link_text?.toLowerCase().includes(search.toLowerCase());

    const matchesStatus =
      statusFilter === "all" ||
      (statusFilter === "broken" && link.is_broken) ||
      (statusFilter === "accessible" && link.is_accessible === true) ||
      (statusFilter === "unchecked" && link.is_accessible === null);

    const matchesType =
      typeFilter === "all" ||
      (typeFilter === "internal" && link.is_internal) ||
      (typeFilter === "external" && !link.is_internal);

    return matchesSearch && matchesStatus && matchesType;
  });

  // Get status badge variant and color
  const getStatusBadge = (link: FoundLink) => {
    if (link.is_accessible === null) {
      return <Badge variant="outline">Unchecked</Badge>;
    }

    if (link.is_broken) {
      return <Badge variant="destructive">{link.status_code || "Error"}</Badge>;
    }

    if (link.is_accessible) {
      return <Badge variant="default">{link.status_code || "OK"}</Badge>;
    }

    return <Badge variant="outline">Unknown</Badge>;
  };

  const brokenCount = links.filter((link) => link.is_broken).length;
  const totalCount = links.length;

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            Links Analysis
            {brokenCount > 0 && (
              <AlertTriangle className="h-5 w-5 text-red-500" />
            )}
          </div>
          <span className="text-sm font-normal text-gray-600">
            {brokenCount > 0 && (
              <span className="text-red-600">{brokenCount} broken · </span>
            )}
            {totalCount} total
          </span>
        </CardTitle>
      </CardHeader>
      <CardContent>
        {totalCount === 0 ? (
          <div className="text-center py-8 text-gray-500">
            No links found on this page
          </div>
        ) : (
          <>
            {/* Filters */}
            <div className="flex flex-col sm:flex-row gap-3 mb-4">
              <div className="relative flex-1">
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-gray-400" />
                <Input
                  placeholder="Search links..."
                  value={search}
                  onChange={(e) => setSearch(e.target.value)}
                  className="pl-10"
                />
              </div>

              <Select value={statusFilter} onValueChange={setStatusFilter}>
                <SelectTrigger className="w-full sm:w-32">
                  <SelectValue placeholder="Status" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">All Status</SelectItem>
                  <SelectItem value="accessible">Accessible</SelectItem>
                  <SelectItem value="broken">Broken</SelectItem>
                  <SelectItem value="unchecked">Unchecked</SelectItem>
                </SelectContent>
              </Select>

              <Select value={typeFilter} onValueChange={setTypeFilter}>
                <SelectTrigger className="w-full sm:w-32">
                  <SelectValue placeholder="Type" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">All Types</SelectItem>
                  <SelectItem value="internal">Internal</SelectItem>
                  <SelectItem value="external">External</SelectItem>
                </SelectContent>
              </Select>
            </div>

            {/* Links Table */}
            <div className="border rounded-lg">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>URL</TableHead>
                    <TableHead>Text</TableHead>
                    <TableHead>Type</TableHead>
                    <TableHead>Status</TableHead>
                    <TableHead className="w-12"></TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {filteredLinks.length === 0 ? (
                    <TableRow>
                      <TableCell
                        colSpan={5}
                        className="text-center py-8 text-gray-500"
                      >
                        No links match your filters
                      </TableCell>
                    </TableRow>
                  ) : (
                    filteredLinks.map((link) => (
                      <TableRow
                        key={link.id}
                        className={link.is_broken ? "bg-red-50" : ""}
                      >
                        <TableCell className="font-medium">
                          <div className="max-w-md">
                            <a
                              href={link.link_url}
                              target="_blank"
                              rel="noopener noreferrer"
                              className="text-blue-600 hover:text-blue-800 hover:underline"
                              title={link.link_url}
                            >
                              {truncateUrl(link.link_url, 50)}
                            </a>
                          </div>
                        </TableCell>

                        <TableCell>
                          {link.link_text ? (
                            <span
                              className="text-gray-700"
                              title={link.link_text}
                            >
                              {truncateUrl(link.link_text, 30)}
                            </span>
                          ) : (
                            <span className="text-gray-400 italic">
                              No text
                            </span>
                          )}
                        </TableCell>

                        <TableCell>
                          <Badge
                            variant={link.is_internal ? "secondary" : "outline"}
                          >
                            {link.is_internal ? "Internal" : "External"}
                          </Badge>
                        </TableCell>

                        <TableCell>
                          {getStatusBadge(link)}
                          {link.error_message && (
                            <div
                              className="text-xs text-red-600 mt-1"
                              title={link.error_message}
                            >
                              {truncateUrl(link.error_message, 30)}
                            </div>
                          )}
                        </TableCell>

                        <TableCell>
                          <Button variant="ghost" size="sm" asChild>
                            <a
                              href={link.link_url}
                              target="_blank"
                              rel="noopener noreferrer"
                              title="Open link in new tab"
                            >
                              <ExternalLink className="h-4 w-4" />
                            </a>
                          </Button>
                        </TableCell>
                      </TableRow>
                    ))
                  )}
                </TableBody>
              </Table>
            </div>

            {/* Summary */}
            <div className="mt-4 text-sm text-gray-600">
              Showing {filteredLinks.length} of {totalCount} links
            </div>
          </>
        )}
      </CardContent>
    </Card>
  );
}

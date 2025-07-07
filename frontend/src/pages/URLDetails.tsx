import { useParams, Link, useNavigate } from "react-router-dom";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { StatusBadge } from "@/components/dashboard/StatusBadge";
import { LinkDistributionChart } from "@/components/charts/LinkDistributionChart";
import { HeadingsChart } from "@/components/charts/HeadingsChart";
import { BrokenLinksTable } from "@/components/url-details/BrokenLinksTable";
import { useURLDetail, useStartCrawl, useDeleteURL } from "@/hooks/useURLs";
import { formatDate, formatDuration, truncateUrl } from "@/lib/utils";
import {
  ArrowLeft,
  ExternalLink,
  Play,
  Trash2,
  RefreshCw,
  Clock,
  Globe,
  FileText,
  Link as LinkIcon,
  AlertTriangle,
  CheckCircle,
  Loader2,
} from "lucide-react";
import { toast } from "sonner";

export function URLDetails() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const urlId = parseInt(id || "0");

  const { data: response, isLoading, error, refetch } = useURLDetail(urlId);
  const startCrawl = useStartCrawl();
  const deleteURL = useDeleteURL();

  const handleStartCrawl = async () => {
    try {
      await startCrawl.mutateAsync(urlId);
      toast.success("Crawl started successfully");
    } catch (error: any) {
      toast.error(
        error.response?.data?.error?.message || "Failed to start crawl"
      );
    }
  };

  const handleDelete = async () => {
    if (
      !confirm(
        "Are you sure you want to delete this URL? This action cannot be undone."
      )
    ) {
      return;
    }

    try {
      await deleteURL.mutateAsync(urlId);
      toast.success("URL deleted successfully");
      navigate("/");
    } catch (error: any) {
      toast.error(
        error.response?.data?.error?.message || "Failed to delete URL"
      );
    }
  };

  if (isLoading) {
    return (
      <div className="space-y-6">
        <div className="flex items-center gap-4">
          <Skeleton className="h-10 w-24" />
          <Skeleton className="h-8 w-64" />
        </div>
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <Skeleton className="h-32" />
          <Skeleton className="h-32" />
          <Skeleton className="h-32" />
        </div>
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <Skeleton className="h-80" />
          <Skeleton className="h-80" />
        </div>
      </div>
    );
  }

  if (error || !response?.data) {
    return (
      <div className="space-y-6">
        <div className="flex items-center gap-4">
          <Button variant="outline" asChild>
            <Link to="/">
              <ArrowLeft className="h-4 w-4 mr-2" />
              Back to Dashboard
            </Link>
          </Button>
        </div>

        <Alert variant="destructive">
          <AlertTriangle className="h-4 w-4" />
          <AlertDescription>
            {error ? "Failed to load URL details" : "URL not found"}
          </AlertDescription>
        </Alert>

        <Button onClick={() => refetch()}>
          <RefreshCw className="h-4 w-4 mr-2" />
          Try Again
        </Button>
      </div>
    );
  }

  const urlDetail = response.data;
  const crawlResult = urlDetail.crawl_result;
  const foundLinks = urlDetail.found_links || [];

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col space-y-4 lg:flex-row lg:items-start lg:justify-between lg:space-y-0">
        <div className="flex flex-col sm:flex-row sm:items-center gap-4">
          <Button variant="outline" asChild>
            <Link to="/">
              <ArrowLeft className="h-4 w-4 mr-2" />
              Back
            </Link>
          </Button>
          <div>
            <h1 className="text-xl md:text-2xl font-bold text-gray-900">
              URL Details
            </h1>
            <p className="text-gray-600 mt-1 text-sm md:text-base">
              Analysis and crawl results
            </p>
          </div>
        </div>

        <div className="flex flex-col sm:flex-row items-stretch sm:items-center gap-2">
          {(urlDetail.status === "queued" || urlDetail.status === "error") && (
            <Button
              onClick={handleStartCrawl}
              disabled={startCrawl.isPending}
              className="w-full sm:w-auto"
            >
              {startCrawl.isPending && (
                <Loader2 className="h-4 w-4 mr-2 animate-spin" />
              )}
              <Play className="h-4 w-4 mr-2" />
              {urlDetail.status === "error" ? "Retry Crawl" : "Start Crawl"}
            </Button>
          )}

          <Button
            variant="outline"
            onClick={() => refetch()}
            className="w-full sm:w-auto"
          >
            <RefreshCw className="h-4 w-4 mr-2" />
            Refresh
          </Button>

          <Button
            variant="destructive"
            onClick={handleDelete}
            disabled={deleteURL.isPending}
            className="w-full sm:w-auto"
          >
            {deleteURL.isPending && (
              <Loader2 className="h-4 w-4 mr-2 animate-spin" />
            )}
            <Trash2 className="h-4 w-4 mr-2" />
            Delete
          </Button>
        </div>
      </div>

      {/* URL Info Card */}
      <Card>
        <CardHeader>
          <div className="flex flex-col space-y-4 lg:flex-row lg:items-center lg:justify-between lg:space-y-0">
            <CardTitle className="flex flex-col sm:flex-row sm:items-center gap-2 break-all">
              <div className="flex items-center gap-2 min-w-0">
                <Globe className="h-5 w-5 flex-shrink-0" />
                <a
                  href={urlDetail.url}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-blue-600 hover:text-blue-800 hover:underline break-all"
                  title={urlDetail.url}
                >
                  {window.innerWidth < 640
                    ? truncateUrl(urlDetail.url, 40)
                    : truncateUrl(urlDetail.url, 80)}
                </a>
                <ExternalLink className="h-4 w-4 text-gray-400 flex-shrink-0" />
              </div>
            </CardTitle>
            <StatusBadge
              status={urlDetail.status}
              className="self-start lg:self-center"
            />
          </div>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
            <div>
              <div className="text-sm text-gray-600">Created</div>
              <div className="font-medium text-sm md:text-base">
                {formatDate(urlDetail.created_at)}
              </div>
            </div>
            <div>
              <div className="text-sm text-gray-600">Last Updated</div>
              <div className="font-medium text-sm md:text-base">
                {formatDate(urlDetail.updated_at)}
              </div>
            </div>
            {crawlResult?.crawled_at && (
              <div>
                <div className="text-sm text-gray-600">Crawled</div>
                <div className="font-medium text-sm md:text-base">
                  {formatDate(crawlResult.crawled_at)}
                </div>
              </div>
            )}
            {crawlResult?.crawl_duration_ms && (
              <div>
                <div className="text-sm text-gray-600">Duration</div>
                <div className="font-medium text-sm md:text-base">
                  {formatDuration(crawlResult.crawl_duration_ms)}
                </div>
              </div>
            )}
          </div>

          {urlDetail.error_message && (
            <Alert variant="destructive" className="mt-4">
              <AlertTriangle className="h-4 w-4" />
              <AlertDescription>{urlDetail.error_message}</AlertDescription>
            </Alert>
          )}
        </CardContent>
      </Card>

      {/* Crawl Results Summary */}
      {crawlResult && (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-5 gap-3 md:gap-4">
          <Card>
            <CardContent className="p-3 md:p-4">
              <div className="flex items-center gap-2">
                <FileText className="h-4 w-4 text-blue-600 flex-shrink-0" />
                <div className="min-w-0">
                  <div className="text-xs md:text-sm text-gray-600">
                    Page Title
                  </div>
                  <div
                    className="font-medium text-sm md:text-base break-words"
                    title={crawlResult.page_title || "No title"}
                  >
                    {crawlResult.page_title
                      ? window.innerWidth < 640
                        ? truncateUrl(crawlResult.page_title, 15)
                        : truncateUrl(crawlResult.page_title, 20)
                      : "No title found"}
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="p-3 md:p-4">
              <div className="flex items-center gap-2">
                <Globe className="h-4 w-4 text-green-600 flex-shrink-0" />
                <div className="min-w-0">
                  <div className="text-xs md:text-sm text-gray-600">
                    HTML Version
                  </div>
                  <div className="font-medium text-sm md:text-base">
                    {crawlResult.html_version || "Unknown"}
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="p-3 md:p-4">
              <div className="flex items-center gap-2">
                <LinkIcon className="h-4 w-4 text-purple-600 flex-shrink-0" />
                <div className="min-w-0">
                  <div className="text-xs md:text-sm text-gray-600">
                    Total Links
                  </div>
                  <div className="font-medium text-sm md:text-base">
                    {crawlResult.total_links}
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="p-3 md:p-4">
              <div className="flex items-center gap-2">
                {crawlResult.has_login_form ? (
                  <CheckCircle className="h-4 w-4 text-green-600 flex-shrink-0" />
                ) : (
                  <Clock className="h-4 w-4 text-gray-400 flex-shrink-0" />
                )}
                <div className="min-w-0">
                  <div className="text-xs md:text-sm text-gray-600">
                    Login Form
                  </div>
                  <div className="font-medium text-sm md:text-base">
                    {crawlResult.has_login_form ? "Detected" : "Not found"}
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="p-3 md:p-4">
              <div className="flex items-center gap-2">
                <AlertTriangle className="h-4 w-4 text-red-600 flex-shrink-0" />
                <div className="min-w-0">
                  <div className="text-xs md:text-sm text-gray-600">
                    Broken Links
                  </div>
                  <div className="font-medium text-sm md:text-base text-red-600">
                    {foundLinks.filter((link) => link.is_broken).length}
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>
      )}

      {/* Charts Section */}
      {crawlResult && (
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <LinkDistributionChart crawlResult={crawlResult} />
          <HeadingsChart crawlResult={crawlResult} />
        </div>
      )}

      {/* Links Analysis */}
      {foundLinks.length > 0 && <BrokenLinksTable links={foundLinks} />}

      {/* No crawl results message */}
      {!crawlResult && urlDetail.status === "queued" && (
        <Card>
          <CardContent className="p-8 text-center">
            <Clock className="h-12 w-12 text-gray-400 mx-auto mb-4" />
            <h3 className="text-lg font-medium text-gray-900 mb-2">
              Waiting to be crawled
            </h3>
            <p className="text-gray-600 mb-4">
              This URL is queued for crawling. Results will appear here once the
              crawl is complete.
            </p>
            <Button onClick={handleStartCrawl} disabled={startCrawl.isPending}>
              {startCrawl.isPending && (
                <Loader2 className="h-4 w-4 mr-2 animate-spin" />
              )}
              <Play className="h-4 w-4 mr-2" />
              Start Crawl Now
            </Button>
          </CardContent>
        </Card>
      )}

      {!crawlResult && urlDetail.status === "running" && (
        <Card>
          <CardContent className="p-8 text-center">
            <Loader2 className="h-12 w-12 text-blue-600 mx-auto mb-4 animate-spin" />
            <h3 className="text-lg font-medium text-gray-900 mb-2">
              Crawl in progress
            </h3>
            <p className="text-gray-600">
              This URL is currently being crawled. Results will appear here once
              the process is complete.
            </p>
          </CardContent>
        </Card>
      )}

      {!crawlResult && urlDetail.status === "error" && (
        <Card>
          <CardContent className="p-8 text-center">
            <AlertTriangle className="h-12 w-12 text-red-500 mx-auto mb-4" />
            <h3 className="text-lg font-medium text-gray-900 mb-2">
              Crawl failed
            </h3>
            <p className="text-gray-600 mb-4">
              The crawl process failed. Check the error message above and try
              again.
            </p>
            <Button onClick={handleStartCrawl} disabled={startCrawl.isPending}>
              {startCrawl.isPending && (
                <Loader2 className="h-4 w-4 mr-2 animate-spin" />
              )}
              <RefreshCw className="h-4 w-4 mr-2" />
              Retry Crawl
            </Button>
          </CardContent>
        </Card>
      )}
    </div>
  );
}

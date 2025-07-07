import type { ReactNode } from "react";
import { Link, useLocation } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { useQueueStatus } from "@/hooks/useURLs";
import { Activity, Database, Home } from "lucide-react";

interface LayoutProps {
  children: ReactNode;
}

export function Layout({ children }: LayoutProps) {
  const location = useLocation();
  const { data: queueStatus } = useQueueStatus();

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-white border-b border-gray-200 sticky top-0 z-40">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-14 md:h-16">
            {/* Logo/Title */}
            <div className="flex items-center space-x-2 md:space-x-4">
              <Database className="h-6 w-6 md:h-8 md:w-8 text-blue-600" />
              <h1 className="text-lg md:text-xl font-semibold text-gray-900">
                <span className="hidden sm:inline">Web Crawler Dashboard</span>
                <span className="sm:hidden">Crawler</span>
              </h1>
            </div>

            {/* Navigation */}
            <nav className="flex items-center space-x-2 md:space-x-4">
              <Button
                variant={location.pathname === "/" ? "default" : "ghost"}
                size="sm"
                asChild
                className="hidden sm:flex"
              >
                <Link to="/">
                  <Home className="h-4 w-4 mr-2" />
                  Dashboard
                </Link>
              </Button>

              {/* Mobile home button */}
              <Button
                variant={location.pathname === "/" ? "default" : "ghost"}
                size="sm"
                asChild
                className="sm:hidden"
              >
                <Link to="/">
                  <Home className="h-4 w-4" />
                </Link>
              </Button>

              {/* Queue Status Indicator */}
              {queueStatus?.data && (
                <div className="hidden md:flex items-center space-x-2 text-sm text-gray-600">
                  <Activity className="h-4 w-4" />
                  <span>Queue:</span>
                  <Badge variant="secondary" className="text-xs">
                    {queueStatus.data.queue_manager.queue_length} /{" "}
                    {queueStatus.data.queue_manager.queue_size}
                  </Badge>
                  <Badge
                    variant={
                      queueStatus.data.queue_manager.is_running
                        ? "default"
                        : "destructive"
                    }
                    className="text-xs"
                  >
                    {queueStatus.data.queue_manager.is_running
                      ? "Running"
                      : "Stopped"}
                  </Badge>
                </div>
              )}

              {/* Mobile Queue Status */}
              {queueStatus?.data && (
                <div className="md:hidden">
                  <Badge
                    variant={
                      queueStatus.data.queue_manager.is_running
                        ? "default"
                        : "destructive"
                    }
                    className="text-xs"
                  >
                    {queueStatus.data.queue_manager.queue_length}
                  </Badge>
                </div>
              )}
            </nav>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4 md:py-8">
        {children}
      </main>

      {/* Footer */}
      <footer className="bg-white border-t border-gray-200 mt-auto">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-3 md:py-4">
          <div className="flex flex-col sm:flex-row justify-between items-center text-xs md:text-sm text-gray-600 space-y-2 sm:space-y-0">
            <p>Web Crawler Dashboard - Full-Stack Test Task</p>
            {queueStatus?.data && (
              <div className="flex items-center space-x-2 md:space-x-4 text-xs">
                <span>
                  Total:{" "}
                  {queueStatus.data.database_stats.queued_count +
                    queueStatus.data.database_stats.running_count +
                    queueStatus.data.database_stats.completed_count +
                    queueStatus.data.database_stats.error_count}
                </span>
                <span className="hidden sm:inline">
                  Completed: {queueStatus.data.database_stats.completed_count}
                </span>
                <span className="sm:hidden">
                  ✓ {queueStatus.data.database_stats.completed_count}
                </span>
              </div>
            )}
          </div>
        </div>
      </footer>
    </div>
  );
}

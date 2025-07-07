import { Badge } from "@/components/ui/badge";
import { Loader2, CheckCircle, Clock, AlertCircle } from "lucide-react";
import type { URLStatus } from "@/types/api";

interface StatusBadgeProps {
  status: URLStatus;
  className?: string;
}

export function StatusBadge({ status, className }: StatusBadgeProps) {
  const getStatusConfig = (status: URLStatus) => {
    switch (status) {
      case "queued":
        return {
          variant: "outline" as const,
          icon: <Clock className="h-3 w-3" />,
          label: "Queued",
          color: "text-blue-600",
        };
      case "running":
        return {
          variant: "secondary" as const,
          icon: <Loader2 className="h-3 w-3 animate-spin" />,
          label: "Running",
          color: "text-orange-600",
        };
      case "completed":
        return {
          variant: "default" as const,
          icon: <CheckCircle className="h-3 w-3" />,
          label: "Completed",
          color: "text-green-600",
        };
      case "error":
        return {
          variant: "destructive" as const,
          icon: <AlertCircle className="h-3 w-3" />,
          label: "Error",
          color: "text-red-600",
        };
      default:
        return {
          variant: "outline" as const,
          icon: null,
          label: status,
          color: "text-gray-600",
        };
    }
  };

  const config = getStatusConfig(status);

  return (
    <Badge variant={config.variant} className={className}>
      <div className="flex items-center gap-1">
        {config.icon}
        <span>{config.label}</span>
      </div>
    </Badge>
  );
}

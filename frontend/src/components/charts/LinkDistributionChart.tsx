import {
  PieChart,
  Pie,
  Cell,
  ResponsiveContainer,
  Tooltip,
  Legend,
} from "recharts";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import type { CrawlResult } from "@/types/api";

interface LinkDistributionChartProps {
  crawlResult: CrawlResult;
}

export function LinkDistributionChart({
  crawlResult,
}: LinkDistributionChartProps) {
  const data = [
    {
      name: "Internal Links",
      value: crawlResult.internal_links_count,
      color: "#3b82f6", // Blue
    },
    {
      name: "External Links",
      value: crawlResult.external_links_count,
      color: "#10b981", // Green
    },
  ];

  // Only show chart if there are links
  const hasLinks =
    crawlResult.internal_links_count > 0 ||
    crawlResult.external_links_count > 0;

  if (!hasLinks) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>Link Distribution</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex items-center justify-center h-48 text-gray-500">
            No links found on this page
          </div>
        </CardContent>
      </Card>
    );
  }

  const renderCustomTooltip = ({ active, payload }: any) => {
    if (active && payload && payload.length) {
      const data = payload[0];
      return (
        <div className="bg-white p-3 border rounded-lg shadow-md">
          <p className="font-medium">{data.name}</p>
          <p className="text-sm text-gray-600">
            {data.value} link{data.value !== 1 ? "s" : ""}
          </p>
          <p className="text-sm text-gray-600">
            {((data.value / crawlResult.total_links) * 100).toFixed(1)}%
          </p>
        </div>
      );
    }
    return null;
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle>Link Distribution</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="h-64">
          <ResponsiveContainer width="100%" height="100%">
            <PieChart>
              <Pie
                data={data}
                cx="50%"
                cy="50%"
                innerRadius={40}
                outerRadius={80}
                paddingAngle={2}
                dataKey="value"
              >
                {data.map((entry, index) => (
                  <Cell key={`cell-${index}`} fill={entry.color} />
                ))}
              </Pie>
              <Tooltip content={renderCustomTooltip} />
              <Legend
                verticalAlign="bottom"
                height={36}
                formatter={(value, entry: any) => (
                  <span style={{ color: entry.color }}>{value}</span>
                )}
              />
            </PieChart>
          </ResponsiveContainer>
        </div>

        {/* Summary stats */}
        <div className="grid grid-cols-2 gap-4 mt-4 pt-4 border-t">
          <div className="text-center">
            <div className="text-2xl font-bold text-blue-600">
              {crawlResult.internal_links_count}
            </div>
            <div className="text-sm text-gray-600">Internal Links</div>
            <div className="text-xs text-gray-500">
              {crawlResult.total_links > 0
                ? `${(
                    (crawlResult.internal_links_count /
                      crawlResult.total_links) *
                    100
                  ).toFixed(1)}%`
                : "0%"}
            </div>
          </div>
          <div className="text-center">
            <div className="text-2xl font-bold text-green-600">
              {crawlResult.external_links_count}
            </div>
            <div className="text-sm text-gray-600">External Links</div>
            <div className="text-xs text-gray-500">
              {crawlResult.total_links > 0
                ? `${(
                    (crawlResult.external_links_count /
                      crawlResult.total_links) *
                    100
                  ).toFixed(1)}%`
                : "0%"}
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}

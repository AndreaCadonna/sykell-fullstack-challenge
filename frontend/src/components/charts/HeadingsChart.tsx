import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
} from "recharts";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import type { CrawlResult } from "@/types/api";

interface HeadingsChartProps {
  crawlResult: CrawlResult;
}

export function HeadingsChart({ crawlResult }: HeadingsChartProps) {
  const data = [
    { name: "H1", count: crawlResult.heading_counts.h1 || 0 },
    { name: "H2", count: crawlResult.heading_counts.h2 || 0 },
    { name: "H3", count: crawlResult.heading_counts.h3 || 0 },
    { name: "H4", count: crawlResult.heading_counts.h4 || 0 },
    { name: "H5", count: crawlResult.heading_counts.h5 || 0 },
    { name: "H6", count: crawlResult.heading_counts.h6 || 0 },
  ];

  const totalHeadings = data.reduce((sum, item) => sum + item.count, 0);

  const renderCustomTooltip = ({ active, payload, label }: any) => {
    if (active && payload && payload.length) {
      const count = payload[0].value;
      return (
        <div className="bg-white p-3 border rounded-lg shadow-md">
          <p className="font-medium">{label} Tags</p>
          <p className="text-sm text-gray-600">
            {count} heading{count !== 1 ? "s" : ""}
          </p>
        </div>
      );
    }
    return null;
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center justify-between">
          Heading Structure
          <span className="text-sm font-normal text-gray-600">
            {totalHeadings} total headings
          </span>
        </CardTitle>
      </CardHeader>
      <CardContent>
        {totalHeadings === 0 ? (
          <div className="flex items-center justify-center h-48 text-gray-500">
            No headings found on this page
          </div>
        ) : (
          <div className="h-48">
            <ResponsiveContainer width="100%" height="100%">
              <BarChart
                data={data}
                margin={{ top: 20, right: 30, left: 20, bottom: 5 }}
              >
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis
                  dataKey="name"
                  tick={{ fontSize: 12 }}
                  tickLine={{ stroke: "#gray" }}
                />
                <YAxis tick={{ fontSize: 12 }} tickLine={{ stroke: "#gray" }} />
                <Tooltip content={renderCustomTooltip} />
                <Bar dataKey="count" fill="#8b5cf6" radius={[2, 2, 0, 0]} />
              </BarChart>
            </ResponsiveContainer>
          </div>
        )}

        {/* Summary grid */}
        {totalHeadings > 0 && (
          <div className="grid grid-cols-6 gap-2 mt-4 pt-4 border-t text-center">
            {data.map((item) => (
              <div key={item.name} className="text-center">
                <div className="text-lg font-bold text-purple-600">
                  {item.count}
                </div>
                <div className="text-xs text-gray-600">{item.name}</div>
              </div>
            ))}
          </div>
        )}
      </CardContent>
    </Card>
  );
}

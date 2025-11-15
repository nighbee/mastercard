import { Card } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Download, Code } from "lucide-react";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, LineChart, Line, PieChart, Pie, Cell } from "recharts";
import { useState } from "react";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog";

interface ResultsViewerProps {
  results: {
    type?: "table" | "chart" | "text" | "error";
    data?: any;
    sql_query?: string | null;
    result_data?: string;
    result_format?: "text" | "table" | "chart" | "error" | null;
    error_message?: string | null;
  };
}

const ResultsViewer = ({ results }: ResultsViewerProps) => {
  const [showSQL, setShowSQL] = useState(false);
  
  // Parse result_data if it's a JSON string
  let parsedData: any[] = [];
  let format = results.result_format || results.type || "text";
  
  if (results.result_data) {
    try {
      parsedData = JSON.parse(results.result_data);
    } catch (e) {
      // If parsing fails, try to use data directly
      parsedData = results.data || [];
    }
  } else if (results.data) {
    parsedData = Array.isArray(results.data) ? results.data : [results.data];
  }

  // Handle error case
  if (format === "error" || results.error_message) {
    return (
      <Card className="p-6 border-destructive">
        <div className="text-destructive">
          <h3 className="text-lg font-semibold mb-2">Error</h3>
          <p>{results.error_message || "An error occurred while executing the query"}</p>
        </div>
      </Card>
    );
  }

  // Handle text result (single value) - show just the value, no heading or card
  if (format === "text" && parsedData.length === 1) {
    const firstRow = parsedData[0];
    const values = Object.values(firstRow);
    if (values.length === 1) {
      return (
        <>
          <div className="mt-3 flex items-center justify-start gap-2">
            <div className="text-3xl font-bold">{String(values[0])}</div>
            {results.sql_query && (
              <Button variant="outline" size="sm" onClick={() => setShowSQL(true)}>
                <Code className="w-4 h-4 mr-2" />
                View SQL
              </Button>
            )}
          </div>
          <Dialog open={showSQL} onOpenChange={setShowSQL}>
            <DialogContent className="max-w-2xl">
              <DialogHeader>
                <DialogTitle>Generated SQL Query</DialogTitle>
              </DialogHeader>
              <pre className="bg-muted p-4 rounded-lg overflow-auto max-h-96">
                <code>{results.sql_query}</code>
              </pre>
            </DialogContent>
          </Dialog>
        </>
      );
    }
  }

  // Handle table/chart results
  if (parsedData.length === 0) {
    return (
      <Card className="p-6">
        <div className="text-center text-muted-foreground">No results found</div>
      </Card>
    );
  }

  // Get column names from first row
  const columns = parsedData.length > 0 ? Object.keys(parsedData[0]) : [];

  // Determine if we should show chart (if we have numeric data)
  const hasNumericData = parsedData.some((row) =>
    columns.some((col) => typeof row[col] === "number")
  );

  return (
    <>
    <Card className="p-6">
      <div className="flex items-center justify-between mb-4">
        <h3 className="text-lg font-semibold">Query Results</h3>
        <div className="flex gap-2">
            {results.sql_query && (
              <Button variant="outline" size="sm" onClick={() => setShowSQL(true)}>
            <Code className="w-4 h-4 mr-2" />
            View SQL
          </Button>
            )}
            <Button variant="outline" size="sm" onClick={() => {
              // Export to CSV
              const csv = [
                columns.join(","),
                ...parsedData.map((row) =>
                  columns.map((col) => JSON.stringify(row[col] || "")).join(",")
                ),
              ].join("\n");
              const blob = new Blob([csv], { type: "text/csv" });
              const url = URL.createObjectURL(blob);
              const a = document.createElement("a");
              a.href = url;
              a.download = "query-results.csv";
              a.click();
            }}>
            <Download className="w-4 h-4 mr-2" />
              Export CSV
          </Button>
        </div>
      </div>

      <Tabs defaultValue="table" className="w-full">
        <TabsList>
          <TabsTrigger value="table">Table</TabsTrigger>
            {hasNumericData && <TabsTrigger value="chart">Chart</TabsTrigger>}
        </TabsList>

        <TabsContent value="table" className="mt-4">
            <div className="border rounded-lg overflow-auto max-h-96">
            <Table>
              <TableHeader>
                <TableRow>
                    {columns.map((col) => (
                      <TableHead key={col}>{col}</TableHead>
                    ))}
                </TableRow>
              </TableHeader>
              <TableBody>
                  {parsedData.map((row, index) => (
                  <TableRow key={index}>
                      {columns.map((col) => (
                        <TableCell key={col}>
                          {row[col] !== null && row[col] !== undefined
                            ? String(row[col])
                            : ""}
                        </TableCell>
                      ))}
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </div>
        </TabsContent>

          {hasNumericData && (
        <TabsContent value="chart" className="mt-4">
          <ResponsiveContainer width="100%" height={300}>
                <BarChart data={parsedData}>
              <CartesianGrid strokeDasharray="3 3" />
                  <XAxis dataKey={columns[0]} />
              <YAxis />
              <Tooltip />
                  {columns.slice(1).map((col, idx) => {
                    const isNumeric = parsedData.some((row) => typeof row[col] === "number");
                    if (isNumeric) {
                      return <Bar key={col} dataKey={col} fill={`hsl(${idx * 60}, 70%, 50%)`} />;
                    }
                    return null;
                  })}
            </BarChart>
          </ResponsiveContainer>
        </TabsContent>
          )}
      </Tabs>
    </Card>

      <Dialog open={showSQL} onOpenChange={setShowSQL}>
        <DialogContent className="max-w-2xl">
          <DialogHeader>
            <DialogTitle>Generated SQL Query</DialogTitle>
          </DialogHeader>
          <pre className="bg-muted p-4 rounded-lg overflow-auto max-h-96">
            <code>{results.sql_query}</code>
          </pre>
        </DialogContent>
      </Dialog>
    </>
  );
};

export default ResultsViewer;

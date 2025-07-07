import { useState } from "react";
import { useAddURL } from "@/hooks/useURLs";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Plus, Loader2 } from "lucide-react";
import { toast } from "sonner";

interface AddURLDialogProps {
  children?: React.ReactNode;
}

export function AddURLDialog({ children }: AddURLDialogProps) {
  const [open, setOpen] = useState(false);
  const [url, setUrl] = useState("");
  const [error, setError] = useState("");

  const addURL = useAddURL();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");

    // Basic URL validation
    if (!url.trim()) {
      setError("URL is required");
      return;
    }

    // Check if URL has protocol
    const urlPattern = /^https?:\/\/.+/;
    if (!urlPattern.test(url.trim())) {
      setError("URL must start with http:// or https://");
      return;
    }

    try {
      await addURL.mutateAsync({ url: url.trim() });
      toast.success("URL added successfully!");
      setUrl("");
      setOpen(false);
    } catch (error: any) {
      const errorMessage =
        error.response?.data?.error?.message || "Failed to add URL";
      setError(errorMessage);
      toast.error(errorMessage);
    }
  };

  const handleOpenChange = (newOpen: boolean) => {
    setOpen(newOpen);
    if (!newOpen) {
      setUrl("");
      setError("");
    }
  };

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <DialogTrigger asChild>
        {children || (
          <Button>
            <Plus className="h-4 w-4 mr-2" />
            Add URL
          </Button>
        )}
      </DialogTrigger>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Add New URL</DialogTitle>
          <DialogDescription>
            Enter a website URL to crawl and analyze. The URL will be queued for
            processing.
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit}>
          <div className="grid gap-4 py-4">
            <div className="space-y-2">
              <label htmlFor="url" className="text-sm font-medium">
                Website URL
              </label>
              <Input
                id="url"
                type="url"
                placeholder="https://example.com"
                value={url}
                onChange={(e) => setUrl(e.target.value)}
                disabled={addURL.isPending}
              />
            </div>

            {error && (
              <Alert variant="destructive">
                <AlertDescription>{error}</AlertDescription>
              </Alert>
            )}
          </div>

          <DialogFooter>
            <Button
              type="button"
              variant="outline"
              onClick={() => setOpen(false)}
              disabled={addURL.isPending}
            >
              Cancel
            </Button>
            <Button type="submit" disabled={addURL.isPending}>
              {addURL.isPending && (
                <Loader2 className="h-4 w-4 mr-2 animate-spin" />
              )}
              Add URL
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}

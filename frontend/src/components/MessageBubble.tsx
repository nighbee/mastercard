import { Card } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { GitBranch, Copy, CheckCircle2 } from "lucide-react";
import { cn } from "@/lib/utils";
import { useState } from "react";

interface Message {
  id: string;
  role: "user" | "assistant";
  content: string;
  timestamp: Date;
}

interface MessageBubbleProps {
  message: Message;
}

const MessageBubble = ({ message }: MessageBubbleProps) => {
  const [copied, setCopied] = useState(false);

  const handleCopy = () => {
    navigator.clipboard.writeText(message.content);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <div
      className={cn(
        "flex gap-3",
        message.role === "user" ? "justify-end" : "justify-start"
      )}
    >
      {message.role === "assistant" && (
        <div className="w-8 h-8 bg-primary rounded-full flex items-center justify-center flex-shrink-0">
          <span className="text-primary-foreground text-sm font-semibold">AI</span>
        </div>
      )}
      <div className={cn("max-w-2xl", message.role === "user" && "order-first")}>
        <Card
          className={cn(
            "p-4",
            message.role === "user"
              ? "bg-primary text-primary-foreground"
              : "bg-card"
          )}
        >
          <p className="whitespace-pre-wrap">{message.content}</p>
          <div className="flex items-center gap-2 mt-3 pt-3 border-t border-border/50">
            <span className="text-xs opacity-70">
              {message.timestamp.toLocaleTimeString()}
            </span>
            <div className="flex gap-1 ml-auto">
              <Button
                variant="ghost"
                size="sm"
                onClick={handleCopy}
                className="h-7 px-2"
              >
                {copied ? (
                  <CheckCircle2 className="w-4 h-4" />
                ) : (
                  <Copy className="w-4 h-4" />
                )}
              </Button>
              {message.role === "assistant" && (
                <Button variant="ghost" size="sm" className="h-7 px-2">
                  <GitBranch className="w-4 h-4" />
                </Button>
              )}
            </div>
          </div>
        </Card>
      </div>
      {message.role === "user" && (
        <div className="w-8 h-8 bg-muted rounded-full flex items-center justify-center flex-shrink-0">
          <span className="text-foreground text-sm font-semibold">U</span>
        </div>
      )}
    </div>
  );
};

export default MessageBubble;

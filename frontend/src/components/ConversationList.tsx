import { ScrollArea } from "@/components/ui/scroll-area";
import { Input } from "@/components/ui/input";
import { Search, MessageSquare, GitBranch } from "lucide-react";
import { cn } from "@/lib/utils";
import { useQuery } from "@tanstack/react-query";
import { api, Conversation } from "@/lib/api";
import { formatDistanceToNow } from "date-fns";

interface ConversationListProps {
  currentConversation: string;
  onSelectConversation: (id: string) => void;
}

const ConversationList = ({ currentConversation, onSelectConversation }: ConversationListProps) => {
  const { data, isLoading, refetch } = useQuery({
    queryKey: ['conversations'],
    queryFn: () => api.getConversations(50, 0),
  });

  const conversations = data?.conversations || [];

  return (
    <div className="flex-1 flex flex-col">
      <div className="p-4">
        <div className="relative">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted-foreground" />
          <Input placeholder="Search conversations..." className="pl-9" />
        </div>
      </div>
      <ScrollArea className="flex-1">
        <div className="p-2 space-y-1">
          {isLoading ? (
            <div className="p-4 text-center text-muted-foreground">Loading conversations...</div>
          ) : conversations.length === 0 ? (
            <div className="p-4 text-center text-muted-foreground">
              No conversations yet. Start a new conversation to begin!
            </div>
          ) : (
            conversations.map((conv: Conversation) => {
              const lastMessage = conv.messages?.[conv.messages.length - 1];
              const timestamp = conv.updated_at
                ? formatDistanceToNow(new Date(conv.updated_at), { addSuffix: true })
                : '';

              return (
            <button
              key={conv.id}
                  onClick={() => onSelectConversation(conv.id.toString())}
              className={cn(
                "w-full text-left p-3 rounded-lg hover:bg-muted transition-colors",
                    currentConversation === conv.id.toString() && "bg-muted"
              )}
            >
              <div className="flex items-start gap-3">
                <MessageSquare className="w-5 h-5 text-primary mt-0.5 flex-shrink-0" />
                <div className="flex-1 min-w-0">
                  <div className="flex items-center gap-2 mb-1">
                        <h4 className="font-medium truncate">
                          {conv.title || `Conversation ${conv.id}`}
                        </h4>
                        {conv.parent_branch_id && (
                      <GitBranch className="w-4 h-4 text-muted-foreground flex-shrink-0" />
                    )}
                  </div>
                      {lastMessage && (
                  <p className="text-sm text-muted-foreground truncate">
                          {lastMessage.user_message}
                  </p>
                      )}
                      {timestamp && (
                        <p className="text-xs text-muted-foreground mt-1">{timestamp}</p>
                      )}
                </div>
              </div>
            </button>
              );
            })
          )}
        </div>
      </ScrollArea>
    </div>
  );
};

export default ConversationList;

import { useState, useEffect } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { ScrollArea } from "@/components/ui/scroll-area";
import { MessageSquare, Send, Mic, Plus, User, Settings, LogOut, BarChart3 } from "lucide-react";
import ConversationList from "@/components/ConversationList";
import MessageBubble from "@/components/MessageBubble";
import ResultsViewer from "@/components/ResultsViewer";
import VoiceInputModal from "@/components/VoiceInputModal";
import { useNavigate } from "react-router-dom";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { useToast } from "@/hooks/use-toast";
import { useAuth } from "@/contexts/AuthContext";
import { api, Message, Conversation } from "@/lib/api";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";

interface DisplayMessage {
  id: string;
  role: "user" | "assistant";
  content: string;
  timestamp: Date;
  analysis?: string | null;
  results?: {
    type?: string;
    data?: any;
    sql_query?: string | null;
    result_data?: string;
    result_format?: string | null;
    error_message?: string | null;
  };
}

const Dashboard = () => {
  const navigate = useNavigate();
  const { toast } = useToast();
  const { user, logout } = useAuth();
  const queryClient = useQueryClient();
  const [inputMessage, setInputMessage] = useState("");
  const [showVoiceInput, setShowVoiceInput] = useState(false);
  const [currentConversationId, setCurrentConversationId] = useState<number | null>(null);
  const [messages, setMessages] = useState<DisplayMessage[]>([]);

  // Fetch current conversation
  const { data: conversationData } = useQuery({
    queryKey: ['conversation', currentConversationId],
    queryFn: () => api.getConversation(currentConversationId!),
    enabled: !!currentConversationId,
  });

  // Update messages when conversation loads
  useEffect(() => {
    if (conversationData?.conversation?.messages) {
      const conversationId = conversationData.conversation.id;
      
      // Only update if this is the current conversation
      if (conversationId === currentConversationId) {
        const displayMessages: DisplayMessage[] = conversationData.conversation.messages.map((msg: Message) => ({
          id: msg.id.toString(),
          role: "user" as const,
          content: msg.user_message,
          timestamp: new Date(msg.created_at),
          analysis: msg.analysis || undefined,
          // Store results separately - will be attached to assistant message
          results: {
            type: msg.result_format || undefined,
            sql_query: msg.sql_query || undefined,
            result_data: msg.result_data || undefined,
            result_format: msg.result_format || undefined,
            error_message: msg.error_message || undefined,
          },
        }));

        // Add assistant responses for successful queries
        const messagesWithResponses: DisplayMessage[] = [];
        displayMessages.forEach((msg) => {
          // User message without results
          messagesWithResponses.push({
            id: msg.id,
            role: "user" as const,
            content: msg.content,
            timestamp: msg.timestamp,
            // Don't attach results to user messages
          });
          
          // Add assistant response with results if available
          if (msg.results && msg.results.result_format !== "error" && !msg.results.error_message) {
            // Build content: analysis first, then "Here are the results:"
            let assistantContent = "";
            if (msg.analysis) {
              assistantContent = msg.analysis;
              if (msg.results.result_data) {
                assistantContent += "\n\nHere are the results:";
              }
            } else {
              assistantContent = "Here are the results:";
            }
            
            messagesWithResponses.push({
              id: `${msg.id}-response`,
              role: "assistant" as const,
              content: assistantContent,
              timestamp: msg.timestamp,
              analysis: msg.analysis || undefined,
              results: msg.results, // Results only in assistant message
            });
          } else if (msg.results?.error_message) {
            messagesWithResponses.push({
              id: `${msg.id}-error`,
              role: "assistant" as const,
              content: `Error: ${msg.results.error_message}`,
              timestamp: msg.timestamp,
              results: msg.results,
            });
          } else if (msg.analysis) {
            // If there's analysis but no results (conversational response)
            messagesWithResponses.push({
              id: `${msg.id}-analysis`,
              role: "assistant" as const,
              content: msg.analysis,
              timestamp: msg.timestamp,
              analysis: msg.analysis,
            });
          }
        });

        setMessages(messagesWithResponses);
      }
    } else if (!currentConversationId) {
      // Welcome message when no conversation is selected
      setMessages([
        {
          id: "welcome",
          role: "assistant",
          content: "Hello! I'm your NLP-to-SQL assistant. Ask me anything about your data. Try queries like 'Total transactions for Silk Pay in Q1 2024' or 'Top 5 merchants by revenue in Kazakhstan'.",
          timestamp: new Date(),
        },
      ]);
    }
  }, [conversationData, currentConversationId]);

  // Query mutation
  const queryMutation = useMutation({
    mutationFn: ({ query, conversationId }: { query: string; conversationId: number | null }) =>
      api.executeQuery({
        query,
        conversation_id: conversationId || undefined,
      }),
    onSuccess: async (response, variables) => {
      // Refresh conversation to get updated messages (this will trigger useEffect to reload)
      const convId = variables.conversationId || currentConversationId;
      if (convId) {
        // Use refetch instead of invalidate to ensure we get fresh data
        await queryClient.refetchQueries({ queryKey: ['conversation', convId] });
      }
      await queryClient.invalidateQueries({ queryKey: ['conversations'] });
    },
    onError: (error) => {
      toast({
        title: "Query failed",
        description: error instanceof Error ? error.message : "Failed to execute query",
        variant: "destructive",
      });
    },
  });

  // Create conversation mutation
  const createConversationMutation = useMutation({
    mutationFn: (title: string) => api.createConversation(title),
    onSuccess: async (response) => {
      await queryClient.invalidateQueries({ queryKey: ['conversations'] });
      setCurrentConversationId(response.conversation.id);
      setMessages([]);
      toast({
        title: "New conversation created",
        description: "You can start asking questions now",
      });
    },
  });

  const handleSendMessage = async () => {
    if (!inputMessage.trim() || queryMutation.isPending) return;

    const queryText = inputMessage.trim();
    setInputMessage(""); // Clear input immediately for better UX

    // Create conversation if none exists
    if (!currentConversationId) {
      try {
        const convResponse = await createConversationMutation.mutateAsync(
          `Query: ${queryText.substring(0, 50)}`
        );
        const newConvId = convResponse.conversation.id;
        setCurrentConversationId(newConvId);
        // Wait a bit for conversation to be set, then send query
        setTimeout(() => {
          queryMutation.mutate({ query: queryText, conversationId: newConvId });
        }, 100);
      } catch (error) {
        toast({
          title: "Failed to create conversation",
          description: error instanceof Error ? error.message : "Unknown error",
          variant: "destructive",
        });
        setInputMessage(queryText); // Restore input on error
      }
    } else {
      queryMutation.mutate({ query: queryText, conversationId: currentConversationId });
    }
  };

  const handleNewConversation = () => {
    setCurrentConversationId(null);
    setMessages([
      {
        id: "welcome",
        role: "assistant",
        content: "Hello! I'm your NLP-to-SQL assistant. Ask me anything about your data.",
        timestamp: new Date(),
      },
    ]);
  };

  const handleSelectConversation = (id: string) => {
    setCurrentConversationId(parseInt(id));
  };

  const handleVoiceInput = (transcript: string) => {
    setInputMessage(transcript);
    setShowVoiceInput(false);
  };

  return (
    <div className="h-screen flex bg-background">
      {/* Sidebar */}
      <div className="w-80 border-r border-border flex flex-col bg-card">
        <div className="p-4 border-b border-border">
          <div className="flex items-center justify-between mb-4">
            <div className="flex items-center gap-2">
              <div className="w-8 h-8 bg-primary rounded-full flex items-center justify-center">
                <MessageSquare className="w-5 h-5 text-primary-foreground" />
              </div>
              <span className="font-semibold">NLP Analytics</span>
            </div>
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="ghost" size="icon">
                  <User className="w-5 h-5" />
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end">
                <DropdownMenuItem onClick={() => navigate("/profile")}>
                  <Settings className="w-4 h-4 mr-2" />
                  Settings
                </DropdownMenuItem>
                {user?.role?.name === "admin" && (
                <DropdownMenuItem onClick={() => navigate("/admin")}>
                  <BarChart3 className="w-4 h-4 mr-2" />
                  Admin Panel
                </DropdownMenuItem>
                )}
                <DropdownMenuSeparator />
                <DropdownMenuItem onClick={logout}>
                  <LogOut className="w-4 h-4 mr-2" />
                  Log Out
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
          <Button className="w-full" onClick={handleNewConversation}>
            <Plus className="w-4 h-4 mr-2" />
            New Conversation
          </Button>
        </div>
        <ConversationList
          currentConversation={currentConversationId?.toString() || ""}
          onSelectConversation={handleSelectConversation}
        />
      </div>

      {/* Main Chat Area */}
      <div className="flex-1 flex flex-col">
        <ScrollArea className="flex-1 p-6">
          <div className="max-w-4xl mx-auto space-y-6">
            {messages.map((message) => (
              <div key={message.id}>
                <MessageBubble message={message} />
                {message.results && (
                  <div className="mt-4">
                    <ResultsViewer results={message.results} />
                  </div>
                )}
              </div>
            ))}
            {queryMutation.isPending && (
              <div className="text-center text-muted-foreground">
                Analyzing your query...
              </div>
            )}
          </div>
        </ScrollArea>

        {/* Input Area */}
        <div className="border-t border-border p-4 bg-card">
          <div className="max-w-4xl mx-auto flex gap-2">
            <Button
              variant="outline"
              size="icon"
              onClick={() => setShowVoiceInput(true)}
            >
              <Mic className="w-5 h-5" />
            </Button>
            <Input
              placeholder="Ask about your data... (e.g., 'Total transactions for Silk Pay in Q1 2024')"
              value={inputMessage}
              onChange={(e) => setInputMessage(e.target.value)}
              onKeyPress={(e) => e.key === "Enter" && !queryMutation.isPending && handleSendMessage()}
              className="flex-1"
              disabled={queryMutation.isPending}
            />
            <Button onClick={handleSendMessage} disabled={queryMutation.isPending || !inputMessage.trim()}>
              <Send className="w-5 h-5" />
            </Button>
          </div>
        </div>
      </div>

      <VoiceInputModal
        open={showVoiceInput}
        onClose={() => setShowVoiceInput(false)}
        onSubmit={handleVoiceInput}
      />
    </div>
  );
};

export default Dashboard;

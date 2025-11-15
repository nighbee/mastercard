import { useState, useEffect } from "react";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Mic, MicOff } from "lucide-react";
import { useToast } from "@/hooks/use-toast";

interface VoiceInputModalProps {
  open: boolean;
  onClose: () => void;
  onSubmit: (transcript: string) => void;
}

const VoiceInputModal = ({ open, onClose, onSubmit }: VoiceInputModalProps) => {
  const [isListening, setIsListening] = useState(false);
  const [transcript, setTranscript] = useState("");
  const { toast } = useToast();

  useEffect(() => {
    if (open) {
      setIsListening(true);
      // Simulate voice input
      setTimeout(() => {
        setTranscript("Total transactions for Silk Pay in Q1 2024");
        setIsListening(false);
      }, 2000);
    } else {
      setTranscript("");
      setIsListening(false);
    }
  }, [open]);

  const handleSubmit = () => {
    if (transcript) {
      onSubmit(transcript);
      setTranscript("");
    }
  };

  return (
    <Dialog open={open} onOpenChange={onClose}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>Voice Input</DialogTitle>
        </DialogHeader>
        <div className="flex flex-col items-center gap-6 py-6">
          <div
            className={`w-24 h-24 rounded-full flex items-center justify-center transition-all ${
              isListening
                ? "bg-primary animate-pulse"
                : "bg-muted"
            }`}
          >
            {isListening ? (
              <Mic className="w-12 h-12 text-primary-foreground" />
            ) : (
              <MicOff className="w-12 h-12 text-muted-foreground" />
            )}
          </div>
          <div className="text-center">
            <p className="text-lg font-medium mb-2">
              {isListening ? "Listening..." : "Ready"}
            </p>
            {transcript && (
              <p className="text-muted-foreground max-w-sm">
                "{transcript}"
              </p>
            )}
          </div>
          <div className="flex gap-2">
            <Button variant="outline" onClick={onClose}>
              Cancel
            </Button>
            <Button onClick={handleSubmit} disabled={!transcript}>
              Submit
            </Button>
          </div>
        </div>
      </DialogContent>
    </Dialog>
  );
};

export default VoiceInputModal;

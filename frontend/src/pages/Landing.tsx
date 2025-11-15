import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { useNavigate } from "react-router-dom";
import { MessageSquare, Database, BarChart3, Shield, Zap, GitBranch } from "lucide-react";

const Landing = () => {
  const navigate = useNavigate();

  const features = [
    {
      icon: MessageSquare,
      title: "Natural Language Queries",
      description: "Ask questions in plain English and get instant SQL-powered answers"
    },
    {
      icon: Database,
      title: "Secure Database Access",
      description: "Enterprise-grade security with role-based access control"
    },
    {
      icon: BarChart3,
      title: "Visual Analytics",
      description: "Transform data into beautiful charts and tables automatically"
    },
    {
      icon: GitBranch,
      title: "Conversation Branching",
      description: "Explore different query paths and maintain conversation history"
    },
    {
      icon: Zap,
      title: "Real-time Processing",
      description: "Lightning-fast query translation and execution"
    },
    {
      icon: Shield,
      title: "Audit & Compliance",
      description: "Complete interaction logging for audit trails and improvements"
    }
  ];

  return (
    <div className="min-h-screen bg-gradient-to-br from-background via-background to-muted">
      <header className="border-b border-border bg-card/50 backdrop-blur-sm sticky top-0 z-50">
        <div className="container mx-auto px-6 py-4 flex items-center justify-between">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 bg-primary rounded-full flex items-center justify-center">
              <MessageSquare className="w-6 h-6 text-primary-foreground" />
            </div>
            <span className="text-xl font-bold">Mastercard NLP Analytics</span>
          </div>
          <div className="flex gap-3">
            <Button variant="ghost" onClick={() => navigate("/login")}>
              Log In
            </Button>
            <Button onClick={() => navigate("/register")}>
              Get Started
            </Button>
          </div>
        </div>
      </header>

      <main className="container mx-auto px-6">
        <section className="py-20 text-center">
          <h1 className="text-5xl md:text-6xl font-bold mb-6 bg-gradient-to-r from-primary to-primary/70 bg-clip-text text-transparent">
            Transform Data Queries into Insights
          </h1>
          <p className="text-xl text-muted-foreground mb-8 max-w-2xl mx-auto">
            Enterprise NLP-to-SQL chatbot platform. Ask questions naturally, get powerful analytics instantly.
          </p>
          <div className="flex gap-4 justify-center">
            <Button size="lg" onClick={() => navigate("/register")}>
              Start Free Trial
            </Button>
            <Button size="lg" variant="outline" onClick={() => navigate("/login")}>
              Sign In
            </Button>
          </div>
        </section>

        <section className="py-16">
          <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
            {features.map((feature) => (
              <Card key={feature.title} className="p-6 hover:shadow-lg transition-shadow border-border">
                <div className="w-12 h-12 bg-primary/10 rounded-lg flex items-center justify-center mb-4">
                  <feature.icon className="w-6 h-6 text-primary" />
                </div>
                <h3 className="text-lg font-semibold mb-2">{feature.title}</h3>
                <p className="text-muted-foreground">{feature.description}</p>
              </Card>
            ))}
          </div>
        </section>

        <section className="py-20 text-center">
          <h2 className="text-3xl font-bold mb-6">Ready to get started?</h2>
          <p className="text-muted-foreground mb-8 max-w-xl mx-auto">
            Join thousands of analysts using natural language to unlock data insights
          </p>
          <Button size="lg" onClick={() => navigate("/register")}>
            Create Account
          </Button>
        </section>
      </main>

      <footer className="border-t border-border py-8">
        <div className="container mx-auto px-6 text-center text-muted-foreground">
          <p>© 2024 Mastercard NLP Analytics. All rights reserved.</p>
        </div>
      </footer>
    </div>
  );
};

export default Landing;

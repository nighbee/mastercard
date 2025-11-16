const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1';

export interface ApiError {
  error: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  password: string;
  full_name: string;
}

export interface LoginResponse {
  message: string;
  access_token: string;
  refresh_token: string;
  user: {
    id: number;
    email: string;
    full_name: string;
    role_id: number | null;
    role?: {
      id: number;
      name: string;
    };
  };
}

export interface User {
  id: number;
  email: string;
  full_name: string;
  role_id: number | null;
  role?: {
    id: number;
    name: string;
  };
  is_active: boolean;
}

export interface Conversation {
  id: number;
  user_id: number;
  title: string | null;
  parent_branch_id: number | null;
  branch_point_message_id: number | null;
  created_at: string;
  updated_at: string;
  messages?: Message[];
}

export interface Message {
  id: number;
  conversation_id: number;
  user_message: string;
  sql_query: string | null;
  result_data: string | null;
  result_format: 'text' | 'table' | 'chart' | 'error' | null;
  error_message: string | null;
  analysis: string | null;
  execution_time_ms: number | null;
  created_at: string;
}

export interface QueryRequest {
  query: string;
  conversation_id?: number;
}

export interface QueryResponse {
  message: Message;
}

class ApiClient {
  private getAuthToken(): string | null {
    return localStorage.getItem('access_token');
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const token = this.getAuthToken();
    const headers: HeadersInit = {
      'Content-Type': 'application/json',
      ...options.headers,
    };

    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
    }

    const response = await fetch(`${API_BASE_URL}${endpoint}`, {
      ...options,
      headers,
    });

    if (!response.ok) {
      const error: ApiError = await response.json().catch(() => ({
        error: `HTTP ${response.status}: ${response.statusText}`,
      }));
      throw new Error(error.error || 'An error occurred');
    }

    return response.json();
  }

  // Auth endpoints
  async register(data: RegisterRequest): Promise<{ message: string; user: User }> {
    return this.request('/auth/register', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  async login(data: LoginRequest): Promise<LoginResponse> {
    return this.request('/auth/login', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  async refreshToken(refreshToken: string): Promise<{ access_token: string }> {
    return this.request('/auth/refresh', {
      method: 'POST',
      body: JSON.stringify({ refresh_token: refreshToken }),
    });
  }

  async getProfile(): Promise<{ user: User }> {
    return this.request('/auth/profile');
  }

  // Query endpoints
  async executeQuery(data: QueryRequest): Promise<QueryResponse> {
    return this.request('/query', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  // Conversation endpoints
  async getConversations(limit = 50, offset = 0): Promise<{
    conversations: Conversation[];
    total: number;
    limit: number;
    offset: number;
  }> {
    return this.request(`/conversations?limit=${limit}&offset=${offset}`);
  }

  async getConversation(id: number): Promise<{ conversation: Conversation }> {
    return this.request(`/conversations/${id}`);
  }

  async createConversation(title: string): Promise<{ conversation: Conversation }> {
    return this.request('/conversations', {
      method: 'POST',
      body: JSON.stringify({ title }),
    });
  }

  async updateConversation(id: number, title: string): Promise<{ conversation: Conversation }> {
    return this.request(`/conversations/${id}`, {
      method: 'PUT',
      body: JSON.stringify({ title }),
    });
  }

  async deleteConversation(id: number): Promise<{ message: string }> {
    return this.request(`/conversations/${id}`, {
      method: 'DELETE',
    });
  }

  async searchConversations(keyword: string, limit = 50, offset = 0): Promise<{
    conversations: Conversation[];
    total: number;
    limit: number;
    offset: number;
  }> {
    return this.request(`/conversations/search?q=${encodeURIComponent(keyword)}&limit=${limit}&offset=${offset}`);
  }

  async createBranch(
    parentId: number,
    branchPointMessageId: number,
    title: string
  ): Promise<{ conversation: Conversation }> {
    return this.request(`/conversations/${parentId}/branch`, {
      method: 'POST',
      body: JSON.stringify({
        title,
        branch_point_message_id: branchPointMessageId,
      }),
    });
  }
}

export const api = new ApiClient();


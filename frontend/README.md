# Mastercard NLP-to-SQL Frontend

React + TypeScript frontend for the Mastercard NLP-to-SQL Analytics Chatbot Platform.

## 🚀 Quick Start

### Prerequisites
- Node.js 18+ and npm
- Backend server running on `http://localhost:8080`

### Setup

1. **Install dependencies:**
   ```bash
   npm install
   ```

2. **Configure API URL (optional):**
   ```bash
   # Create .env file if needed
   echo "VITE_API_URL=http://localhost:8080/api/v1" > .env
   ```

3. **Start development server:**
   ```bash
npm run dev
```

The frontend will start on `http://localhost:5173` (or the port shown in terminal)

## 📁 Project Structure

```
frontend/
├── src/
│   ├── components/          # Reusable UI components
│   │   ├── ui/             # shadcn/ui components
│   │   ├── ConversationList.tsx
│   │   ├── MessageBubble.tsx
│   │   ├── ResultsViewer.tsx
│   │   └── VoiceInputModal.tsx
│   ├── contexts/           # React contexts
│   │   └── AuthContext.tsx # Authentication context
│   ├── hooks/              # Custom React hooks
│   ├── lib/                # Utilities and API client
│   │   ├── api.ts         # API client
│   │   └── utils.ts       # Utility functions
│   ├── pages/             # Page components
│   │   ├── Dashboard.tsx
│   │   ├── Login.tsx
│   │   ├── Register.tsx
│   │   └── ...
│   ├── App.tsx            # Main app component
│   └── main.tsx           # Entry point
└── package.json
```

## 🔧 Features

- **Authentication**: Login, register, and JWT token management
- **Query Interface**: Natural language query input with voice support
- **Conversation Management**: Create, view, and manage conversations
- **Results Display**: Table and chart visualization of query results
- **Real-time Updates**: React Query for data fetching and caching
- **Protected Routes**: Route protection based on authentication status

## 🔌 API Integration

The frontend uses the API client in `src/lib/api.ts` to communicate with the backend:

- **Authentication**: `/api/v1/auth/*`
- **Queries**: `/api/v1/query`
- **Conversations**: `/api/v1/conversations/*`

All API requests include JWT tokens from localStorage automatically.

## 🎨 UI Components

Built with:
- **React 18** - UI framework
- **TypeScript** - Type safety
- **Tailwind CSS** - Styling
- **shadcn/ui** - UI component library
- **Recharts** - Data visualization
- **React Query** - Data fetching and caching
- **React Router** - Navigation

## 📝 Environment Variables

Create a `.env` file in the frontend directory:

```env
VITE_API_URL=http://localhost:8080/api/v1
```

## 🛠️ Development

### Build for production:
```bash
npm run build
```

### Preview production build:
```bash
npm run preview
```

### Lint:
```bash
npm run lint
```

## 🔍 Troubleshooting

1. **API connection errors**: 
   - Check that backend is running on `http://localhost:8080`
   - Verify `VITE_API_URL` in `.env` matches backend URL
   - Check browser console for CORS errors

2. **Authentication issues**:
   - Clear localStorage: `localStorage.clear()`
   - Check that tokens are being stored correctly
   - Verify backend JWT configuration

3. **Query execution fails**:
   - Check browser console for error messages
   - Verify Gemini API key is set in backend
   - Check network tab for API request/response

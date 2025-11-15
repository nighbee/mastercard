User Flow Document: Mastercard NLP-to-SQL Chatbot Platform
1. Entry Points
Welcome/Landing Page

    Branding/logo and tagline

    Options: Log In, Register, Learn More, (optional: Demo or FAQ)

    Redirects unauthenticated users to Login/Register

2. Authentication
Login Page

    Fields: Email, Password

    Options: “Forgot password?”, “Register”

    Buttons: Log In

    Social Login (optional)

    Validation messages for errors

Register Page

    Fields: Full Name, Email, Password, Confirm Password, (optional: Organization/role)

    Password strength meter

    Accept Terms checkbox

    Button: Register

    “Already have an account? Log In”

Password Reset Flow

    Forgot Password: enter email, get reset link

    Password Reset: enter new password, confirm

3. Main App (Post-Login)
Layout (consistent throughout):

    Top Bar: Logo, profile/settings dropdown, log out button

    Left Sidebar: Conversations list (with search), New Conversation button, history/branch tree navigation

A. Chat Interface (Core Window)

    Message input (text and voice mic icon)

    Display previous chat (with timestamps, sender/AI distinction)

    Message content: User query, system chat, tables, charts, text answers

    Action buttons: Export (CSV/Excel), Copy, View SQL (optional)

    Loading/status indicators (“Analyzing your query…”)

    Example prompts suggestions when chat is empty

    Error messages and fallback guidance

B. Conversation Management

    Conversation List: Shows all previous sessions (search/filter by date/title)

    Conversation Details: Shows messages, results, visualizations

    Create Branch: Fork previous message to explore alternative path

    Rename/Delete Conversation

    View Conversation Tree: graphical branch navigation

C. Voice Input Modal

    Opens when mic icon is clicked

    Shows waveform, “Listening…” indicator

    Text preview of recognized speech

    Buttons: Submit/edit, Cancel

D. Query Results Windows

    Table Results: Paginated, sortable, copy/export features

    Chart Results: Bar/Line/Pie/Area chart as appropriate, interactive, with legend

    Text Summary: For direct answers or stats

    SQL View: Show/hide the raw generated SQL for transparency

E. User Settings/Profile

    Edit profile: name, email, password

    (Optional) Set preferences: language, default export format

    View API usage/statistics

    Log out

F. Admin Panel (Role: Admin only)

    User management: List, search, assign roles

    Audit logs: View/filter/export security logs

    System metrics: Usage stats, active users, query volume

4. Special/Support Windows

    Error Pages: 404, 500, network/server issues

    System notifications (above chat or as toast messages): “Query failed,” “Access denied,” etc.

    Onboarding Tooltips: First-time user guidance overlays

5. State/Navigation Flows

    After login: redirect to last active conversation or new chat window

    Switching between conversations/branches: auto-load relevant history

    Logout: clear session, redirect to login

    Inactivity timeout: show session expiry warning, auto-log out

6. Chat Branch/History Flow

    At any previous message: “Branch here” option

    Creates new conversation branch, visible in sidebar tree

    Maintain context and results up to that branch point

For a visual diagram, the main windows (“nodes”) would be:

    Welcome / Landing

    Login

    Register

    Password Reset (2-step)

    Main Chat

    Conversation List

    Conversation Details/Branches

    Voice Input Modal

    Result Viewer (Table/Chart/Text/SQL)

    Profile / Settings

    Admin Dashboard (if admin)

    Error / Notification Windows

Each screen is directly accessible through logical navigation from the previous (except Error/Notifications which are modal or overlay).
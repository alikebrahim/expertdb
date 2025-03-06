import { Routes, Route, Navigate } from "react-router-dom";
import ProtectedRoute from "./components/ProtectedRoute";
import Login from "./pages/Login";
import { useAuth } from "./context/AuthContext";

function App() {
  const { token } = useAuth();

  return (
    <Routes>
      <Route path="/login" element={<Login />} />
      <Route path="/search" element={<ProtectedRoute><div>Search Page</div></ProtectedRoute>} />
      <Route path="/requests" element={<ProtectedRoute><div>Requests Page</div></ProtectedRoute>} />
      <Route path="/statistics" element={<ProtectedRoute><div>Statistics Page</div></ProtectedRoute>} />
      <Route path="/admin" element={<ProtectedRoute requireAdmin><div>Admin Panel</div></ProtectedRoute>} />
      <Route path="/" element={token ? <Navigate to="/search" replace /> : <Navigate to="/login" replace />} />
    </Routes>
  );
}

export default App;

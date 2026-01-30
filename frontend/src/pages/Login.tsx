import { useEffect } from "react";
import { useNavigate } from "react-router";
import { useAuth } from "../hooks/useAuth";
import { API_BASE_URL } from "../libs/api";

export default function Login() {
  const { user, loading } = useAuth();
  const navigate = useNavigate();

  useEffect(() => {
    if (!loading && user) {
      navigate("/dashboard");
    }
  }, [user, loading, navigate]);

  const handleLogin = () => {
    window.location.href = `${API_BASE_URL}/login`;
  };

  if (loading) return null;

  return <button onClick={handleLogin}>Login with Google</button>;
}

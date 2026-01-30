import { useEffect, useState } from "react";
import { useNavigate } from "react-router";
import { useAuth } from "../hooks/useAuth";
import { API_BASE_URL } from "../libs/api";

export default function Dashboard() {
  const { user, loading } = useAuth();
  const navigate = useNavigate();
  const [status, setStatus] = useState<string | null>(null);

  useEffect(() => {
    if (!loading && !user) {
      navigate("/");
    }
  }, [user, loading, navigate]);

  useEffect(() => {
    console.log(status);
  }, [status])

  const handleStart = async () => {
    setStatus(null);
    try {
      const res = await fetch(`${API_BASE_URL}/api/minecraft/start`, {
        method: "POST",
        credentials: "include",
      });
      console.log(res);
      const text = await res.text();
      const data = text ? JSON.parse(text) : {};
      setStatus(res.ok ? data.message ?? "Server started" : data.error ?? "Failed to start server");
    } catch (error) {
      console.log(error);
      setStatus("Request failed with error: " + error);
    }
  };

  const handleStop = async () => {
    setStatus(null);
    try {
      const res = await fetch(`${API_BASE_URL}/api/minecraft/stop`, {
        method: "POST",
        credentials: "include",
      });
      console.log(res);
      const text = await res.text();
      const data = text ? JSON.parse(text) : {};
      setStatus(res.ok ? data.message ?? "Server stopped" : data.error ?? "Failed to stop server");
    } catch (error) {
      console.log(error);
      setStatus("Request failed with error: " + error);
    }
  };

  if (loading) return null;

  return (
    <div>
      <button onClick={handleStart}>Start Server</button>
      <button onClick={handleStop}>Stop Server</button>
      {status && <p>{status}</p>}
    </div>
  );
}

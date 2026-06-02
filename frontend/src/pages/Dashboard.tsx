import { useState } from "react";
import { useAuth } from "../hooks/useAuth";
import { API_BASE_URL } from "../libs/api";

export default function Dashboard() {
  const { user, loading, refetch } = useAuth();
  const [status, setStatus] = useState<string | null>(null);

  const handleStart = async () => {
    setStatus(null);
    try {
      const res = await fetch(`${API_BASE_URL}/api/minecraft/start`, {
        method: "POST",
        credentials: "include",
      });
      const text = await res.text();
      const data = text ? JSON.parse(text) : {};
      setStatus(res.ok ? data.message ?? "Server started" : data.error ?? "Failed to start server");
    } catch (error) {
      setStatus("Request failed with error: " + (error as Error).message);
    }
  };

  const handleStop = async () => {
    setStatus(null);
    try {
      const res = await fetch(`${API_BASE_URL}/api/minecraft/stop`, {
        method: "POST",
        credentials: "include",
      });
      const data = await res.json();
      setStatus(res.ok ? data.message ?? "Server stopped" : data.error ?? "Failed to stop server");
    } catch (error) {
      setStatus("Request failed with error: " + (error as Error).message);
    }
  };

  const handleLogin = () => {
    window.location.href = `${API_BASE_URL}/login`;
  };

  const handleLogout = async () => {
    await fetch(`${API_BASE_URL}/logout`, { credentials: "include" });
    refetch();
  };

  if (loading) return null;

  return (
    <div style={{ padding: "16px" }}>
      <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center" }}>
        <div style={{ display: "flex", gap: "8px", alignItems: "center" }}>
          <button onClick={handleStart} disabled={!user}>Start Server</button>
          <button onClick={handleStop} disabled={!user}>Stop Server</button>
        </div>
        {user ? (
          <div style={{ display: "flex", alignItems: "center", gap: "8px" }}>
            <span>{user.name}</span>
            <button onClick={handleLogout}>Logout</button>
          </div>
        ) : (
          <button onClick={handleLogin}>Login with Google</button>
        )}
      </div>
      {status && <p style={{ marginTop: "16px" }}>{status}</p>}
    </div>
  );
}

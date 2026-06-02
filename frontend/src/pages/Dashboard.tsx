import { useState } from "react";
import { useAuth } from "../hooks/useAuth";
import { API_BASE_URL } from "../libs/api";
import { InvGrid, McButton } from "../components/minecraft-ui";

// TODO: Temporary placeholder data. Will be replaced with live inventory
// fetched from the backend (typed at the API boundary), at which point these
// hardcoded item IDs go away.
const DEMO_ITEMS: (string | null)[] = [
  "Diamond_Helmet",    "Diamond_Chestplate", "Diamond_Leggings",
  "Diamond_Boots",     "Shield",             "Bow",
  "Arrow",             "Golden_Apple",       "Ender_Pearl",
  "Diamond_Sword",     "Diamond_Pickaxe",    "Diamond_Axe",
  "Torch",             "Cobblestone",         "Oak_Planks",
  "Bread",             "Cooked_Beef",              "Water_Bucket",
  "Iron_Ingot",        "Gold_Ingot",         "Redstone",
  "Lapis_Lazuli",      "Coal",               "Emerald",
  "String",            "Bone",               "Gunpowder",
];

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
      const data = await res.json();
      setStatus(res.ok ? data.message ?? "Server started" : data.error ?? "Failed to start server");
    } catch (error) {
      setStatus("Request failed with error: " + String(error));
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
      setStatus("Request failed with error: " + String(error));
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
    <div className="mc-page">
      <h1 className="mc-title">Kraft Bois</h1>

      <div style={{ position: "absolute", top: 12, right: 12, display: "flex", alignItems: "center", gap: 8, zIndex: 1 }}>
        {user ? (
          <>
            <span className="mc-status" style={{ fontSize: 14 }}>{user.name}</span>
            <McButton onClick={handleLogout}>Logout</McButton>
          </>
        ) : (
          <McButton onClick={handleLogin}>Login with Google</McButton>
        )}
      </div>

      <div style={{ display: "flex", gap: 12 }}>
        <McButton onClick={handleStart} disabled={!user}>Start Server</McButton>
        <McButton onClick={handleStop} disabled={!user}>Stop Server</McButton>
      </div>

      {status && <p className="mc-status">{status}</p>}

      <div className="mc-panel">
        <span className="mc-panel-title">Inventory</span>
        <InvGrid rows={3} cols={9} items={DEMO_ITEMS} />
      </div>
    </div>
  );
}

import { useState } from "react";
import { useAuth } from "../hooks/useAuth";
import { useOnlinePlayers } from "../hooks/useOnlinePlayers";
import { apiFetch, apiUrl } from "../libs/api";
import { InvGrid, McButton, PlayerList } from "../components/minecraft-ui";

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
  const { players, loading: playersLoading, error: playersError } = useOnlinePlayers();

  const handleStart = async () => {
    setStatus(null);
    try {
      const data = await apiFetch<{ message?: string }>("/api/minecraft/start", { method: "POST" });
      setStatus(data.message ?? "Server started");
    } catch (error) {
      setStatus(error instanceof Error ? error.message : "Failed to start server");
    }
  };

  const handleStop = async () => {
    setStatus(null);
    try {
      const data = await apiFetch<{ message?: string }>("/api/minecraft/stop", { method: "POST" });
      setStatus(data.message ?? "Server stopped");
    } catch (error) {
      setStatus(error instanceof Error ? error.message : "Failed to stop server");
    }
  };

  const handleLogin = () => {
    window.location.href = apiUrl("/login");
  };

  const handleLogout = async () => {
    try {
      await apiFetch("/logout");
    } finally {
      refetch();
    }
  };

  if (loading) return null;

  return (
    <div className="mc-page">
      <h1 className="mc-title">Kraft Bois</h1>

      <div className="mc-topbar">
        <div className="mc-topbar-user">
          {user ? (
            <>
              <span className="mc-status mc-username">{user.name}</span>
              <McButton onClick={handleLogout}>Logout</McButton>
            </>
          ) : (
            <McButton onClick={handleLogin}>Login with Google</McButton>
          )}
        </div>
        <PlayerList players={players} loading={playersLoading} error={playersError} />
      </div>

      <div className="mc-actions">
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

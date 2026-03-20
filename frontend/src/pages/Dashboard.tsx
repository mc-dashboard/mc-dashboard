import { useEffect, useState } from "react";
import { useNavigate } from "react-router";
import { useAuth } from "../hooks/useAuth";
import { API_BASE_URL } from "../libs/api";
import { InvGrid, McButton } from "../components/minecraft-ui";

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
  const { user, loading } = useAuth();
  const navigate = useNavigate();
  const [status, setStatus] = useState<string | null>(null);

  useEffect(() => {
    if (!loading && !user) {
      navigate("/");
    }
  }, [user, loading, navigate]);

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
      const text = await res.text();
      const data = text ? JSON.parse(text) : {};
      setStatus(res.ok ? data.message ?? "Server stopped" : data.error ?? "Failed to stop server");
    } catch (error) {
      setStatus("Request failed with error: " + error);
    }
  };

  if (loading) return null;

  return (
    <div className="mc-page">
      <h1 className="mc-title">Kraft Bois</h1>

      <div style={{ display: "flex", gap: 12 }}>
        <McButton onClick={handleStart}>Start Server</McButton>
        <McButton onClick={handleStop}>Stop Server</McButton>
      </div>

      {status && <p className="mc-status">{status}</p>}

      <div className="mc-panel">
        <span className="mc-panel-title">Inventory</span>
        <InvGrid rows={3} cols={9} items={DEMO_ITEMS} />
      </div>
    </div>
  );
}

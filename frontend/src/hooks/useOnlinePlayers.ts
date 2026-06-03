import { useCallback, useEffect, useState } from "react";
import { apiFetch } from "../libs/api";

export interface PlayerInfo {
  name: string;
  uuid: string;
}

interface PlayerListResponse {
  players: PlayerInfo[];
  count: number;
}

const POLL_INTERVAL_MS = 10_000;

/**
 * Polls the backend for the list of currently connected players. Returns the
 * latest list along with loading/error state so the UI can render gracefully
 * while the server is offline or unreachable.
 */
export function useOnlinePlayers() {
  const [players, setPlayers] = useState<PlayerInfo[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchPlayers = useCallback(async () => {
    try {
      const data = await apiFetch<PlayerListResponse>("/api/minecraft/players");
      setPlayers(data.players ?? []);
      setError(null);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to load players");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchPlayers();
    const id = setInterval(fetchPlayers, POLL_INTERVAL_MS);
    return () => clearInterval(id);
  }, [fetchPlayers]);

  return { players, loading, error };
}

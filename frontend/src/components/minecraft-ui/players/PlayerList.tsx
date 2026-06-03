import "../minecraft-ui.css";
import type { PlayerInfo } from "../../../hooks/useOnlinePlayers";

export interface PlayerListProps {
  players: PlayerInfo[];
  loading: boolean;
  error: string | null;
}

export default function PlayerList({ players, loading, error }: PlayerListProps) {
  return (
    <div className="mc-panel">
      <span className="mc-panel-title">
        Players Online ({loading ? "…" : players.length})
      </span>

      {error ? (
        <span className="mc-player-empty">{error}</span>
      ) : players.length === 0 ? (
        <span className="mc-player-empty">
          {loading ? "Loading…" : "No players online"}
        </span>
      ) : (
        <ul className="mc-player-list">
          {players.map((player) => (
            <li key={player.uuid || player.name} className="mc-player-row">
              <img
                className="mc-player-head"
                src={`https://mc-heads.net/avatar/${encodeURIComponent(player.name)}/24`}
                alt=""
                aria-hidden="true"
              />
              <span className="mc-player-name">{player.name}</span>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}

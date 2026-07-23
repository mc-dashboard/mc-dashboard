import InvSlot from "./InvSlot";
import "../minecraft-ui.css";

export interface InvGridProps {
  rows: number;
  cols: number;
  items?: (string | null)[];
  className?: string;
}

export default function InvGrid({ rows, cols, items, className }: InvGridProps) {
  return (
    <div className={`mcui ${className ?? ""}`}>
      {Array.from({ length: rows }, (_, row) => (
        <div key={row} className="mcui-row">
          {Array.from({ length: cols }, (_, col) => {
            const idx = row * cols + col;
            return <InvSlot key={idx} item={items?.[idx] ?? null} />;
          })}
        </div>
      ))}
    </div>
  );
}

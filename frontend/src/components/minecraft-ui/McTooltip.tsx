import { useEffect, useState } from "react";
import { createPortal } from "react-dom";
import "./minecraft-ui.css";

interface TooltipState {
  text: string;
  x: number;
  y: number;
}

type SetTooltip = React.Dispatch<React.SetStateAction<TooltipState | null>>;

let setGlobalTooltip: SetTooltip | null = null;

export function showTooltip(text: string, x: number, y: number) {
  setGlobalTooltip?.({ text, x, y });
}

export function moveTooltip(x: number, y: number) {
  setGlobalTooltip?.((prev) => (prev ? { ...prev, x, y } : null));
}

export function hideTooltip() {
  setGlobalTooltip?.(null);
}

export default function McTooltipPortal() {
  const [tooltip, setTooltip] = useState<TooltipState | null>(null);

  useEffect(() => {
    setGlobalTooltip = setTooltip;
    return () => {
      setGlobalTooltip = null;
    };
  }, []);

  if (!tooltip) return null;

  return createPortal(
    <div
      className="mc-tooltip"
      style={{
        left: tooltip.x + 12,
        top: tooltip.y - 12,
        transform: "translateY(-100%)",
      }}
    >
      {tooltip.text}
    </div>,
    document.body
  );
}

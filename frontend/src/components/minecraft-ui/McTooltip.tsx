import { useEffect, useState } from "react";
import { createPortal } from "react-dom";
import { bindTooltipSetter, type TooltipState } from "./McTooltipController";
import "./minecraft-ui.css";

export default function McTooltipPortal() {
  const [tooltip, setTooltip] = useState<TooltipState | null>(null);

  useEffect(() => {
    bindTooltipSetter(setTooltip);
    return () => {
      bindTooltipSetter(null);
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

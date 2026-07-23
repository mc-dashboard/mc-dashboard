import { useEffect, useState } from "react";
import { createPortal } from "react-dom";
import { tooltipController, type TooltipState } from "./McTooltipController";
import "../minecraft-ui.css";

export default function McTooltipPortal() {
  const [tooltip, setTooltip] = useState<TooltipState | null>(null);

  useEffect(() => {
    tooltipController.bind(setTooltip);
    return () => {
      tooltipController.bind(null);
    };
  }, []);

  if (!tooltip) return null;

  return createPortal(
    <div
      className="mc-tooltip"
      style={{
        left: tooltip.x + 12,
        top: tooltip.y - 12,
      }}
    >
      {tooltip.text}
    </div>,
    document.body
  );
}

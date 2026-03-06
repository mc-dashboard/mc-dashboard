import { showTooltip, moveTooltip, hideTooltip } from "./McTooltipController";
import "./minecraft-ui.css";

export interface InvSlotProps {
  item?: string | null;
  className?: string;
}

function itemDisplayName(item: string): string {
  return item.replace(/_/g, " ");
}

function itemImageUrl(item: string): string {
  return `https://minecraft.wiki/images/Invicon_${item}.png`;
}

export default function InvSlot({ item, className }: InvSlotProps) {
  return (
    <span className={`invslot ${className ?? ""}`}>
      {item && (
        <span
          className="invslot-item"
          onMouseEnter={(e) =>
            showTooltip(itemDisplayName(item), e.clientX, e.clientY)
          }
          onMouseMove={(e) => moveTooltip(e.clientX, e.clientY)}
          onMouseLeave={hideTooltip}
        >
          <img
            src={itemImageUrl(item)}
            alt={itemDisplayName(item)}
            className="invslot-item-image"
            draggable={false}
          />
        </span>
      )}
    </span>
  );
}

import { tooltipController } from "./McTooltipController";
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
            tooltipController.show(itemDisplayName(item), e.clientX, e.clientY)
          }
          onMouseMove={(e) => tooltipController.move(e.clientX, e.clientY)}
          onMouseLeave={tooltipController.hide}
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

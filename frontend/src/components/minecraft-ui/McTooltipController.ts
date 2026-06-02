import type { Dispatch, SetStateAction } from "react";

export interface TooltipState {
  text: string;
  x: number;
  y: number;
}

type SetTooltip = Dispatch<SetStateAction<TooltipState | null>>;

class TooltipController {
  private setter: SetTooltip | null = null;

  // Arrow fields so methods stay bound when passed as event handlers
  // (e.g. onMouseLeave={tooltipController.hide}).
  bind = (setter: SetTooltip | null) => {
    this.setter = setter;
  };

  show = (text: string, x: number, y: number) => {
    this.setter?.({ text, x, y });
  };

  move = (x: number, y: number) => {
    this.setter?.((prev) => (prev ? { ...prev, x, y } : null));
  };

  hide = () => {
    this.setter?.(null);
  };
}

export const tooltipController = new TooltipController();

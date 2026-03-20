import type { Dispatch, SetStateAction } from "react";

export interface TooltipState {
  text: string;
  x: number;
  y: number;
}

type SetTooltip = Dispatch<SetStateAction<TooltipState | null>>;

let setGlobalTooltip: SetTooltip | null = null;

export function bindTooltipSetter(setter: SetTooltip | null) {
  setGlobalTooltip = setter;
}

export function showTooltip(text: string, x: number, y: number) {
  setGlobalTooltip?.({ text, x, y });
}

export function moveTooltip(x: number, y: number) {
  setGlobalTooltip?.((prev) => (prev ? { ...prev, x, y } : null));
}

export function hideTooltip() {
  setGlobalTooltip?.(null);
}

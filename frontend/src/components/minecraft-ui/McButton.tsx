import "./minecraft-ui.css";

interface McButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  children: React.ReactNode;
}

export default function McButton({ children, className, ...props }: McButtonProps) {
  return (
    <button className={`mc-btn ${className ?? ""}`} {...props}>
      {children}
    </button>
  );
}

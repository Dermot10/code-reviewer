import { Play, Sun, Moon, Settings, LogOut } from "lucide-react";
import { ReviewStatus } from "../types"; 
type Props = {
  theme: "vs-dark" | "light";
  setTheme: (theme: "vs-dark" | "light") => void;
  handleReview: () => void;
  reviewStatus: ReviewStatus;
  handleLogout: () => void;
};

export default function TopBar({ theme, setTheme, handleReview, reviewStatus, handleLogout }: Props) {
  return (
    <header className="ide-header">
      <div className="header-left">
        <h1 className="logo">Code Reviewer</h1>
      </div>

      <div className="header-center">
        <button
          onClick={handleReview}
          disabled={reviewStatus === "processing"}
          className="btn-primary"
        >
          <Play size={16} />
          {reviewStatus === "processing" ? "Processing..." : "Review Code"}
        </button>
      </div>

      <div className="header-right">
        <button
          onClick={() => setTheme(theme === "vs-dark" ? "light" : "vs-dark")}
          className="btn-icon"
          title="Toggle theme"
        >
          {theme === "vs-dark" ? <Sun size={18} /> : <Moon size={18} />}
        </button>

        <button className="btn-icon" title="Settings">
          <Settings size={18} />
        </button>

        <button onClick={handleLogout} className="btn-icon" title="Logout">
          <LogOut size={18} />
        </button>
      </div>
    </header>
  );
}

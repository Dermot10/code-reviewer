import { AlertCircle, CheckCircle } from "lucide-react";
import { ReviewStatus } from "../types";

type Props = {
  reviewStatus: ReviewStatus;
  reviewResult: any;
  isCollapsed: boolean;
  setIsCollapsed: (val: boolean) => void;
  panelHeight: number;
  setPanelHeight: (val: number) => void;
  reviewId: number | null;
};

export default function ReviewPanel({ reviewStatus, reviewResult, isCollapsed, setIsCollapsed, panelHeight, reviewId }: Props) {
  return (
    <div
      className={`review-panel ${isCollapsed ? "collapsed" : "expanded"}`}
      style={{ height: isCollapsed ? 30 : panelHeight }}
    >
      <button onClick={() => setIsCollapsed(!isCollapsed)} className="btn-icon">
        {isCollapsed ? "▲" : "▼"}
      </button>

      {!isCollapsed && (
        <div className="review-content-wrapper">
          {reviewStatus === "processing" && (
            <div className="status-message">
              <div className="spinner" />
              <span>Analysing code (Review #{reviewId})...</span>
            </div>
          )}

          {reviewStatus === "completed" && reviewResult && (
            <div className="review-body">
              <CheckCircle size={20} className="success-icon" />
              <pre>{JSON.stringify(reviewResult, null, 2)}</pre>
            </div>
          )}

          {reviewStatus === "failed" && (
            <div className="status-message error">
              <AlertCircle size={20} />
              <span>Review failed. Please try again.</span>
            </div>
          )}
        </div>
      )}
    </div>
  );
}

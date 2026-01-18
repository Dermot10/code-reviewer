// app/(protected)/editor/page.tsx
"use client";

import { useState, useEffect } from "react";
import { Editor } from "@monaco-editor/react";

export default function EditorPage() {
  const [code, setCode] = useState("");
  const [reviewId, setReviewId] = useState<number | null>(null);
  const [reviewStatus, setReviewStatus] = useState<"idle" | "processing" | "completed" | "failed">("idle");
  const [reviewResult, setReviewResult] = useState<any>(null);

  // Submit code for review
  async function handleReview() {
    const token = localStorage.getItem("token");
    if (!token) return;

    try {
      setReviewStatus("processing");
      console.log("Submitting code:", code, "to URL:", `${process.env.NEXT_PUBLIC_API_URL}/review-code`);

      const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/review-code`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          "Authorization": `Bearer ${token}`,
        },
        body: JSON.stringify({ code }),
      });

      if (res.status === 429) {
        alert("Rate limit exceeded. Please wait.");
        setReviewStatus("idle");
        return;
      }

      if (!res.ok) throw new Error("Review submission failed");

      const data = await res.json();
      setReviewId(data.review_id);
      
      // Start polling for results
      startPolling(data.review_id, token);
    } catch (err) {
      console.error(err);
      setReviewStatus("failed");
    }
  }

  // Poll for review results
  function startPolling(id: number, token: string) {
    const interval = setInterval(async () => {
      try {
        const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/reviews/${id}`, {
          headers: { "Authorization": `Bearer ${token}` },
        });

        if (!res.ok) throw new Error("Failed to fetch review");

        const data = await res.json();

        if (data.status === "completed") {
          clearInterval(interval);
          setReviewResult(data);
          setReviewStatus("completed");
        } else if (data.status === "failed") {
          clearInterval(interval);
          setReviewStatus("failed");
        }
      } catch (err) {
        console.error(err);
        clearInterval(interval);
        setReviewStatus("failed");
      }
    }, 2000); // Poll every 2 seconds
  }

  return (
    <div className="ide-container">
      {/* Header */}
      <header className="ide-header">
        <div className="header-left">
          <h1>Code Reviewer IDE</h1>
        </div>
        <div className="header-right">
          <button 
            onClick={handleReview} 
            disabled={reviewStatus === "processing" || !code}
            className="review-button"
          >
            {reviewStatus === "processing" ? "Processing..." : "Review Code"}
          </button>
          <button onClick={() => {
            localStorage.removeItem("token");
            window.location.href = "/login";
          }}>
            Logout
          </button>
        </div>
      </header>

      {/* Main editor area */}
      <div className="ide-main">
        {/* File tree sidebar (placeholder for now) */}
        <aside className="file-tree">
          <h3>Files</h3>
          <div className="file-item active">📄 main.py</div>
          <button className="new-file-btn">+ New File</button>
        </aside>

        {/* Editor */}
        <div className="editor-area">
          <Editor
            height="100%"
            defaultLanguage="python"
            theme="vs-dark"
            value={code}
            onChange={(value) => setCode(value || "")}
            options={{
              minimap: { enabled: false },
              fontSize: 14,
            }}
          /> 
        </div>
        

      </div>

      {/* Review panel (slides up when results ready) */}
      {reviewStatus !== "idle" && (
        <div className={`review-panel ${reviewStatus === "completed" ? "expanded" : ""}`}>
          {reviewStatus === "processing" && (
            <div className="review-loading">
              <span className="spinner">⏳</span>
              Processing review... (ID: {reviewId})
            </div>
          )}

          {reviewStatus === "completed" && reviewResult && (
            <div className="review-results">
              <div className="review-header">
                <h3>Review Complete</h3>
                <button onClick={() => setReviewStatus("idle")}>✕ Close</button>
              </div>
              <div className="review-content">
                <pre>{JSON.stringify(reviewResult, null, 2)}</pre>
              </div>
            </div>
          )}

          {reviewStatus === "failed" && (
            <div className="review-error">
              Review failed. Please try again.
              <button onClick={() => setReviewStatus("idle")}>Close</button>
            </div>
          )}
        </div>
      )}
    </div>
  );
}
// app/(protected)/editor/page.tsx
"use client";

import { useState } from "react";
import TopBar from "./components/topbar";
import Sidebar from "./components/sidebar";
import EditorContainer from "./components/editorContainer";
import ReviewPanel from "./components/reviewPanel";
import Resizer from "./components/resizer";
import { EditorFile, ReviewStatus } from "./types";

export default function EditorPage() {

  const [theme, setTheme] = useState<"vs-dark" | "light">("vs-dark");

  const [files, setFiles] = useState<EditorFile[]>([
    { id: "1", name: "main.py", content: "# Write your Python code here\nprint('Hello, world')" },
    { id: "2", name: "utils.py", content: "def helper():\n    pass" },
  ]);
  const [activeFileId, setActiveFileId] = useState("1");
  const activeFile = files.find(f => f.id === activeFileId);

  const [reviewStatus, setReviewStatus] = useState<ReviewStatus>("idle");
  const [reviewResult, setReviewResult] = useState<any>(null);
  const [reviewId, setReviewId] = useState<number | null>(null);

  const [isReviewCollapsed, setIsReviewCollapsed] = useState(false);
  const [panelHeight, setPanelHeight] = useState(200);

  
  const handleNewFile = () => {
    const id = crypto.randomUUID();
    const name = `file${files.length + 1}.py`;
    setFiles(prev => [...prev, { id, name, content: "" }]);
    setActiveFileId(id);
  };

  const handleLogout = () => {
    localStorage.removeItem("token");
    window.location.href = "/login";
  };

  const handleReview = async () => {
    const token = localStorage.getItem("token");
    if (!token) return;

    try {
      setReviewStatus("processing");
      const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/review-code`, {
        method: "POST",
        headers: { "Content-Type": "application/json", Authorization: `Bearer ${token}` },
        body: JSON.stringify({ code: activeFile?.content }),
      });

      if (res.status === 429) {
        alert("Rate limit exceeded. Try again in an hour.");
        setReviewStatus("idle");
        return;
      }

      if (!res.ok) throw new Error("Review failed");
      const data = await res.json();
      setReviewId(data.review_id);
      startPolling(data.review_id, token);
    } catch (err) {
      console.error(err);
      setReviewStatus("failed");
    }
  };

  const startPolling = (id: number, token: string) => {
    const interval = setInterval(async () => {
      try {
        const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/reviews/${id}`, {
          headers: { Authorization: `Bearer ${token}` },
        });

        if (!res.ok) throw new Error("Failed to fetch");
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
    }, 2000);
  };

  const initDrag = (e: React.MouseEvent) => {
    const startY = e.clientY;
    const startHeight = panelHeight;

    function onMouseMove(e: MouseEvent) {
      setPanelHeight(startHeight - (e.clientY - startY));
    }
    function onMouseUp() {
      window.removeEventListener("mousemove", onMouseMove);
      window.removeEventListener("mouseup", onMouseUp);
    }

    window.addEventListener("mousemove", onMouseMove);
    window.addEventListener("mouseup", onMouseUp);
  };

  
  return (
    <div className="ide-container">
      <TopBar
        theme={theme}
        setTheme={setTheme}
        handleReview={handleReview}
        reviewStatus={reviewStatus}
        handleLogout={handleLogout}
      />

      <div className="ide-main">
        <Sidebar
          files={files}
          activeFileId={activeFileId}
          setActiveFileId={setActiveFileId}
          handleNewFile={handleNewFile}
        />

        <EditorContainer
          activeFile={activeFile}
          setFiles={setFiles}
          activeFileId={activeFileId}
          theme={theme}
        />
      </div>

      <Resizer initDrag={initDrag} />

      <ReviewPanel
        reviewStatus={reviewStatus}
        reviewResult={reviewResult}
        isCollapsed={isReviewCollapsed}
        setIsCollapsed={setIsReviewCollapsed}
        panelHeight={panelHeight}
        setPanelHeight={setPanelHeight}
        reviewId={reviewId}
      />
    </div>
  );
}

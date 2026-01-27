"use client";

import { useState } from "react";
import TopBar from "./components/topbar";
import Sidebar from "./components/sidebar";
import EditorContainer from "./components/editorContainer";
import ReviewPanel from "./components/reviewPanel";
import Resizer from "./components/resizer";
import { useEditorFiles } from "./hooks/useEditorFiles";
import { useReview } from "./hooks/useReview";

export default function EditorPage() {

  const [theme, setTheme] = useState<"vs-dark" | "light">("vs-dark");
  const { files, setFiles, activeFile, activeFileId, setActiveFileId, handleNewFile } =
    useEditorFiles();

  const { reviewStatus, reviewResult, reviewId, handleReview } =
    useReview(activeFile);
 
  const [isReviewCollapsed, setIsReviewCollapsed] = useState(false);
  const [panelHeight, setPanelHeight] = useState(200);


  const handleLogout = () => {
    localStorage.removeItem("token");
    window.location.href = "/login";
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

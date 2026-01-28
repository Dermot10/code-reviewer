"use client";

import { useState } from "react";
import TopBar from "./components/topbar";
import Sidebar from "./components/sidebar";
import EditorContainer from "./components/editorContainer";
import ReviewPanel from "./components/reviewPanel";
import Resizer from "./components/resizer";
import { useEditorFiles } from "./hooks/useEditorFiles";
import { useReview } from "./hooks/useReview";
import Tabs from "./components/tabs";

export default function EditorPage() {

  const [theme, setTheme] = useState<"vs-dark" | "light">("vs-dark");
  const [dirtyFiles, setDirtyFiles] = useState<Set<number>>(new Set());
  const [openFiles, setOpenFiles] = useState<number[]>([]);
  const { files, setFiles, activeFile, activeFileId, setActiveFileId, handleNewFile, handleSaveFile, handleDeleteFile } =
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

  function openFile(fileId: number) {
    setOpenFiles(prev => 
      prev.includes(fileId) ? prev : [...prev, fileId]
    );
    setActiveFileId(fileId)
  }

  function closeFile(fileId: number) { 
    setOpenFiles(prev => {
      const newOpenFiles = prev.filter(id => id !== fileId); 

      if (fileId === activeFileId) {
        setActiveFileId(newOpenFiles[0] ?? null);
      }

      return newOpenFiles;
    });
  }

  
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
          dirtyFiles={dirtyFiles}
          activeFileId={activeFileId}
          openFile={openFile}
          handleNewFile={handleNewFile}
          handleDeleteFile={handleDeleteFile}
        />

        <div className="editor-container">
          <Tabs
            files={files}
            openFiles={openFiles}
            activeFileId={activeFileId}
            dirtyFiles={dirtyFiles}
            setActiveFileId={setActiveFileId}
            closeFile={closeFile}
          />

          <EditorContainer
            theme={theme}
            activeFile={activeFile}
            setFiles={setFiles}
            activeFileId={activeFileId}
            dirtyFiles={dirtyFiles}       // <-- pass the Set itself
            setDirtyFiles={setDirtyFiles}
            handleSaveFile={handleSaveFile}
            files={files}
            openFiles={openFiles}
            setActiveFileId={setActiveFileId}
            closeFile={closeFile}
          />
        </div>
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

"use client";

import { useEffect } from "react";
import { Editor } from "@monaco-editor/react";
import { EditorFile } from "../types";

type Props = {
  activeFile: EditorFile | undefined;
  setFiles: React.Dispatch<React.SetStateAction<EditorFile[]>>;
  activeFileId: number | null;
  theme: "vs-dark" | "light";
  dirtyFiles: Set<number>;  // <-- ADD THIS
  setDirtyFiles: React.Dispatch<React.SetStateAction<Set<number>>>;
  handleSaveFile: (fileId: number, content: string) => Promise<void>;
  files: EditorFile[];
  openFiles: number[];
  setActiveFileId: (id: number) => void;
  closeFile: (id: number) => void;
};


export default function EditorContainer({
  activeFile,
  activeFileId,
  theme,
  setFiles,
  dirtyFiles,
  setDirtyFiles,
  handleSaveFile,
  files,
  openFiles,
  setActiveFileId,
  closeFile,
}: Props) {
  // Auto-save after 2s of inactivity
  useEffect(() => {
    if (!activeFile || activeFileId === null) return;

    const timer = setTimeout(() => {
      handleSaveFile(activeFileId, activeFile.content);
      setDirtyFiles(prev => {
        const next = new Set(prev);
        next.delete(activeFileId);
        return next;
      });
    }, 2000);

    return () => clearTimeout(timer);
  }, [activeFile?.content, activeFileId, handleSaveFile, setDirtyFiles]);

  // No file open fallback
  if (!activeFile || activeFileId === null) {
    return (
      <main className="editor-container empty-editor">
        <p style={{ color: "#888", textAlign: "center", marginTop: "2rem" }}>
          No file open
        </p>
      </main>
    );
  }

  const content = activeFile.content;

  return (
    <main className="editor-container">
      {/* Save button toolbar */}
      <div className="editor-toolbar">
        <button
          className="btn-primary"
          onClick={() => {
            if (!activeFile || activeFileId === null) return;

            handleSaveFile(activeFileId, content).then(() => {
              setDirtyFiles(prev => {
                const next = new Set(prev);
                next.delete(activeFileId);
                return next;
              });
            });
          }}
        >
          Save (Ctrl+S)
        </button>
      </div>

      {/* Tabs */}
      <div className="editor-actions-top">
        {openFiles.map(fileId => {
          const file = files.find(f => f.id === fileId);
          if (!file) return null;

          const isActive = file.id === activeFileId;

          return (
            <div
              key={file.id}
              className={`tab ${isActive ? "active" : ""}`}
              onClick={() => setActiveFileId(file.id)}
            >
              <span className="tab-label">
                {file.name}
                {dirtyFiles.has(file.id) && (
                  <span className="dirty-indicator">*</span>
                )}
              </span>
              <button
                className="tab-close"
                onClick={e => {
                  e.stopPropagation();
                  closeFile(file.id);
                }}
              >
                ×
              </button>
            </div>
          );
        })}
      </div>

      {/* Monaco Editor */}
      <Editor
        height="100%"
        language="python"
        theme={theme}
        value={content}
        onChange={value => {
          if (activeFileId === null) return;

          setFiles(prev =>
            prev.map(f =>
              f.id === activeFileId ? { ...f, content: value || "" } : f
            )
          );

          setDirtyFiles(prev => new Set(prev).add(activeFileId));
        }}
        options={{
          minimap: { enabled: true },
          fontSize: 14,
          fontFamily: "'JetBrains Mono', 'Fira Code', 'Courier New', monospace",
          lineNumbers: "on",
          rulers: [80],
          renderWhitespace: "selection",
          scrollBeyondLastLine: false,
          automaticLayout: true,
        }}
      />
    </main>
  );
}

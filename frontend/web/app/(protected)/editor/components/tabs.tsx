import { EditorFile } from "../types";


type Props = {
    files: EditorFile[];
    openFiles: number[];
    activeFileId: number | null;
    dirtyFiles: Set<number>;
    setActiveFileId: (id: number) => void;
    closeFile: (id: number) => void;
}

export default function Tabs({
  files,
  openFiles,
  activeFileId,
  dirtyFiles,
  setActiveFileId,
  closeFile,
}: Props) {
  return (
    <div className="editor-tabs">
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
              onClick={(e) => {
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
  );
}

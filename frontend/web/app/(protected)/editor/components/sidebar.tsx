import { FileText, FolderOpen, Plus} from "lucide-react";
import { EditorFile } from "../types";


type Props = {
    files: EditorFile[];
    activeFileId: number | null;
    openFile: (id: number) => void;
    handleNewFile:() => void; 
    handleDeleteFile: (id: number) => void;
    dirtyFiles: Set<number>;
}

export default function Sidebar({ 
  files, 
  activeFileId, 
  openFile, 
  handleNewFile, 
  handleDeleteFile,
  dirtyFiles,
}: Props) {
    return (
        <aside className="sidebar">
          <div className="sidebar-header">
            <FolderOpen size={16} />
            <span>Files</span>
          </div>

          <div className="file-list">
            {files.map((file) => (
              <div
                key={file.id}
                className={`file-item ${file.id === activeFileId ? "active" : ""}`}
              >
                <div 
                  className="file-info"
                  onClick={() => openFile(file.id)}
                >
                  <FileText size={14} />
                  <span>
                    {file.name}
                    {dirtyFiles.has(file.id) && (
                    <span className="dirty-indicator">*</span>
                    )}
                  </span>
                </div>
                <button
                  onClick={(e) => {
                    e.stopPropagation();  
                    handleDeleteFile(file.id);
                  }}
                  className="delete-btn"
                >
                  ×
                </button>
              </div>
            ))}
          </div>


          <button className="btn-secondary sidebar-btn" onClick={handleNewFile}>
            <Plus size={14} />
            New File
          </button>
        </aside>
    );
}
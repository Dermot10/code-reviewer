import { FileText, FolderOpen, Plus} from "lucide-react";
import {File} from "@/editor/types";


type Props = {
    files: File[];
    activeFileId: string;
    setActiveFileId: (id: string) => void;
    handleNewFile:() => void; 
}

export default function Sidebar({ files, activeFileId, setActiveFileId, handleNewFile}: Props){
    return (
        <aside className="sidebar">
          <div className="sidebar-header">
            <FolderOpen size={16} />
            <span>Files</span>
          </div>

          <div className="file-list">
            {files.map(file => (
              <div
                key={file.id}
                className={`file-item ${file.id === activeFileId ? "active" : ""}`}
                onClick={() => setActiveFileId(file.id)}
              >
                <FileText size={14}/>
                <span>{file.name}</span>
              </div>
            ))}
          </div>
          <button className="btn-secondary sidebar-btn" onClick={handleNewFile}>
            <Plus size={14} />
            New File
          </button>
        </aside>
    )

}
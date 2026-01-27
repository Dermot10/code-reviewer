import { useState } from "react";
import { EditorFile } from "../types";


export function useEditorFiles(){
    const [files, setFiles] = useState<EditorFile[]>([
        { id: "1", name: "main.py", content: "# Write your Python code here\nprint('Hello, world')" },
        { id: "2", name: "utils.py", content: "def helper():\n    pass" },
    ]);

    const [activeFileId, setActiveFileId] = useState("1");

    const activeFile = files.find(f => f.id === activeFileId)

    const handleNewFile = () => {
        const id = crypto.randomUUID();
        const name = `files${files.length + 1 }.py`;
        setFiles(prev => [...prev, { id, name, content: ""}]);
        setActiveFileId(id);
    };
    return {
        files, 
        setFiles,
        activeFileId, 
        setActiveFileId, 
        activeFile, 
        handleNewFile,
    };
}
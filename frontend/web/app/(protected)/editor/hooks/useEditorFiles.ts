import { useState, useEffect } from "react";
import { EditorFile } from "../types";


const API_URL = process.env.NEXT_PUBLIC_API_URL


export function useEditorFiles(){
    const [files, setFiles] = useState<EditorFile[]>([]);
    const [activeFileId, setActiveFileId] = useState<number | null>(null);
    const [loading, setLoading] = useState(true); 

    // will run on mount, and fallback case for empty file tree
    useEffect(() => {
        fetchFiles();
    }, []); 

    async function fetchFiles(){
        // send token to authenticated endpoint before parsing response data and setting file state
        const token = localStorage.getItem("token");
        if (!token) return;

        try {
            const res = await fetch(`${API_URL}/api/files`, {
                headers: {Authorization: `Bearer ${token}`}, 
            });

            if (!res.ok) throw new Error("Failed to fetch files"); 

            const data: EditorFile[] = await res.json();

            // Debug
            console.log("API Response:", data);
            console.log("Is Array?", Array.isArray(data));
            
            // Safety check
            if (!Array.isArray(data)) {
                console.error("Expected array, got:", typeof data);
                setFiles([]);
                return;
            }
            
            setFiles(data); 

            if (data.length > 0) {
                setActiveFileId(data[0].id);
            } 

        } catch (err) {
            console.error("Failed to fetch files:", err);
        } finally {
            setLoading(false);
        }
    }

    async function handleNewFile() {
        const token = localStorage.getItem("token"); 
        if (!token) return;

        const fileName = prompt("File name:");
        if (!fileName) return;

        try {
            const res = await fetch(`${API_URL}/api/files`, {
                method: "POST", 
                headers: {
                    "Content-Type": "application/json", 
                    Authorization: `Bearer ${token}`, 
                }, 
                body: JSON.stringify({name: fileName, content: ""}), 
            });

            if (!res.ok) throw new Error("Failed to create file");

            const newFile = await res.json();
            setFiles([...files, newFile]);
            setActiveFileId(newFile.id);
        }catch (err){
            console.error("Failed to create file:", err); 
        }
    }

    async function handleSaveFile(fileId: number, content: string){
        const token = localStorage.getItem("token");
        if (!token) return;

        try {
            const res = await fetch(`${API_URL}/api/files/${fileId}`, {
                method: "PUT", 
                headers: {
                    "Content-Type": "application/json", 
                    Authorization: `Bearer ${token}`,
                }, 
                body: JSON.stringify({content}),
            });

            if (!res.ok) throw new Error("Failed to save file");

            const updatedFile = await res.json();
            setFiles(files.map(f => f.id === fileId ? updatedFile : f));
        }catch (err) {
            console.error("Failed to save file:", err); 
        }
    }


    async function handleDeleteFile(fileId: number) {
        const token = localStorage.getItem("token");
        if (!token) return; 

        if (!confirm("Delete this file?")) return; 

        try {
            const res = await fetch(`${API_URL}/api/files/${fileId}`, {
                method: "DELETE", 
                headers: { Authorization: `Bearer ${token}`},
           });

           if (!res.ok) throw new Error("Failed to delete"); 

           setFiles(files.filter(f => f.id !== fileId));

           if (activeFileId === fileId) {
            setActiveFileId(null);
           }
        }catch (err){
            console.error("Failed to delete file:", err);
        }
    }

    // derived from existing state for synchronicity 
    const activeFile = (files && files.length > 0) 
        ? files.find(f => f.id === activeFileId) ?? null
        : null;

    return {
    files: files || [],
    setFiles,
    activeFile: activeFile ?? undefined,
    activeFileId,
    setActiveFileId,
    handleNewFile,
    handleSaveFile,
    handleDeleteFile,
    loading,
  };
}
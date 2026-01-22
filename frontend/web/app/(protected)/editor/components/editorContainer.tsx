import { Editor } from "@monaco-editor/react";
import { EditorFile } from "../types";

type Props = {
  activeFile: EditorFile | undefined;
  setFiles: React.Dispatch<React.SetStateAction<EditorFile[]>>;
  activeFileId: string;
  theme: "vs-dark" | "light";
};

export default function EditorContainer({ activeFile, setFiles, activeFileId, theme }: Props) {
  return (
    <main className="editor-container">
      <div className="editor-tabs">
        <div className="tab active">
          <span>{activeFile?.name || "Untitled"}</span>
        </div>
      </div>

      <Editor
        height="100%"
        language="python"
        theme={theme}
        value={activeFile?.content || ""}
        onChange={(value) =>
          setFiles(prev =>
            prev.map(f => f.id === activeFileId ? { ...f, content: value || "" } : f)
          )
        }
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

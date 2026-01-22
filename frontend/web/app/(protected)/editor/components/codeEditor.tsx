"use client";


import dynamic from "next/dynamic";

const Editor = dynamic(() => import("@monaco-editor/react"), {
  ssr: false,
});

type CodeEditorProps = {
  value: string;
  onChange?:(value: string) => void;
  readOnly?: boolean;
}

export default function CodeEditor({ value, onChange, readOnly = false}: CodeEditorProps) {
  return (
    <Editor
      height="400px"
      width="500px"
      defaultLanguage="python"
      value={value}
      onChange={(newValue) => onChange?.(newValue ?? "")}
      theme="vs-dark"
      options={{
        readOnly: readOnly,
        minimap: { enabled: false },
        fontSize: 14,
        scrollBeyondLastLine: false,
      }}
    />
  );
}

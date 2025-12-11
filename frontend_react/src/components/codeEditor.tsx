import Editor from "@monaco-editor/react";


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

import Editor from "@monaco-editor/react";


type CodeEditorProps = {
  value: string;
  onChange:(value: string) => void;
}


export default function CodeEditor({ value, onChange }: CodeEditorProps) {
  return (
    <Editor
      height="400px"
      width="600px"
      defaultLanguage="python"
      value={value}
      onChange={(newValue) =>onChange(newValue ?? "")}
      theme="vs-dark"
    />
  );
}

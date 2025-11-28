import Editor from "@monaco-editor/react";

export default function CodeEditor({ language, value, onChange }: any) {
  return (
    <Editor
      height="400px"
      width="600px"
      defaultLanguage={language}
      defaultValue={value}
      onChange={onChange}
      theme="vs-dark"
    />
  );
}

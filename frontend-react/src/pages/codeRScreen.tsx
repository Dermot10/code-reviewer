import CodeEditor from "../components/codeEditor";
import ResultsPanel from "../components/resultsPanel";
import ErrorPanel from "../components/errorPanel";
import { useState } from 'react';
import { useEffect } from "react";
import "../index.css";

type AppState = (typeof AppStates)[number];
type MockResponse = {
  ok: Boolean;
  json: () => Promise<{ feedback: string; issues: any[]}>;
};
const AppStates = ["idle", "submitting", "results", "error"] as const; // ensure they types and not str

export default function MainScreen() {
  const [currentState, setCurrentState] = useState<AppState>("idle");
  const [code, setCode] = useState("");
  const [result, setResult] = useState("");
  const [copied, setCopied] = useState(false);
  const [theme, setTheme] = useState("light");

  useEffect(() => {
    if (currentState === "results"){
      const panel = document.getElementById("results-panel");
      panel?.scrollIntoView({behavior: "smooth"})
    }
  }, [currentState, result])

  useEffect(() => {
    document.documentElement.setAttribute("data-theme", theme);
  })

  async function handleSubmit(){
    setCurrentState("submitting");
    const res = await submitCode(code);
    if (!res.ok) setCurrentState("error");
    else {
      const data = await res.json();
      setResult(data.feedback);
      setCurrentState("results");
    }
  }

  function handleCopy(e: React.MouseEvent<HTMLButtonElement>){
    navigator.clipboard.writeText(result);
    e.currentTarget.blur();
    setTimeout(() => setCopied(false), 1000)
  }

  function handleFileUpload(e: React.ChangeEvent<HTMLInputElement>){
    const file = e.target.files?.[0]; // if file, chained conditional 
    if (!file) return;
    const reader = new FileReader();
    reader.onload = (event) => {
      if (event.target?.result) setCode(event.target.result as string)
    };
    reader.readAsText(file);
  }

  return (
  <div className="app-root">
    <header className="app-header">
      <h1>AI Code Reviewer</h1>

      <button 
        className="theme-toggle"
        onClick={() => setTheme(theme === "light" ? "dark" : "light")}
      >
        {theme === "light" ? "🌙" : "☀️"}
      </button>
    </header>

    <main className="main-container">
      <CodeEditor value={code} onChange={setCode} />

      <div className="right-panel">
        {currentState === "results" && (
          <ResultsPanel 
            result={result} 
            onCopy={handleCopy} 
            copied={copied} 
          />
        )}

        {currentState === "error" && (
          <ErrorPanel onRetry={handleSubmit} />
        )}

        {currentState === "submitting" && (
          <div className="spinner" />
        )}
      </div>
    </main>

    <footer>
      <button 
        onClick={handleSubmit} 
        disabled={currentState === "submitting"}
      >
        Review Code
      </button>
    </footer>
  </div>
);

}


async function submitCode(code: string): Promise<MockResponse>{
  return new Promise((resolve) => {
    setTimeout(() => {
      resolve({ 
        ok: true, 
        json: async () => ({ feedback: "Looks good!", issues: []}),
      });
    }, 1000);  // simulated network delay
  });
}
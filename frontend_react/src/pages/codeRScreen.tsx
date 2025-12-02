import CodeEditor from "../components/codeEditor";
import ResultsPanel from "../components/resultsPanel";
import ErrorPanel from "../components/errorPanel";
import { useState, useEffect } from "react";
import "../index.css";

type AppState = (typeof AppStates)[number];
type ReviewResponse = { feedback: string; issues: any[]};
const AppStates = ["idle", "submitting", "results", "error"] as const; // ensures types, not strings

export default function MainScreen() {
  const [currentState, setCurrentState] = useState<AppState>("idle");
  const [code, setCode] = useState("");
  const [result, setResult] = useState("");
  const [copied, setCopied] = useState(false);
  const [theme, setTheme] = useState("light");

  useEffect(() => {
    if (currentState === "results") {
      const panel = document.getElementById("results-panel");
      panel?.scrollIntoView({ behavior: "smooth" });
    }
  }, [currentState, result]);

  useEffect(() => {
    document.documentElement.setAttribute("data-theme", theme);
  }, [theme]);

  async function handleSubmit() {
    try {
      setCurrentState("submitting");
      const data = await submitCode(code);
      setResult(data.feedback);
      setCurrentState("results");
    } catch (err) {
      console.error("Submit failed:", err);
      setCurrentState("error");
    }
  }

  function handleCopy(e: React.MouseEvent<HTMLButtonElement>) {
    navigator.clipboard.writeText(result);
    e.currentTarget.blur();
    setTimeout(() => setCopied(false), 1000);
  }

  function handleFileUpload(e: React.ChangeEvent<HTMLInputElement>) {
    const file = e.target.files?.[0];
    if (!file) return;
    const reader = new FileReader();
    reader.onload = (event) => {
      if (event.target?.result) setCode(event.target.result as string);
    };
    reader.readAsText(file);
  }

  return (
    <div className="app-root">
      <header className="app-header">
        <h1>Python Reviewer</h1>
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
            <ResultsPanel result={result} onCopy={handleCopy} copied={copied} />
          )}

          {currentState === "error" && <ErrorPanel onRetry={handleSubmit} />}

          {currentState === "submitting" && <div className="spinner" />}
        </div>
      </main>

      <footer>
        <button onClick={handleSubmit} disabled={currentState === "submitting"}>
          Review Code
        </button>
      </footer>
    </div>
  );
}

// submitCode stays outside the component
async function submitCode(code: string): Promise<ReviewResponse> {
  const res = await fetch("http://localhost:8080/review-code", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ submitted_code:code }),
  });

  if (!res.ok) {
    throw new Error(`Request failed: ${res.status}`);
  }
  return await res.json();
}

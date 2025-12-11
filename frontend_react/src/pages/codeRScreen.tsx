import CodeEditor from "../components/codeEditor";
import {ReviewPanel} from "../components/resultsPanel";
import { useState, useEffect } from "react";
import "../index.css";

const AppStates = ["idle", "submitting", "results", "error"] as const;
type AppState = (typeof AppStates)[number];
type ReviewResponse = { feedback: string; issues: any[] };
type EnhancedResponse = {output: string};

export default function MainScreen() {
  const [currentState, setCurrentState] = useState<AppState>("idle");
  const [code, setCode] = useState("");
  const [review, setReview] = useState<ReviewResponse | null>(null);
  const [enhancedCode, setEnhancedCode] = useState<EnhancedResponse| null>(null);
  const [copied, setCopied] = useState(false);
  const [exportType, setExportType] = useState("md");
  const [theme, setTheme] = useState("dark");

  useEffect(() => {
    if (currentState === "results") {
      const panel = document.getElementById("results-panel");
      panel?.scrollIntoView({ behavior: "smooth" });
    }
  }, [currentState, review]); // useEffect is called when dependency array changes in value 
  

  useEffect(() => {
    document.documentElement.setAttribute("data-theme", theme);
  }, [theme]);

  async function handleReviewSubmit() {
    try {
      setCurrentState("submitting");
      const data = await submitReviewCode(code);
      setReview(data);
      setCurrentState("results");
    } catch (err) {
      console.error("Submit failed:", err);
      setCurrentState("error");
    }
  }                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                             

  async function handleEnhanceSubmit() {
    try {
      setCurrentState("submitting");
      const data = await submitEnhancedCode(code);
      console.log("enhanced:", data);
      setEnhancedCode(data);
      setCurrentState("results");
    } catch (err) {
      console.error("Submit failed:", err);
      setCurrentState("error");
    }
  }

  function handleCopy(text: string, e: React.MouseEvent<HTMLButtonElement>) {
      navigator.clipboard.writeText(text);
      e.currentTarget.blur();
      setCopied(true);
      setTimeout(() => setCopied(false), 1000);
    }

  function handleFileUpload(e: React.ChangeEvent<HTMLInputElement>) {
    try {
      const file = e.target.files?.[0];
      if (!file) return;

      const reader = new FileReader();
      reader.onload = (event) => {
        if (event.target?.result) setCode(event.target.result as string);
      };
      reader.readAsText(file);
    } catch (err) {
      console.error("File upload failed:", err);
      setCurrentState("error");
    }
  }

  async function handleExport(type: string) {
    try {
      setCurrentState("submitting");

      const res = await fetch(
        `http://localhost:8080/review-code/download?type=${type}`,
        {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify(review),
        }
      );

      if (!res.ok) {
        throw new Error(`Export failed: ${res.status}`);
      }

      const blob = await res.blob();
      const url = window.URL.createObjectURL(blob);

      const a = document.createElement("a");
      a.href = url;
      a.download = `code_review.${type}`;
      a.click();

      window.URL.revokeObjectURL(url);

      setCurrentState("results");
    } catch (err) {
      console.error("Export failed:", err);
      setCurrentState("error");
    }
  }

  return (
    <div className="app-root">
      <header className="app-header">
        <div className="header-left">
          <h1>Python Reviewer</h1>
        </div>
        <div className="header-right">
          <button
            className="theme-toggle"
            onClick={() => setTheme(theme === "light" ? "dark" : "light")}
          >
            {theme === "light" ? "🌙" : "☀️"}
          </button>
          <input
            type="file"
            id="file-upload"
            style={{ display: "none" }}
            onChange={handleFileUpload}
          />
          <button
            className="upload-button"
            onClick={() => document.getElementById("file-upload")?.click()}
          >
            📁 Upload
          </button>

          <div className="export-container">
            <button
              onClick={() => handleExport(exportType)}
              disabled={currentState !== "results"}
            >
              Export
            </button>
            <select
              value={exportType}
              onChange={(e) => setExportType(e.target.value)}
            >
              <option value="md">Markdown</option>
              <option value="txt">Text</option>
              <option value="json">JSON</option>
              <option value="csv">CSV</option>
                    </select>
                  </div>
                </div>
              </header>

              <main className="main-container">
            {/* Input Code Editor */}
            <div className="editor-panel">
              <h2>Input Code</h2>
              <CodeEditor value={code} onChange={setCode} />
            </div>

            {/* Enhanced Code Editor */}
            <div className="editor-panel">
              <h2>Enhanced Code</h2>
              <CodeEditor
                value={enhancedCode?.output || ""}
                onChange={() => {}}
                readOnly={true}
              />
            </div>

            {/* Review & Issues Panel */}
            <div className="editor-panel review-panel-wrapper">
              <h2>Review / Issues</h2>
              {review ? (
                <>
                  <ReviewPanel result={review.feedback} onCopy={handleCopy} copied={copied} />
                  {review.issues.length > 0 && (
                    <div className="issues-panel">
                      <h3>Issues Found:</h3>
                      <ul>
                        {review.issues.map((issue, idx) => (
                          <li key={idx}>
                            <strong>Line {issue.line}</strong> [{issue.type}]: {issue.description}
                          </li>
                        ))}
                      </ul>
                    </div>
                  )}
                </>
              ) : (
                <div className="placeholder">No review data yet</div>
              )}
            </div>
        </main>


      <footer>
        <button onClick={handleReviewSubmit} disabled={currentState === "submitting"}>
          Review Code
        </button>
        <button
            onClick={(e) => {
              e.preventDefault(); // optional, good habit
              handleEnhanceSubmit();
            }}
            disabled={currentState === "submitting"}
          >
            Enhance Code
      </button>
      </footer>
    </div>
  );
}

async function submitReviewCode(code: string): Promise<ReviewResponse> {
  const res = await fetch("http://localhost:8080/review-code", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ submitted_code: code }),
  });

  if (!res.ok) {
    throw new Error(`Request failed: ${res.status}`);
  }
  return await res.json();
}


async function submitEnhancedCode(code: string): Promise<EnhancedResponse> {
  const res = await fetch("http://localhost:8080/enhance-code", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ submitted_code: code }),
  });

  if (!res.ok) {
    throw new Error(`Request failed: ${res.status}`);
  }
  return await res.json();
}

"use client";

import CodeEditor from "../components/codeEditor";
import { ReviewPanel } from "../components/resultsPanel";
import { useState, useEffect } from "react";
import { ReviewResult, EnhancedResult, Issue } from "../interfaces/responses";
import { CodeRequest } from "../interfaces/requests";

const AppStates = ["idle", "submitting", "results", "error"] as const;
type AppState = (typeof AppStates)[number];

export default function MainScreen() {
  const [currentState, setCurrentState] = useState<AppState>("idle");
  const [code, setCode] = useState<string>("");
  const [review, setReview] = useState<ReviewResult | null>(null);
  const [enhancedCode, setEnhancedCode] = useState<EnhancedResult | null>(null);
  const [copied, setCopied] = useState<boolean>(false);
  const [exportType, setExportType] = useState<string>("md");
  const [theme, setTheme] = useState<"light" | "dark">("dark");

  useEffect(() => {
    if (currentState === "results") {
      document.getElementById("results-panel")?.scrollIntoView({ behavior: "smooth" });
    }
  }, [currentState, review, enhancedCode]);

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
      setEnhancedCode(data);
      setCurrentState("results");
    } catch (err) {
      console.error("Enhance failed:", err);
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
    const file = e.target.files?.[0];
    if (!file) return;

    const reader = new FileReader();
    reader.onload = (event) => {
      if (event.target?.result) setCode(event.target.result as string);
    };
    reader.readAsText(file);
  }

  async function handleExport(type: string) {
    if (!review) return;

    try {
      setCurrentState("submitting");

      const payload = { ...review, enhanced_code: enhancedCode?.enhanced_code || null };

      const res = await fetch(
        `${process.env.NEXT_PUBLIC_API_URL}/review-code/download?type=${type}`,
        {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify(payload),
        }
      );

      if (!res.ok) throw new Error(`Export failed: ${res.status}`);

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
          <button className="theme-toggle" onClick={() => setTheme(theme === "light" ? "dark" : "light")}>
            {theme === "light" ? "🌙" : "☀️"}
          </button>

          <input type="file" id="file-upload" style={{ display: "none" }} onChange={handleFileUpload} />
          <button className="upload-button" onClick={() => document.getElementById("file-upload")?.click()}>
            📁 Upload
          </button>

          <div className="export-container">
            <button onClick={() => handleExport(exportType)} disabled={currentState !== "results"}>
              Export
            </button>
            <select value={exportType} onChange={(e) => setExportType(e.target.value)}>
              <option value="md">Markdown</option>
              <option value="txt">Text</option>
              <option value="json">JSON</option>
              <option value="csv">CSV</option>
            </select>
          </div>
        </div>
      </header>

      <main className="main-container">
        <div className="editor-panel">
          <h2>Input Code</h2>
          <CodeEditor value={code} onChange={setCode} />
        </div>

        <div className="editor-panel">
          <h2>Enhanced Code</h2>
          <CodeEditor value={enhancedCode?.enhanced_code || ""} onChange={() => {}} readOnly />
        </div>

        <div className="editor-panel review-panel-wrapper" id="results-panel">
          <h2>Review / Issues</h2>
          {review ? (
            <>
              <ReviewPanel result={review.feedback} onCopy={handleCopy} copied={copied} />
              {review.issues.length > 0 && (
                <div className="issues-panel">
                  <h3>Issues Found:</h3>
                  <ul>
                    {review.issues.map((issue: Issue, idx: number) => (
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
        <button onClick={handleEnhanceSubmit} disabled={currentState === "submitting"}>
          Enhance Code
        </button>
      </footer>
    </div>
  );
}

async function submitReviewCode(code: string): Promise<ReviewResult> {
  const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/review-code`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ submitted_code: code } as CodeRequest),
  });
  if (!res.ok) throw new Error(`Request failed: ${res.status}`);
  return await res.json();
}

async function submitEnhancedCode(code: string): Promise<EnhancedResult> {
  const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/enhance-code`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ submitted_code: code } as CodeRequest),
  });
  if (!res.ok) throw new Error(`Request failed: ${res.status}`);
  return await res.json();
}

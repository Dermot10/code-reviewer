interface PanelProps {
  result: string;
  onCopy: (text: string, e: React.MouseEvent<HTMLButtonElement>) => void;
  copied: boolean;
}

export function ReviewPanel({ result, onCopy, copied }: PanelProps) {
  return (
    <div className="panel results-panel">
      <h2>Results</h2>
      <pre className="output">{result}</pre>
      <button onClick={(e) => onCopy(result, e)}>
        {copied ? "Copied" : "Copy All"}
      </button>
    </div>
  );
}

export function EnhancedCodePanel({ result, onCopy, copied }: PanelProps) {
  return (
    <div className="panel results-panel">
      <h2>Results</h2>
      <pre className="output">{result}</pre>
      <button onClick={(e) => onCopy(result, e)}>
        {copied ? "Copied" : "Copy All"}
      </button>
    </div>
  );
}

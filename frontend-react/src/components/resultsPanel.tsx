interface ResultsPanelProps {
    result: string;
    onCopy:(e: React.MouseEvent<HTMLButtonElement>) => void;
    copied: boolean;
}

export default function ResultsPanel({result, onCopy, copied}: ResultsPanelProps) {
    return(
        <div className="panel results-panel">
            <h2>Results</h2>
            <pre className="output">{result}</pre>
            <button 
                onClick={(e) => {
                    onCopy(e);
                    e.currentTarget.blur();
                }}>
                {copied ? "Copied": "Copy All"}  
            </button>
        </div>
    );
}
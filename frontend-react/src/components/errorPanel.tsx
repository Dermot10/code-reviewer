

interface ErrorPanelProps {
    onRetry:(e: React.MouseEvent<HTMLButtonElement>) => void;
}

export default function ErrorPanel({onRetry}: ErrorPanelProps) {
    return(
        <div className="panel error-panel">
            <h2>Error</h2>
            
            <button 
                onClick={(e) => {
                    onRetry(e); 
                    e.currentTarget.blur()
                }}
            >
                Resubmit   
            </button>
        </div>
    );
}
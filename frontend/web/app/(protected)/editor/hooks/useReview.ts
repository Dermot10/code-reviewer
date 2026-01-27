import { useState } from "react";
import { ReviewStatus } from "../types";
import { EditorFile } from "../types";


export function useReview(activeFile?: EditorFile){
    const [reviewStatus, setReviewStatus] = useState<ReviewStatus>("idle");
    const [reviewResult, setReviewResult] = useState<any>(null);
    const [reviewId, setReviewId] = useState<number | null>(null); 

    const handleReview = async () => {
        if (!activeFile) return;

        const token = localStorage.getItem("token");
        if (!token) return;

        try {
            setReviewStatus("processing");

            const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/review-code`, {
                method: "POST", 
                headers: {
                    "Content-Type": "application/json", 
                    Authorization: `Bearer ${token}`,
                },
                body: JSON.stringify({code : activeFile.content}),
            });
            
            if (!res.ok) throw new Error("Review failed");

            const data = await res.json();
            setReviewId(data.review_id);
            startPolling(data.review_id, token); 
        } catch {
            setReviewStatus("failed");
        }
    };
    
        const startPolling = (id: number, token: string) => {
            const interval = setInterval(async () => {
                try {
                    const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/reviews/${id}`, {
                        headers: {Authorization: `Bearer ${token}`},
                    });

                    const data = await res.json();

                    if (data.status === "completed") {
                        clearInterval(interval);
                        setReviewResult(data);
                        setReviewStatus("completed");
                    }

                    if (data.status === "failed"){
                        clearInterval(interval);
                        setReviewResult("failed");
                    }
                } catch {
                    clearInterval(interval);
                    setReviewStatus("failed");
                }
            }, 2000);
        };
    return {
        reviewStatus, 
        reviewResult, 
        reviewId, 
        handleReview,
    };
}
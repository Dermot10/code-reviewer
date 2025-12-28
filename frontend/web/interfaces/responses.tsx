export interface Issue{ 
    line: number;
    type: "bug" | "security" | "style" | "other";
    description: string;
}

export interface ReviewResult{
    review_id: string;
    project_id?: string;
    user_id?: string; 
    feedback: string;
    issues: Issue[];
    ai_model_version: string; 
    created_at: string; 
}

export interface EnhancedResult{
    review_id: string;
    enhanced_code: string; 
    ai_model_version: string; 
    created_at: string; 
}
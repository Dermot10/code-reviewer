export type ReviewStatus = "idle" | "processing" | "completed" | "failed";
export type EditorFile = {
  id: number;
  name: string;
  content: string;
  created_at: string;
  updated_at: string;
};


"use client";

import { useRouter } from "next/navigation";
import { useEffect } from "react";

export default function ProtectedRoot() {
  const router = useRouter();
  
  useEffect(() => {
    router.push("/editor");
  }, [router]);

  return <div>Redirecting...</div>;
}
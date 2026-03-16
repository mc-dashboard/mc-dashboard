import { useCallback, useEffect, useState } from "react";
import { API_BASE_URL } from "../libs/api";

interface User {
  email: string;
  name: string;
}

export function useAuth() {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);

  const refetch = useCallback(() => {
    fetch(`${API_BASE_URL}/api/user`, { credentials: "include" })
      .then((res) => (res.ok ? res.json() : null))
      .then((data) => setUser(data?.email ? data : null));
  }, []);

  useEffect(() => {
    fetch(`${API_BASE_URL}/api/user`, { credentials: "include" })
      .then((res) => (res.ok ? res.json() : null))
      .then((data) => setUser(data?.email ? data : null))
      .finally(() => setLoading(false));
  }, []);

  return { user, loading, refetch };
}

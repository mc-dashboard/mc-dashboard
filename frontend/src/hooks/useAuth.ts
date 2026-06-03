import { useCallback, useEffect, useState } from "react";
import { apiFetch } from "../libs/api";

interface User {
  email: string;
  name: string;
}

interface AuthState {
  user: User | null;
  loading: boolean;
}

export function useAuth() {
  const [{ user, loading }, setAuth] = useState<AuthState>({ user: null, loading: true });

  const refetch = useCallback(() => {
    apiFetch<User>("/api/user")
      .then((data) => setAuth({ user: data?.email ? data : null, loading: false }))
      .catch(() => setAuth({ user: null, loading: false }));
  }, []);

  useEffect(() => {
    let cancelled = false;
    apiFetch<User>("/api/user")
      .then((data) => { if (!cancelled) setAuth({ user: data?.email ? data : null, loading: false }); })
      .catch(() => { if (!cancelled) setAuth({ user: null, loading: false }); });
    return () => { cancelled = true; };
  }, []);

  return { user, loading, refetch };
}

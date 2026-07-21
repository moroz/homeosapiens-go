import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

import { api } from "~/lib/api";

/** `GET /api/admin/session` — the currently signed in user. */
export function useGetSessionQuery() {
  return useQuery({
    queryKey: ["getSession"],
    queryFn: async () => {
      const { data } = await api.GET("/session");
      return data;
    },
  });
}

async function signOut() {
  return fetch("/sign-out", {
    method: "DELETE",
    credentials: "include",
    headers: { Accept: "application/json" },
  });
}

export function useSignOutMutation() {
  const client = useQueryClient();

  return useMutation({
    mutationFn: signOut,
    onSuccess() {
      // Set (not just invalidate) so the stale user doesn't linger: a failed
      // background refetch keeps the last-known-good `data` in the cache,
      // which would leave `useGetSessionQuery` reporting a signed-in user.
      client.setQueryData(["getSession"], null);
    },
  });
}

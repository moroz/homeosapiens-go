import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

import { api } from "~/lib/api";

/** `GET /api/admin/session` — the currently signed in user. */
export function useGetSessionQuery() {
  return useQuery({
    queryKey: ["getSession"],
    queryFn: async () => {
      try {
        const { data } = await api.GET("/session");
        return data!;
      } catch (e) {
        return null;
      }
    },
    // Force a refetch on focus regardless of staleTime, so signing out in
    // another tab is picked up as soon as this one regains focus.
    refetchOnWindowFocus: "always",
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

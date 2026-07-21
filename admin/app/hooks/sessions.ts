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
  });
}

export function useSignOutMutation() {
  const client = useQueryClient();

  return useMutation({
    mutationFn: signOut,
    async onSuccess() {
      await client.invalidateQueries({
        queryKey: ["getSession"],
      });
    },
  });
}

import { useQuery } from "@tanstack/react-query";

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

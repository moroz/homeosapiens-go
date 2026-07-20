import { useQuery } from "@tanstack/react-query";

import { api } from "~/lib/api";

/** Example query hook: `GET /api/admin/health`. */
export function useHealth() {
  return useQuery({
    queryKey: ["health"],
    queryFn: async () => {
      const { data } = await api.GET("/health");
      return data;
    },
  });
}

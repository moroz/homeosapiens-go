import { useQuery } from "@tanstack/react-query";

import { api } from "~/lib/api";

/** `GET /api/admin/hosts` — a page of hosts, ordered by name. */
export function useListHostsQuery(page = 1, perPage = 100) {
  return useQuery({
    queryKey: ["listHosts", page, perPage],
    queryFn: async () => {
      const { data } = await api.GET("/hosts", { params: { query: { page, perPage } } });
      return data;
    },
  });
}
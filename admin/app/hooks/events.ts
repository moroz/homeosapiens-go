import { useQuery } from "@tanstack/react-query";

import { api } from "~/lib/api";

/** `GET /api/admin/events` — a page of events, newest first. */
export function useListEventsQuery(page = 1, perPage = 20) {
  return useQuery({
    queryKey: ["listEvents", page, perPage],
    queryFn: async () => {
      const { data } = await api.GET("/events", { params: { query: { page, perPage } } });
      return data;
    },
  });
}

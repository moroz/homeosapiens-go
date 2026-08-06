import { useQuery } from "@tanstack/react-query";
import { api } from "~/lib/api";
import type { components } from "~/lib/api-types";
import type { UUID } from "~/lib/interfaces";

export type EventAttendant = components["schemas"]["EventAttendant"];

/** `GET /api/admin/events/{id}/attendants` — a page of event attendants (subset of user data), newest first */
export function useListEventAttendantsQuery(eventId: UUID, page = 1, perPage = 20) {
  return useQuery({
    queryKey: ["listEventAttendants", eventId],
    queryFn: async () => {
      const { data } = await api.GET("/events/{id}/attendants", {
        params: { path: { id: eventId }, query: { page, perPage } },
      });
      return data;
    },
  });
}

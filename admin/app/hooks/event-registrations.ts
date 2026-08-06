import { useQuery } from "@tanstack/react-query";
import { api } from "~/lib/api";
import type { components } from "~/lib/api-types";
import type { UUID } from "~/lib/interfaces";

export type EventAttendant = components["schemas"]["EventAttendant"];

interface UseListEventAttendantsQueryParams {
  eventId: UUID;
  page?: number;
  perPage?: number;
}

/** `GET /api/admin/events/{id}/attendants` — a page of event attendants (subset of user data), newest first */
export function useListEventAttendantsQuery({
  eventId,
  page = 1,
  perPage = 20,
}: UseListEventAttendantsQueryParams) {
  return useQuery({
    queryKey: ["listEventAttendants", eventId],
    queryFn: async () => {
      const result = await api.GET("/events/{id}/attendants", {
        params: { path: { id: eventId }, query: { page, perPage } },
      });
      return result.data;
    },
  });
}

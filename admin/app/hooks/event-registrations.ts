import { keepPreviousData, useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { api } from "~/lib/api";
import type { components } from "~/lib/api-types";
import type { UUID } from "~/lib/interfaces";

export type EventAttendant = components["schemas"]["EventAttendant"];
export type EnrollStudentForEventInput = components["schemas"]["EnrollStudentForEventInput"];

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
    queryKey: ["listEventAttendants", { eventId, page, perPage }],
    queryFn: async () => {
      const { data } = await api.GET("/events/{id}/attendants", {
        params: { path: { id: eventId }, query: { page, perPage } },
      });
      return data;
    },
  });
}

interface UseListEligibleUsersForEventQueryParams {
  eventId: UUID;
  searchTerm: string;
}

export function useListEligibleUsersForEventQuery({
  eventId,
  searchTerm,
}: UseListEligibleUsersForEventQueryParams) {
  return useQuery({
    queryKey: ["listEligibleUsersForEvent", eventId, searchTerm],
    queryFn: async () => {
      const { data } = await api.GET("/events/{id}/eligible-users", {
        params: { path: { id: eventId }, query: { q: searchTerm } },
      });
      return data;
    },
    placeholderData: keepPreviousData,
  });
}

export function useEnrollStudentForEventMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (params: EnrollStudentForEventInput) => {
      const { data } = await api.POST("/event-registrations", { body: params });
      return data;
    },
    onSuccess(_, { eventId }) {
      queryClient.invalidateQueries({
        queryKey: ["listEventAttendants", { eventId }],
        exact: false,
      });
    },
  });
}

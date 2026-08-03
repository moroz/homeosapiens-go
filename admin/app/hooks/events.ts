import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

import { api } from "~/lib/api";
import type { components } from "~/lib/api-types";
import type { UUID } from "~/lib/interfaces";

export type EventInput = components["schemas"]["EventInput"];
export type PatchEventInput = components["schemas"]["PatchEventInput"];
export type EventDetails = components["schemas"]["EventDetails"];

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

/** `GET /api/admin/events/{id}` — a single event by primary key. */
export function useGetEventQuery(id: string | undefined) {
  return useQuery({
    queryKey: ["getEvent", id],
    enabled: id != null,
    queryFn: async () => {
      const { data } = await api.GET("/events/{id}", { params: { path: { id: id! } } });
      return data;
    },
  });
}

/** `POST /api/admin/events` — create an event. Throws {@link ApiError} on failure (see `~/lib/api`). */
export function useCreateEventMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (body: EventInput) => {
      const { data } = await api.POST("/events", { body });
      return data!;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["listEvents"] });
    },
  });
}

interface UpdateEventMutationParams {
  id: UUID;
  params: PatchEventInput;
}

/** `PATCH /api/admin/events/{id}` — selective update. Throws {@link ApiError} on failure (see `~/lib/api`). */
export function useUpdateEventMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({ id, params }: UpdateEventMutationParams) => {
      const { data } = await api.PATCH("/events/{id}", {
        params: { path: { id } },
        body: params,
      });
      return data!;
    },
    onSuccess(event) {
      queryClient.invalidateQueries({ queryKey: ["listEvents"] });
      queryClient.invalidateQueries({ queryKey: ["getEvent", event.id] });
    },
  });
}

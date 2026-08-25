import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

import { api } from "~/lib/api";
import type { components } from "~/lib/api-types";
import type { UUID } from "~/lib/interfaces";

export type Host = components["schemas"]["Host"];
export type HostInput = components["schemas"]["HostInput"];

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

/** `GET /api/admin/hosts/{id}` — a single host by primary key. */
export function useGetHostQuery(id: string | undefined) {
  return useQuery({
    queryKey: ["getHost", id],
    enabled: id != null,
    queryFn: async () => {
      const { data } = await api.GET("/hosts/{id}", { params: { path: { id: id! } } });
      return data;
    },
  });
}

/** `POST /api/admin/hosts` — create a host. Throws {@link ApiError} on failure (see `~/lib/api`). */
export function useCreateHostMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (body: HostInput) => {
      const { data } = await api.POST("/hosts", { body });
      return data!;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["listHosts"] });
    },
  });
}

interface UpdateHostMutationParams {
  id: UUID;
  params: HostInput;
}

/** `PATCH /api/admin/hosts/{id}` — replace the editable fields of a host. */
export function useUpdateHostMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({ id, params }: UpdateHostMutationParams) => {
      const { data } = await api.PATCH("/hosts/{id}", {
        params: { path: { id } },
        body: params,
      });
      return data!;
    },
    onSuccess(host) {
      queryClient.invalidateQueries({ queryKey: ["listHosts"] });
      queryClient.invalidateQueries({ queryKey: ["getHost", host.id] });
    },
  });
}

/**
 * `DELETE /api/admin/hosts/{id}` — delete a host. The server answers 409 when
 * the host is still referenced by an event or a video.
 */
export function useDeleteHostMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (id: UUID) => {
      await api.DELETE("/hosts/{id}", { params: { path: { id } } });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["listHosts"] });
    },
  });
}

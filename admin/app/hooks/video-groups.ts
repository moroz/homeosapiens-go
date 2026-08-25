import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

import { api } from "~/lib/api";
import type { components } from "~/lib/api-types";
import type { UUID } from "~/lib/interfaces";

export type VideoGroup = components["schemas"]["VideoGroup"];
export type VideoGroupInput = components["schemas"]["VideoGroupInput"];
export type PatchVideoGroupInput = components["schemas"]["VideoGroupPatch"];

/** `GET /api/admin/video-groups` — a page of video groups, newest first. */
export function useListVideoGroupsQuery(page = 1, perPage = 20) {
  return useQuery({
    queryKey: ["listVideoGroups", page, perPage],
    queryFn: async () => {
      const { data } = await api.GET("/video-groups", { params: { query: { page, perPage } } });
      return data;
    },
  });
}

/** `GET /api/admin/video-groups/{id}` — a single video group by primary key. */
export function useGetVideoGroupQuery(id: string | undefined) {
  return useQuery({
    queryKey: ["getVideoGroup", id],
    enabled: id != null,
    queryFn: async () => {
      const { data } = await api.GET("/video-groups/{id}", { params: { path: { id: id! } } });
      return data;
    },
  });
}

/** `GET /api/admin/video-groups/{id}/videos` — the group's videos, in playback order. */
export function useListVideosInVideoGroupQuery(id: string | undefined) {
  return useQuery({
    queryKey: ["listVideosInVideoGroup", id],
    enabled: id != null,
    queryFn: async () => {
      const { data } = await api.GET("/video-groups/{id}/videos", {
        params: { path: { id: id! } },
      });
      return data;
    },
  });
}

interface ReplaceVideosMutationParams {
  id: UUID;
  videoIds: UUID[];
}

/**
 * `PUT /api/admin/video-groups/{id}/videos` — set the group's videos and their
 * order in one go. Throws {@link ApiError} on failure (see `~/lib/api`).
 */
export function useReplaceVideosInVideoGroupMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({ id, videoIds }: ReplaceVideosMutationParams) => {
      const { data } = await api.PUT("/video-groups/{id}/videos", {
        params: { path: { id } },
        body: { videoIds },
      });
      return data!;
    },
    onSuccess(_data, { id }) {
      queryClient.invalidateQueries({ queryKey: ["listVideosInVideoGroup", id] });
      queryClient.invalidateQueries({ queryKey: ["listVideoGroups"] });
      queryClient.invalidateQueries({ queryKey: ["getVideoGroup", id] });
    },
  });
}

/** `POST /api/admin/video-groups` — create a group. Throws {@link ApiError} on failure (see `~/lib/api`). */
export function useCreateVideoGroupMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (body: VideoGroupInput) => {
      const { data } = await api.POST("/video-groups", { body });
      return data!;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["listVideoGroups"] });
    },
  });
}

interface UpdateVideoGroupMutationParams {
  id: UUID;
  params: PatchVideoGroupInput;
}

/**
 * `PATCH /api/admin/video-groups/{id}` — selective update. Pricing a group is
 * what puts it behind the paywall. Throws {@link ApiError} on failure (see `~/lib/api`).
 */
export function useUpdateVideoGroupMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({ id, params }: UpdateVideoGroupMutationParams) => {
      const { data } = await api.PATCH("/video-groups/{id}", {
        params: { path: { id } },
        body: params,
      });
      return data!;
    },
    onSuccess(group) {
      queryClient.invalidateQueries({ queryKey: ["listVideoGroups"] });
      queryClient.invalidateQueries({ queryKey: ["getVideoGroup", group.id] });
    },
  });
}

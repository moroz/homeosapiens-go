import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

import { api } from "~/lib/api";
import type { components } from "~/lib/api-types";
import type { UUID } from "~/lib/interfaces";

export type Video = components["schemas"]["Video"];
export type VideoDetails = components["schemas"]["VideoDetails"];
export type UpdateVideoInput = components["schemas"]["UpdateVideoInput"];
export type VideoProvider = components["schemas"]["VideoProvider"];

/** `GET /api/admin/videos` — a page of every video, newest first. */
export function useListVideosQuery(page = 1, perPage = 100) {
  return useQuery({
    queryKey: ["listVideos", page, perPage],
    queryFn: async () => {
      const { data } = await api.GET("/videos", { params: { query: { page, perPage } } });
      return data;
    },
  });
}

/** `GET /api/admin/videos/{id}` — a single video, with everything the edit form needs. */
export function useGetVideoQuery(id: string | undefined) {
  return useQuery({
    queryKey: ["getVideo", id],
    enabled: id != null,
    queryFn: async () => {
      const { data } = await api.GET("/videos/{id}", { params: { path: { id: id! } } });
      return data;
    },
  });
}

interface UpdateVideoMutationParams {
  id: UUID;
  params: UpdateVideoInput;
}

/**
 * `PATCH /api/admin/videos/{id}` — replace the editable fields of a video.
 * Throws {@link ApiError} on failure (see `~/lib/api`).
 */
export function useUpdateVideoMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({ id, params }: UpdateVideoMutationParams) => {
      const { data } = await api.PATCH("/videos/{id}", {
        params: { path: { id } },
        body: params,
      });
      return data!;
    },
    onSuccess(video) {
      queryClient.invalidateQueries({ queryKey: ["listVideos"] });
      queryClient.invalidateQueries({ queryKey: ["getVideo", video.id] });
      // A retitled video shows up inside its playlists, too.
      queryClient.invalidateQueries({ queryKey: ["listVideosInVideoGroup"] });
    },
  });
}

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { api } from "~/lib/api";
import type { components } from "~/lib/api-types";
import type { UUID } from "~/lib/interfaces";

export type BlogPost = components["schemas"]["BlogPost"];
export type CreateBlogPostInput = components["schemas"]["CreateBlogPostInput"];
export type UpdateBlogPostInput = components["schemas"]["UpdateBlogPostInput"];

export function useListBlogPostsQuery() {
  return useQuery({
    queryKey: ["listBlogPosts"],
    queryFn: async () => {
      const { data } = await api.GET("/blog-posts");
      return data;
    },
  });
}

export function useCreateBlogPostMutation() {
  return useMutation({
    mutationFn: async (params: CreateBlogPostInput) => {
      const { data } = await api.POST("/blog-posts", {
        body: params,
      });
      return data!;
    },
  });
}

export function useGetBlogPostQuery(id: UUID) {
  return useQuery({
    queryKey: ["getBlogPost", id],
    queryFn: async () => {
      const { data } = await api.GET("/blog-posts/{id}", {
        params: { path: { id } },
      });
      return data;
    },
  });
}

interface UpdateBlogPostMutationParams {
  id: UUID;
  params: UpdateBlogPostInput;
}

export function useUpdateBlogPostMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({ id, params }: UpdateBlogPostMutationParams) => {
      const { data } = await api.PATCH("/blog-posts/{id}", {
        params: { path: { id } },
        body: params,
      });
      return data!;
    },
    onSuccess(blogPost) {
      queryClient.invalidateQueries({ queryKey: ["listBlogPosts"] });
      queryClient.invalidateQueries({ queryKey: ["getBlogPost", blogPost.id] });
    },
  });
}

/** `POST /api/admin/blog-posts/{id}/publish` */
export function usePublishBlogPostMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (id: UUID) => {
      await api.POST("/blog-posts/{id}/publish", { params: { path: { id } } });
    },
    onSuccess(_data, id) {
      queryClient.invalidateQueries({ queryKey: ["listBlogPosts"] });
      queryClient.invalidateQueries({ queryKey: ["getBlogPost", id] });
    },
  });
}

/** `POST /api/admin/blog-posts/{id}/unpublish` */
export function useUnpublishBlogPostMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (id: UUID) => {
      await api.POST("/blog-posts/{id}/unpublish", { params: { path: { id } } });
    },
    onSuccess(_data, id) {
      queryClient.invalidateQueries({ queryKey: ["listBlogPosts"] });
      queryClient.invalidateQueries({ queryKey: ["getBlogPost", id] });
    },
  });
}

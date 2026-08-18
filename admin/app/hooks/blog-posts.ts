import { useMutation, useQuery } from "@tanstack/react-query";
import { api } from "~/lib/api";
import type { components } from "~/lib/api-types";
import type { UUID } from "~/lib/interfaces";

export type BlogPost = components["schemas"]["BlogPost"];
export type CreateBlogPostInput = components["schemas"]["CreateBlogPostInput"];

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

import { useQuery } from "@tanstack/react-query";
import { api } from "~/lib/api";
import type { components } from "~/lib/api-types";

export type BlogPost = components["schemas"]["BlogPost"];

export function useListBlogPostsQuery() {
  return useQuery({
    queryKey: ["listBlogPosts"],
    queryFn: async () => {
      const { data } = await api.GET("/blog-posts");
      return data;
    },
  });
}

import { useQuery } from "@tanstack/react-query";
import { api } from "~/lib/api";
import type { components } from "~/lib/api-types";
import type { UUID } from "~/lib/interfaces";

export type User = components["schemas"]["User"];

interface ListUsersQueryParams {
  page?: number;
  perPage?: number;
  searchTerm?: string;
}

export function useListUsersQuery(params: ListUsersQueryParams = {}) {
  const { page, perPage, searchTerm } = params;

  return useQuery({
    queryFn: async () => {
      const { data } = await api.GET("/users", {
        params: { query: { page, perPage, search: searchTerm } },
      });
      return data;
    },
    queryKey: ["listUsers", { page, perPage, searchTerm }],
  });
}

export function useGetUserQuery(id: UUID) {
  return useQuery({
    queryKey: ["getUser", id],
    queryFn: async () => {
      const { data } = await api.GET("/users/{id}", { params: { path: { id } } });
      return data;
    },
    enabled: Boolean(id),
  });
}

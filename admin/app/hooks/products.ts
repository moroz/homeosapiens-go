import { useQuery } from "@tanstack/react-query";

import { api } from "~/lib/api";
import type { components } from "~/lib/api-types";

export type Product = components["schemas"]["Product"];
export type ProductType = components["schemas"]["ProductType"];

/** `GET /api/admin/products` — a page of products, newest first. */
export function useListProductsQuery(page = 1, perPage = 20) {
  return useQuery({
    queryKey: ["listProducts", page, perPage],
    queryFn: async () => {
      const { data } = await api.GET("/products", { params: { query: { page, perPage } } });
      return data;
    },
  });
}

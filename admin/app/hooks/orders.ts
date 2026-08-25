import { useQuery } from "@tanstack/react-query";

import { api } from "~/lib/api";
import type { components } from "~/lib/api-types";

export type Order = components["schemas"]["Order"];
export type OrderDetails = components["schemas"]["OrderDetails"];
export type OrderLineItem = components["schemas"]["OrderLineItem"];
export type OrderStatus = components["schemas"]["OrderStatus"];

/** `GET /api/admin/orders` — a page of orders, newest first. */
export function useListOrdersQuery(page = 1, perPage = 20) {
  return useQuery({
    queryKey: ["listOrders", page, perPage],
    queryFn: async () => {
      const { data } = await api.GET("/orders", { params: { query: { page, perPage } } });
      return data;
    },
  });
}

/** `GET /api/admin/orders/{id}` — one order with its line items and billing address. */
export function useGetOrderQuery(id: string | undefined) {
  return useQuery({
    queryKey: ["getOrder", id],
    enabled: id != null,
    queryFn: async () => {
      const { data } = await api.GET("/orders/{id}", { params: { path: { id: id! } } });
      return data;
    },
  });
}

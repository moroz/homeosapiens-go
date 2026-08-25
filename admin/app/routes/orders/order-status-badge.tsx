import { Badge } from "~/components/ui/badge";
import type { OrderStatus } from "~/hooks";

const STATUS_BADGES: Record<
  OrderStatus,
  { label: string; variant: "default" | "secondary" | "destructive" }
> = {
  paid: { label: "paid", variant: "default" },
  pending: { label: "pending", variant: "secondary" },
  cancelled: { label: "cancelled", variant: "destructive" },
};

interface Props {
  status: OrderStatus;
}

export function OrderStatusBadge({ status }: Props) {
  const { label, variant } = STATUS_BADGES[status];
  return <Badge variant={variant}>{label}</Badge>;
}

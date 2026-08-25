import React from "react";
import { Link, useParams } from "react-router";

import { AdminLayout } from "~/components/admin-layout";
import { BackButton } from "~/components/back-button";
import { PageTitle } from "~/components/page-title";
import { DetailsTable, DetailsTableField as Field } from "~/components/ui/details-table";
import {
  Table,
  TableBody,
  TableCell,
  TableFooter,
  TableHead,
  TableHeader,
  TableRow,
} from "~/components/ui/table";
import { useGetOrderQuery, useGetUserQuery } from "~/hooks";
import { formatPrice } from "~/lib/money";
import { formatInstant } from "~/lib/time";

import { OrderStatusBadge } from "./order-status-badge";
import { buttonVariants } from "~/components/ui/button";
import { cn } from "~/lib/utils";

interface Props {}

export const OrderDetails: React.FC<Props> = () => {
  const { id } = useParams();
  const { data: order, isPending: orderPending, isError } = useGetOrderQuery(id!);
  const { data: user, isPending: userPending } = useGetUserQuery(order?.userId ?? "");

  const dataLoading = userPending || orderPending;

  return (
    <AdminLayout title="Order details">
      <BackButton href="/orders">Back to list</BackButton>

      {dataLoading ? (
        <p className="text-muted-foreground">Loading&hellip;</p>
      ) : isError || !order ? (
        <p>Not found</p>
      ) : (
        <div className="grid gap-6">
          <PageTitle subtitle={`${order.givenName} ${order.familyName}`}>
            Order #{order.orderNumber}
          </PageTitle>

          <DetailsTable>
            <Field label="ID" copy monospace>
              {order.id}
            </Field>
            <Field label="Status">
              <OrderStatusBadge status={order.status} />
            </Field>
            <Field label="Grand total">{formatPrice(order.grandTotal, order.currency)}</Field>
            <Field label="Placed at">{formatInstant(order.insertedAt)}</Field>
            <Field label="Paid at">{order.paidAt ? formatInstant(order.paidAt) : ""}</Field>
            <Field label="Cancelled at">
              {order.cancelledAt ? formatInstant(order.cancelledAt) : ""}
            </Field>
            <Field label="Checkout locale">{order.preferredLocale}</Field>
            <Field label="Stripe session" copy monospace>
              {order.stripeCheckoutSessionId}
            </Field>
            <Field label="User">
              <Link to={`/users/${order.userId}`} className="link">
                {user?.givenName} {user?.familyName} &lt;{user?.email}&gt;
              </Link>
            </Field>
          </DetailsTable>

          <section className="grid gap-2">
            <h3 className="text-lg font-semibold">Billing details</h3>
            <DetailsTable>
              <Field label="Given name">{order.givenName}</Field>
              <Field label="Family name">{order.familyName}</Field>
              <Field label="Email" copy>
                {order.email}
              </Field>
              <Field label="Phone">{order.phone}</Field>
              <Field label="Address line 1">{order.addressLine1}</Field>
              <Field label="Address line 2">{order.addressLine2}</Field>
              <Field label="City">{order.city}</Field>
              <Field label="Postal code">{order.postalCode}</Field>
              <Field label="Country">{order.billingCountry}</Field>
              <Field label="Tax ID">{order.taxId}</Field>
            </DetailsTable>
          </section>

          <section className="grid gap-2">
            <h3 className="text-lg font-semibold">Line items</h3>
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Product</TableHead>
                  <TableHead className="w-24 text-right">Quantity</TableHead>
                  <TableHead className="w-32 text-right">Unit price</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {order.lineItems.map((item) => (
                  <TableRow key={item.id}>
                    <TableCell>{item.productTitle}</TableCell>
                    <TableCell className="text-right">{item.quantity}</TableCell>
                    <TableCell className="text-right">
                      {formatPrice(item.productPrice, item.productPriceCurrency)}
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
              <TableFooter>
                <TableRow>
                  <TableCell colSpan={2}>Grand total</TableCell>
                  <TableCell className="text-right font-medium">
                    {formatPrice(order.grandTotal, order.currency)}
                  </TableCell>
                </TableRow>
              </TableFooter>
            </Table>
          </section>
        </div>
      )}
    </AdminLayout>
  );
};

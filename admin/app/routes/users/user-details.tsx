import React from "react";
import { AdminLayout } from "~/components/admin-layout";
import { PageTitle } from "~/components/page-title";
import { Link, useLocation, useParams } from "react-router";
import { useGetUserQuery } from "~/hooks";
import { BackButton } from "~/components/back-button";
import {
  DetailsTable,
  DetailsTableField as Field,
} from "~/components/ui/details-table";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "~/components/ui/table";
import { formatInstant } from "~/lib/time";

interface Props {}

export const UserDetails: React.FC<Props> = () => {
  const { id } = useParams();
  const { search } = useLocation();
  const { data: user, isPending, isError } = useGetUserQuery(id!);

  return (
    <AdminLayout title="User details" searchFormAction="/users">
      <BackButton href={`/users${search}`}>Back to list</BackButton>

      {isPending ? (
        <p className="text-muted-foreground">Loading&hellip;</p>
      ) : isError || !user ? (
        <p>Not found</p>
      ) : (
        <div className="grid gap-6">
          <PageTitle subtitle="User details">
            {user.givenName} {user.familyName}
          </PageTitle>
          <DetailsTable>
            <Field label="ID" copy monospace>
              {user.id}
            </Field>
            <Field label="Given name">{user.givenName}</Field>
            <Field label="Family name">{user.familyName}</Field>
            <Field label="Email" copy monospace>
              {user.email}
            </Field>
            <Field label="Account created at">
              {formatInstant(user.insertedAt)}
            </Field>
            <Field label="Email verified at">
              {user.emailConfirmedAt
                ? formatInstant(user.emailConfirmedAt)
                : "Not verified"}
            </Field>
          </DetailsTable>

          <section className="grid gap-2">
            <h3 className="text-lg font-semibold">Product access</h3>
            {user.productAccess.length === 0 ? (
              <p className="text-muted-foreground">No products.</p>
            ) : (
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Product</TableHead>
                    <TableHead className="w-32">Type</TableHead>
                    <TableHead className="w-40">Granted by</TableHead>
                    <TableHead className="w-40">Granted at</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {user.productAccess.map((access) => (
                    <TableRow key={access.productId}>
                      <TableCell>{access.titleEn}</TableCell>
                      <TableCell>{access.productType}</TableCell>
                      <TableCell>
                        {access.orderId ? (
                          <Link
                            to={`/orders/${access.orderId}`}
                            className="link"
                          >
                            Order
                          </Link>
                        ) : access.grantedByUserId ? (
                          <Link
                            to={`/users/${access.grantedByUserId}`}
                            className="link"
                          >
                            {access.grantedByName}
                          </Link>
                        ) : (
                          "Imported from legacy system"
                        )}
                      </TableCell>
                      <TableCell>{formatInstant(access.grantedAt)}</TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            )}
          </section>

          <section className="grid gap-2">
            <h3 className="text-lg font-semibold">Event registrations</h3>
            {user.eventRegistrations.length === 0 ? (
              <p className="text-muted-foreground">No registrations.</p>
            ) : (
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Event</TableHead>
                    <TableHead className="w-40">Starts at</TableHead>
                    <TableHead className="w-40">Registered at</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {user.eventRegistrations.map((registration) => (
                    <TableRow key={registration.eventId}>
                      <TableCell>
                        <Link
                          to={`/events/${registration.eventId}`}
                          className="link"
                        >
                          {registration.titleEn}
                        </Link>
                      </TableCell>
                      <TableCell>
                        {formatInstant(registration.startsAt)}
                      </TableCell>
                      <TableCell>
                        {formatInstant(registration.registeredAt)}
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            )}
          </section>
        </div>
      )}
    </AdminLayout>
  );
};

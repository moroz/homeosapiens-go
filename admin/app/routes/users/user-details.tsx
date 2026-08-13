import React from "react";
import { AdminLayout } from "~/components/admin-layout";
import { PageTitle } from "~/components/page-title";
import { useParams } from "react-router";
import { useGetUserQuery } from "~/hooks";
import { BackButton } from "~/components/back-button";
import { DetailsTable, DetailsTableField as Field } from "~/components/ui/details-table";
import { formatInstant } from "~/lib/time";

interface Props {}

export const UserDetails: React.FC<Props> = () => {
  const { id } = useParams();
  const { data: user, isPending, isError } = useGetUserQuery(id!);

  return (
    <AdminLayout title="User details">
      <BackButton href="/users">Back to list</BackButton>

      {isPending ? (
        <p className="text-muted-foreground">Loading&hellip;</p>
      ) : isError || !user ? (
        <p>Not found</p>
      ) : (
        <div className="grid gap-4">
          <PageTitle subtitle="User details">
            {user.givenName} {user.familyName}
          </PageTitle>
          <DetailsTable>
            <Field label="ID" copy>
              {user.id}
            </Field>
            <Field label="Given name">{user.givenName}</Field>
            <Field label="Family name">{user.familyName}</Field>
            <Field label="Email" copy>
              {user.email}
            </Field>
            <Field label="Account created at">{formatInstant(user.insertedAt)}</Field>
            <Field label="Email verified at">
              {user.emailConfirmedAt ? formatInstant(user.emailConfirmedAt) : "Not verified"}
            </Field>
          </DetailsTable>
        </div>
      )}
    </AdminLayout>
  );
};

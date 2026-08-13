import React from "react";
import { AdminLayout } from "~/components/admin-layout";
import { PageTitle } from "~/components/page-title";
import { useParams } from "react-router";
import { useGetUserQuery } from "~/hooks";

interface Props {}

export const UserDetails: React.FC<Props> = () => {
  const { id } = useParams();
  const { data } = useGetUserQuery(id!);

  return (
    <AdminLayout title="User details">
      <PageTitle subtitle="User details">
        {data ? `${data.givenName} ${data.familyName}` : ""}
      </PageTitle>
    </AdminLayout>
  );
};

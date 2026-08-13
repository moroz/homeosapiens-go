import React from "react";
import { AdminLayout } from "~/components/admin-layout";
import { PageTitle } from "~/components/page-title";
import { useParams } from "react-router";

interface Props {}

export const UserDetails: React.FC<Props> = () => {
  const { id } = useParams();

  return (
    <AdminLayout title="User details">
      <PageTitle subtitle="User details"></PageTitle>
    </AdminLayout>
  );
};

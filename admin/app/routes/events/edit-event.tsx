import React from "react";
import { AdminLayout } from "~/components/admin-layout";
import { useParams } from "react-router";
import { useGetEventQuery } from "~/hooks";

interface Props {}

export const EditEvent: React.FC<Props> = () => {
  const { id } = useParams();
  const { data: event, isPending, isError } = useGetEventQuery(id);

  return <AdminLayout>Edit event</AdminLayout>;
};

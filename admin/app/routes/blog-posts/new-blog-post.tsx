import React from "react";
import { AdminLayout } from "~/components/admin-layout";
import { useForm } from "react-hook-form";

interface Props {}

export const NewBlogPost: React.FC<Props> = () => {
  const form = useForm();

  return (
    <AdminLayout title="New blog post">
      <form></form>
    </AdminLayout>
  );
};

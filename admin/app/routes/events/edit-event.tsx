import React, { useEffect } from "react";
import { AdminLayout } from "~/components/admin-layout";
import { useParams } from "react-router";
import { useGetEventQuery } from "~/hooks";
import { FormProvider, useForm } from "react-hook-form";
import type { EventFormValues } from "./interfaces";
import { FormFields } from "./form-fields";

interface Props {}

export const EditEvent: React.FC<Props> = () => {
  const { id } = useParams();
  const { data: event, isPending, isError } = useGetEventQuery(id);
  const form = useForm<EventFormValues>({ defaultValues: { eventType: "" } });

  useEffect(() => {
    if (isPending || !event) return;

    form.reset({
      ...event,
      subtitleEn: event.subtitleEn ?? undefined,
      subtitlePl: event.subtitlePl ?? undefined,
      price: event.price ?? undefined,
      currency: event.currency ?? undefined,
    });
  }, [event, isPending]);

  return (
    <AdminLayout title="Edit event">
      <FormProvider {...form}>
        <FormFields />
      </FormProvider>
    </AdminLayout>
  );
};

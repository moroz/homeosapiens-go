import React, { useEffect } from "react";
import { AdminLayout } from "~/components/admin-layout";
import { useParams } from "react-router";
import { useGetEventQuery } from "~/hooks";
import { FormProvider, useForm } from "react-hook-form";
import type { EventFormValues } from "./interfaces";
import { FormFields } from "./form-fields";
import { ISO8601ToDatetimeLocalValue } from "~/lib/time";
import { PageTitle } from "~/components/page-title";

interface Props {}

export const EditEvent: React.FC<Props> = () => {
  const { id } = useParams();
  const { data: event, isPending, isError } = useGetEventQuery(id);
  const form = useForm<EventFormValues>({
    defaultValues: { eventType: "", isFree: false, isVirtual: false, currency: "", hostIds: [] },
  });

  useEffect(() => {
    if (isPending || !event) return;

    form.reset({
      ...event,
      subtitleEn: event.subtitleEn ?? undefined,
      subtitlePl: event.subtitlePl ?? undefined,
      price: event.price ?? undefined,
      currency: event.currency ?? undefined,
      startsAt: ISO8601ToDatetimeLocalValue(event.startsAt),
      endsAt: ISO8601ToDatetimeLocalValue(event.endsAt),
      venueNameEn: event.venueNameEn ?? undefined,
      venueNamePl: event.venueNamePl ?? undefined,
      venueStreet: event.venueStreet ?? undefined,
      venueCityEn: event.venueCityEn ?? undefined,
      venueCityPl: event.venueCityPl ?? undefined,
      venuePostalCode: event.venuePostalCode ?? undefined,
      venueCountryCode: event.venueCountryCode ?? undefined,
    });
  }, [event, isPending]);

  return (
    <AdminLayout title="Edit event">
      {isPending ? (
        <p className="text-muted-foreground">Loading&hellip;</p>
      ) : isError || !event ? (
        <p className="text-desctructive">Event not found.</p>
      ) : (
        <FormProvider {...form}>
          <PageTitle subtitle="Edit event">{event.titleEn}</PageTitle>
          <FormFields />
        </FormProvider>
      )}
    </AdminLayout>
  );
};

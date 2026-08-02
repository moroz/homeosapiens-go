import React, { useCallback, useEffect } from "react";
import { AdminLayout } from "~/components/admin-layout";
import { Outlet, useNavigate, useParams } from "react-router";
import { useGetEventQuery, useUpdateEventMutation } from "~/hooks";
import { FormProvider, type Path, useForm } from "react-hook-form";
import type { EventFormValues } from "./interfaces";
import { FormFields } from "./form-fields";
import { ISO8601ToDatetimeLocalValue } from "~/lib/time";
import { PageTitle } from "~/components/page-title";
import { ApiError, isValidationErrorBody } from "~/lib/api";
import { Button } from "~/components/ui/button";

interface Props {}

export const EditEvent: React.FC<Props> = () => {
  const { id } = useParams();
  const { data: event, isPending, isError } = useGetEventQuery(id);
  const form = useForm<EventFormValues>({
    defaultValues: { eventType: "", isFree: false, isVirtual: false, currency: "", hostIds: [] },
  });
  const mutation = useUpdateEventMutation();
  const navigate = useNavigate();

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
      hostIds: event.hosts.map(({ id }) => id),
    });
  }, [event, isPending]);

  const onSubmit = useCallback(
    async (params: EventFormValues) => {
      try {
        const data = await mutation.mutateAsync({ id: id!, params });
        navigate(`/events/${data.id}`);
      } catch (err) {
        if (err instanceof ApiError && err.status === 422 && isValidationErrorBody(err.body)) {
          for (const [field, message] of Object.entries(err.body.errors)) {
            form.setError(field as Path<EventFormValues>, { message });
          }
        }
        return;
      }
    },
    [mutation, id],
  );

  return (
    <AdminLayout title="Edit event">
      {isPending ? (
        <p className="text-muted-foreground">Loading&hellip;</p>
      ) : isError || !event ? (
        <p className="text-desctructive">Event not found.</p>
      ) : (
        <FormProvider {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)}>
            <PageTitle subtitle="Edit event">{event.titleEn}</PageTitle>
            <FormFields />
            <Outlet />
            <div className="flex gap-2">
              <Button type="submit" disabled={mutation.isPending}>
                {mutation.isPending ? "Updating…" : "Update event"}
              </Button>
            </div>
          </form>
        </FormProvider>
      )}
    </AdminLayout>
  );
};

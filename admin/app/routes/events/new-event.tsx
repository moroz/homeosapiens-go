import { useState } from "react";
import { Controller, FormProvider, useForm } from "react-hook-form";
import { useNavigate } from "react-router";

import { AdminLayout } from "~/components/admin-layout";
import { Button } from "~/components/ui/button";
import { Label } from "~/components/ui/label";
import { Select } from "~/components/forms";
import { Switch } from "~/components/ui/switch";
import { ApiError, isValidationErrorBody } from "~/lib/api";
import { useCreateEventMutation } from "~/hooks";
import { InputField } from "~/components/forms/input-field";
import { FieldError } from "~/components/forms/field-error";
import { type EventFormValues, toEventInput } from "./interfaces";
import { InputGroup } from "~/components/forms/input-group";
import { FormFields } from "./form-fields";
import { HostMultiSelect } from "./host-multi-select";

const defaultValues: Partial<EventFormValues> = {
  eventType: "webinar",
  isVirtual: true,
  isFree: true,
  currency: "PLN",
  hostIds: [],
};

/** Server field names line up 1:1 with `FormValues` keys, so validation errors map straight onto form fields. */
const FORM_FIELDS = new Set<string>(Object.keys(defaultValues));

function isFormField(field: string): field is keyof EventFormValues {
  return FORM_FIELDS.has(field);
}

export default function NewEvent() {
  const navigate = useNavigate();
  const createEvent = useCreateEventMutation();
  const [formError, setFormError] = useState<string | null>(null);

  const form = useForm<EventFormValues>({ defaultValues });

  const { handleSubmit, setError } = form;

  async function onSubmit(values: EventFormValues) {
    setFormError(null);
    try {
      const event = await createEvent.mutateAsync(toEventInput(values));
      navigate(`/events/${event.id}`);
    } catch (err) {
      if (err instanceof ApiError && err.status === 422 && isValidationErrorBody(err.body)) {
        for (const [field, message] of Object.entries(err.body.errors)) {
          if (isFormField(field)) {
            setError(field, { message });
          }
        }
        setFormError("Please fix the errors below.");
        return;
      }
      setFormError("Something went wrong creating the event. Please try again.");
    }
  }

  return (
    <AdminLayout title="New event">
      <FormProvider {...form}>
        <form onSubmit={handleSubmit(onSubmit)} className="flex max-w-2xl flex-col gap-6">
          <FieldError message={formError ?? undefined} />

          <FormFields />

          <div className="flex gap-2">
            <Button type="submit" disabled={createEvent.isPending}>
              {createEvent.isPending ? "Creating…" : "Create event"}
            </Button>
            <Button type="button" variant="ghost" onClick={() => navigate("/events")}>
              Cancel
            </Button>
          </div>
        </form>
      </FormProvider>
    </AdminLayout>
  );
}

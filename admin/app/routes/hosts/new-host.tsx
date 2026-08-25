import React, { useCallback, useState } from "react";
import { FormProvider, type Path, useForm } from "react-hook-form";
import { useNavigate } from "react-router";

import { AdminLayout } from "~/components/admin-layout";
import { BackButton } from "~/components/back-button";
import { Notification } from "~/components/notification";
import { PageTitle } from "~/components/page-title";
import { Button } from "~/components/ui/button";
import { useCreateHostMutation } from "~/hooks";
import { ApiError, isValidationErrorBody } from "~/lib/api";

import { FormFields } from "./form-fields";
import { toHostInput, type HostFormValues } from "./interfaces";

interface Props {}

export const NewHost: React.FC<Props> = () => {
  const form = useForm<HostFormValues>({
    defaultValues: { salutation: "", givenName: "", familyName: "", country: "" },
  });
  const mutation = useCreateHostMutation();
  const [formError, setFormError] = useState<string | null>(null);
  const navigate = useNavigate();

  const { handleSubmit, setError } = form;

  const onSubmit = useCallback(
    async (values: HostFormValues) => {
      setFormError(null);
      try {
        await mutation.mutateAsync(toHostInput(values));
        navigate("/hosts");
      } catch (err) {
        if (err instanceof ApiError && err.status === 422 && isValidationErrorBody(err.body)) {
          for (const [field, message] of Object.entries(err.body.errors)) {
            setError(field as Path<HostFormValues>, { message });
          }
          setFormError("Please fix the errors below.");
          return;
        }
        setFormError("Something went wrong creating the host. Please try again.");
      }
    },
    [mutation, navigate, setError],
  );

  return (
    <AdminLayout title="New host">
      <BackButton href="/hosts">Back to list</BackButton>
      <FormProvider {...form}>
        <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-4">
          <PageTitle>Create a host</PageTitle>
          {formError ? (
            <Notification
              title="An error has prevented this host from being saved."
              variant="destructive"
            >
              {formError}
            </Notification>
          ) : null}

          <FormFields />

          <div className="flex gap-2">
            <Button type="submit" disabled={mutation.isPending}>
              {mutation.isPending ? "Creating…" : "Create host"}
            </Button>
          </div>
        </form>
      </FormProvider>
    </AdminLayout>
  );
};

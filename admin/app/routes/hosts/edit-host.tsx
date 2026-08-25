import React, { useCallback, useEffect, useState } from "react";
import { FormProvider, type Path, useForm } from "react-hook-form";
import { useNavigate, useParams } from "react-router";

import { AdminLayout } from "~/components/admin-layout";
import { BackButton } from "~/components/back-button";
import { Notification } from "~/components/notification";
import { PageTitle } from "~/components/page-title";
import { Button } from "~/components/ui/button";
import { useDeleteHostMutation, useGetHostQuery, useUpdateHostMutation } from "~/hooks";
import { ApiError } from "~/lib/api";
import { isValidationErrorBody } from "~/lib/api";

import { FormFields } from "./form-fields";
import { toHostInput, type HostFormValues } from "./interfaces";

interface Props {}

export const EditHost: React.FC<Props> = () => {
  const { id } = useParams();
  const { data: host, isPending, isError } = useGetHostQuery(id!);
  const form = useForm<HostFormValues>({
    defaultValues: { salutation: "", givenName: "", familyName: "", country: "" },
  });
  const mutation = useUpdateHostMutation();
  const deleteMutation = useDeleteHostMutation();
  const [formError, setFormError] = useState<string | null>(null);
  const navigate = useNavigate();

  useEffect(() => {
    if (isPending || !host) return;

    form.reset({
      salutation: host.salutation ?? "",
      givenName: host.givenName,
      familyName: host.familyName,
      country: host.country ?? "",
    });
  }, [host, isPending, form]);

  const onSubmit = useCallback(
    async (values: HostFormValues) => {
      setFormError(null);
      try {
        await mutation.mutateAsync({ id: id!, params: toHostInput(values) });
        navigate("/hosts");
      } catch (err) {
        if (err instanceof ApiError && err.status === 422 && isValidationErrorBody(err.body)) {
          for (const [field, message] of Object.entries(err.body.errors)) {
            form.setError(field as Path<HostFormValues>, { message });
          }
          setFormError("Please fix the errors below.");
          return;
        }
        setFormError("Something went wrong updating the host. Please try again.");
      }
    },
    [mutation, id, navigate, form],
  );

  // A host that appears on an event or a video cannot be deleted; the server
  // answers 409 and the host stays as it is.
  const onDelete = useCallback(async () => {
    if (!window.confirm("Delete this host? This cannot be undone.")) return;

    setFormError(null);
    try {
      await deleteMutation.mutateAsync(id!);
      navigate("/hosts");
    } catch (err) {
      if (err instanceof ApiError && err.status === 409) {
        setFormError("This host is used by an event or a video and cannot be deleted.");
        return;
      }
      setFormError("Something went wrong deleting the host. Please try again.");
    }
  }, [deleteMutation, id, navigate]);

  return (
    <AdminLayout title="Edit host">
      <BackButton href="/hosts">Back to list</BackButton>
      {isPending ? (
        <p className="text-muted-foreground">Loading&hellip;</p>
      ) : isError || !host ? (
        <p className="text-destructive">Host not found.</p>
      ) : (
        <FormProvider {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="flex flex-col gap-4">
            <PageTitle subtitle="Edit host">
              {host.givenName} {host.familyName}
            </PageTitle>
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
                {mutation.isPending ? "Updating…" : "Update host"}
              </Button>
              <Button
                type="button"
                variant="destructive"
                onClick={onDelete}
                disabled={deleteMutation.isPending}
              >
                {deleteMutation.isPending ? "Deleting…" : "Delete host"}
              </Button>
            </div>
          </form>
        </FormProvider>
      )}
    </AdminLayout>
  );
};

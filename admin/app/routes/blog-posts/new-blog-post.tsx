import React, { useCallback, useState } from "react";
import { AdminLayout } from "~/components/admin-layout";
import { FormProvider, type Path, useForm } from "react-hook-form";
import { useCreateBlogPostMutation } from "~/hooks";
import { ApiError, isValidationErrorBody } from "~/lib/api";
import { PageTitle } from "~/components/page-title";
import { Notification } from "~/components/notification";
import { Button } from "~/components/ui/button";
import { slugify } from "~/lib/slugify";
import { useNavigate } from "react-router";
import { BackButton } from "~/components/back-button";
import { type BlogPostFormValues, toCreateBlogPostInput } from "./interfaces";
import { FormFields } from "./form-fields";

interface Props {}

export const NewBlogPost: React.FC<Props> = () => {
  const form = useForm<BlogPostFormValues>({
    defaultValues: { language: "en" },
  });
  const mutation = useCreateBlogPostMutation();
  const [formError, setFormError] = useState<string | null>(null);
  const navigate = useNavigate();

  const { handleSubmit, setError, setValue, getValues } = form;

  const onSubmit = useCallback(
    async (values: BlogPostFormValues) => {
      setFormError(null);
      try {
        const blogPost = await mutation.mutateAsync(toCreateBlogPostInput(values));
        navigate(`/blog-posts/${blogPost.id}`);
      } catch (err) {
        if (err instanceof ApiError && err.status === 422 && isValidationErrorBody(err.body)) {
          for (const [field, message] of Object.entries(err.body.errors)) {
            setError(field as Path<BlogPostFormValues>, { message });
          }
          setFormError("Please fix the errors below.");
          return;
        }
        setFormError("Something went wrong creating the blog post. Please try again.");
      }
    },
    [mutation, navigate, setError],
  );

  const onTitleBlur: React.ChangeEventHandler<HTMLInputElement> = useCallback(
    (e) => {
      const slug = getValues("slug");
      if (slug || !e.currentTarget.value) return;
      setValue("slug", slugify(e.currentTarget.value));
    },
    [getValues, setValue],
  );

  return (
    <AdminLayout title="New blog post">
      <BackButton href="/blog-posts">Back to list</BackButton>
      <FormProvider {...form}>
        <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-4">
          <PageTitle>Create a blog post</PageTitle>
          {formError ? (
            <Notification
              title="An error has prevented this blog post from being saved."
              variant="destructive"
            >
              {formError}
            </Notification>
          ) : null}

          <FormFields onTitleBlur={onTitleBlur} />

          <div className="flex gap-2">
            <Button type="submit" disabled={mutation.isPending}>
              {mutation.isPending ? "Creating…" : "Create blog post"}
            </Button>
          </div>
        </form>
      </FormProvider>
    </AdminLayout>
  );
};

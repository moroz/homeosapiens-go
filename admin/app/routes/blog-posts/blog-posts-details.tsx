import React from "react";
import { AdminLayout } from "~/components/admin-layout";
import { PageTitle } from "~/components/page-title";
import { useParams } from "react-router";
import { useGetBlogPostQuery } from "~/hooks";
import { DetailsTable, DetailsTableField as Field } from "~/components/ui/details-table";
import { formatInstant } from "~/lib/time";
import { BackButton } from "~/components/back-button";
import Markdown from "react-markdown";

interface Props {}

export const BlogPostDetails: React.FC<Props> = ({}: Props) => {
  const { id } = useParams();
  const { data: blogPost, isPending, isError } = useGetBlogPostQuery(id!);

  return (
    <AdminLayout title="Blog post details">
      <BackButton href="/blog-posts">Back to list</BackButton>
      {isPending ? (
        <p>Loading...</p>
      ) : isError || !blogPost ? (
        <p>No post found.</p>
      ) : (
        <div className="grid gap-4">
          <PageTitle subtitle="Blog post details">{blogPost.title}</PageTitle>
          <DetailsTable>
            <Field label="ID" copy>
              {blogPost.id}
            </Field>
            <Field label="Title">{blogPost.title}</Field>
            <Field label="Slug" className="font-mono" copy>
              {blogPost.slug}
            </Field>
            <Field label="Created at">{formatInstant(blogPost.insertedAt)}</Field>
            <Field label="Body">
              {blogPost.body ? (
                <Markdown>{blogPost.body}</Markdown>
              ) : (
                <span className="text-muted-foreground">(No data)</span>
              )}
            </Field>
          </DetailsTable>
        </div>
      )}
    </AdminLayout>
  );
};

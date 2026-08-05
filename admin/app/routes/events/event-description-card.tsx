import React from "react";
import { Card, CardAction, CardContent, CardHeader, CardTitle } from "~/components/ui/card";
import { Link } from "react-router";
import { buttonVariants } from "~/components/ui/button";
import { PencilIcon } from "@phosphor-icons/react/ssr";
import Markdown from "react-markdown";
import type { EventDetails } from "~/hooks";
import { Notification } from "~/components/notification";
import { cn } from "~/lib/utils";

interface Props {
  language: "pl" | "en";
  event: EventDetails;
}

const FIELD_BY_LOCALE = { en: "descriptionEn", pl: "descriptionPl" } as const;

export const EventDescriptionCard: React.FC<Props> = ({ event, language }) => {
  const description = event[FIELD_BY_LOCALE[language]] ?? "";

  return (
    <Card className="min-w-[70ch] flex-1">
      <CardHeader>
        <CardTitle>{new Intl.DisplayNames("en", { type: "language" }).of(language)}</CardTitle>
        <CardAction>
          <Link
            to={`description/${language}`}
            className={buttonVariants({ variant: "outline", size: "sm" })}
          >
            <PencilIcon className="w-4" />
            Edit
          </Link>
        </CardAction>
      </CardHeader>
      <CardContent
        className={cn(
          "prose text-foreground dark:prose-invert",
          description ? "max-w-[65ch]" : "w-full",
        )}
      >
        {description ? (
          <Markdown>{description}</Markdown>
        ) : (
          <span className="text-muted-foreground">[No description.]</span>
        )}
      </CardContent>
    </Card>
  );
};

import React from "react";
import { Badge } from "~/components/ui/badge";

interface Props {
  value?: any;
}

export const BooleanBadge: React.FC<Props> = ({ value }) => {
  return <Badge variant={value ? "default" : "secondary"}>{value ? "Yes" : "No"}</Badge>;
};

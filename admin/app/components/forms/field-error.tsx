import React from "react";

interface Props {
  message: string | null | undefined;
}

export const FieldError: React.FC<Props> = ({ message }) => {
  if (!message) return null;

  return <p className="text-sm text-destructive">{message}</p>;
};
